package logzio

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	testAccExpectedAlertChannelName string
	testAccExpectedApplicationName  string
	testAccProviderFactories        = map[string]func() (*schema.Provider, error){
		"kubernetes": func() (*schema.Provider, error) {
			return Provider(), nil
		},
	}
	testAccProvider *schema.Provider
)

func init() {
	testAccProvider = Provider()
	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		"logzio": func() (*schema.Provider, error) {
			return Provider(), nil
		},
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
