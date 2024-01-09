package deployment

import (
  "fmt"
  "github.com/sirupsen/logrus"
)

func (t *Tier) DeployArtifacts(d *Deployment) {
  logrus.Info("TIER -> DEPLOY ARTIFACTS STARTED")
  defer logrus.Info("TIER -> DEPLOY ARTIFACTS ENDED")

  nodeStatus := GetNodeStatus()
  nodeStatus.ArtifactsDeployStatus = make([]ArtifactDeployStatus, len(t.Artifacts))

  for index, artifact := range t.Artifacts {
    deployStatus := ArtifactDeployStatus{
      Started:         GetTimeNow(),
      Source:          t.formatSourceString(d.Secrets, artifact),
      ControlSequence: t.LatestControlSequence.CreatedAt,
      State:           DEPLOY_ARTIFACTS_STATE_IN_PROGRESS,
    }

    nodeStatus.ArtifactsDeployStatus[index] = deployStatus

    if err := artifact.ExecuteDeployment(d.Secrets); err != nil {
      nodeStatus.ArtifactsDeployStatus[index].State = DEPLOY_ARTIFACTS_STATE_FAILED
      nodeStatus.ArtifactsDeployStatus[index].Message = err.Error()
    } else {
      nodeStatus.ArtifactsDeployStatus[index].State = DEPLOY_ARTIFACTS_STATE_DONE
    }

    nodeStatus.ArtifactsDeployStatus[index].Ended = GetTimeNow()
    nodeStatus.ArtifactsDeployStatus[index].DurationSeconds = int64(nodeStatus.ArtifactsDeployStatus[index].Ended.Sub(nodeStatus.ArtifactsDeployStatus[index].Started).Seconds())
  }
}

func (t *Tier) formatSourceString(secrets map[string]any, artifact Artifact) string {
  var sourceStrings []string
  for _, source := range artifact.Source {
    sourceS3Bucket := TransformValuePlaceholderIntoValue(secrets, source.From.S3.Bucket)
    sourceS3Object := TransformValuePlaceholderIntoValue(secrets, source.From.S3.Object)
    sourceS3Region := TransformValuePlaceholderIntoValue(secrets, source.From.S3.Region)

    sourceString := fmt.Sprintf("s3://%s/%s --region %s", sourceS3Bucket, sourceS3Object, sourceS3Region)
    sourceStrings = append(sourceStrings, sourceString)
  }

  return fmt.Sprint(sourceStrings)
}
