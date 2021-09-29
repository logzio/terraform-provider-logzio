package utils

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
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

func IdFromResourceData(d *schema.ResourceData) (int64, error) {
	return strconv.ParseInt(d.Id(), BASE_10, BITSIZE_64)
}

func ReadFixtureFromFile(fileName string) string {
	content, err := ioutil.ReadFile("testdata/fixtures/" + fileName)
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%s", content)
}

func ReadResourceFromFile(resourceName string, fileName string) string {
	return fmt.Sprintf(ReadFixtureFromFile(fileName), resourceName)
}
