package main

import (
	"github.com/jonboydell/logzio_client"
	"testing"
)

func TestValidateOperation(t *testing.T) {
	validOperations := []string{
		logzio_client.OperatorGreaterThanOrEquals,
		logzio_client.OperatorLessThanOrEquals,
		logzio_client.OperatorGreaterThan,
		logzio_client.OperatorLessThan,
		logzio_client.OperatorEquals,
		logzio_client.OperatorNotEquals,
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