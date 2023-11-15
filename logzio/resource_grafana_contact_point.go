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

	grafanaContactPointPagerduty                 = "pagerduty"
	grafanaContactPointPagerdutyClass            = "class"
	grafanaContactPointPagerdutyComponent        = "component"
	grafanaContactPointPagerdutyGroup            = "group"
	grafanaContactPointPagerdutyIntegrationKey   = "integration_key"
	grafanaContactPointPagerdutySummary          = "summary"
	grafanaContactPointPagerdutySeverity         = "severity"
	grafanaContactPointPagerdutySeverityInfo     = "info"
	grafanaContactPointPagerdutySeverityWarning  = "warning"
	grafanaContactPointPagerdutySeverityError    = "error"
	grafanaContactPointPagerdutySeverityCritical = "critical"

	grafanaContactPointSlack                      = "slack"
	grafanaContactPointSlackEndpointUrl           = "endpoint_url"
	grafanaContactPointSlackMentionChannel        = "mention_channel"
	grafanaContactPointSlackMentionChannelHere    = "here"
	grafanaContactPointSlackMentionChannelChannel = "channel"
	grafanaContactPointSlackMentionChannelDisable = ""
	grafanaContactPointSlackMentionGroups         = "mention_groups"
	grafanaContactPointSlackMentionUsers          = "mention_users"
	grafanaContactPointSlackRecipient             = "recipient"
	grafanaContactPointSlackText                  = "text"
	grafanaContactPointSlackTitle                 = "title"
	grafanaContactPointSlackToken                 = "token"
	grafanaContactPointSlackUrl                   = "url"
	grafanaContactPointSlackUsername              = "username"

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
	googleChatNotifier{},
	opsGenieNotifier{},
	pagerDutyNotifier{},
	slackNotifier{},
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
	for k, v := range pt.Settings {
		if v == "" {
			delete(pt.Settings, k)
		}
	}
	return pt
}
