package deployment

import (
	"encoding/json"
	"fmt"
	"github.com/DawnBreather/go-commons/aws"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
)

// TODO: print warning on the start of the application
// please store KMS keys in the region you set in the cont below
const KMS_KEY_REGION = "us-east-1"

var (
	//SSM_REFERENCE_FOR_GIT_CREDENTIALS   = "arn:aws:secretsmanager:us-east-1:010987917155:secret:git"
	SSM_REFERENCE_FOR_REDIS_CREDENTIALS = "arn:aws:secretsmanager:us-east-1:010987917155:secret:redis"
	//SSM_REFERENCE_FOR_FIREBASE_CREDENTIALS = "arn:aws:secretsmanager:us-east-1:010987917155:secret:firebase"
)

var (
	//GIT_CREDENTIALS   = GitLabCredentials{}
	REDIS_CREDENTIALS = RedisCredentials{}
)

func PullCredentials() {
	//getFromSsmAndUnmarshal(SSM_REFERENCE_FOR_GIT_CREDENTIALS, &GIT_CREDENTIALS)
	getFromSsmAndUnmarshal(SSM_REFERENCE_FOR_REDIS_CREDENTIALS, &REDIS_CREDENTIALS)
}

type GitLabCredentials struct {
	RepositoryHttpsUrl string `json:"repository_https_url"`
	Username           string `json:"username"`
	Token              string `json:"token"`
}

func (gc *GitLabCredentials) GetBasicAuth() *http.BasicAuth {
	return &http.BasicAuth{
		Username: gc.Username,
		Password: gc.Token,
	}
}

type RedisCredentials struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
}

func getFromSsmAndUnmarshal(key string, out any) (err error) {
	value, ok := aws.SecretsManager.GetSecret(key)
	if ok {
		err = json.Unmarshal([]byte(value), out)
	}
	if !ok && err != nil {
		logrus.Fatalf("[E] Extracting credentials from SSM { %s }: %v", key, err)
	}
	return nil
}

//var (
//	REDIS_HOST     = "deployer-redis.general.cohero-health.com"
//	REDIS_PORT     = "6379"
//	REDIS_PASSWORD = "Ooghee7diegohqua"
//)

var (
	//CONFIGURATION_STORAGE_GIT_URL                    = `https://gitlab.com/aptar-digital-health/web/devops_cicd.git`
	//CONFIGURATION_STORAGE_GIT_USERNAME               = "dawnbreather"
	//CONFIGURATION_STORAGE_GIT_TOKEN                  = "arn:aws:secretsmanager:us-east-1:010987917155:secret:git/token"
	CONFIGURATION_STORAGE_GIT_CLONE_DESTINATION_PATH = filepath.FromSlash("./config")
	//CONFIGURATION_STORAGE_GIT_TOKEN_VALUE, _         = aws.SecretsManager.GetSecret(CONFIGURATION_STORAGE_GIT_TOKEN)
	//CONFIGURATIONS_STORAGE_GIT_BASIC_AUTH_OBJECT = &http.BasicAuth{
	//	Username: CONFIGURATION_STORAGE_GIT_USERNAME,
	//	HubAdminPassword: CONFIGURATION_STORAGE_GIT_TOKEN_VALUE,
	//}
	CONFIGURATION_STORAGE_TIERS_PATH             = filepath.FromSlash(fmt.Sprintf("%s/%s/%s.yaml", CONFIGURATION_STORAGE_GIT_CLONE_DESTINATION_PATH, "tiers", strings.ReplaceAll(os.Getenv(ABZ_DEPLOYER_AGENT_TIER_ID), "::", ".")))
	CONFIGURATION_STORAGE_CONTROL_SEQUENCES_PATH = filepath.FromSlash(fmt.Sprintf("%s/%s/%s.yaml", CONFIGURATION_STORAGE_GIT_CLONE_DESTINATION_PATH, "control_sequences", strings.ReplaceAll(os.Getenv(ABZ_DEPLOYER_AGENT_TIER_ID), "::", ".")))
	CONFIGURATION_STORAGE_SECRETS_PATH           = filepath.FromSlash(fmt.Sprintf("%s/%s", CONFIGURATION_STORAGE_GIT_CLONE_DESTINATION_PATH, "secrets.yaml"))
	CONFIGURATION_STORAGE_GIT_BRANCH_NAME        = fmt.Sprintf("deployer/%s", os.Getenv(ABZ_DEPLOYER_AGENT_ENVIRONMENT_ID))
	CONFIGURATION_STORAGE_GIT_PULL_PERIOD        = 30
)
