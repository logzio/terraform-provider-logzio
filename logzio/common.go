package logzio

import "strings"

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
