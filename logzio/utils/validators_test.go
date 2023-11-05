package utils

import (
	"github.com/logzio/logzio_terraform_client/alerts_v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
