package main

import (
	"deployer/apps/auto_deployments_invoker"
	"flag"
	"fmt"
	"go.uber.org/zap"
	"os"
	"regexp"
	"time"

	. "deployer/apps/deployments_status_tracker"
	"deployer/commons/models/deployment"
)

func main() {

	//TODO: pull credentials
	deployment.PullCredentials()

	environmentName := flag.String("environment", "some-environment.cohero-health.com", "name of the target environment")
	deploymentStatusReportJsonOutputFile := flag.String("deployment-status-report-json", "some-environment.cohero-health.com.deployment.json", "name of the destination file for storing deployment status")
	flag.Parse()

	Logger.Info(fmt.Sprintf("tracking { %s } environment", *environmentName))

	timeLimitInSeconds := 300
	timeCounter := 0
	interval := 1 * time.Second

	statusOk :=
		func() (statusOk bool) {
			for {
				Logger.Info(fmt.Sprintf("starting { %d } seconds tracking cycle", timeLimitInSeconds))
				//TODO: pull entries from Redis by environment name
				redisKeysMask := fmt.Sprintf("environments/%s/status/*/artifacts_deploy/*", *environmentName)
				Logger.Info(fmt.Sprintf("pulling redis keys by mask { %s }", redisKeysMask))
				deploymentsStatusesEntriesNames :=
					func() []string {
						keys, err := deployment.GetRedisKeysByMask(redisKeysMask)
						if err != nil {
							Logger.Fatal(fmt.Sprintf("getting keys from redis by mask: { %s }", err.Error()))
						}
						return keys
					}()

				Logger.Info(fmt.Sprintf("collecting statuses of the deployments corresponding to { %s } control sequence", auto_deployments_invoker.GenerateValueForCreatedAtPropertyForControlSequenceObjectInContextOfGitlabCiPipline()))
				//TODO: filter entries by control_sequence
				var deploymentsStatuses = func(entriesNames []string) []deployment.ArtifactDeployStatus {
					var deploymentsStatuses []deployment.ArtifactDeployStatus
					for _, entryName := range entriesNames {
						var deployStatus deployment.ArtifactDeployStatus
						err := deployment.GetJsonEntryFromRedis(entryName, &deployStatus)
						if err != nil {
							Logger.Error(fmt.Sprintf("getting value from Redis for entry { %s } failed", entryName), zap.Error(err))
						} else {
							if deployStatus.ControlSequence == auto_deployments_invoker.GenerateValueForCreatedAtPropertyForControlSequenceObjectInContextOfGitlabCiPipline() {
								deploymentsStatuses = append(deploymentsStatuses, deployStatus)
							}
						}
					}
					return deploymentsStatuses
				}(deploymentsStatusesEntriesNames)

				Logger.Info(fmt.Sprintf("collected { %d } statuses of the deployments corresponding to { %s } control sequence", len(deploymentsStatuses), auto_deployments_invoker.GenerateValueForCreatedAtPropertyForControlSequenceObjectInContextOfGitlabCiPipline()))

				Logger.Info(fmt.Sprintf("analyzing each deployment status for accomplishment"))
				//TODO: wait for 300 seconds whether the job finished (can be DONE or FAILED)
				var finished bool
				finished, statusOk =
					func(deploymentStatuses []deployment.ArtifactDeployStatus) (finished, ok bool) {
						finished = true
						for _, deploymentStatus := range deploymentStatuses {
							Logger.Info(fmt.Sprintf("deployment status: %+v", deploymentStatus))
							if deploymentStatus.State == deployment.DEPLOY_ARTIFACTS_STATE_IN_PROGRESS {
								finished = false
							} else {
								if deploymentStatus.State == deployment.DEPLOY_ARTIFACTS_STATE_DONE || deploymentStatus.State == deployment.DEPLOY_ARTIFACTS_STATE_FAILED {
									finished = true
									ok = deploymentStatus.State == deployment.DEPLOY_ARTIFACTS_STATE_DONE
								} else {
									if deploymentStatus.State == deployment.DEPLOY_ARTIFACTS_STATE_FAILED {
										Logger.Warn(fmt.Sprintf("failed deployment status identified: %v", deploymentStatus))
									}
								}
							}
						}
						return
					}(deploymentsStatuses)

				if finished {
					break
				}

				if timeCounter < timeLimitInSeconds {
					Logger.Info(fmt.Sprintf("Passed { %d } timeout seconds, moving to the next cycle, sleeping for { %.0f } seconds", timeCounter, interval.Seconds()))
					time.Sleep(interval)
					timeCounter++
				} else {
					break
				}
			}
			return statusOk
		}()

	//TODO: create status report JSON (status, version)
	err := os.WriteFile(*deploymentStatusReportJsonOutputFile, []byte(fmt.Sprintf(`{"status": "%s", "version": "%s"}`, func() string {
		if statusOk {
			return deployment.DEPLOY_ARTIFACTS_STATE_DONE
		} else {
			return deployment.DEPLOY_ARTIFACTS_STATE_FAILED
		}
	}(), GetSemanticVersionBasedOnGitBranchName(os.Getenv("CI_COMMIT_BRANCH")))), 0x644)
	if err != nil {
		Logger.Error(fmt.Sprintf("saving deployment status into { %s }", *deploymentStatusReportJsonOutputFile), zap.Error(err))
	}
}

func GetSemanticVersionBasedOnGitBranchName(inputString string) (version string) {
	semverPattern := `([0-9]+)(\.([0-9]+))*`
	Logger.Info(fmt.Sprintf("extracting semver by regexp pattern { %s } from branch name { %s }", semverPattern, inputString))
	re := regexp.MustCompile(semverPattern)
	matches := re.FindStringSubmatch(inputString)
	Logger.Info(fmt.Sprintf("extracted semver by regexp pattern { %s } from branch name { %s }: { %s }", semverPattern, inputString, func() string {
		if len(matches) > 0 {
			return matches[0]
		}
		return ""
	}()))
	if len(matches) > 0 {
		return matches[0]
	}
	return ""
}
