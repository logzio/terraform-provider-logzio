package logzio

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/logzio/logzio_terraform_client/grafana_contact_points"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
	"github.com/stoewer/go-strcase"
)

const (
	grafanaContactPointName                  = "name"
	grafanaContactPointUid                   = "uid"
	grafanaContactPointDisableResolveMessage = "disable_resolve_message"
	grafanaContactPointType                  = "type"

	grafanaContactPointEmail            = "email"
	grafanaContactPointEmailAddresses   = "addresses"
	grafanaContactPointEmailSingleEmail = "single_email"
	grafanaContactPointEmailMessage     = "message"

	grafanaContactPointGoogleChat        = "googlechat"
	grafanaContactPointGoogleChatUrl     = "url"
	grafanaContactPointGoogleChatMessage = "message"

	grafanaContactPointOpsgenie                 = "opsgenie"
	grafanaContactPointOpsgenieApiUrl           = "api_url"
	grafanaContactPointOpsgenieApiKey           = "api_key"
	grafanaContactPointOpsgenieAutoClose        = "auto_close"
	grafanaContactPointOpsgenieOverridePriority = "override_priority"
	grafanaContactPointOpsgenieSendTagsAs       = "send_tags_as"
	grafanaContactPointOpsgenieSendTagsTags     = "tags"
	grafanaContactPointOpsgenieSendTagsDetails  = "details"
	grafanaContactPointOpsgenieSendTagsBoth     = "both"

	grafanaContactPointPagerduty               = "pagerduty"
	grafanaContactPointPagerdutyClass          = "class"
	grafanaContactPointPagerdutyComponent      = "component"
	grafanaContactPointPagerdutyGroup          = "group"
	grafanaContactPointPagerdutyIntegrationKey = "integration_key"
	grafanaContactPointPagerdutySeverity       = "severity"
	grafanaContactPointPagerdutySummary        = "summary"

	grafanaContactPointSlack               = "slack"
	grafanaContactPointSlackEndpointUrl    = "endpoint_url"
	grafanaContactPointSlackMentionChannel = "mention_channel"
	grafanaContactPointSlackMentionGroups  = "mention_groups"
	grafanaContactPointSlackMentionUsers   = "mention_users"
	grafanaContactPointSlackRecipient      = "recipient"
	grafanaContactPointSlackText           = "text"
	grafanaContactPointSlackTitle          = "title"
	grafanaContactPointSlackToken          = "token"
	grafanaContactPointSlackUrl            = "url"
	grafanaContactPointSlackUsername       = "username"

	grafanaContactPointMicrosoftTeams        = "teams"
	grafanaContactPointMicrosoftTeamsMessage = "message"
	grafanaContactPointMicrosoftTeamsUrl     = "url"

	grafanaContactPointVictorops            = "victorops"
	grafanaContactPointVictoropsMessageType = "message_type"
	grafanaContactPointVictoropsUrl         = "url"

	grafanaContactPointWebhook           = "webhook"
	grafanaContactPointWebhookHttpMethod = "http_method"
	grafanaContactPointWebhookMaxAlerts  = "max_alerts"
	grafanaContactPointWebhookPassword   = "password"
	grafanaContactPointWebhookUrl        = "url"
	grafanaContactPointWebhookUsername   = "username"
)

