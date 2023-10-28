package auto_deployments_invoker

import (
	"deployer/commons/models/deployment"
	"deployer/commons/utils/aws"
	"deployer/commons/utils/gitwrapper"
	"fmt"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

type EnvironmentsConfigurations struct {
  deployments                  Deployments
  artifactReferenceReplacement ArtifactReferenceReplacement
  gitWrapper                   *gitwrapper.GitWrapper
}

func (ec *EnvironmentsConfigurations) Initialize() *EnvironmentsConfigurations {

  ec.artifactReferenceReplacement.InitializeForGitlabCiWorkflow()

  Logger.Info("Pulling credentials")
  deployment.PullCredentials()

  ec.gitWrapper = newGitWrapper()
  ec.gitWrapper.CleanClone()

  ec.deployments.pullDeploymentsConfigurationFromRedis()

  return ec
}

func newGitWrapper() *gitwrapper.GitWrapper {
  url, ok := aws.SecretsManager.GetSecret(os.Getenv(AWS_GIT_URL_WITH_CREDENTIALS_SECRETS_MANAGER_ARN))
  if !ok {
    Logger.Fatal("failed pulling { git URL with credentials } from AWS secrets manager", zap.String("arn", os.Getenv(AWS_GIT_URL_WITH_CREDENTIALS_SECRETS_MANAGER_ARN)))
  }

  return gitwrapper.New(url, filepath.FromSlash(fmt.Sprintf("%s/%s", os.TempDir(), "environmentsConfigurations")), true, Logger)
}

func (ec *EnvironmentsConfigurations) SubmitAutodeploy() (exitCode int) {

  exitCode = 1

  if ok := ec.deployments.submitAutodeploy(ec.artifactReferenceReplacement, ec.gitWrapper); ok {
    exitCode = 0
  }

  return
}

func NewEnvironmentsConfigurations() *EnvironmentsConfigurations {
  return &EnvironmentsConfigurations{
    deployments:                  Deployments{},
    artifactReferenceReplacement: ArtifactReferenceReplacement{},
  }
}
