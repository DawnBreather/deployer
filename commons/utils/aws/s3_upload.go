package aws

import (
  "github.com/DawnBreather/s3sync"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/sirupsen/logrus"
)

func S3Upload(region, from, to string) {
  // Creates an AWS session
  sess, _ := session.NewSession(&aws.Config{
    Region: aws.String(region),
  })

  syncManager := s3sync.New(sess)

  //// Sync from s3 to local
  //syncManager.Sync("s3://yourbucket/path/to/dir", "local/path/to/dir")

  // Sync from local to s3
  err := syncManager.Sync(from, to)
  if err != nil {
    logrus.Errorf("[E] uploading objects from { %s } path to { %s } s3 bucket: %v", from, to, err)
  }
}
