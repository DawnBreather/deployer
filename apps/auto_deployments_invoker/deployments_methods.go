package auto_deployments_invoker

import (
	"bytes"
	"deployer/commons/models/deployment"
	"deployer/commons/utils/gitwrapper"
	"encoding/json"
	"flag"
	"fmt"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func (d Deployments) initializeDeploymentObjectsIfMissingInMap(deploymentName string) *deployment.Deployment {
	if _, ok := d[deploymentName]; !ok {
		d[deploymentName] = &deployment.Deployment{
			Metadata: deployment.Metadata{Name: deploymentName},
		}
	}

	return d[deploymentName]
}

func (d Deployments) pullDeploymentsConfigurationFromRedis() {
	var mask = "configuration/*/tiers/*"
	configurationKeys, err := deployment.GetRedisKeysByMask(mask)
	if err != nil {
		Logger.Error("Pulling { keys of configuration entries } by { mask } from redis failed", zap.NamedError("error", err), zap.Any("details", map[string]string{
			"mask": mask,
		}))
		return
	}

	tiersPerDeployment := map[string][]string{}
	for _, configurationKey := range configurationKeys {
		deploymentName, tierName := extractDeploymentAndTierNamesFromConfigurationKey(configurationKey)
		if _, ok := tiersPerDeployment[deploymentName]; !ok {
			tiersPerDeployment[deploymentName] = []string{}

		}

		tiersPerDeployment[deploymentName] = append(tiersPerDeployment[deploymentName], tierName)
	}

	for deploymentName := range tiersPerDeployment {
		Logger.Info("Pulling deployments configurations from Redis", zap.Any("details", map[string]string{
			"environment": deploymentName,
			"tiers":       strings.Join(tiersPerDeployment[deploymentName], ", "),
		}))
		d.initializeDeploymentObjectsIfMissingInMap(deploymentName).PullConfigurationFromRedisWithoutDecryptionForMultipleTiers(tiersPerDeployment[deploymentName])
	}

}

func (d Deployments) submitAutodeploy(artifactRefReplacement ArtifactReferenceReplacement, versionRefReplacement VersionReferenceReplacement, git *gitwrapper.GitWrapper) (anyOk bool) {

	resultingReport := map[string]bool{}

	for deploymentName, deploymentObject := range d {

		if shouldAdjustDeployment, tiersForAutodeploy := filterTiersApplicableForAutodeploy(deploymentObject); shouldAdjustDeployment {
			if resultingReport[deploymentName] = adjustDeployment(deploymentName, deploymentObject, artifactRefReplacement, versionRefReplacement, tiersForAutodeploy, git); resultingReport[deploymentName] {
				anyOk = true
			}
		}
	}

	Logger.Info("Automated deployment results", zap.Any("status", resultingReport))

	generateDynamicJobsForGitlabCiPipeline(resultingReport)

	return anyOk
}

func filterTiersApplicableForAutodeploy(deploymentObject *deployment.Deployment) (shouldAdjustDeployment bool, tiersForAutodeploy map[string]*deployment.Tier) {
	tiersForAutodeploy = map[string]*deployment.Tier{}

	for tierName, tierObject := range deploymentObject.Tiers {
		if tierObject.LatestControlSequence.AutoDeploy {
			shouldAdjustDeployment = true
			tiersForAutodeploy[tierName] = tierObject
		}
	}

	return
}

func adjustDeployment(deploymentName string, deploymentObject *deployment.Deployment, artifactRefReplacement ArtifactReferenceReplacement, versionRefReplacement VersionReferenceReplacement, tiersForAutodeploy map[string]*deployment.Tier, git *gitwrapper.GitWrapper) (ok bool) {

	secretsYamlBytes := getAdjustedSecrets(deploymentName, deploymentObject.Secrets, artifactRefReplacement, versionRefReplacement)
	if len(secretsYamlBytes) == 0 {
		return false
	}

	// TODO: publish secrets to Git repository
	if ok := commitSecretsToGitRepository(deploymentName, git, secretsYamlBytes); !ok {
		return false
	}

	//publishAdjustedSecretsToRedisChannel(deploymentName, deploymentObject, secretsYamlBytes)

	Logger.Info("Waiting for 3 seconds for secrets decryption on deployer agents")
	time.Sleep(3 * time.Second)

	for tierName, tierObject := range tiersForAutodeploy {

		controlSequenceJsonBytes := prepareNewControlSequence(tierObject, deploymentName, tierName)
		if len(controlSequenceJsonBytes) == 0 {
			continue
		}

		// TODO: save control sequence in Git
		commitControlSequenceToGit(deploymentName, git, tierName, controlSequenceJsonBytes)

		//publishControlSequenceToRedisChannel(deploymentName, tierName, deploymentObject, controlSequenceJsonBytes)
	}

	return git.PullAndPush()
}

func getAdjustedSecrets(deploymentName string, secrets deployment.Secrets, artifactRefReplacement ArtifactReferenceReplacement, versionRefReplacement VersionReferenceReplacement) (secretsYamlBytes []byte) {
	secretsYamlBytes, err := yaml.Marshal(secrets)
	if err != nil {
		Logger.Error("YAML marshalling of { secret } failed", zap.NamedError("error", err), zap.String("environment", deploymentName))
		return
	}
	Logger.Info("Replacing { old values } in { secrets YAML }", zap.Any("details", map[string]string{
		"environment": deploymentName,
	}))
	secretsYamlAdjustedBytes := versionRefReplacement.ReplaceAll(artifactRefReplacement.ReplaceAll(secretsYamlBytes))
	if err != nil {
		Logger.Error("Replacing { artifacts references } in { secret } failed", zap.NamedError("error", err), zap.String("environment", deploymentName))
		return
	}

	if !bytes.Equal(secretsYamlAdjustedBytes, secretsYamlBytes) {
		return secretsYamlAdjustedBytes
	} else {
		Logger.Warn("Replacing { artifacts references } in { secrets } warning", zap.NamedError("warning", fmt.Errorf("YAML remains the same after replacing artifacts references")), zap.String("environment", deploymentName))
	}

	return nil
}

func commitSecretsToGitRepository(deploymentName string, git *gitwrapper.GitWrapper, secretsYamlBytes []byte) (ok bool) {
	targetBranch := fmt.Sprintf("deployer/%s", deploymentName)
	Logger.Info("Publishing adjusted { secrets YAML } to git { repository }", zap.Any("details", map[string]string{
		"environment": deploymentName,
		"repository":  git.Repository(),
		"branch":      targetBranch,
	}))
	if git.RemoteBranchExists(targetBranch) {
		git.Checkout(fmt.Sprintf("deployer/%s", deploymentName))
		secretsFilePath, exists := git.YamlFileExists("secrets")
		if exists {
			if _, ok := git.ChangeFileAndCommit(secretsFilePath, string(secretsYamlBytes)); ok {
				//git.PullAndPush()
			} else {
				return false
			}
		} else {
			Logger.Error("adjusting deployment's secret file in git repository failed", zap.NamedError("err", fmt.Errorf("secrets file doesn't exist in the repository")), zap.String("environment", deploymentName))
			return false
		}
	} else {
		Logger.Error("adjusting deployment's secret file in git repository failed", zap.NamedError("err", fmt.Errorf("branch doesn't exist in the repository")), zap.String("environment", deploymentName))
		return false
	}

	return true
}

func publishAdjustedSecretsToRedisChannel(deploymentName string, deploymentObject *deployment.Deployment, secretsYamlBytes []byte) {
	Logger.Info("Publishing adjusted { secrets YAML } to redis { channel }", zap.Any("details", map[string]string{
		"environment": deploymentName,
		"channel":     deploymentObject.SecretsRedisChannel(),
	}))
	err := deployment.PublishMessageToRedisChannel(deploymentObject.SecretsRedisChannel(), string(secretsYamlBytes))
	if err != nil {
		Logger.Error("Publishing { secrets } to redis channel failed", zap.NamedError("error", err), zap.Any("details", map[string]string{
			"channel":     deploymentObject.SecretsRedisChannel(),
			"environment": deploymentName,
		}))
	}
}

func prepareNewControlSequence(tierObject *deployment.Tier, deploymentName string, tierName string) (controlSequenceJsonBytes []byte) {
	timestamp := deployment.GetTimeNowString()
	controlSequence := tierObject.LatestControlSequence
	//controlSequence.CreatedAt = timestamp
	// TODO: Getting branch name and commit short SHA hash
	controlSequence.CreatedAt = fmt.Sprintf("%s/%s", os.Getenv("CI_COMMIT_BRANCH"), os.Getenv("CI_COMMIT_SHORT_SHA"))

	Logger.Info("Creating { control sequence YAML } for publishing to redis channel", zap.Any("details", map[string]string{
		"environment": deploymentName,
		"tier":        tierName,
		"timestamp":   timestamp,
	}))

	controlSequenceJsonBytes, err := json.Marshal(controlSequence)
	if err != nil {
		Logger.Error("JSON marshalling of { control sequence } failed", zap.NamedError("error", err), zap.Any("details", map[string]string{
			"environment": deploymentName,
			"tier":        tierName,
		}))
	}

	return controlSequenceJsonBytes
}

func publishControlSequenceToRedisChannel(deploymentName string, tierName string, deploymentObject *deployment.Deployment, controlSequenceJsonBytes []byte) {
	Logger.Info("Publishing { control sequence YAML } to redis { channel }", zap.Any("details", map[string]string{
		"environment": deploymentName,
		"tier":        tierName,
		"channel":     deploymentObject.ControlSequenceRedisChannel(tierName),
	}))
	err := deployment.PublishMessageToRedisChannel(deploymentObject.ControlSequenceRedisChannel(tierName), string(controlSequenceJsonBytes))
	if err != nil {
		// TODO: log error
		Logger.Error("Publishing { control sequence } to redis channel failed", zap.NamedError("error", err), zap.Any("details", map[string]string{
			"channel":     deploymentObject.ControlSequenceRedisChannel(tierName),
			"environment": deploymentName,
			"tier":        tierName,
		}))
	}
}

func commitControlSequenceToGit(deploymentName string, git *gitwrapper.GitWrapper, tierName string, controlSequenceYamlBytes []byte) (ok bool) {
	targetBranch := fmt.Sprintf("deployer/%s", deploymentName)
	Logger.Info("Publishing adjusted { secrets YAML } to git { repository }", zap.Any("details", map[string]string{
		"environment": deploymentName,
		"repository":  git.Repository(),
		"branch":      targetBranch,
	}))
	if git.RemoteBranchExists(targetBranch) {
		git.Checkout(fmt.Sprintf("deployer/%s", deploymentName))
		controlSequenceFilePath, exists := git.YamlFileExists(fmt.Sprintf("control_sequences/%s", tierName))
		if exists {
			if _, ok := git.ChangeFileAndCommit(controlSequenceFilePath, string(controlSequenceYamlBytes)); !ok {
				return false
			}
		} else {
			if _, ok := git.ChangeFileAndCommit(fmt.Sprintf("control_sequences/%s.yaml", tierName), string(controlSequenceYamlBytes)); !ok {
				return false
			}
		}
		//git.PullAndPush()
	} else {
		Logger.Error("adjusting deployment's secret file in git repository failed", zap.NamedError("err", fmt.Errorf("secrets file doesn't exist in the repository")), zap.String("environment", deploymentName))
		return false
	}

	return true
}

//func exportDeploymentInitializationReportToFilesystem(report map[string]bool) {
//
//	finalReportFileName := "auto_deployment_initialization_report.json"
//
//	bytes, err := json.Marshal(report)
//	if err != nil {
//		Logger.Error("Marhsalling final report", zap.Error(err))
//	}
//	err = os.WriteFile(finalReportFileName, bytes, 0x644)
//	if err != nil {
//		Logger.Error(fmt.Sprintf("Saving { %s } final report", finalReportFileName), zap.Error(err))
//	}
//}

func generateDynamicJobsForGitlabCiPipeline(report map[string]bool) (dynamicJobsTemplate string) {

	outputFile := flag.String("dynamic-jobs-output-file", "dynamic_jobs.yml", "Destination file for generated dynamic jobs for Gitlab CI pipeline")
	isAqaJobsEnabled := flag.String("is-aqa-jobs-enabled", "true", "Enabled / disable post-deployment jobs for executing AQA API tests")
	flag.Parse()

	Logger.Info("[configuration] destination file into which we will store dynamically generated CI jobs definitions", zap.String("dynamic-jobs-output-file", *outputFile))
	Logger.Info("[configuration] whether we enable AQA jobs", zap.String("is-aqa-jobs-enabled", *isAqaJobsEnabled))

	Logger.Info("Printing current project working directory", zap.String("CI_PROJECT_DIR", os.Getenv("CI_PROJECT_DIR")))
	*outputFile = filepath.FromSlash(fmt.Sprintf("%s/%s", os.Getenv("CI_PROJECT_DIR"), *outputFile))

	var listOfEnvironments = []string{}
	var controlSequence = GenerateValueForCreatedAtPropertyForControlSequenceObjectInContextOfGitlabCiPipline()
	for environmentName, status := range report {
		if status {
			listOfEnvironments = append(listOfEnvironments, fmt.Sprintf("\"%s\"", environmentName))
		}
	}

	dynamicJobsTemplate = strings.ReplaceAll(DYNAMIC_JOBS_GITLAB_CI_TEMPLATE, "{{LIST_OF_APPLICABLE_ENVIRONMENTS}}", strings.Join(listOfEnvironments, ","))
	dynamicJobsTemplate = strings.ReplaceAll(dynamicJobsTemplate, "{{CONTROL_SEQUENCE}}", controlSequence)
	dynamicJobsTemplate = strings.ReplaceAll(dynamicJobsTemplate, "{{IS_AQA_JOBS_ENABLED}}", *isAqaJobsEnabled)

	err := os.WriteFile(*outputFile, []byte(dynamicJobsTemplate), 0x644)
	if err != nil {
		Logger.Error(fmt.Sprintf("writing dynamic jobs file { %s } for Gitlab CI pipeline", *outputFile), zap.NamedError("err", fmt.Errorf(err.Error())))
	} else {
		path, err := os.Getwd()
		if err != nil {
			Logger.Warn("getting current { working directory }", zap.Error(err))
		}
		Logger.Info(fmt.Sprintf("writing dynamic jobs file { %s } for Gitlab CI pipeline to { %s } location", *outputFile, path), zap.String("dynanmic job content", fmt.Sprintf("%s", dynamicJobsTemplate)))
	}

	return
}

const (
	DYNAMIC_JOBS_GITLAB_CI_TEMPLATE = `
variables:
  IS_AQA_JOBS_ENABLED: "{{IS_AQA_JOBS_ENABLED}}"
  

stages:
  - deployments.tracking
  - aqa.assets_pulling
  - aqa.api_tests_execution

deployments.track:
  image: public.ecr.aws/r8d1t0k0/aptar-digital-health/deployer:tracker-alpine
  stage: deployments.tracking
  script:
    - echo "deployer.deployments_status_tracker --environment=${ENVIRONMENT_NAME} --deployment-status-report-json=${ENVIRONMENT_NAME}.deployment.json"
    - deployer.deployments_status_tracker --environment=${ENVIRONMENT_NAME} --deployment-status-report-json=${ENVIRONMENT_NAME}.deployment.json
  tags:
    - deployer
  parallel:
    matrix:
      - ENVIRONMENT_NAME: [{{LIST_OF_APPLICABLE_ENVIRONMENTS}}]
  artifacts:
    untracked: true
    expire_in: 1 hour

aqa.pull_assets:
  image:
    name: bitnami/git
    entrypoint: [""]
  stage: aqa.assets_pulling
  script:
    - git clone ${AQA_REPOSITORY_URL_WITH_CREDENTIALS} assets
    - ls assets
  artifacts:
    paths:
      - ./
    expire_in: 1 hour
  only:
    variables:
      - $IS_AQA_JOBS_ENABLED == "true"

aqa.execute_api_tests:
  image:
    name: public.ecr.aws/r8d1t0k0/aptar-digital-health/deployer:aqa-alpine
    entrypoint: [""]
  stage: aqa.api_tests_execution
  script:
    - echo "${ENVIRONMENT_NAME}.deployment.json"
    - cat "${ENVIRONMENT_NAME}.deployment.json"
    - status=$(jq -r '.status' "${ENVIRONMENT_NAME}.deployment.json")
    - 'if [ "$status" != "DONE" ]; then exit 250; fi'
    - postman_collection_version=$(jq -r '.version' "${ENVIRONMENT_NAME}.deployment.json")
    - cd assets
    - echo "newman run ${postman_collection_version}.api.postman_collection.json -g workspace.postman_globals.json -e ${ENVIRONMENT_NAME}.postman_environment.json"
    - newman run ${postman_collection_version}.api.postman_collection.json -g workspace.postman_globals.json -e ${ENVIRONMENT_NAME}.postman_environment.json
		- echo "newman run ${postman_collection_version}.caregiver.postman_collection.json -g workspace.postman_globals.json -e ${ENVIRONMENT_NAME}.postman_environment.json"
    - newman run ${postman_collection_version}.caregiver.postman_collection.json -g workspace.postman_globals.json -e ${ENVIRONMENT_NAME}.postman_environment.json
		- echo "newman run ${postman_collection_version}.doctor.postman_collection.json -g workspace.postman_globals.json -e ${ENVIRONMENT_NAME}.postman_environment.json"
    - newman run ${postman_collection_version}.doctor.postman_collection.json -g workspace.postman_globals.json -e ${ENVIRONMENT_NAME}.postman_environment.json
		- echo "newman run ${postman_collection_version}.admin.postman_collection.json -g workspace.postman_globals.json -e ${ENVIRONMENT_NAME}.postman_environment.json"
    - newman run ${postman_collection_version}.admin.postman_collection.json -g workspace.postman_globals.json -e ${ENVIRONMENT_NAME}.postman_environment.json
		- echo "newman run ${postman_collection_version}.signalrcore.postman_collection.json -g workspace.postman_globals.json -e ${ENVIRONMENT_NAME}.postman_environment.json"
    - newman run ${postman_collection_version}.signalrcore.postman_collection.json -g workspace.postman_globals.json -e ${ENVIRONMENT_NAME}.postman_environment.json
  parallel:
    matrix:
      - ENVIRONMENT_NAME: [{{LIST_OF_APPLICABLE_ENVIRONMENTS}}]
  allow_failure:
    exit_codes:
      - 250
  only:
    variables:
      - $IS_AQA_JOBS_ENABLED == "true"
`
)
