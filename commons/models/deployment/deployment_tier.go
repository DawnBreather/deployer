package deployment

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

func (t *Tier) DeployArtifacts(d *Deployment) {
	logrus.Infof("TIER -> DEPLOY ARTIFACTS STARTED")
	defer logrus.Infof("TIER -> DEPLOY ARTIFACTS ENDED")

	GetNodeStatus().ArtifactsDeployStatus = []ArtifactDeployStatus{}
	for index, artifact := range t.Artifacts {
		sourceString := ""

		for _, source := range artifact.Source {
			sourceS3Bucket := TransformValuePlaceholderIntoValue(d.Secrets, source.From.S3.Bucket)
			sourceS3Object := TransformValuePlaceholderIntoValue(d.Secrets, source.From.S3.Object)
			sourceS3Region := TransformValuePlaceholderIntoValue(d.Secrets, source.From.S3.Region)
			sourceString += fmt.Sprintf("s3://%s/%s --region %s", sourceS3Bucket, sourceS3Object, sourceS3Region)
		}

		GetNodeStatus().ArtifactsDeployStatus = append(GetNodeStatus().ArtifactsDeployStatus, ArtifactDeployStatus{
			Started:         GetTimeNow(),
			Source:          sourceString,
			ControlSequence: t.LatestControlSequence.CreatedAt,
			State:           DEPLOY_ARTIFACTS_STATE_IN_PROGRESS,
		})

		//PutIntoRedis(fmt.Sprintf(string(RDB_ARTIFACTS_DEPLOY_STATUS_PATH), d.Name(), d.tierName(), agentId, index), ArtifactDeployStatus{
		//	Source:  sourceString,
		//	State:   DEPLOY_ARTIFACTS_STATE_IN_PROGRESS,
		//	Message: "",
		//	Log:     "",
		//}, 6000*time.Second)

		if err := artifact.ExecuteDeployment(d.Secrets); err != nil {
			//PutIntoRedis(fmt.Sprintf(string(RDB_ARTIFACTS_DEPLOY_STATUS_PATH), d.Name(), d.tierName(), agentId, index), ArtifactDeployStatus{
			//	Source:  sourceString,
			//	State:   DEPLOY_ARTIFACTS_STATE_FAILED,
			//	Message: err.Error(),
			//	Log:     "",
			//}, 6000*time.Second)
			GetNodeStatus().ArtifactsDeployStatus[index].State = DEPLOY_ARTIFACTS_STATE_FAILED
			GetNodeStatus().ArtifactsDeployStatus[index].Message = err.Error()
		} else {
			//PutIntoRedis(fmt.Sprintf(string(RDB_ARTIFACTS_DEPLOY_STATUS_PATH), d.Name(), d.tierName(), agentId, index), ArtifactDeployStatus{
			//	Source:  sourceString,
			//	State:   DEPLOY_ARTIFACTS_STATE_DONE,
			//	Message: "",
			//	Log:     "",
			//}, 600000*time.Second)
			GetNodeStatus().ArtifactsDeployStatus[index].State = DEPLOY_ARTIFACTS_STATE_DONE
		}

		GetNodeStatus().ArtifactsDeployStatus[index].Ended = GetTimeNow()
		GetNodeStatus().ArtifactsDeployStatus[index].DurationSeconds = int64(GetNodeStatus().ArtifactsDeployStatus[index].Ended.Sub(GetNodeStatus().ArtifactsDeployStatus[index].Started).Seconds())

	}
}
