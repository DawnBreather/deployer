package main

import (
	"deployer/commons/models/deployment"
	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
	"github.com/sirupsen/logrus"
	"github.com/thoas/go-funk"
	"os"
	"regexp"
	"strings"
	"time"
)

// environments/deployer.local.cohero-health.com/status/api.admin/VKs-MacBook-Air-local/artifacts_deploy/0
func main() {
	commitHash := os.Getenv("CI_COMMIT")

	deployStatuses := map[string]deployment.ArtifactDeployStatus{}
	interestingKeys := []string{}
	// Create a regular expression pattern to match the numeric suffix after the last /
	pattern := regexp.MustCompile(`/\d+$`)

	// TODO: collect artifacts_deploy redis keys
	keyMask := "*/artifacts_deploy/*"
	keys, err := deployment.GetRedisKeysByMask(keyMask)
	if err != nil {
		logrus.Errorf("[E] getting Redis key by mask { %s }: %v", keyMask, err)
	}

	// TODO: collect entries and its values (artifacts_depoys)
	for _, key := range keys {
		deployStatus := deployment.ArtifactDeployStatus{}
		err := deployment.GetJsonEntryFromRedis(key, &deployStatus)
		if err != nil {
			logrus.Errorf("[E] getting Json Entry from redis: %v", err)
		}
		// TODO: filter artifact_deploys: keep only relevant by {CI_COMMIT}
		if strings.Contains(deployStatus.Source, commitHash) {
			deployStatuses[key] = deployStatus
			if !funk.Contains(interestingKeys, pattern.ReplaceAllString(key, "")) {
				interestingKeys = append(interestingKeys, pattern.ReplaceAllString(key, ""))
			}
		}
	}

	// TODO: extract {environment}, {tier}, {nodeName} from artifact_deploys
	// TODO: wait for completion of artifact_deploys
	finalKeys := funk.FilterString(keys, func(key string) bool {
		for _, interestingKey := range interestingKeys {
			if strings.HasPrefix(key, interestingKey) {
				return true
			}
		}
		return false
	})

	bar := progressbar.NewOptions(len(finalKeys),
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(false),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetDescription("[cyan][1/1][reset] Waiting for deployment completion..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	for _, key := range finalKeys {
		go func() {
			for {
				deployStatus := deployment.ArtifactDeployStatus{}
				err := deployment.GetJsonEntryFromRedis(key, &deployStatus)
				if err != nil {
					logrus.Errorf("[E] getting Json Entry from redis: %v", err)
				}
				if deployStatus.State == deployment.DEPLOY_ARTIFACTS_STATE_DONE {
					bar.Add(1)
				}
				if deployStatus.State == deployment.DEPLOY_ARTIFACTS_STATE_FAILED {
					logrus.Errorf("[E] deploying { %s } for { %s } in { %s }: %s", deployStatus.Source, deployStatus.Message)
				}
				time.Sleep(3 * time.Second)
			}
		}()
	}

	// TODO: generate gitlab.ci.yaml

	// TODO: collect all artifcats which are deploying for the node
}
