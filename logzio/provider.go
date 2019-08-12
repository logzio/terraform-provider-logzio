package logzio

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

const (
	providerApiToken     = "api_token"
	providerBaseUrl      = "base_url"
	resourceAlertType    = "logzio_alert"
	resourceEndpointType = "logzio_endpoint"
	resourceUserType     = "logzio_user"
	envLogzioApiToken    = "LOGZIO_API_TOKEN"
	envLogzioBaseURL     = "LOGZIO_BASE_URL"
	envLogzioAccountId   = "LOGZIO_ACCOUNT_ID"

	defaultBaseUrl = "https://api.logz.io"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			providerApiToken: {
				Type:        schema.TypeString,
				Required:    true,
				Description: descriptions[providerApiToken],
				DefaultFunc: schema.EnvDefaultFunc(envLogzioApiToken, nil),
				Sensitive:   true,
			},
			providerBaseUrl: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions[providerBaseUrl],
				DefaultFunc: schema.EnvDefaultFunc(envLogzioBaseURL, defaultBaseUrl),
				Sensitive:   false,
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			resourceAlertType:    dataSourceAlert(),
			resourceEndpointType: dataSourceEndpoint(),
			resourceUserType:     dataSourceUser(),
		},
		ResourcesMap: map[string]*schema.Resource{
			resourceAlertType:    resourceAlert(),
			resourceEndpointType: resourceEndpoint(),
			resourceUserType:     resourceUser(),
		},
		ConfigureFunc: providerConfigure,
	}
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{providerApiToken: "Your API token", providerBaseUrl: "API base URL"}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	apiToken, ok := d.GetOk(providerApiToken)
	baseUrl, ok := d.GetOk(providerBaseUrl)
	if !ok {
		return nil, fmt.Errorf("can't find the %s, either set it in the provider or set the %s env var", providerApiToken, envLogzioApiToken)
	}
	config := Config{
		apiToken: apiToken.(string),
		baseUrl:  baseUrl.(string),
	}
	return config, nil
}
