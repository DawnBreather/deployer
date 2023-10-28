package deployment

import (
  "fmt"
  "github.com/zabawaba99/firego"
)

const (
  // RDB stands for (Firebase) (R)eal-(T)ime (D)atabase
  RDB_ENVIRONMENT_PATH                       rdbPathTemplate = `environments/%s`
  RDB_ENVIRONMENT_SECRETS_PATH               rdbPathTemplate = `environments/%s/secrets`
  RDB_ENVIRONMENT_METADATA_PATH              rdbPathTemplate = `environments/%s/metadata`
  RDB_TIER_PATH                              rdbPathTemplate = `environments/%s/tiers/%s`
  RDB_INSTALLATION_STATUS_PATH               rdbPathTemplate = `environments/%s/status/%s/%s/installation`
  RDB_AGENT_HEARTBEAT_PATH                   rdbPathTemplate = `environments/%s/status/%s/%s/agent_heartbeat_at`
  RDB_LATEST_CONTROL_SEQUENCE_PATH           rdbPathTemplate = `environments/%v/tiers/%s/latest_control_sequence`
  RDB_LATEST_CONTROL_SEQUENCE_TIMESTAMP_PATH rdbPathTemplate = `environments/%s/tiers/%s/latest_control_sequence/created_at`
  RDB_LATEST_CONTROL_SEQUENCE_STATUS_PATH    rdbPathTemplate = `environments/%v/status/%s/%s/latest_control_sequence`
  RDB_ARTIFACTS_DEPLOY_STATUS_PATH           rdbPathTemplate = `environments/%s/status/%s/%s/artifacts_deploy/%d`
  RDB_TCP_HEALTH_CHECKS_STATUS_PATH          rdbPathTemplate = `environments/%s/status/%s/%s/health_checks/tcp/%d`
  RDB_UDP_HEALTH_CHECKS_STATUS_PATH          rdbPathTemplate = `environments/%s/status/%s/%s/health_checks/udp/%d`
  RDB_MSSQL_HEALTH_CHECKS_STATUS_PATH        rdbPathTemplate = `environments/%s/status/%s/%s/health_checks/mssql/%d`
  RDB_REDIS_HEALTH_CHECKS_STATUS_PATH        rdbPathTemplate = `environments/%s/status/%s/%s/health_checks/redis/%d`
  RDB_HTTP_HEALTH_CHECKS_STATUS_PATH         rdbPathTemplate = `environments/%s/status/%s/%s/health_checks/http/%d`
)

type rdbPathTemplate string

func (rpt rdbPathTemplate) getRefPath(args ...any) string {
  return fmt.Sprintf(string(rpt), args...)
}

func (rpt rdbPathTemplate) getRef(args ...any) (*firego.Firebase, error) {
  return f.Ref(rpt.getRefPath(args...))
}

func GetEnvironmentRef(environmentName string) (*firego.Firebase, error) {
  return RDB_ENVIRONMENT_PATH.getRef(environmentName)
}

func GetEnvironmentMetadataRef(environmentName string) (*firego.Firebase, error) {
  return RDB_ENVIRONMENT_METADATA_PATH.getRef(environmentName)
}

func GetEnvironmentSecretsRef(environmentName string) (*firego.Firebase, error) {
  return RDB_ENVIRONMENT_SECRETS_PATH.getRef(environmentName)
}

func GetTierRef(environmentName, tierName string) (*firego.Firebase, error) {
  return RDB_TIER_PATH.getRef(environmentName, tierName)
}

func GetLatestControlSequenceStatusRef(environmentName, tierName, agentName string) (*firego.Firebase, error) {
  return RDB_LATEST_CONTROL_SEQUENCE_STATUS_PATH.getRef(environmentName, tierName, agentName)
}

func GetInstallationStatusRef(environmentName, tierName, agentName string) (*firego.Firebase, error) {
  return RDB_INSTALLATION_STATUS_PATH.getRef(environmentName, tierName, agentName)
}

func GetAgentHeartBeatRef(environmentName, tierName, agentName string) (*firego.Firebase, error) {
  return RDB_AGENT_HEARTBEAT_PATH.getRef(environmentName, tierName, agentName)
}

func GetLatestControlSequenceRef(environmentName, tierName string) (*firego.Firebase, error) {
  return RDB_LATEST_CONTROL_SEQUENCE_PATH.getRef(environmentName, tierName)
}

func GetLatestControlSequenceTimestampRef(environmentName, tierName string) (*firego.Firebase, error) {
  return RDB_LATEST_CONTROL_SEQUENCE_TIMESTAMP_PATH.getRef(environmentName, tierName)
}

func GetArtifactsDeployStatusRef(environmentName, tierName, agentName string, artifactIndex int) (*firego.Firebase, error) {
  return RDB_ARTIFACTS_DEPLOY_STATUS_PATH.getRef(environmentName, tierName, agentName, artifactIndex)
}

func GetTCPHealthChecksStatusRef(environmentName, tierName, agentName string, tcpHealthCheckIndex int) (*firego.Firebase, error) {
  return RDB_TCP_HEALTH_CHECKS_STATUS_PATH.getRef(environmentName, tierName, agentName, tcpHealthCheckIndex)
}

func GetUDPHealthChecksStatusRef(environmentName, tierName, agentName string, udpHealthCheckIndex int) (*firego.Firebase, error) {
  return RDB_UDP_HEALTH_CHECKS_STATUS_PATH.getRef(environmentName, tierName, agentName, udpHealthCheckIndex)
}

func GetMSSQLHealthChecksStatusRef(environmentName, tierName, agentName string, mssqlHealthCheckIndex int) (*firego.Firebase, error) {
  return RDB_MSSQL_HEALTH_CHECKS_STATUS_PATH.getRef(environmentName, tierName, agentName, mssqlHealthCheckIndex)
}

func GetRedisHealthChecksStatusRef(environmentName, tierName, agentName string, redisHealthCheckIndex int) (*firego.Firebase, error) {
  return RDB_REDIS_HEALTH_CHECKS_STATUS_PATH.getRef(environmentName, tierName, agentName, redisHealthCheckIndex)
}

func GetHTTPHealthChecksStatusRef(environmentName, tierName, agentName string, httpHealthCheckIndex int) (*firego.Firebase, error) {
  return RDB_HTTP_HEALTH_CHECKS_STATUS_PATH.getRef(environmentName, tierName, agentName, httpHealthCheckIndex)
}
