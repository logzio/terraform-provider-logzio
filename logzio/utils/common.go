package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/avast/retry-go"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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

func ConvertToStrings[T ~string](items []T) []string {
	result := make([]string, len(items))
	for i, v := range items {
		result[i] = string(v)
	}
	return result
}

// ReadUntilConsistent retries the readFunc until the isConsistent function returns true or the retryAttempts is exhausted.
func ReadUntilConsistent(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
	retryAttempts int,
	operation string,
	readFunc func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics,
	isConsistent func() bool) diag.Diagnostics {
	err := retry.Do(
		func() error {
			diags := readFunc(ctx, d, m)
			if diags != nil && len(diags) > 0 {
				return fmt.Errorf("failed to read after %s: %v", operation, diags)
			}
			if !isConsistent() {
				return fmt.Errorf("resource state not consistent after %s", operation)
			}
			return nil
		},
		retry.RetryIf(func(err error) bool {
			return err != nil
		}),
		retry.Attempts(uint(retryAttempts)),
		retry.DelayType(retry.BackOffDelay),
	)

	if err != nil {
		tflog.Warn(ctx, fmt.Sprintf("Failed to achieve consistency after %s: %v", operation, err))
	}

	return nil
}

// GetOptionalInt32Pointer returns a pointer to the numeric value from the config, or nil if it was not set.
// We don't use d.Get since it returns 0 for nil, when the field is unset.
// And we don't use d.GetOk because it returns false for 0, even if it's explicitly set.
func GetOptionalInt32Pointer(d *schema.ResourceData, key string) *int32 {
	val, diags := d.GetRawConfigAt(cty.GetAttrPath(key))
	if diags.HasError() || val.IsNull() {
		return nil
	} else {
		int64Val, _ := val.AsBigFloat().Int64()
		int32Val := int32(int64Val)
		return &int32Val
	}
}

func GetOptionalFloat32Pointer(d *schema.ResourceData, key string) *float32 {
	val, diags := d.GetRawConfigAt(cty.GetAttrPath(key))
	if diags.HasError() || val.IsNull() {
		return nil
	}
	floatVal, _ := val.AsBigFloat().Float32()
	return &floatVal
}
