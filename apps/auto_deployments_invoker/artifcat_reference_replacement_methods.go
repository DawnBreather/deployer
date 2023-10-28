package auto_deployments_invoker

import (
	"fmt"
	"go.uber.org/zap"
	"os"
	"regexp"
)

func (arr *ArtifactReferenceReplacement) InitializeForGitlabCiWorkflow() *ArtifactReferenceReplacement {
	var repository = os.Getenv(GITLAB_CI_PROJECT_NAME)         // i.e. breathesmart
	var branch = os.Getenv(GITLAB_CI_COMMIT_BRANCH)            // i.e. release-3.6-stallergenes
	var commitShortSha = os.Getenv(GITLAB_CI_COMMIT_SHORT_SHA) // i.e. 17b51953
	var prefix = os.Getenv(ARTIFACT_OBJECT_REFERENCE_PREFIX)

	return arr.initializeGeneric(prefix, repository, branch, commitShortSha)
}

func (arr *ArtifactReferenceReplacement) initializeGeneric(prefix, repository, branch, commitShortSha string) *ArtifactReferenceReplacement {
	var matchPattern = fmt.Sprintf("%s/%s/%s/build-\\b[0-9a-f]{8}\\b.deployer.zip", prefix, repository, branch)
	var newValue = fmt.Sprintf("%s/%s/%s/build-%s.deployer.zip", prefix, repository, branch, commitShortSha)

	Logger.Info("Initializing { artifact reference replacement }", zap.Any("details", map[string]string{
		"oldValuePattern": matchPattern,
		"newValue":        newValue,
	}))

	arr.matchPattern = regexp.MustCompile(matchPattern)
	arr.newValue = newValue

	return arr
}

func (arr *ArtifactReferenceReplacement) ReplaceAll(input []byte) []byte {
	//inputString := string(input)
	//fmt.Printf("", inputString)
	//found := arr.matchPattern.FindAllString(string(input), -1)
	//fmt.Println("Found: ", found)
	return arr.matchPattern.ReplaceAll(input, []byte(arr.newValue))
}
