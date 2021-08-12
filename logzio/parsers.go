package logzio

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
)

// This function receives a string, tries to parse the string, and returns the string as
// the first type it managed to parse
func parseFromStringToType(value string) interface{} {
	// will try to parse in this order: json, float, int, bool, string
	var returnObject map[string]interface{}
	err := json.Unmarshal([]byte(value), &returnObject)
	if err == nil {
		return returnObject
	}

	returnFloat, err := strconv.ParseFloat(value, 64)
	if err == nil {
		return returnFloat
	}

	returnInt, err := strconv.Atoi(value)
	if err == nil {
		return returnInt
	}

	returnBool, err := strconv.ParseBool(value)
	if err == nil {
		return returnBool
	}

	return value
}

// This function receives an object and returns it as a string
func parseObjectToString(value interface{}) string {
	switch value.(type) {
	case map[string]interface{}:
		byteArray, _ := json.Marshal(value)
		return string(byteArray)
	case []map[string]interface{}:
		if len(value.([]map[string]interface{})) > 0 {
			byteArray, _ := json.Marshal(value)
			return string(byteArray)
		}
		return ""
	case map[string]string:
		byteArray, _ := json.Marshal(value)
		return string(byteArray)
	case string:
		return value.(string)
	default:
		return fmt.Sprintf("%v", value)
	}
}

func parseStringToMapList(value string) []map[string]interface{} {
	var returnObject []map[string]interface{}
	err := json.Unmarshal([]byte(value), &returnObject)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	return returnObject
}
