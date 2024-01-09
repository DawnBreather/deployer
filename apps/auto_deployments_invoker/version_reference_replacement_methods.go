package auto_deployments_invoker

import (
	"fmt"
	"go.uber.org/zap"
	"os"
	"regexp"
)

func (vrr *VersionReferenceReplacement) InitializeForGitlabCiWorkflow() *VersionReferenceReplacement {
	var commitShortSha = os.Getenv(GITLAB_CI_COMMIT_SHORT_SHA) // i.e. 17b51953

	return vrr.initializeGeneric(commitShortSha)
}

func (vrr *VersionReferenceReplacement) initializeGeneric(commitShortSha string) *VersionReferenceReplacement {
	var matchPattern = fmt.Sprintf("build-number: .*")
	var newValue = fmt.Sprintf("build-number: %s", commitShortSha)

	Logger.Info("Initializing { version reference replacement }", zap.Any("details", map[string]string{
		"oldValuePattern": matchPattern,
		"newValue":        newValue,
	}))

	vrr.matchPattern = regexp.MustCompile(matchPattern)
	vrr.newValue = newValue

	return vrr
}

func (vrr *VersionReferenceReplacement) ReplaceAll(input []byte) []byte {
	//inputString := string(input)
	//fmt.Printf("", inputString)
	//found := vrr.matchPattern.FindAllString(string(input), -1)
	//fmt.Println("Found: ", found)
	return vrr.matchPattern.ReplaceAll(input, []byte(vrr.newValue))
}
