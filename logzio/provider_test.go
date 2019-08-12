package logzio

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var (
	testAccExpectedAlertChannelName string
	testAccExpectedApplicationName  string
	testAccProviders                map[string]terraform.ResourceProvider
	testAccProvider                 *schema.Provider
)

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"logzio": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProviderImpl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv(envLogzioApiToken); v == "" {
		t.Log(v)
		t.Fatalf("%s must be set for acceptance tests", envLogzioApiToken)
	}
}
