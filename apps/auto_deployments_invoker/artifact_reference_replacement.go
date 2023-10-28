package auto_deployments_invoker

import (
	"regexp"
)

type ArtifactReferenceReplacement struct {
	matchPattern *regexp.Regexp
	newValue     string
}
