package deployment

import (
	"fmt"
	"github.com/DawnBreather/go-commons/app/cicd_envsubst"
	"github.com/DawnBreather/go-commons/path"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
)

func (a *Artifact) ExecuteDeployment(secrets Secrets) (err error) {
	logrus.Infof("ARTIFACT -> DEPLOY STARTED")
	defer logrus.Infof("ARTIFACT -> DEPLOY ENDED")
	for _, source := range a.Source {
		err = source.Extract(secrets)
		if err != nil {
			return err
		}
	}

	err = a.Middleware.Envsubst.Execute(secrets)
	if err != nil {
		return err
	}

	for _, deploy := range a.Deploy {
		err = deploy.Deploy(secrets)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Source) Extract(secrets Secrets) (err error) {

	os.Setenv("AWS_ACCESS_KEY", TransformValuePlaceholderIntoValue(secrets, SECRETS_PATH_AWS_TOKEN))
	os.Setenv("AWS_SECRET_KEY", TransformValuePlaceholderIntoValue(secrets, SECRETS_PATH_AWS_SECRET))
	os.Setenv("AWS_REGION", TransformValuePlaceholderIntoValue(secrets, SECRETS_PATH_AWS_REGION))

	bucket := TransformValuePlaceholderIntoValue(secrets, s.From.S3.Bucket)
	item := TransformValuePlaceholderIntoValue(secrets, s.From.S3.Object)
	region := TransformValuePlaceholderIntoValue(secrets, s.From.S3.Region)

	logrus.Infof("SOURCE -> EXTRACT { s3://%s/%s --region %s} STARTED", bucket, item, region)
	defer logrus.Infof("SOURCE -> { s3://%s/%s --region %s} EXTRACT STARTED", bucket, item, region)

	filename := item
	if strings.Contains(item, "/") {
		_, filename = filepath.Split(item)
	}

	dstDir := path.New(s.To)
	if !dstDir.Exists() {
		dstDir.MkdirAll(os.ModePerm)
	}

	dstFile := path.New(os.TempDir() + string(filepath.Separator) + filename)

	//Unzip(dstFile.GetPath(), s.To)

	//TODO: uncomment
	if dstFile.Exists() {
		err := os.Remove(dstFile.GetPath())
		if err != nil {
			logrus.Errorf("[E] removing file: %v", err)
			return err
		}
	}

	file, err := os.Create(dstFile.GetPath())
	if err != nil {
		logrus.Errorf("[E] creating file: %v", err)
		return err
	}
	defer file.Close()

	sess, _ := session.NewSession(&aws.Config{Region: aws.String(region)})
	downloader := s3manager.NewDownloader(sess)
	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(item),
		})
	if err != nil {
		fmt.Println(err)
		logrus.Errorf("[E] Downloading file: %v", err)
		return err
	}

	logrus.Infof("[I] Downloaded %s %d bytes", file.Name(), numBytes)

	err = Unzip(dstFile.GetPath(), s.To)
	if err != nil {
		return err
	}

	return nil
}

func (d *Deploy) Deploy(secrets Secrets) (err error) {

	from := TransformValuePlaceholderIntoValue(secrets, d.From)
	to := TransformValuePlaceholderIntoValue(secrets, d.To)

	logrus.Infof("DEPLOY -> DEPLOY from { %s } to { %s } STARTED", from, to)
	defer logrus.Infof("DEPLOY -> DEPLOY from { %s } to { %s } ENDED", from, to)

	if err := d.Scripts.BeforeDeploy.Execute(secrets); err != nil {
		return err
	}

	if d.Mode == ENUM_DEPLOY_MODE_DELETE {
		logrus.Infof("[I] Removing objects in { %s } directory", to)
		removeContentsOfDestinationLocation(to)
	}

	logrus.Infof("[I] Copying objects from { %s } to { %s } directory", from, to)
	copy(from, to)

	if err := d.Scripts.AfterDeploy.Execute(secrets); err != nil {
		return err
	}

	return nil
}

// TODO: implement error return for envsubst
func (e *Envsubst) Execute(secrets Secrets) (err error) {
	for key, value := range e.Variables {
		os.Setenv(key, TransformValuePlaceholderIntoValue(secrets, value))
	}

	if !cliOptionsForEnvsubstInitialized {
		cicd_envsubst.ReadCliOptionsEnvsubst()
		cliOptionsForEnvsubstInitialized = true
	}

	cicd_envsubst.ProcessPaths(e.Target)
	return nil
}
