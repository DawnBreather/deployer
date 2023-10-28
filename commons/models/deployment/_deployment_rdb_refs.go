package deployment

import (
  "github.com/sirupsen/logrus"
  "github.com/zabawaba99/firego"
  "strings"
)

func (d *Deployment) environmentRef() (res *firego.Firebase) {
  return handleRefError(GetEnvironmentRef(d.Name()))
}

func (d *Deployment) environmentMetadataRef() (res *firego.Firebase) {
  return handleRefError(GetEnvironmentMetadataRef(d.Name()))
}

func (d *Deployment) environmentSecretsRef() (res *firego.Firebase) {
  return handleRefError(GetEnvironmentSecretsRef(d.Name()))
}

func (d *Deployment) tierRef() (res *firego.Firebase) {
  return handleRefError(GetTierRef(d.Name(), d.tierName()))
}

func (d *Deployment) installationStatusRef() (res *firego.Firebase) {
  return handleRefError(GetInstallationStatusRef(d.Name(), d.tierName(), agentId))
}

func (d *Deployment) latestControlSequenceStatusRef() (res *firego.Firebase) {
  return handleRefError(GetLatestControlSequenceStatusRef(d.Name(), d.tierName(), agentId))
}

func (d *Deployment) agentHeartBeatRef() (res *firego.Firebase) {
  return handleRefError(GetAgentHeartBeatRef(d.Name(), d.tierName(), agentId))
}

func (d *Deployment) latestControlSequenceRef() (res *firego.Firebase) {
  return handleRefError(GetLatestControlSequenceRef(d.Name(), d.tierName()))
}

func (d *Deployment) latestControlSequenceTimestampRef() (res *firego.Firebase) {
  return handleRefError(GetLatestControlSequenceTimestampRef(d.Name(), d.tierName()))
}

func (d *Deployment) artifactsDeployStatusRef(index int) (res *firego.Firebase) {
  return handleRefError(GetArtifactsDeployStatusRef(d.Name(), d.tierName(), agentId, index))
}

func (d *Deployment) tcpHealthCheckStatus(index int) (res *firego.Firebase) {
  return handleRefError(GetTCPHealthChecksStatusRef(d.Name(), d.tierName(), agentId, index))
}

func (d *Deployment) udpHealthCheckStatus(index int) (res *firego.Firebase) {
  return handleRefError(GetUDPHealthChecksStatusRef(d.Name(), d.tierName(), agentId, index))
}

func (d *Deployment) mssqlHealthCheckStatus(index int) (res *firego.Firebase) {
  return handleRefError(GetMSSQLHealthChecksStatusRef(d.Name(), d.tierName(), agentId, index))
}

func (d *Deployment) redisHealthCheckStatus(index int) (res *firego.Firebase) {
  return handleRefError(GetRedisHealthChecksStatusRef(d.Name(), d.tierName(), agentId, index))
}

func (d *Deployment) httpHealthCheckStatus(index int) (res *firego.Firebase) {
  return handleRefError(GetHTTPHealthChecksStatusRef(d.Name(), d.tierName(), agentId, index))
}

func handleRefError(res *firego.Firebase, err error) *firego.Firebase {
  if err != nil {
    logrus.Fatalf("[E] referring to { %s } in firebase: %v", strings.TrimPrefix(res.URL(), FIREBASE_RTDB_URL), err)
  }
  return res
}
