package logzio

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	providerApiToken                 = "api_token"
	providerBaseUrl                  = "base_url"
	providerRegion                   = "region"
	resourceAlertType                = "logzio_alert"
	resourceAlertV2Type              = "logzio_alert_v2"
	resourceEndpointType             = "logzio_endpoint"
	resourceUserType                 = "logzio_user"
	resourceSubAccountType           = "logzio_subaccount"
	resourceLogShippingTokenType     = "logzio_log_shipping_token"
	resourceDropFilterType           = "logzio_drop_filter"
	resourceArchiveLogsType          = "logzio_archive_logs"
	resourceRestoreLogsType          = "logzio_restore_logs"
	resourceAuthenticationGroupsType = "logzio_authentication_groups"
	resourceKibanaObjectType         = "logzio_kibana_object"
	envLogzioApiToken                = "LOGZIO_API_TOKEN"
	envLogzioRegion                  = "LOGZIO_REGION"
	envLogzioBaseURL                 = "LOGZIO_BASE_URL"
	envLogzioAccountId               = "LOGZIO_ACCOUNT_ID"
	envLogzioS3Path                  = "S3_PATH"
	envLogzioAwsAccessKey            = "AWS_ACCESS_KEY"
	envLogzioAwsSecretKey            = "AWS_SECRET_KEY"
	envLogzioAwsArn                  = "AWS_ARN"
	envLogzioAzureAccountName        = "AZURE_ACCOUNT_NAME"
	envLogzioAzureClientId           = "AZURE_CLIENT_ID"
	envLogzioAzureClientSecret       = "AZURE_CLIENT_SECRET"
	envLogzioAzureContainerName      = "AZURE_CONTAINER_NAME"
	envLogzioAzureTenantId           = "AZURE_TENANT_ID"
	envLogzioAzurePath               = "BLOB_PATH"

	baseUrl = "https://api%s.logz.io"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			providerApiToken: {
				Type:        schema.TypeString,
				Required:    true,
				Description: descriptions[providerApiToken],
				DefaultFunc: schema.EnvDefaultFunc(envLogzioApiToken, nil),
				Sensitive:   true,
			},
			providerRegion: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions[providerRegion],
				DefaultFunc: schema.EnvDefaultFunc(envLogzioRegion, ""),
				Sensitive:   false,
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			resourceEndpointType:             dataSourceEndpoint(),
			resourceUserType:                 dataSourceUser(),
			resourceSubAccountType:           dataSourceSubAccount(),
			resourceAlertV2Type:              dataSourceAlertV2(),
			resourceLogShippingTokenType:     dataSourceLogShippingToken(),
			resourceDropFilterType:           dataSourceDropFilter(),
			resourceArchiveLogsType:          dataSourceArchiveLogs(),
			resourceRestoreLogsType:          dataSourceRestoreLogs(),
			resourceAuthenticationGroupsType: dataSourceAuthenticationGroups(),
			resourceKibanaObjectType:         dataSourceKibanaObject(),
		},
		ResourcesMap: map[string]*schema.Resource{
			resourceEndpointType:             resourceEndpoint(),
			resourceUserType:                 resourceUser(),
			resourceSubAccountType:           resourceSubAccount(),
			resourceAlertV2Type:              resourceAlertV2(),
			resourceLogShippingTokenType:     resourceLogShippingToken(),
			resourceDropFilterType:           resourceDropFilter(),
			resourceArchiveLogsType:          resourceArchiveLogs(),
			resourceRestoreLogsType:          resourceRestoreLogs(),
			resourceAuthenticationGroupsType: resourceAuthenticationGroups(),
			resourceKibanaObjectType:         resourceKibanaObject(),
		},
		ConfigureContextFunc: providerConfigureWrapper,
	}
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{providerApiToken: "Your API token", providerRegion: "Your logz.io region"}
}

func providerConfigure(d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	apiToken, ok := d.GetOk(providerApiToken)
	if !ok {
		return nil, diag.Errorf("can't find the %s, either set it in the provider or set the %s env var", providerApiToken, envLogzioApiToken)
	}
	region := d.Get(providerRegion).(string)
	regionCode := ""
	if region != "" && region != "us" {
		regionCode = fmt.Sprintf("-%s", region)
	}
	apiUrl := fmt.Sprintf(baseUrl, regionCode)

	config := Config{
		apiToken: apiToken.(string),
		baseUrl:  apiUrl,
	}
	return config, diag.Diagnostics{}
}

func providerConfigureWrapper(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return providerConfigure(d)
}
