package logzio

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	testAccExpectedAlertChannelName string
	testAccExpectedApplicationName  string
	testAccProviders                map[string]*schema.Provider
	testAccProvider                 *schema.Provider
)

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"logzio": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProviderImpl(t *testing.T) {
	var _ *schema.Provider = Provider()
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
