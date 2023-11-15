package logzio

import (
	"context"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/logzio/logzio_terraform_client/grafana_contact_points"
	"reflect"
	"strings"
)

const (
	grafanaContactPointName                  = "name"
	grafanaContactPointUid                   = "uid"
	grafanaContactPointDisableResolveMessage = "disable_resolve_message"
	grafanaContactPointSettings              = "settings"

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

	grafanaContactPointEmailAddressSeparator = ";"

	grafanaContactPointRetryAttempts = 8
)

var notifiers = []grafanaContactPointNotifier{
	emailNotifier{},
	//googleChatNotifier{},
	//opsGenieNotifier{},
	//pagerDutyNotifier{},
	//slackNotifier{},
	//teamsNotifier{},
	//victorOpsNotifier{},
	//webhookNotifier{},
}

func resourceGrafanaContactPoint() *schema.Resource {
	resource := &schema.Resource{
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
		},
	}

	// Build list of available notifier fields, at least one has to be specified
	notifierFields := make([]string, len(notifiers))
	for i, n := range notifiers {
		notifierFields[i] = n.meta().field
	}

	for _, n := range notifiers {
		resource.Schema[n.meta().field] = &schema.Schema{
			Type:         schema.TypeList,
			Optional:     true,
			Elem:         n.schema(),
			ExactlyOneOf: notifierFields,
		}
	}

	return resource
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
	contactPoint, err := grafanaContactPointClient(m).GetGrafanaContactPointByUid(d.Id())

	if err != nil {
		tflog.Error(ctx, err.Error())
		if strings.Contains(err.Error(), "missing grafana contact point") {
			// If we were not able to find the resource - delete from state
			d.SetId("")
			return diag.Diagnostics{}
		} else {
			return diag.FromErr(err)
		}
	}

	err = setGrafanaContactPoint(d, contactPoint)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceGrafanaContactPointUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	updateContactPoint, err := getGrafanaContactPointFromSchema(d)
	if err != nil {
		return diag.FromErr(err)
	}

	err = grafanaContactPointClient(m).UpdateContactPoint(updateContactPoint)
	if err != nil {
		return diag.FromErr(err)
	}

	var diagRet diag.Diagnostics
	readErr := retry.Do(
		func() error {
			diagRet = resourceGrafanaContactPointRead(ctx, d, m)
			if diagRet.HasError() {
				return fmt.Errorf("received error from read grafana contact point")
			}

			return nil
		},
		retry.RetryIf(
			// Retry ONLY if the resource was not updated yet
			func(err error) bool {
				if err != nil {
					return false
				} else {
					// Check if the update shows on read
					// if not updated yet - retry
					grafanaContactPointFromSchema, _ := getGrafanaContactPointFromSchema(d)
					return !reflect.DeepEqual(updateContactPoint, grafanaContactPointFromSchema)
				}
			}),
		retry.DelayType(retry.BackOffDelay),
		retry.Attempts(grafanaContactPointRetryAttempts),
	)

	if readErr != nil {
		tflog.Error(ctx, "could not update schema")
		return diagRet
	}

	return nil
}

func resourceGrafanaContactPointDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := grafanaContactPointClient(m).DeleteGrafanaContactPoint(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func setGrafanaContactPoint(d *schema.ResourceData, contactPoint grafana_contact_points.GrafanaContactPoint) error {
	d.Set(grafanaContactPointName, contactPoint.Name)
	d.Set(grafanaContactPointUid, contactPoint.Uid)
	d.Set(grafanaContactPointDisableResolveMessage, contactPoint.DisableResolveMessage)
	for _, n := range notifiers {
		if contactPoint.Type == n.meta().typeStr {
			packed, err := n.getGrafanaContactPointFromObject(d, contactPoint)
			if err != nil {
				return err
			}
			d.Set(n.meta().field, []interface{}{packed})
			return nil
		}
	}

	return fmt.Errorf("could not find notifier")
}

func getGrafanaContactPointFromSchema(d *schema.ResourceData) (grafana_contact_points.GrafanaContactPoint, error) {
	for _, notifier := range notifiers {
		if point, ok := d.GetOk(notifier.meta().field); ok {
			uid := ""
			if v, okUid := d.GetOk(grafanaContactPointUid); okUid {
				uid = v.(string)
			}
			return unpackPointConfig(notifier,
				point.([]interface{}),
				d.Get(grafanaContactPointName).(string),
				d.Get(grafanaContactPointDisableResolveMessage).(bool),
				uid), nil
		}
	}

	return grafana_contact_points.GrafanaContactPoint{}, fmt.Errorf("could not find notifier")
}

func unpackPointConfig(n grafanaContactPointNotifier, data []interface{}, name string, disableResolveMessage bool, uid string) grafana_contact_points.GrafanaContactPoint {
	pt := n.getGrafanaContactPointFromSchema(data, name, disableResolveMessage, uid)
	// Treat settings like `omitempty`
	for k, v := range pt.Settings {
		if v == "" {
			delete(pt.Settings, k)
		}
	}
	return pt
}

func parseAddressStringToList(addressString string) []interface{} {
	arrStr := strings.Split(addressString, grafanaContactPointEmailAddressSeparator)
	var interfaceSlice = make([]interface{}, len(arrStr))
	for i, v := range arrStr {
		interfaceSlice[i] = v
	}

	return interfaceSlice
}

//package logzio
//
//import (
//	"context"
//	"fmt"
//	"github.com/avast/retry-go"
//	"github.com/hashicorp/terraform-plugin-log/tflog"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
//	"github.com/logzio/logzio_terraform_client/grafana_contact_points"
//	"github.com/logzio/logzio_terraform_provider/logzio/utils"
//	"github.com/stoewer/go-strcase"
//	"reflect"
//	"strings"
//)
//
//const (
//	grafanaContactPointName                  = "name"
//	grafanaContactPointUid                   = "uid"
//	grafanaContactPointDisableResolveMessage = "disable_resolve_message"
//	grafanaContactPointType                  = "type"
//
//	grafanaContactPointEmail            = "email"
//	grafanaContactPointEmailAddresses   = "addresses"
//	grafanaContactPointEmailSingleEmail = "single_email"
//	grafanaContactPointEmailMessage     = "message"
//
//	grafanaContactPointGoogleChat        = "googlechat"
//	grafanaContactPointGoogleChatUrl     = "url"
//	grafanaContactPointGoogleChatMessage = "message"
//
//	grafanaContactPointOpsgenie                 = "opsgenie"
//	grafanaContactPointOpsgenieApiUrl           = "api_url"
//	grafanaContactPointOpsgenieApiKey           = "api_key"
//	grafanaContactPointOpsgenieAutoClose        = "auto_close"
//	grafanaContactPointOpsgenieOverridePriority = "override_priority"
//	grafanaContactPointOpsgenieSendTagsAs       = "send_tags_as"
//	grafanaContactPointOpsgenieSendTagsTags     = "tags"
//	grafanaContactPointOpsgenieSendTagsDetails  = "details"
//	grafanaContactPointOpsgenieSendTagsBoth     = "both"
//
//	grafanaContactPointPagerduty               = "pagerduty"
//	grafanaContactPointPagerdutyClass          = "class"
//	grafanaContactPointPagerdutyComponent      = "component"
//	grafanaContactPointPagerdutyGroup          = "group"
//	grafanaContactPointPagerdutyIntegrationKey = "integration_key"
//	grafanaContactPointPagerdutySeverity       = "severity"
//	grafanaContactPointPagerdutySummary        = "summary"
//
//	grafanaContactPointSlack               = "slack"
//	grafanaContactPointSlackEndpointUrl    = "endpoint_url"
//	grafanaContactPointSlackMentionChannel = "mention_channel"
//	grafanaContactPointSlackMentionGroups  = "mention_groups"
//	grafanaContactPointSlackMentionUsers   = "mention_users"
//	grafanaContactPointSlackRecipient      = "recipient"
//	grafanaContactPointSlackText           = "text"
//	grafanaContactPointSlackTitle          = "title"
//	grafanaContactPointSlackToken          = "token"
//	grafanaContactPointSlackUrl            = "url"
//	grafanaContactPointSlackUsername       = "username"
//
//	grafanaContactPointMicrosoftTeams        = "teams"
//	grafanaContactPointMicrosoftTeamsMessage = "message"
//	grafanaContactPointMicrosoftTeamsUrl     = "url"
//
//	grafanaContactPointVictorops            = "victorops"
//	grafanaContactPointVictoropsMessageType = "message_type"
//	grafanaContactPointVictoropsUrl         = "url"
//
//	grafanaContactPointWebhook           = "webhook"
//	grafanaContactPointWebhookHttpMethod = "http_method"
//	grafanaContactPointWebhookMaxAlerts  = "max_alerts"
//	grafanaContactPointWebhookPassword   = "password"
//	grafanaContactPointWebhookUrl        = "url"
//	grafanaContactPointWebhookUsername   = "username"
//
//	grafanaContactPointEmailAddressSeparator = ";"
//
//	grafanaContactPointRetryAttempts = 8
//)
//
//func resourceGrafanaContactPoint() *schema.Resource {
//	return &schema.Resource{
//		CreateContext: resourceGrafanaContactPointCreate,
//		ReadContext:   resourceGrafanaContactPointRead,
//		UpdateContext: resourceGrafanaContactPointUpdate,
//		DeleteContext: resourceGrafanaContactPointDelete,
//		Importer: &schema.ResourceImporter{
//			StateContext: schema.ImportStatePassthroughContext,
//		},
//
//		Schema: map[string]*schema.Schema{
//			grafanaContactPointName: {
//				Type:     schema.TypeString,
//				Required: true,
//			},
//			grafanaContactPointUid: {
//				Type:     schema.TypeString,
//				Computed: true,
//			},
//			grafanaContactPointDisableResolveMessage: {
//				Type:     schema.TypeBool,
//				Optional: true,
//				Default:  false,
//			},
//			grafanaContactPointType: {
//				Type:             schema.TypeString,
//				Required:         true,
//				ForceNew:         true,
//				ValidateDiagFunc: utils.ValidateGrafanaContactPointType,
//			},
//			grafanaContactPointEmail: {
//				Type:     schema.TypeSet,
//				Optional: true,
//				MinItems: 1,
//				MaxItems: 1,
//				Elem: &schema.Resource{
//					Schema: map[string]*schema.Schema{
//						grafanaContactPointEmailAddresses: {
//							Type:     schema.TypeList,
//							Required: true,
//							Elem: &schema.Schema{
//								Type:         schema.TypeString,
//								ValidateFunc: validation.StringIsNotEmpty,
//							},
//						},
//						grafanaContactPointEmailSingleEmail: {
//							Type:     schema.TypeBool,
//							Optional: true,
//							Default:  false,
//						},
//						grafanaContactPointEmailMessage: {
//							Type:     schema.TypeString,
//							Optional: true,
//							Default:  "",
//						},
//					},
//				},
//			},
//			grafanaContactPointGoogleChat: {
//				Type:     schema.TypeSet,
//				Optional: true,
//				MinItems: 1,
//				MaxItems: 1,
//				Elem: &schema.Resource{
//					Schema: map[string]*schema.Schema{
//						grafanaContactPointGoogleChatUrl: {
//							Type:      schema.TypeString,
//							Required:  true,
//							Sensitive: true,
//						},
//						grafanaContactPointGoogleChatMessage: {
//							Type:     schema.TypeString,
//							Optional: true,
//						},
//					},
//				},
//			},
//			grafanaContactPointOpsgenie: {
//				Type:     schema.TypeSet,
//				Optional: true,
//				MinItems: 1,
//				MaxItems: 1,
//				Elem: &schema.Resource{
//					Schema: map[string]*schema.Schema{
//						grafanaContactPointOpsgenieApiUrl: {
//							Type:     schema.TypeString,
//							Optional: true,
//						},
//						grafanaContactPointOpsgenieApiKey: {
//							Type:      schema.TypeString,
//							Required:  true,
//							Sensitive: true,
//						},
//						grafanaContactPointOpsgenieAutoClose: {
//							Type:     schema.TypeBool,
//							Optional: true,
//						},
//						grafanaContactPointOpsgenieOverridePriority: {
//							Type:     schema.TypeBool,
//							Optional: true,
//						},
//						grafanaContactPointOpsgenieSendTagsAs: {
//							Type:     schema.TypeString,
//							Optional: true,
//							ValidateFunc: validation.StringInSlice([]string{grafanaContactPointOpsgenieSendTagsTags,
//								grafanaContactPointOpsgenieSendTagsDetails,
//								grafanaContactPointOpsgenieSendTagsBoth}, false),
//						},
//					},
//				},
//			},
//			grafanaContactPointPagerduty: {
//				Type:     schema.TypeSet,
//				Optional: true,
//				MinItems: 1,
//				MaxItems: 1,
//				Elem: &schema.Resource{
//					Schema: map[string]*schema.Schema{
//						grafanaContactPointPagerdutyClass: {
//							Type:     schema.TypeString,
//							Optional: true,
//						},
//						grafanaContactPointPagerdutyComponent: {
//							Type:     schema.TypeString,
//							Optional: true,
//						},
//						grafanaContactPointPagerdutyGroup: {
//							Type:     schema.TypeString,
//							Optional: true,
//						},
//						grafanaContactPointPagerdutyIntegrationKey: {
//							Type:      schema.TypeString,
//							Required:  true,
//							Sensitive: true,
//						},
//						grafanaContactPointPagerdutySeverity: {
//							Type:     schema.TypeString,
//							Optional: true,
//						},
//						grafanaContactPointPagerdutySummary: {
//							Type:     schema.TypeString,
//							Optional: true,
//						},
//					},
//				},
//			},
//			grafanaContactPointSlack: {
//				Type:     schema.TypeSet,
//				Optional: true,
//				MinItems: 1,
//				MaxItems: 1,
//				Elem: &schema.Resource{
//					Schema: map[string]*schema.Schema{
//						grafanaContactPointSlackEndpointUrl: {
//							Type:     schema.TypeString,
//							Optional: true,
//						},
//						grafanaContactPointSlackMentionChannel: {
//							Type:     schema.TypeString,
//							Optional: true,
//						},
//						grafanaContactPointSlackMentionGroups: {
//							Type:     schema.TypeString,
//							Optional: true,
//						},
//						grafanaContactPointSlackMentionUsers: {
//							Type:     schema.TypeString,
//							Optional: true,
//						},
//						grafanaContactPointSlackRecipient: {
//							Type:     schema.TypeString,
//							Optional: true,
//						},
//						grafanaContactPointSlackText: {
//							Type:     schema.TypeString,
//							Optional: true,
//						},
//						grafanaContactPointSlackTitle: {
//							Type:     schema.TypeString,
//							Optional: true,
//						},
//						grafanaContactPointSlackToken: {
//							Type:      schema.TypeString,
//							Optional:  true,
//							Sensitive: true,
//						},
//						grafanaContactPointSlackUrl: {
//							Type:      schema.TypeString,
//							Optional:  true,
//							Sensitive: true,
//						},
//						grafanaContactPointSlackUsername: {
//							Type:     schema.TypeString,
//							Optional: true,
//						},
//					},
//				},
//			},
//			grafanaContactPointMicrosoftTeams: {
//				Type:     schema.TypeSet,
//				Optional: true,
//				MinItems: 1,
//				MaxItems: 1,
//				Elem: &schema.Resource{
//					Schema: map[string]*schema.Schema{
//						grafanaContactPointMicrosoftTeamsMessage: {
//							Type:     schema.TypeString,
//							Optional: true,
//						},
//						grafanaContactPointMicrosoftTeamsUrl: {
//							Type:      schema.TypeString,
//							Required:  true,
//							Sensitive: true,
//						},
//					},
//				},
//			},
//			grafanaContactPointVictorops: {
//				Type:     schema.TypeSet,
//				Optional: true,
//				MinItems: 1,
//				MaxItems: 1,
//				Elem: &schema.Resource{
//					Schema: map[string]*schema.Schema{
//						grafanaContactPointVictoropsMessageType: {
//							Type:     schema.TypeString,
//							Optional: true,
//						},
//						grafanaContactPointVictoropsUrl: {
//							Type:     schema.TypeString,
//							Required: true,
//						},
//					},
//				},
//			},
//			grafanaContactPointWebhook: {
//				Type:     schema.TypeSet,
//				Optional: true,
//				MinItems: 1,
//				MaxItems: 1,
//				Elem: &schema.Resource{
//					Schema: map[string]*schema.Schema{
//						grafanaContactPointWebhookHttpMethod: {
//							Type:     schema.TypeString,
//							Optional: true,
//						},
//						grafanaContactPointWebhookMaxAlerts: {
//							Type:     schema.TypeString,
//							Optional: true,
//						},
//						grafanaContactPointWebhookPassword: {
//							Type:      schema.TypeString,
//							Optional:  true,
//							Sensitive: true,
//						},
//						grafanaContactPointWebhookUrl: {
//							Type:     schema.TypeString,
//							Required: true,
//						},
//						grafanaContactPointWebhookUsername: {
//							Type:     schema.TypeString,
//							Optional: true,
//						},
//					},
//				},
//			},
//		},
//	}
//}
//
//func grafanaContactPointClient(m interface{}) *grafana_contact_points.GrafanaContactPointClient {
//	var client *grafana_contact_points.GrafanaContactPointClient
//	client, _ = grafana_contact_points.New(m.(Config).apiToken, m.(Config).baseUrl)
//	return client
//}
//
//func resourceGrafanaContactPointCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
//	createContactPoint, err := getGrafanaContactPointFromSchema(d)
//	if err != nil {
//		return diag.FromErr(err)
//	}
//
//	contactPoint, err := grafanaContactPointClient(m).CreateGrafanaContactPoint(createContactPoint)
//	if err != nil {
//		return diag.FromErr(err)
//	}
//
//	d.SetId(contactPoint.Uid)
//	// when using GET, sensitive fields return as "[REDACTED]" so we can't set them from read, we need to do it at this point
//	setSensitiveFields(d, contactPoint)
//	return resourceGrafanaContactPointRead(ctx, d, m)
//}
//
//func resourceGrafanaContactPointRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
//	contactPoint, err := grafanaContactPointClient(m).GetGrafanaContactPointByUid(d.Id())
//
//	if err != nil {
//		tflog.Error(ctx, err.Error())
//		if strings.Contains(err.Error(), "missing grafana contact point") {
//			// If we were not able to find the resource - delete from state
//			d.SetId("")
//			return diag.Diagnostics{}
//		} else {
//			return diag.FromErr(err)
//		}
//	}
//
//	setGrafanaContactPoint(d, contactPoint)
//	return nil
//}
//
//func resourceGrafanaContactPointUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
//	updateContactPoint, err := getGrafanaContactPointFromSchema(d)
//	if err != nil {
//		return diag.FromErr(err)
//	}
//
//	err = grafanaContactPointClient(m).UpdateContactPoint(updateContactPoint)
//	if err != nil {
//		return diag.FromErr(err)
//	}
//
//	var diagRet diag.Diagnostics
//	readErr := retry.Do(
//		func() error {
//			diagRet = resourceGrafanaContactPointRead(ctx, d, m)
//			if diagRet.HasError() {
//				return fmt.Errorf("received error from read grafana contact point")
//			}
//
//			return nil
//		},
//		retry.RetryIf(
//			// Retry ONLY if the resource was not updated yet
//			func(err error) bool {
//				if err != nil {
//					return false
//				} else {
//					// Check if the update shows on read
//					// if not updated yet - retry
//					grafanaContactPointFromSchema, _ := getGrafanaContactPointFromSchema(d)
//					return !reflect.DeepEqual(updateContactPoint, grafanaContactPointFromSchema)
//				}
//			}),
//		retry.DelayType(retry.BackOffDelay),
//		retry.Attempts(grafanaContactPointRetryAttempts),
//	)
//
//	if readErr != nil {
//		tflog.Error(ctx, "could not update schema")
//		return diagRet
//	}
//
//	return nil
//}
//
//func resourceGrafanaContactPointDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
//	err := grafanaContactPointClient(m).DeleteGrafanaContactPoint(d.Id())
//	if err != nil {
//		return diag.FromErr(err)
//	}
//
//	d.SetId("")
//	return nil
//}
//
//func setSensitiveFields(d *schema.ResourceData, contactPoint grafana_contact_points.GrafanaContactPoint) {
//	var sensitiveFields []string
//	switch contactPoint.Type {
//	case grafanaContactPointGoogleChat:
//		sensitiveFields = []string{grafanaContactPointGoogleChatUrl}
//	case grafanaContactPointOpsgenie:
//		sensitiveFields = []string{grafanaContactPointOpsgenieApiKey}
//	case grafanaContactPointPagerduty:
//		sensitiveFields = []string{grafanaContactPointPagerdutyIntegrationKey}
//	case grafanaContactPointSlack:
//		sensitiveFields = []string{grafanaContactPointSlackToken, grafanaContactPointSlackUrl}
//	case grafanaContactPointMicrosoftTeams:
//		sensitiveFields = []string{grafanaContactPointMicrosoftTeamsUrl}
//	case grafanaContactPointWebhook:
//		sensitiveFields = []string{grafanaContactPointWebhookPassword}
//	default:
//		return
//	}
//
//	prefix := fmt.Sprintf("%s.0.", contactPoint.Type)
//	setFieldsFromApiKey(d, prefix, sensitiveFields, contactPoint.Settings)
//}
//
//func setGrafanaContactPoint(d *schema.ResourceData, contactPoint grafana_contact_points.GrafanaContactPoint) {
//	d.Set(grafanaContactPointName, contactPoint.Name)
//	d.Set(grafanaContactPointUid, contactPoint.Uid)
//	d.Set(grafanaContactPointDisableResolveMessage, contactPoint.DisableResolveMessage)
//	d.Set(grafanaContactPointType, contactPoint.Type)
//
//	var fieldsToSet []string
//	switch contactPoint.Type {
//	case grafanaContactPointEmail:
//		fieldsToSet = []string{grafanaContactPointEmailAddresses, grafanaContactPointEmailSingleEmail, grafanaContactPointEmailMessage}
//	case grafanaContactPointGoogleChat:
//		fieldsToSet = []string{grafanaContactPointGoogleChatMessage}
//	case grafanaContactPointOpsgenie:
//		fieldsToSet = []string{grafanaContactPointOpsgenieApiUrl,
//			grafanaContactPointOpsgenieAutoClose,
//			grafanaContactPointOpsgenieOverridePriority,
//			grafanaContactPointOpsgenieSendTagsAs,
//		}
//	case grafanaContactPointPagerduty:
//		fieldsToSet = []string{
//			grafanaContactPointPagerdutyClass,
//			grafanaContactPointPagerdutyComponent,
//			grafanaContactPointPagerdutyGroup,
//			grafanaContactPointPagerdutySeverity,
//			grafanaContactPointPagerdutySummary,
//		}
//	case grafanaContactPointSlack:
//		fieldsToSet = []string{
//			grafanaContactPointSlackEndpointUrl,
//			grafanaContactPointSlackMentionChannel,
//			grafanaContactPointSlackMentionGroups,
//			grafanaContactPointSlackMentionUsers,
//			grafanaContactPointSlackRecipient,
//			grafanaContactPointSlackText,
//			grafanaContactPointSlackTitle,
//			grafanaContactPointSlackUsername,
//		}
//	case grafanaContactPointMicrosoftTeams:
//		fieldsToSet = []string{grafanaContactPointMicrosoftTeamsMessage}
//	case grafanaContactPointVictorops:
//		fieldsToSet = []string{grafanaContactPointVictoropsMessageType, grafanaContactPointVictoropsUrl}
//	case grafanaContactPointWebhook:
//		fieldsToSet = []string{
//			grafanaContactPointWebhookHttpMethod,
//			grafanaContactPointWebhookMaxAlerts,
//			grafanaContactPointWebhookUrl,
//			grafanaContactPointWebhookUsername,
//		}
//	default:
//		panic("unidentified Grafana Contact Point type!")
//
//	}
//
//	prefix := fmt.Sprintf("%s.0.", contactPoint.Type)
//	setFieldsFromApiKey(d, prefix, fieldsToSet, contactPoint.Settings)
//}
//
//func setFieldsFromApiKey(d *schema.ResourceData, prefix string, fieldsToSet []string, settings map[string]interface{}) {
//	for _, fieldToSet := range fieldsToSet {
//		apiKey := strcase.LowerCamelCase(fieldToSet)
//		if val, ok := settings[apiKey]; ok {
//			switch fieldToSet {
//			case grafanaContactPointEmailAddresses:
//				d.Set(prefix+fieldToSet, parseAddressStringToList(val.(string)))
//			default:
//				d.Set(prefix+fieldToSet, val)
//			}
//		}
//	}
//}
//
//func getGrafanaContactPointFromSchema(d *schema.ResourceData) (grafana_contact_points.GrafanaContactPoint, error) {
//	contactPoint := grafana_contact_points.GrafanaContactPoint{
//		Name:                  d.Get(grafanaContactPointName).(string),
//		Type:                  d.Get(grafanaContactPointType).(string),
//		DisableResolveMessage: d.Get(grafanaContactPointDisableResolveMessage).(bool),
//	}
//
//	if uid, ok := d.GetOk(grafanaContactPointUid); ok {
//		contactPoint.Uid = uid.(string)
//	}
//
//	settings, err := utils.ParseTypeSetToMap(d, contactPoint.Type)
//	if err != nil {
//		return contactPoint, err
//	}
//
//	// in tf we use snake case for keys, but in the api uses lower camel case, so we need to convert the relevant fields
//	var convertKeys []string
//	switch contactPoint.Type {
//	case grafanaContactPointEmail:
//		convertKeys = []string{grafanaContactPointEmailSingleEmail}
//		if val, ok := settings[grafanaContactPointEmailAddresses]; ok {
//			settings[grafanaContactPointEmailAddresses] = parseAddressListToString(val.([]interface{}))
//		}
//	case grafanaContactPointOpsgenie:
//		convertKeys = []string{grafanaContactPointOpsgenieApiUrl,
//			grafanaContactPointOpsgenieApiKey,
//			grafanaContactPointOpsgenieAutoClose,
//			grafanaContactPointOpsgenieOverridePriority,
//			grafanaContactPointOpsgenieSendTagsAs,
//		}
//	case grafanaContactPointPagerduty:
//		convertKeys = []string{grafanaContactPointPagerdutyIntegrationKey}
//	case grafanaContactPointSlack:
//		convertKeys = []string{grafanaContactPointSlackEndpointUrl,
//			grafanaContactPointSlackMentionChannel,
//			grafanaContactPointSlackMentionGroups,
//			grafanaContactPointSlackMentionUsers,
//		}
//	case grafanaContactPointVictorops:
//		convertKeys = []string{grafanaContactPointVictoropsMessageType}
//	case grafanaContactPointWebhook:
//		convertKeys = []string{grafanaContactPointWebhookHttpMethod,
//			grafanaContactPointWebhookMaxAlerts,
//		}
//	default:
//		panic("unidentified Grafana Contact Point type!")
//	}
//
//	for _, key := range convertKeys {
//		convertSettingsMapToApiKeys(settings, key)
//	}
//
//	contactPoint.Settings = settings
//	return contactPoint, nil
//}
//
//func convertSettingsMapToApiKeys(settings map[string]interface{}, schemaKey string) {
//	if val, ok := settings[schemaKey]; ok {
//		apiKey := strcase.LowerCamelCase(schemaKey)
//		settings[apiKey] = val
//		delete(settings, schemaKey)
//	}
//}
//
//func parseAddressListToString(addressList []interface{}) string {
//	strArr := make([]string, len(addressList))
//	for i, v := range addressList {
//		strArr[i] = fmt.Sprintf("%v", v)
//	}
//
//	return strings.Join(strArr, grafanaContactPointEmailAddressSeparator)
//}
//
//func parseAddressStringToList(addressString string) []interface{} {
//	arrStr := strings.Split(addressString, grafanaContactPointEmailAddressSeparator)
//	var interfaceSlice = make([]interface{}, len(arrStr))
//	for i, v := range arrStr {
//		interfaceSlice[i] = v
//	}
//
//	return interfaceSlice
//}
