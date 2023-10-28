package main

// deployer environment update --name cohero/adhereit/demo/backend --path artifacts/0/source/0/from --value s3://breathesmart-builds/CloudFormation/PreProd/build-286.zip

import (
	"deployer/commons/models/deployment"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

// TODO
// 1.

func main() {

	logrus.Infof("Version 1.0.10")

	logrus.SetFormatter(&logrus.TextFormatter{
		DisableQuote: true,
		ForceColors:  true,
	})

	environmentID := os.Getenv(deployment.ABZ_DEPLOYER_AGENT_ENVIRONMENT_ID)
	if environmentID == "" {
		logrus.Fatalf("[E] { ABZ_DEPLOYER_AGENT_ENVIRONMENT_ID } environment variable is not set: please provide Environment ID.")
	}

	var d = deployment.Deployment{
		Metadata: deployment.Metadata{
			Name: os.Getenv(deployment.ABZ_DEPLOYER_AGENT_ENVIRONMENT_ID),
		},
	}

	deployment.PullCredentials()

	//d.StartPullingConfigurationFromGitRepository()
	//d.StartListeningForTierConfiguration()
	//d.StartListeningForTierConfigurationFromGit()
	//d.WaitForInitialization()
	//d.Tier().HealthChecks.StartValidatingHealthChecks(&d)
	//
	//for {
	//	time.Sleep(1 * time.Second)
	//}

	d.PullConfigurationFromRedisWithDecryptionAndDeploy()
	//d.StartPullingConfigurationFromGitRepository()
	//deployment.StartKeepingFirebaseAuthenticated()
	//d.WaitForFirebaseAuthentication()

	// # Enable
	//d.WaitForConfigurationFromGitInitialRetrieval()

	//d.StartListeningForMetadataConfiguration()
	//d.StartListeningForSecretsConfiguration()
	//d.StartListeningForSecretsConfiguration()

	// # Enable
	//d.StartListeningForSecretsConfigurationFromGit()
	//d.StartListeningForTierConfigurationFromGit()

	// =====================> d.StartListeningForTierConfiguration()

	d.Install()

	d.WaitForInitialization()

	// # Enable
	//d.StartListeningForControlSequenceConfigurationFromGit()

	//d.StartListeningForSecretsConfigurationFromRedisChannel()
	d.StartListeningForSecretsConfigurationFromRedisChannel()
	d.StartListeningForTierConfigurationFromRedisChannel()
	d.StartListeningForControlSequenceConfigurationFromRedisChannel()
	deployment.GetNodeStatus().StartSubmittingToRedis(&d)

	// =====================> d.StartListeningForLatestControlSequence()
	////d.Tier.DeployArtifacts()
	d.Tier().HealthChecks.StartValidatingHealthChecks(&d)

	if d.Tier().Entrypoint.Autostart {
		d.Tier().StartRunningEntrypoint()
	}

	//go func() {
	//	d.Tier.StartRunningEntrypoint()
	//	//L:
	//	//	for {
	//	//		select {
	//	//		case <-stop:
	//	//			fmt.Println("stopping")
	//	//			cancel()
	//	//			break L
	//	//		default:
	//	//			d.Tier.StartRunningEntrypoint(ctx)
	//	//		}
	//	//	}
	//}()
	//
	//go func() {
	//	time.Sleep(2 * time.Second)
	//	fmt.Println("Stopping")
	//	cancel()
	//	//stop <- struct{}{}
	//	//close(stop)
	//}()

	//go func() {
	//	time.Sleep(2 * time.Second)
	//	fmt.Println("Stopping")
	//	deployment.EntrypointCtxCancel()
	//	//stop <- struct{}{}
	//	//close(stop)
	//}()

	for {
		//fmt.Println("running main process")
		time.Sleep(1 * time.Second)
	}

	//deployment.RemoveContentsOfDirectory("./webapi")

	//// TODO: CAREFUL, don't put slashes at the end of { destination }
	//deployment.Unzip(`C:\Users\won\AppData\Local\Temp\build-307.zip`, "webapi")

	//for {
	//	time.Sleep(1 * time.Second)
	//}

	//logrus.Error("Hello\nHello")

	//
	//data, _ := os.ReadFile("/tmp/webapi/build-307.zip")
	//buffer := bytes.NewBuffer(data)
	//err := extract.Zip(context.TODO(), buffer, "unzip-test/dst/2/3/4", nil)
	//if err != nil {
	//  logrus.Errorf("[E] extracting zip: %v", err)
	//}

	//err := Unzip("/tmp/webapi/build-307.zip", "unzip-test/dst")
	//if err != nil {
	//  logrus.Errorf("[E] extracting zip: %v", err)
	//}

	//uz := unzip.New("/tmp/webapi/build-307.zip", "unzip-test/dst/")
	//err := uz.Extract()
	//if err != nil {
	//  logrus.Errorf("[E] extracting zip file { %s } to { %s }")
	//}

	//err := cp.Copy("copy-test/src", "copy-test/dst/src")
	//if err != nil {
	//  logrus.Errorf("[E] copying files: %v", err)
	//}

	//var d = Deployment{
	//  Name: "demo.adhereit.cohero-health.com",
	//}
	//
	//d.StartListeningForTierConfiguration()
	//d.WaitForInitialization()
	//
	//d.Tier.Artifacts[0].Source[0].Extract(d.Tier.Credentials)
	//for {
	//  time.Sleep(1 * time.Second)
	//}

	//f := firego.New("https://test-7f148-default-rtdb.firebaseio.com/", nil)
	//
	//var data map[string]interface{}
	//err := json.Unmarshal([]byte(sampleSchema), &data)
	//if err != nil {
	//  fmt.Println("error unmarshalling json", err)
	//}
	//
	//err = f.Set(data)
	//if err != nil {
	//  fmt.Println("error sending data", err)
	//}
	//
	//for {
	//  time.Sleep(1 * time.Second)
	//}

}
