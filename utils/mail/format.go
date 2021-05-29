package mailUtils

import "regexp"

func CheckMailFormat(email string) bool {

	mailCompile := regexp.MustCompile("^(.*)@(.*)\\.(.*)$")

	r := mailCompile.FindSubmatch([]byte(email))
	if len(r) != 4 {
		return false
	}
	return true
}

func IgnoreSensMatch(value string) bool {

	valueCompile := regexp.MustCompile(`(?i)`+value)

	r := valueCompile.MatchString(value)
	if !r {
		return false
	}
	return true
}
