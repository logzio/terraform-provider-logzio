package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
)

// This function receives a string, tries to parse the string, and returns the string as
// the first type it managed to parse
func ParseFromStringToType(value string) interface{} {
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
func ParseObjectToString(value interface{}) string {
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

func ParseStringToMapList(value string) []map[string]interface{} {
	var returnObject []map[string]interface{}
	err := json.Unmarshal([]byte(value), &returnObject)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	return returnObject
}

// ParseInterfaceSliceToStringSlice receives slice of interface, and returns a slice of string
func ParseInterfaceSliceToStringSlice(interfaceSlice []interface{}) []string {
	stringSlice := make([]string, 0, len(interfaceSlice))
	for _, s := range interfaceSlice {
		val, ok := s.(string)
		if !ok {
			val = ""
		}
		stringSlice = append(stringSlice, val)
	}
	return stringSlice
}
