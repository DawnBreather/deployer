package deployment

const (
	// RDB stands for (Firebase) (R)eal-(T)ime (D)atabase
	RDB_ENVIRONMENT_PATH                       rdbPathTemplate = `environments/%s`
	RDB_ENVIRONMENT_SECRETS_PATH               rdbPathTemplate = `environments/%s/secrets`
	RDB_ENVIRONMENT_METADATA_PATH              rdbPathTemplate = `environments/%s/metadata`
	RDB_TIER_PATH                              rdbPathTemplate = `environments/%s/tiers/%s`
	RDB_INSTALLATION_STATUS_PATH               rdbPathTemplate = `environments/%s/status/%s/%s/installation`
	RDB_AGENT_STATUS_PATH                      rdbPathTemplate = `environments/%s/status/%s/%s/status`
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
