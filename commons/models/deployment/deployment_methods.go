package deployment

import (
  "encoding/json"
  _ "github.com/denisenkom/go-mssqldb"
  "github.com/sirupsen/logrus"
  "github.com/tidwall/pretty"
  "io"
  "os"
  "os/exec"
  "strings"
  "time"
)

const (
  ABZ_DEPLOYER_AGENT_ENVIRONMENT_ID                   = "ABZ_DEPLOYER_AGENT_ENVIRONMENT_ID"
  ABZ_DEPLOYER_AGENT_TIER_ID                          = "ABZ_DEPLOYER_AGENT_TIER_ID"
  ABZ_DEPLOYER_AGENT_INSTALLATION_ENABLED             = "ABZ_DEPLOYER_AGENT_INSTALLATION_ENABLED"
  ABZ_DEPLOYER_AGENT_FIREBASE_AUTHENTICATION_JSON_URL = "ABZ_DEPLOYER_AGENT_FIREBASE_AUTHENTICATION_JSON_URL"
)

const (
  DEPLOY_ARTIFACTS_STATE_IN_PROGRESS = "IN PROGRESS"
  DEPLOY_ARTIFACTS_STATE_DONE        = "DONE"
  DEPLOY_ARTIFACTS_STATE_FAILED      = "FAILED"
)

const (
  ENUM_DEPLOY_MODE_DELETE = "delete"
)

const (
  COMMAND_SEQUENCE_DEPLOY             = "deploy"
  COMMAND_SEQUENCE_RESTART_ENTRYPOINT = "restart-entrypoint"
  COMMAND_SEQUENCE_STOP_ENTRYPOINT    = "stop-entrypoint"
  COMMAND_SEQUENCE_START_ENTRYPOINT   = "start-entrypoint"
)

const (
  SECRETS_PATH_AWS_REGION = `${aws.region}`
  SECRETS_PATH_AWS_TOKEN  = `${aws.token}`
  SECRETS_PATH_AWS_SECRET = `${aws.secret}`
)

var (
  agentId                                   = AgentName()
  IsTierConfigurationInitialized            = false
  cliOptionsForEnvsubstInitialized          = false
  PauseConfigurationListening               = false
  IsFirebaseAuthenticated                   = false
  IsConfigurationFromGitRepositoryRetrieved = false
)

var (
  entrypointCommand *exec.Cmd
  entrypointStdin   io.WriteCloser
)

func (d *Deployment) WaitForConfigurationFromGitInitialRetrieval() {
  for !IsConfigurationFromGitRepositoryRetrieved {
    logrus.Infof("[I] waiting for configuration from Git initial retrieval")
    time.Sleep(3000 * time.Millisecond)
  }
}

func (d *Deployment) WaitForFirebaseAuthentication() {
  for !IsFirebaseAuthenticated {
    logrus.Infof("[I] waiting for Firebase authentication")
    time.Sleep(3000 * time.Millisecond)
  }
}

func (d *Deployment) Name() string {
  if d.Metadata.Name == "" {
    return os.Getenv(ABZ_DEPLOYER_AGENT_ENVIRONMENT_ID)
  }
  return d.Metadata.Name
}

func (d *Deployment) tierName() string {
  var tierName = strings.ReplaceAll(os.Getenv(ABZ_DEPLOYER_AGENT_TIER_ID), "::", ".")
  if tierName == "" {
    logrus.Fatalf("[E] { ABZ_DEPLOYER_AGENT_TIER_ID } environment variable is not set. Please provide the Tier ID.")
  }
  return tierName
}

func (d *Deployment) Tier() *Tier {
  return d.Tiers[d.tierName()]
}

func (d *Deployment) ExportDeployment() {
  jsonStringBase64, err := json.Marshal(*d)
  if err != nil {
    logrus.Errorf("[E] marshalling { d } into JSON: %v", err)
  } else {
    pretty.Color(jsonStringBase64, nil)
  }
}
