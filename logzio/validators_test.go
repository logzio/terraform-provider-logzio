package logzio

import (
	"github.com/jonboydell/logzio_client/alerts"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateOperation(t *testing.T) {
	validOperations := []string{
		alerts.OperatorGreaterThanOrEquals,
		alerts.OperatorLessThanOrEquals,
		alerts.OperatorGreaterThan,
		alerts.OperatorLessThan,
		alerts.OperatorEquals,
		alerts.OperatorNotEquals,
	}

	for _, s := range validOperations {
		_, errors := validateOperation(s, "operation")
		if len(errors) > 0 {
			t.Fatalf("%q should be a valid operation: %v", s, errors)
		}
	}

	invalidNames := []string{
		"",
		"this is not a valid operation",
	}

	for _, s := range invalidNames {
		_, errors := validateOperation(s, "operation")
		if len(errors) == 0 {
			t.Fatalf("%q should not be a valid operations: %v", s, errors)
		}
	}
}

func TestValidUrl(t *testing.T) {
	str := "https://some.url"
	_, errors := validateUrl(str, "url")
	assert.Len(t, errors, 0)
}
