package logzio

import (
	"strings"
	"unicode"
)

const (
	BASE_10            int    = 10
	BITSIZE_64         int    = 64
	VALIDATE_URL_REGEX string = "^http(s):\\/\\/"
)

func findStringInArray(v string, values []string) bool {
	for i := 0; i < len(values); i++ {
		value := values[i]
		if strings.EqualFold(v, value) {
			return true
		}
	}
	return false
}

func stripAllWhitespace(inputString string) string {
	var b strings.Builder
	b.Grow(len(inputString))
	for _, ch := range inputString {
		if !unicode.IsSpace(ch) {
			b.WriteRune(ch)
		}
	}
	return b.String()
}