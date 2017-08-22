package tools

import "regexp"

// RemoveAllSpaces removes all spaces of a given string.
func RemoveAllSpaces(str string) string {
	rgx, _ := regexp.Compile("\\s\\s+")
	trimmedString := rgx.ReplaceAllString(str, "")

	return trimmedString
}
