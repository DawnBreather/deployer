package aws

import (
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/s3"
  "github.com/aws/aws-sdk-go/service/s3/s3manager"
  "github.com/sirupsen/logrus"
)

func S3EmptyBucket(region, bucketName string) {

  sess, _ := session.NewSession(&aws.Config{
    Region: aws.String(region)},
  )

  // Create S3 service client
  svc := s3.New(sess)

  // Setup BatchDeleteIterator to iterate through a list of objects.
  iter := s3manager.NewDeleteListIterator(svc, &s3.ListObjectsInput{
    Bucket: aws.String(bucketName),
  })

  // Traverse iterator deleting each object
  if err := s3manager.NewBatchDeleteWithClient(svc).Delete(aws.BackgroundContext(), iter); err != nil {
    logrus.Errorf("[E] removing object from { %s } s3 bucket: %v", bucketName, err)
  }
}