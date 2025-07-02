package logzio

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	envLogzioAccountId          = "LOGZIO_ACCOUNT_ID"
	envLogzioWarmAccountId      = "LOGZIO_WARM_ACCOUNT_ID"
	envLogzioS3Path             = "S3_PATH"
	envLogzioAwsAccessKey       = "AWS_ACCESS_KEY"
	envLogzioAwsSecretKey       = "AWS_SECRET_KEY"
	envLogzioAwsArn             = "AWS_ARN"
	envLogzioAwsArnS3Fetcher    = "AWS_ARN_S3_FETCHER"
	envLogzioAzureAccountName   = "AZURE_ACCOUNT_NAME"
	envLogzioAzureClientId      = "AZURE_CLIENT_ID"
	envLogzioAzureClientSecret  = "AZURE_CLIENT_SECRET"
	envLogzioAzureContainerName = "AZURE_CONTAINER_NAME"
	envLogzioAzureTenantId      = "AZURE_TENANT_ID"
	envLogzioAzurePath          = "BLOB_PATH"
	envLogzioApiTokenWarm       = "LOGZIO_WARM_API_TOKEN"
)

var (
	testAccExpectedAlertChannelName string
	testAccExpectedApplicationName  string
	testAccProviderFactories        = map[string]func() (*schema.Provider, error){
		"kubernetes": func() (*schema.Provider, error) {
			return Provider(), nil
		},
	}
	testAccWarmProviderFactories = map[string]func() (*schema.Provider, error){
		"kubernetes": func() (*schema.Provider, error) {
			return ProviderWithEnvVar(envLogzioApiTokenWarm), nil
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
	testAccWarmProviderFactories = map[string]func() (*schema.Provider, error){
		"logzio": func() (*schema.Provider, error) {
			return ProviderWithEnvVar(envLogzioApiTokenWarm), nil
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
func testAccPreCheckApiTokenWarm(t *testing.T) {
	testAccPreCheckEnv(t, envLogzioApiTokenWarm)
}
func testAccPreCheckAccountId(t *testing.T) {
	testAccPreCheckEnv(t, envLogzioAccountId)
}
func testAccPreCheckEmail(t *testing.T) {
	testAccPreCheckEnv(t, envLogzioEmail)
}
