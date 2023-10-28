package url

import "regexp"

func CleanUrlFromCredentials(url string) string {
	return regexp.MustCompile(`^([a-z]*?://)?([^:/]+:[^@]+@)(.*)$`).ReplaceAllString(url, `$1$3`)
}