func resourceGrafanaContactPoint() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGrafanaContactPointCreate,
		ReadContext:   resourceGrafanaContactPointRead,
		UpdateContext: resourceGrafanaContactPointUpdate,
		DeleteContext: resourceGrafanaContactPointDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			grafanaContactPointName: {
				Type:     schema.TypeString,
				Required: true,
			},
			grafanaContactPointUid: {
				Type:     schema.TypeString,
				Computed: true,
			},
			grafanaContactPointDisableResolveMessage: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			grafanaContactPointType: {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: utils.ValidateGrafanaContactPointType,
			},
			grafanaContactPointEmail: {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						grafanaContactPointEmailAddresses: {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
						grafanaContactPointEmailSingleEmail: {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						grafanaContactPointEmailMessage: {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
					},
				},
			},
			grafanaContactPointGoogleChat: {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						grafanaContactPointGoogleChatUrl: {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true, // TODO - note set
						},
						grafanaContactPointGoogleChatMessage: {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			grafanaContactPointOpsgenie: {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						grafanaContactPointOpsgenieApiUrl: {
							Type:     schema.TypeString,
							Optional: true,
						},
						grafanaContactPointOpsgenieApiKey: {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
						grafanaContactPointOpsgenieAutoClose: {
							Type:     schema.TypeBool,
							Optional: true,
						},
						grafanaContactPointOpsgenieOverridePriority: {
							Type:     schema.TypeBool,
							Optional: true,
						},
						grafanaContactPointOpsgenieSendTagsAs: {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{grafanaContactPointOpsgenieSendTagsTags,
								grafanaContactPointOpsgenieSendTagsDetails,
								grafanaContactPointOpsgenieSendTagsBoth}, false),
						},
					},
				},
			},
			grafanaContactPointPagerduty: {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						grafanaContactPointPagerdutyClass: {
							Type:     schema.TypeString,
							Optional: true,
						},
						grafanaContactPointPagerdutyComponent: {
							Type:     schema.TypeString,
							Optional: true,
						},
						grafanaContactPointPagerdutyGroup: {
							Type:     schema.TypeString,
							Optional: true,
						},
						grafanaContactPointPagerdutyIntegrationKey: {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
						grafanaContactPointPagerdutySeverity: {
							Type:     schema.TypeString,
							Optional: true,
						},
						grafanaContactPointPagerdutySummary: {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			grafanaContactPointSlack: {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						grafanaContactPointSlackEndpointUrl: {
							Type:     schema.TypeString,
							Optional: true,
						},
						grafanaContactPointSlackMentionChannel: {
							Type:     schema.TypeString,
							Optional: true,
						},
						grafanaContactPointSlackMentionGroups: {
							Type:     schema.TypeString,
							Optional: true,
						},
						grafanaContactPointSlackMentionUsers: {
							Type:     schema.TypeString,
							Optional: true,
						},
						grafanaContactPointSlackRecipient: {
							Type:     schema.TypeString,
							Optional: true,
						},
						grafanaContactPointSlackText: {
							Type:     schema.TypeString,
							Optional: true,
						},
						grafanaContactPointSlackTitle: {
							Type:     schema.TypeString,
							Optional: true,
						},
						grafanaContactPointSlackToken: {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						grafanaContactPointSlackUrl: {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						grafanaContactPointSlackUsername: {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			grafanaContactPointMicrosoftTeams: {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						grafanaContactPointMicrosoftTeamsMessage: {
							Type:     schema.TypeString,
							Optional: true,
						},
						grafanaContactPointMicrosoftTeamsUrl: {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
					},
				},
			},
			grafanaContactPointVictorops: {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						grafanaContactPointVictoropsMessageType: {
							Type:     schema.TypeString,
							Optional: true,
						},
						grafanaContactPointVictoropsUrl: {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			grafanaContactPointWebhook: {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						grafanaContactPointWebhookHttpMethod: {
							Type:     schema.TypeString,
							Optional: true,
						},
						grafanaContactPointWebhookMaxAlerts: {
							Type:     schema.TypeInt,
							Optional: true,
						},
						grafanaContactPointWebhookPassword: {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						grafanaContactPointWebhookUrl: {
							Type:     schema.TypeString,
							Required: true,
						},
						grafanaContactPointWebhookUsername: {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func grafanaContactPointClient(m interface{}) *grafana_contact_points.GrafanaContactPointClient {
	var client *grafana_contact_points.GrafanaContactPointClient
	client, _ = grafana_contact_points.New(m.(Config).apiToken, m.(Config).baseUrl)
	return client
}

func resourceGrafanaContactPointCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	createContactPoint, err := getGrafanaContactPointFromSchema(d)
	if err != nil {
		return diag.FromErr(err)
	}

	contactPoint, err := grafanaContactPointClient(m).CreateGrafanaContactPoint(createContactPoint)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(contactPoint.Uid)
	return resourceGrafanaContactPointRead(ctx, d, m)
}

func resourceGrafanaContactPointRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := d.Id()
}

func getGrafanaContactPointFromSchema(d *schema.ResourceData) (grafana_contact_points.GrafanaContactPoint, error) {
	contactPoint := grafana_contact_points.GrafanaContactPoint{
		Name:                  d.Get(grafanaContactPointName).(string),
		Type:                  d.Get(grafanaContactPointType).(string),
		DisableResolveMessage: d.Get(grafanaContactPointDisableResolveMessage).(bool),
	}

	settings, err := utils.ParseTypeSetToMap(d, contactPoint.Type)
	if err != nil {
		return contactPoint, err
	}

	// in tf we use snake case for keys, but in the api uses lower camel case, so we need to convert the relevant fields
	var convertKeys []string
	switch contactPoint.Type {
	case grafanaContactPointEmail:
		convertKeys = []string{grafanaContactPointEmailSingleEmail}
	case grafanaContactPointOpsgenie:
		convertKeys = []string{grafanaContactPointOpsgenieApiUrl,
			grafanaContactPointOpsgenieApiKey,
			grafanaContactPointOpsgenieAutoClose,
			grafanaContactPointOpsgenieOverridePriority,
			grafanaContactPointOpsgenieSendTagsAs,
		}
	case grafanaContactPointPagerduty:
		convertKeys = []string{grafanaContactPointPagerdutyIntegrationKey}
	case grafanaContactPointSlack:
		convertKeys = []string{grafanaContactPointSlackEndpointUrl,
			grafanaContactPointSlackMentionChannel,
			grafanaContactPointSlackMentionGroups,
			grafanaContactPointSlackMentionUsers,
		}
	case grafanaContactPointVictorops:
		convertKeys = []string{grafanaContactPointVictoropsMessageType}
	case grafanaContactPointWebhook:
		convertKeys = []string{grafanaContactPointWebhookHttpMethod,
			grafanaContactPointWebhookMaxAlerts,
		}
	}

	for _, key := range convertKeys {
		convertSettingsMapToApiKeys(settings, key)
	}

	contactPoint.Settings = settings
	return contactPoint, nil
}

func convertSettingsMapToApiKeys(settings map[string]interface{}, schemaKey string) {
	if val, ok := settings[schemaKey]; ok {
		apiKey := strcase.LowerCamelCase(schemaKey)
		settings[apiKey] = val
		delete(settings, schemaKey)
	}
}
