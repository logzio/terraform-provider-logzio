package logzio

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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

func testAccPreCheckEnv(t *testing.T, env string) {
	if v := os.Getenv(env); v == "" {
		t.Errorf("%s must be set for acceptance tests", env)
	}
}

func testAccPreCheckApiToken(t *testing.T) {
	testAccPreCheckEnv(t, envLogzioApiToken)
}
func testAccPreCheckAccountId(t *testing.T) {
	testAccPreCheckEnv(t, envLogzioAccountId)
}
func testAccPreCheckEmail(t *testing.T) {
	testAccPreCheckEnv(t, envLogzioEmail)
}
