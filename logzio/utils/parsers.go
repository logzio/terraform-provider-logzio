package utils

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

func ParseTypeSetToMap(d *schema.ResourceData, key string) (map[string]interface{}, error) {
	if v, ok := d.GetOk(key); ok {
		rawMappings := v.(*schema.Set).List()
		for i := 0; i < len(rawMappings); i++ {
			x := rawMappings[i]
			y := x.(map[string]interface{})
			return y, nil
		}
	}

	return nil, fmt.Errorf("can't load mapping for key %s", key)
}
