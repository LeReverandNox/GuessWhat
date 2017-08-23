package tools

import (
	"regexp"
	"strings"

	"github.com/texttheater/golang-levenshtein/levenshtein"
)

// RemoveAllSpaces removes all spaces of a given string.
func RemoveAllSpaces(str string) string {
	rgx, _ := regexp.Compile("\\s\\s*")
	trimmedString := rgx.ReplaceAllString(str, "")

	return trimmedString
}

func ToLowerCase(str string) string {
	return strings.ToLower(str)
}

func IsAlphaHyphen(str string) bool {
	rgx, _ := regexp.Compile("^[a-zA-Z]+(-?)+[a-zA-Z]+$")
	isMatching := rgx.MatchString(str)

	return isMatching
}

func Sanitize(str string) string {
	trimmedString := strings.TrimSpace(str)
	return trimmedString
}

func Distance(str1 string, str2 string) int {
	dist := levenshtein.DistanceForStrings([]rune(str1), []rune(str2), levenshtein.DefaultOptions)
	return dist
}
