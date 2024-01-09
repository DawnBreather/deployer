package deployment

import (
  "deployer/commons/utils/aws"
  "encoding/json"
  "fmt"
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
  CONFIGURATION_STORAGE_GIT_CLONE_DESTINATION_PATH = filepath.FromSlash("./config")
  CONFIGURATION_STORAGE_TIERS_PATH                 = filepath.FromSlash(fmt.Sprintf("%s/%s/%s.yaml", CONFIGURATION_STORAGE_GIT_CLONE_DESTINATION_PATH, "tiers", strings.ReplaceAll(os.Getenv(ABZ_DEPLOYER_AGENT_TIER_ID), "::", ".")))
  CONFIGURATION_STORAGE_CONTROL_SEQUENCES_PATH     = filepath.FromSlash(fmt.Sprintf("%s/%s/%s.yaml", CONFIGURATION_STORAGE_GIT_CLONE_DESTINATION_PATH, "control_sequences", strings.ReplaceAll(os.Getenv(ABZ_DEPLOYER_AGENT_TIER_ID), "::", ".")))
  CONFIGURATION_STORAGE_SECRETS_PATH               = filepath.FromSlash(fmt.Sprintf("%s/%s", CONFIGURATION_STORAGE_GIT_CLONE_DESTINATION_PATH, "secrets.yaml"))
  CONFIGURATION_STORAGE_GIT_BRANCH_NAME            = fmt.Sprintf("deployer/%s", os.Getenv(ABZ_DEPLOYER_AGENT_ENVIRONMENT_ID))
  CONFIGURATION_STORAGE_GIT_PULL_PERIOD            = 30

  SSM_REFERENCE_FOR_REDIS_CREDENTIALS = "arn:aws:secretsmanager:us-east-1:010987917155:secret:redis"

  REDIS_CREDENTIALS = RedisCredentials{}
)

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

func PullCredentials() {
  if err := fetchAndUnmarshalCredentials(SSM_REFERENCE_FOR_REDIS_CREDENTIALS, &REDIS_CREDENTIALS); err != nil {
    logrus.Fatalf("[E] Pulling Redis credentials: %v", err)
  }
}

func fetchAndUnmarshalCredentials(ssmKey string, out any) error {
  value, err := fetchCredentialsFromSSM(ssmKey)
  if err != nil {
    return err
  }

  return json.Unmarshal([]byte(value), out)
}

func fetchCredentialsFromSSM(ssmKey string) (string, error) {
  value, ok := aws.SecretsManager.GetSecret(ssmKey)
  if !ok {
    return "", fmt.Errorf("credentials not found for key %s", ssmKey)
  }
  return value, nil
}
