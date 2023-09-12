package utils

import (
	"github.com/logzio/logzio_terraform_client/alerts_v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateOperation(t *testing.T) {
	validOperations := []string{
		alerts_v2.OperatorGreaterThanOrEquals,
		alerts_v2.OperatorLessThanOrEquals,
		alerts_v2.OperatorGreaterThan,
		alerts_v2.OperatorLessThan,
		alerts_v2.OperatorEquals,
		alerts_v2.OperatorNotEquals,
	}

	for _, s := range validOperations {
		_, errors := ValidateOperation(s, "operation")
		if len(errors) > 0 {
			t.Fatalf("%q should be a validd operation: %v", s, errors)
		}
	}

	invalidNames := []string{
		"",
		"this is not a valid operation",
	}

	for _, s := range invalidNames {
		_, errors := ValidateOperation(s, "operation")
		if len(errors) == 0 {
			t.Fatalf("%q should not be a valid operations: %v", s, errors)
		}
	}
}

func TestValidUrl(t *testing.T) {
	str := "https://some.url"
	_, errors := ValidateUrl(str, "url")
	assert.Len(t, errors, 0)
}

func TestValidateOutputTypes(t *testing.T) {
	validOptions := []string{
		alerts_v2.OutputTypeJson,
		alerts_v2.OutputTypeTable,
	}

	for _, s := range validOptions {
		_, errors := ValidateOutputType(s, "output_type")
		assert.Empty(t, errors)
	}

	invalidNames := []string{
		"",
		"this is not a valid type",
	}

	for _, s := range invalidNames {
		_, errors := ValidateOutputType(s, "output_type")
		assert.NotEmpty(t, errors)
	}
}

func TestValidateSortTypes(t *testing.T) {
	validTypes := []string{
		alerts_v2.SortAsc,
		alerts_v2.SortDesc,
	}

	for _, s := range validTypes {
		_, errors := ValidateSortTypes(s, "sort")
		assert.Empty(t, errors)
	}

	invalidNames := []string{
		"",
		"this is not a valid type",
	}

	for _, s := range invalidNames {
		_, errors := ValidateSortTypes(s, "sort")
		assert.NotEmpty(t, errors)
	}
}
