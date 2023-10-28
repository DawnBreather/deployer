package main

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go/aws/awserr"
	ecsmetadata "github.com/brunoscheufler/aws-ecs-metadata-go"
	ec2metadata "github.com/travisjeffery/go-ec2-metadata"

	"github.com/sirupsen/logrus"
	"net/http"
)

const (
	WINDOWS_SERVICE_STATE_RUNNING  = 4
	WINDOWS_SERVICE_STATE_STOPPED  = 1
	WINDOWS_SERVICE_STATE_STARTING = 2
)

const (
	RDB_TIER_PATH string = `environments/%s/tiers/%s`
)

func main() {

	TestTransformValuePlaceholderIntoValue()
	//fmt.Println(fmt.Sprintf(RDB_TIER_PATH, "1", "2"))
	//TestGetValueByPathFromMap()
}

func TestTransformValuePlaceholderIntoValue() {
	var nestedMap = map[string]any{
		"k1": "v1",
		"k2": map[string]any{
			"nestedK1": "nestedV1",
			"nestedK2": "nestedV2",
			"nestedK3": map[string]any{
				"superNestedK1": "FOUND!!!",
			},
		},
	}

	fmt.Println(transformValuePlaceholderIntoValue(nestedMap, "${k1}.${k2.nestedK3.superNestedK1}"))
}

func transformValuePlaceholderIntoValue(data map[string]any, path string) string {
	re := regexp.MustCompile(`\${.*?}`)
	result := path
	for _, placeholder := range re.FindAllString(path, -1) {
		valuePath := strings.TrimSuffix(strings.TrimPrefix(placeholder, `${`), `}`)
		if value, isFound := GetValueByPathFromMap(data, valuePath, ""); isFound {
			result = strings.ReplaceAll(result, placeholder, value.(string))
		}
	}
	return result

	//if strings.HasPrefix(path, `${`) && strings.HasSuffix(path, "}") {
	//	path = strings.TrimSuffix(strings.TrimPrefix(path, `${`), `}`)
	//	if result, isFound := getValueByPathFromMap(data, path, ""); isFound {
	//		return result.(string)
	//	}
	//}
	//return path
}

func TestGetValueByPathFromMap() {
	var nestedMap = map[string]any{
		"k1": "v1",
		"k2": map[string]any{
			"nestedK1": "nestedV1",
			"nestedK2": "nestedV2",
			"nestedK3": map[string]any{
				"superNestedK1": "FOUND!!!",
			},
		},
	}

	searchKey := "k2.nestedK3.superNestedK1"

	if res, isFound := GetValueByPathFromMap(nestedMap, searchKey, ""); isFound {
		fmt.Println(res)
	} else {
		fmt.Println(searchKey)
	}

}

func GetValueByPathFromMap(data map[string]any, key string, passedKey string) (result any, found bool) {

	a := strings.SplitN(key, ".", 2)
	currentKey := a[0]
	if passedKey != "" {
		passedKey = passedKey + "." + currentKey
	} else {
		passedKey = currentKey
	}

	if _, isKeyExistInData := data[currentKey]; !isKeyExistInData {
		logrus.Errorf("[E] key path { %s } not found", passedKey)
		return
	} else {

		if len(a) > 1 {
			remainingPath := a[1]
			switch data[currentKey].(type) {
			case map[string]any:
				if result, found = GetValueByPathFromMap(data[currentKey].(map[string]any), remainingPath, passedKey); found {
					return
				}
			}
		} else {
			return data[currentKey], true
		}
	}

	return nil, false
}

func GettingMetadataFromEc2Ecs() {

	//os.Setenv("AWS_ACCESS_KEY_ID", "AKIATAMBS6UH654REHFT")
	//os.Setenv("AWS_SECRET_KEY", "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")
	//os.Setenv("AWS_DEFAULT_REGION", "us-east-1")

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		logrus.Fatalf("")
	}

	svc := sts.NewFromConfig(cfg)
	input := &sts.GetCallerIdentityInput{}

	result, err := svc.GetCallerIdentity(context.TODO(), input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(*result.Account)
	fmt.Println(*result.Arn)
	fmt.Println(*result.UserId)

	client := imds.NewFromConfig(cfg)
	fmt.Println(client.GetRegion(context.TODO(), nil))

	meta, err := ecsmetadata.Get(context.Background(), &http.Client{})
	if err != nil {
		panic(err)
	}

	switch m := meta.(type) {
	case *ecsmetadata.TaskMetadataV3:
		fmt.Println(m.TaskARN)
	case *ecsmetadata.TaskMetadataV4:
		fmt.Println(m.TaskARN)
	}

	ec2metadata.InstanceID()

}

//func PrototypeManagingWindowsServices(){
//	m, err := mgr.Connect()
//	if err != nil {
//		logrus.Errorf("[E] Connecting to services.msc: %v", err)
//	} else {
//		svc, err := m.OpenService("abz.deployer")
//		if err != nil {
//			logrus.Errorf("[E] Opening { abz.deployer } service: %v", err)
//			return
//		} else {
//			if svc != nil {
//				//err = svc.Start()
//				//if err != nil {
//				//	logrus.Errorf("[E] Starting { abz.deployer } service: %v", err)
//				//	return
//				//}
//
//				_, err = svc.Control(windows.SERVICE_STOPPED)
//				if err != nil {
//					logrus.Errorf("[E] Stopping service { abz.deployer }: %v", err)
//					return
//				}
//
//				for {
//					status, err := svc.Query()
//					if err != nil {
//						logrus.Errorf("[E] Querying status of { abz.deployer }: %v", err)
//						return
//					}
//
//					//err = svc.Close()
//					//if err != nil {
//					//	logrus.Errorf("[E] Closing service { abz.deployer }: %v", err)
//					//	return
//					//}
//
//					fmt.Println(status.State)
//					time.Sleep(2 * time.Second)
//				}
//
//			} else {
//				logrus.Warnf("[W] Service { abz.deployer } not found")
//				return
//			}
//		}
//	}
//
//}
