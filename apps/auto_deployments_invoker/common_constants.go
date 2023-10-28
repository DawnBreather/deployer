package auto_deployments_invoker

import "regexp"

const (
	GITLAB_CI_PROJECT_NAME                           = "CI_PROJECT_NAME"
	GITLAB_CI_COMMIT_BRANCH                          = "CI_COMMIT_BRANCH"
	GITLAB_CI_COMMIT_SHORT_SHA                       = "CI_COMMIT_SHORT_SHA"
	ARTIFACT_OBJECT_REFERENCE_PREFIX                 = "ARTIFACT_OBJECT_REFERENCE_PREFIX"
	AWS_GIT_URL_WITH_CREDENTIALS_SECRETS_MANAGER_ARN = "AWS_GIT_URL_WITH_CREDENTIALS_SECRETS_MANAGER_ARN"
)

var patternForEnvironmentConfigurationKeyInRedis = regexp.MustCompile(`configuration/(.*?)/tiers/(.*)`)
