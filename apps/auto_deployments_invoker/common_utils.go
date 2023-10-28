package auto_deployments_invoker

import (
	"fmt"
	"os"
)

func extractDeploymentAndTierNamesFromConfigurationKey(redisEntryKey string) (deploymentName, tierName string) {

	match := patternForEnvironmentConfigurationKeyInRedis.FindStringSubmatch(redisEntryKey)

	return match[1], match[2]
}

func GenerateValueForCreatedAtPropertyForControlSequenceObjectInContextOfGitlabCiPipline() string {
	return fmt.Sprintf("%s/%s", os.Getenv("CI_COMMIT_BRANCH"), os.Getenv("CI_COMMIT_SHORT_SHA"))
}
