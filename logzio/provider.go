package logzio

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	providerApiToken                      = "api_token"
	providerCustomApiUrl                  = "custom_api_url"
	providerBaseUrl                       = "base_url"
	providerRegion                        = "region"
	resourceAlertType                     = "logzio_alert"
	resourceAlertV2Type                   = "logzio_alert_v2"
	resourceEndpointType                  = "logzio_endpoint"
	resourceUserType                      = "logzio_user"
	resourceSubAccountType                = "logzio_subaccount"
	resourceMetricsAccountType            = "logzio_metrics_account"
	resourceLogShippingTokenType          = "logzio_log_shipping_token"
	resourceDropFilterType                = "logzio_drop_filter"
	resourceArchiveLogsType               = "logzio_archive_logs"
	resourceRestoreLogsType               = "logzio_restore_logs"
	resourceAuthenticationGroupsType      = "logzio_authentication_groups"
	resourceKibanaObjectType              = "logzio_kibana_object"
	resourceS3FetcherType                 = "logzio_s3_fetcher"
	resourceGrafanaDashboardType          = "logzio_grafana_dashboard"
	resourceGrafanaFolderType             = "logzio_grafana_folder"
	resourceGrafanaAlertRuleType          = "logzio_grafana_alert_rule"
	resourceGrafanaNotificationPolicyType = "logzio_grafana_notification_policy"
	resourceGrafanaContactPointType       = "logzio_grafana_contact_point"

	envLogzioApiToken = "LOGZIO_API_TOKEN"
	envLogzioRegion   = "LOGZIO_REGION"
	envLogzioCustomApiUrl = "LOGZIO_CUSTOM_API_URL"

	baseUrl = "https://api%s.logz.io"
)

func Provider() *schema.Provider {
	return ProviderWithEnvVar(envLogzioApiToken)
}

func ProviderWithEnvVar(apiTokenEnvVar string) *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			providerApiToken: {
				Type:        schema.TypeString,
				Required:    true,
				Description: descriptions[providerApiToken],
				DefaultFunc: schema.EnvDefaultFunc(apiTokenEnvVar, nil),
				Sensitive:   true,
			},
			providerRegion: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions[providerRegion],
				DefaultFunc: schema.EnvDefaultFunc(envLogzioRegion, ""),
				Sensitive:   false,
			},
			providerCustomApiUrl: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom API URL to override the default Logz.io API endpoint.",
				DefaultFunc: schema.EnvDefaultFunc(envLogzioCustomApiUrl, ""),
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			resourceEndpointType:             dataSourceEndpoint(),
			resourceUserType:                 dataSourceUser(),
			resourceSubAccountType:           dataSourceSubAccount(),
			resourceMetricsAccountType:       dataSourceMetricsAccount(),
			resourceAlertV2Type:              dataSourceAlertV2(),
			resourceLogShippingTokenType:     dataSourceLogShippingToken(),
			resourceDropFilterType:           dataSourceDropFilter(),
			resourceArchiveLogsType:          dataSourceArchiveLogs(),
			resourceRestoreLogsType:          dataSourceRestoreLogs(),
			resourceAuthenticationGroupsType: dataSourceAuthenticationGroups(),
			resourceKibanaObjectType:         dataSourceKibanaObject(),
			resourceS3FetcherType:            dataSourceS3Fetcher(),
			resourceGrafanaDashboardType:     dataSourceGrafanaDashboard(),
			resourceGrafanaFolderType:        dataSourceGrafanaFolder(),
		},
		ResourcesMap: map[string]*schema.Resource{
			resourceEndpointType:                  resourceEndpoint(),
			resourceUserType:                      resourceUser(),
			resourceSubAccountType:                resourceSubAccount(),
			resourceMetricsAccountType:            resourceMetricsAccount(),
			resourceAlertV2Type:                   resourceAlertV2(),
			resourceLogShippingTokenType:          resourceLogShippingToken(),
			resourceDropFilterType:                resourceDropFilter(),
			resourceArchiveLogsType:               resourceArchiveLogs(),
			resourceRestoreLogsType:               resourceRestoreLogs(),
			resourceAuthenticationGroupsType:      resourceAuthenticationGroups(),
			resourceKibanaObjectType:              resourceKibanaObject(),
			resourceS3FetcherType:                 resourceS3Fetcher(),
			resourceGrafanaDashboardType:          resourceGrafanaDashboard(),
			resourceGrafanaFolderType:             resourceGrafanaFolder(),
			resourceGrafanaAlertRuleType:          resourceGrafanaAlertRule(),
			resourceGrafanaNotificationPolicyType: resourceGrafanaNotificationPolicy(),
			resourceGrafanaContactPointType:       resourceGrafanaContactPoint(),
		},
		ConfigureContextFunc: providerConfigureWrapper,
	}
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		providerApiToken: "Your API token",
		providerRegion: "Your logz.io region",
		providerCustomApiUrl: "Custom API URL to override the default Logz.io API endpoint. Useful for routing through internal gateways/proxies.",
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	apiToken, ok := d.GetOk(providerApiToken)
	if !ok {
		return nil, diag.Errorf("can't find the %s, either set it in the provider or set the %s env var", providerApiToken, envLogzioApiToken)
	}
	region := d.Get(providerRegion).(string)
	customApiUrl := d.Get(providerCustomApiUrl).(string)
	var apiUrl string
	if customApiUrl != "" {
		apiUrl = customApiUrl
	} else {
		regionCode := ""
		if region != "" && region != "us" {
			regionCode = fmt.Sprintf("-%s", region)
		}
		apiUrl = fmt.Sprintf(baseUrl, regionCode)
	}

	config := Config{
		apiToken: apiToken.(string),
		baseUrl:  apiUrl,
	}
	return config, diag.Diagnostics{}
}

func providerConfigureWrapper(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return providerConfigure(d)
}
