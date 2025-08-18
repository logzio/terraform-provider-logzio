package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
	content, err := os.ReadFile("testdata/fixtures/" + fileName)
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%s", content)
}

func ReadResourceFromFile(resourceName string, fileName string) string {
	return fmt.Sprintf(ReadFixtureFromFile(fileName), resourceName)
}

func SleepAfterTest() {
	time.Sleep(5 * time.Second)
}

func InterfaceToMapOfStrings(original interface{}) map[string]string {
	res := map[string]string{}
	originalToMap := original.(map[string]interface{})
	for k, v := range originalToMap {
		res[k] = v.(string)
	}
	return res
}

// MakeKibanaObjectDataUnique modifies JSON data to include unique IDs for testing
func MakeKibanaObjectDataUnique(jsonData string) string {
	uniqueSuffix := strconv.FormatInt(time.Now().UnixNano(), 10)

	var dataObj map[string]interface{}
	err := json.Unmarshal([]byte(jsonData), &dataObj)
	if err != nil {
		log.Fatal(err)
	}

	// Update _id field
	if id, ok := dataObj["_id"].(string); ok {
		dataObj["_id"] = id + "-" + uniqueSuffix
	}

	// Update _source.id field
	if source, ok := dataObj["_source"].(map[string]interface{}); ok {
		if id, ok := source["id"].(string); ok {
			source["id"] = id + "-" + uniqueSuffix
		}

		// Update type-specific id and title fields
		if sourceType, ok := source["type"].(string); ok {
			if typeObj, ok := source[sourceType].(map[string]interface{}); ok {
				if id, ok := typeObj["id"].(string); ok {
					typeObj["id"] = id + "-" + uniqueSuffix
				}
				if title, ok := typeObj["title"].(string); ok {
					typeObj["title"] = title + " " + uniqueSuffix
				}
			}
		}
	}

	result, err := json.Marshal(dataObj)
	if err != nil {
		log.Fatal(err)
	}

	return string(result)
}
