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

	grafanaContactPointVictorops                    = "victorops"
	grafanaContactPointVictoropsMessageType         = "message_type"
	grafanaContactPointVictoropsMessageTypeCritical = "CRITICAL"
	grafanaContactPointVictoropsMessageTypeWarning  = "WARNING"
	grafanaContactPointVictoropsMessageTypeNone     = ""
	grafanaContactPointVictoropsUrl                 = "url"

	grafanaContactPointWebhook           = "webhook"
	grafanaContactPointWebhookHttpMethod = "http_method"
	grafanaContactPointWebhookHttpPut    = "PUT"
	grafanaContactPointWebhookHttpPost   = "POST"
	grafanaContactPointWebhookMaxAlerts  = "max_alerts"
	grafanaContactPointWebhookPassword   = "password"
	grafanaContactPointWebhookUrl        = "url"
	grafanaContactPointWebhookUsername   = "username"

	grafanaContactPointEmailAddressSeparator = ";"
	grafanaContactPointUidsSeparator         = ";"

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
			AtLeastOneOf: notifierFields,
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
	createContactPoints, err := getGrafanaContactPointsFromSchema(d)
	if err != nil {
		return diag.FromErr(err)
	}

	uids := make([]string, 0, len(createContactPoints))
	for _, cp := range createContactPoints {
		contactPoint, err := grafanaContactPointClient(m).CreateGrafanaContactPoint(cp)
		if err != nil {
			return diag.FromErr(err)
		}
		uids = append(uids, contactPoint.Uid)
	}

	d.SetId(createUid(uids))
	return resourceGrafanaContactPointRead(ctx, d, m)
}

func resourceGrafanaContactPointRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	uidsToFetch := getUidsToFetch(d.Id())
	contactPoints := []grafana_contact_points.GrafanaContactPoint{}
	for _, uid := range uidsToFetch {
		contactPoint, err := grafanaContactPointClient(m).GetGrafanaContactPointByUid(uid)

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
		contactPoints = append(contactPoints, contactPoint)
	}

	err := setGrafanaContactPoints(d, contactPoints)
	if err != nil {
		return diag.FromErr(err)
	}

	uids := make([]string, 0, len(contactPoints))
	for _, p := range contactPoints {
		uids = append(uids, p.Uid)
	}

	d.SetId(createUid(uids))

	return nil
}

func resourceGrafanaContactPointUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	updateContactPoints, err := getGrafanaContactPointsFromSchema(d)
	if err != nil {
		return diag.FromErr(err)
	}
	existingUIDs := getUidsToFetch(d.Id())
	unprocessedUIDs := toUIDSet(existingUIDs)
	newUIDs := make([]string, 0, len(updateContactPoints))

	for _, contactPointToUpdate := range updateContactPoints {
		delete(unprocessedUIDs, contactPointToUpdate.Uid)
		err = grafanaContactPointClient(m).UpdateContactPoint(contactPointToUpdate)
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
						grafanaContactPointFromSchema, _ := getGrafanaContactPointsFromSchema(d)
						return !reflect.DeepEqual(updateContactPoints, grafanaContactPointFromSchema)
					}
				}),
			retry.DelayType(retry.BackOffDelay),
			retry.Attempts(grafanaContactPointRetryAttempts),
		)

		if readErr != nil {
			tflog.Error(ctx, "could not update schema")
			return diagRet
		}

		newUIDs = append(newUIDs, contactPointToUpdate.Uid)
	}

	// Any UIDs still left in the state that we haven't seen must map to deleted receivers.
	// Delete them on the server and drop them from state.
	for u := range unprocessedUIDs {
		if err := grafanaContactPointClient(m).DeleteGrafanaContactPoint(u); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(createUid(newUIDs))

	return nil
}

func resourceGrafanaContactPointDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	uids := getUidsToFetch(d.Id())
	for _, uid := range uids {
		err := grafanaContactPointClient(m).DeleteGrafanaContactPoint(uid)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId("")
	return nil
}

func setGrafanaContactPoints(d *schema.ResourceData, contactPoints []grafana_contact_points.GrafanaContactPoint) error {
	pointsPerNotifier := map[grafanaContactPointNotifier][]interface{}{}
	for _, contactPoint := range contactPoints {
		d.Set(grafanaContactPointName, contactPoint.Name)
		for _, n := range notifiers {
			if contactPoint.Type == n.meta().typeStr {
				packed, err := n.getGrafanaContactPointFromObject(d, contactPoint)
				if err != nil {
					return err
				}
				pointsPerNotifier[n] = append(pointsPerNotifier[n], packed)
				continue
			}
		}
	}

	for n, pts := range pointsPerNotifier {
		d.Set(n.meta().field, pts)
	}

	return nil
}

func getGrafanaContactPointsFromSchema(d *schema.ResourceData) ([]grafana_contact_points.GrafanaContactPoint, error) {
	contactPoints := make([]grafana_contact_points.GrafanaContactPoint, 0)
	for _, notifier := range notifiers {
		if points, ok := d.GetOk(notifier.meta().field); ok {
			for _, p := range points.([]interface{}) {
				cp := unpackPointConfig(notifier, p, d.Get(grafanaContactPointName).(string))
				contactPoints = append(contactPoints, cp)
			}

		}
	}

	return contactPoints, nil
}

func unpackPointConfig(n grafanaContactPointNotifier, data interface{}, name string) grafana_contact_points.GrafanaContactPoint {
	pt := n.getGrafanaContactPointFromSchema(data, name)
	for k, v := range pt.Settings {
		if v == "" {
			delete(pt.Settings, k)
		}
	}
	return pt
}

func createUid(uids []string) string {
	return strings.Join(uids, grafanaContactPointUidsSeparator)
}

func getUidsToFetch(uidsStr string) []string {
	return strings.Split(uidsStr, grafanaContactPointUidsSeparator)
}

func toUIDSet(uids []string) map[string]bool {
	set := map[string]bool{}
	for _, uid := range uids {
		set[uid] = true
	}
	return set
}
