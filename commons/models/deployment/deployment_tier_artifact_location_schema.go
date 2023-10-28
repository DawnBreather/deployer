package deployment

import (
	"github.com/sirupsen/logrus"
	"regexp"
	"strings"
)

const (
	FILESYSTEM_SCHEMA = iota
	S3_SCHEMA
)

const (
	SCHEMA_PREFIX_FOR_S3_BUCKET        = "s3://"
	S3_BUCKET_LOCATION_GENERAL_PATTERN = `^s3://.*[ ]*.*\d`
)

// i.e. s3://general-bucket-123624/path/to/some/object.json --region us-east-2
func parseS3BucketLocation(location string) (region, bucketName, objectPath, regionFlag string) {
	regExp := regexp.MustCompile(S3_BUCKET_LOCATION_GENERAL_PATTERN)
	if len(regExp.FindString(location)) > 0 {

		bucketPathAndRegionFlag := strings.SplitN(location, " ", 2)
		bucketPath := bucketPathAndRegionFlag[0]
		if len(bucketPathAndRegionFlag) > 1 {
			regionFlag = bucketPathAndRegionFlag[1]
		}

		bucketNameAndObjectPath := strings.SplitN(strings.TrimPrefix(bucketPath, SCHEMA_PREFIX_FOR_S3_BUCKET), "/", 2)
		bucketName = bucketNameAndObjectPath[0]
		if len(bucketNameAndObjectPath) > 1 {
			objectPath = bucketNameAndObjectPath[1]
		}

		region = strings.TrimSpace(strings.TrimPrefix(regionFlag, "--region "))
	} else {
		logrus.Errorf("[E] parsing { %s } as s3 bucket location: not corresponding to regexp pattern { %s }", location, S3_BUCKET_LOCATION_GENERAL_PATTERN)
	}

	return

}

func getLocationSchema(location string) int {
	switch {
	case isSchemaS3Bucket(location):
		return S3_SCHEMA
	default:
		return FILESYSTEM_SCHEMA
	}
}

func isSchemaS3Bucket(location string) bool {
	return strings.HasPrefix(location, SCHEMA_PREFIX_FOR_S3_BUCKET)
}
