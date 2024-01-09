package auto_deployments_invoker

import "regexp"

type VersionReferenceReplacement struct {
	matchPattern *regexp.Regexp
	newValue     string
}
