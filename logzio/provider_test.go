package logzio

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	envLogzioAccountId            = "LOGZIO_ACCOUNT_ID"
	envLogzioWarmAccountId        = "LOGZIO_WARM_ACCOUNT_ID"
	envLogzioMetricsAccountId     = "LOGZIO_METRICS_ACCOUNT_ID"
	envLogzioS3Path               = "S3_PATH"
	envLogzioAwsAccessKey         = "AWS_ACCESS_KEY"
	envLogzioAwsSecretKey         = "AWS_SECRET_KEY"
	envLogzioAwsArn               = "AWS_ARN"
	envLogzioAwsArnS3Fetcher      = "AWS_ARN_S3_FETCHER"
	envLogzioAzureAccountName     = "AZURE_ACCOUNT_NAME"
	envLogzioAzureClientId        = "AZURE_CLIENT_ID"
	envLogzioAzureClientSecret    = "AZURE_CLIENT_SECRET"
	envLogzioAzureContainerName   = "AZURE_CONTAINER_NAME"
	envLogzioAzureTenantId        = "AZURE_TENANT_ID"
	envLogzioAzurePath            = "BLOB_PATH"
	envLogzioApiTokenWarm         = "LOGZIO_WARM_API_TOKEN"
	envLogzioApiTokenConsumption  = "ֿLOGZIO_CONSUMPTION_API_TOKEN"
	envLogzioConsumptionAccountId = "ֿLOGZIO_CONSUMPTION_ACCOUNT_ID"
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
		"logzio": func() (*schema.Provider, error) {
			return ProviderWithEnvVar(envLogzioApiTokenWarm), nil
		},
	}
	testAccConsumptionProviderFactories = map[string]func() (*schema.Provider, error){
		"logzio": func() (*schema.Provider, error) {
			return ProviderWithEnvVar(envLogzioApiTokenConsumption), nil
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
	testAccConsumptionProviderFactories = map[string]func() (*schema.Provider, error){
		"logzio": func() (*schema.Provider, error) {
			return ProviderWithEnvVar(envLogzioApiTokenConsumption), nil
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

func TestProvider_BaseUrlResolution(t *testing.T) {
	testCases := []struct {
		name            string
		region          string
		customApiUrl    string
		expectedBaseUrl string
	}{
		{
			name:            "Only region set (eu)",
			region:          "eu",
			expectedBaseUrl: "https://api-eu.logz.io",
		},
		{
			name:            "Only custom_api_url set",
			customApiUrl:    "https://custom.example.com/api",
			expectedBaseUrl: "https://custom.example.com/api",
		},
		{
			name:            "Both region and custom_api_url set (custom wins)",
			region:          "eu",
			customApiUrl:    "https://custom.example.com/api",
			expectedBaseUrl: "https://custom.example.com/api",
		},
		{
			name:            "Neither region nor custom_api_url set (default US)",
			expectedBaseUrl: "https://api.logz.io",
		},
		{
			name:            "Region set to us (default US)",
			region:          "us",
			expectedBaseUrl: "https://api.logz.io",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			provider := Provider()
			attrs := map[string]interface{}{
				"api_token": "dummy-token",
			}
			if tc.region != "" {
				attrs["region"] = tc.region
			}
			if tc.customApiUrl != "" {
				attrs["custom_api_url"] = tc.customApiUrl
			}
			resourceData := schema.TestResourceDataRaw(t, provider.Schema, attrs)
			cfg, diags := providerConfigure(resourceData)
			if len(diags) > 0 {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}
			config := cfg.(Config)
			if config.baseUrl != tc.expectedBaseUrl {
				t.Errorf("expected baseUrl to be %s, got %s", tc.expectedBaseUrl, config.baseUrl)
			}
		})
	}
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
func testAccPreCheckApiTokenConsumption(t *testing.T) {
	testAccPreCheckEnv(t, envLogzioApiTokenConsumption)
}
func testAccPreCheckAccountId(t *testing.T) {
	testAccPreCheckEnv(t, envLogzioAccountId)
}
func testAccPreCheckEmail(t *testing.T) {
	testAccPreCheckEnv(t, envLogzioEmail)
}

func testAccPreCheckMetricsAccountId(t *testing.T) {
	testAccPreCheckEnv(t, envLogzioMetricsAccountId)
}
