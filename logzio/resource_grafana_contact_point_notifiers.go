package logzio

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/logzio/logzio_terraform_client/grafana_contact_points"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
	"github.com/stoewer/go-strcase"
	"strings"
)

type grafanaContactPointNotifier interface {
	meta() grafanaContactPointNotifierMeta
	schema() *schema.Resource
	getGrafanaContactPointFromObject(d *schema.ResourceData, contactPoint grafana_contact_points.GrafanaContactPoint) (interface{}, error)
	getGrafanaContactPointFromSchema(raw []interface{}, name string, disableResolveMessage bool, uid string) grafana_contact_points.GrafanaContactPoint
}

type grafanaContactPointNotifierMeta struct {
	field        string
	typeStr      string
	secureFields []string
}

func getAddressFromSchema(addresses []interface{}) string {
	strSlice := utils.ParseInterfaceSliceToStringSlice(addresses)
	return strings.Join(strSlice, grafanaContactPointEmailAddressSeparator)
}

func addressStringToStringList(addresses string) []string {
	return strings.Split(addresses, grafanaContactPointEmailAddressSeparator)
}

func getSecuredFieldsFromSchema(notifier map[string]interface{}, secureFields []string, typeStr string, d *schema.ResourceData) {
	for _, tfKey := range secureFields {
		if conf, ok := d.GetOk(typeStr); ok && conf != nil {
			notifier[tfKey] = conf.([]interface{})[0].(map[string]interface{})[tfKey]
		}
	}
}

type emailNotifier struct{}

var _ grafanaContactPointNotifier = (*emailNotifier)(nil)

func (e emailNotifier) meta() grafanaContactPointNotifierMeta {
	return grafanaContactPointNotifierMeta{
		field:   grafanaContactPointEmail,
		typeStr: grafanaContactPointEmail,
	}
}

func (e emailNotifier) schema() *schema.Resource {
	r := &schema.Resource{
		Schema: map[string]*schema.Schema{},
	}
	r.Schema[grafanaContactPointEmailAddresses] = &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		Elem: &schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
	r.Schema[grafanaContactPointEmailSingleEmail] = &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	}
	r.Schema[grafanaContactPointEmailMessage] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Default:  "",
	}
	return r
}

func (e emailNotifier) getGrafanaContactPointFromObject(d *schema.ResourceData, contactPoint grafana_contact_points.GrafanaContactPoint) (interface{}, error) {
	notifier := make(map[string]interface{}, 0)
	if v, ok := contactPoint.Settings[grafanaContactPointEmailAddresses]; ok && v != nil {
		notifier[grafanaContactPointEmailAddresses] = addressStringToStringList(v.(string))
	}
	if v, ok := contactPoint.Settings[strcase.LowerCamelCase(grafanaContactPointEmailSingleEmail)]; ok && v != nil {
		notifier[grafanaContactPointEmailSingleEmail] = v.(bool)
	}
	if v, ok := contactPoint.Settings[grafanaContactPointEmailMessage]; ok && v != nil {
		notifier[grafanaContactPointEmailMessage] = v.(string)
	}

	return notifier, nil
}

func (e emailNotifier) getGrafanaContactPointFromSchema(raw []interface{}, name string, disableResolveMessage bool, uid string) grafana_contact_points.GrafanaContactPoint {
	json := raw[0].(map[string]interface{})
	settings := make(map[string]interface{})

	addresses := getAddressFromSchema(json[grafanaContactPointEmailAddresses].([]interface{}))
	settings[grafanaContactPointEmailAddresses] = addresses
	if v, ok := json[grafanaContactPointEmailSingleEmail]; ok && v != nil {
		settings[strcase.LowerCamelCase(grafanaContactPointEmailSingleEmail)] = v.(bool)
	}
	if v, ok := json[grafanaContactPointEmailMessage]; ok && v != nil {
		settings[grafanaContactPointEmailMessage] = v.(string)
	}

	return grafana_contact_points.GrafanaContactPoint{
		Uid:                   uid,
		Name:                  name,
		Type:                  e.meta().typeStr,
		DisableResolveMessage: disableResolveMessage,
		Settings:              settings,
	}
}

type googleChatNotifier struct{}

var _ grafanaContactPointNotifier = (*googleChatNotifier)(nil)

func (g googleChatNotifier) meta() grafanaContactPointNotifierMeta {
	return grafanaContactPointNotifierMeta{
		field:   grafanaContactPointGoogleChat,
		typeStr: grafanaContactPointGoogleChat,
	}
}

func (g googleChatNotifier) schema() *schema.Resource {
	r := &schema.Resource{
		Schema: map[string]*schema.Schema{},
	}
	r.Schema[grafanaContactPointGoogleChatUrl] = &schema.Schema{
		Type:      schema.TypeString,
		Required:  true,
		Sensitive: true,
	}
	r.Schema[grafanaContactPointGoogleChatMessage] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	return r
}

func (g googleChatNotifier) getGrafanaContactPointFromObject(d *schema.ResourceData, contactPoint grafana_contact_points.GrafanaContactPoint) (interface{}, error) {
	notifier := make(map[string]interface{}, 0)
	if v, ok := contactPoint.Settings[grafanaContactPointGoogleChatUrl]; ok && v != nil {
		notifier[grafanaContactPointGoogleChatUrl] = v.(string)
	}
	if v, ok := contactPoint.Settings[grafanaContactPointGoogleChatMessage]; ok && v != nil {
		notifier[grafanaContactPointGoogleChatMessage] = v.(string)
	}
	return notifier, nil
}

func (g googleChatNotifier) getGrafanaContactPointFromSchema(raw []interface{}, name string, disableResolveMessage bool, uid string) grafana_contact_points.GrafanaContactPoint {
	json := raw[0].(map[string]interface{})
	settings := make(map[string]interface{})

	settings[grafanaContactPointGoogleChatUrl] = json[grafanaContactPointGoogleChatUrl].(string)
	if v, ok := json[grafanaContactPointGoogleChatMessage]; ok && v != nil {
		settings[grafanaContactPointGoogleChatMessage] = v.(string)
	}
	return grafana_contact_points.GrafanaContactPoint{
		Uid:                   uid,
		Name:                  name,
		Type:                  g.meta().typeStr,
		DisableResolveMessage: disableResolveMessage,
		Settings:              settings,
	}
}

type opsGenieNotifier struct{}

var _ grafanaContactPointNotifier = (*opsGenieNotifier)(nil)

func (o opsGenieNotifier) meta() grafanaContactPointNotifierMeta {
	return grafanaContactPointNotifierMeta{
		field:        grafanaContactPointOpsgenie,
		typeStr:      grafanaContactPointOpsgenie,
		secureFields: []string{grafanaContactPointOpsgenieApiKey},
	}
}

func (o opsGenieNotifier) schema() *schema.Resource {
	r := &schema.Resource{
		Schema: map[string]*schema.Schema{},
	}
	r.Schema[grafanaContactPointOpsgenieApiUrl] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	r.Schema[grafanaContactPointOpsgenieApiKey] = &schema.Schema{
		Type:      schema.TypeString,
		Required:  true,
		Sensitive: true,
	}
	r.Schema[grafanaContactPointOpsgenieAutoClose] = &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
	}
	r.Schema[grafanaContactPointOpsgenieOverridePriority] = &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
	}
	r.Schema[grafanaContactPointOpsgenieSendTagsAs] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		ValidateFunc: validation.StringInSlice(
			[]string{grafanaContactPointOpsgenieSendTagsTags,
				grafanaContactPointOpsgenieSendTagsDetails,
				grafanaContactPointOpsgenieSendTagsBoth},
			false),
	}

	return r
}

func (o opsGenieNotifier) getGrafanaContactPointFromObject(d *schema.ResourceData, contactPoint grafana_contact_points.GrafanaContactPoint) (interface{}, error) {
	notifier := make(map[string]interface{}, 0)

	if v, ok := contactPoint.Settings[strcase.LowerCamelCase(grafanaContactPointOpsgenieApiUrl)]; ok && v != nil {
		notifier[grafanaContactPointOpsgenieApiUrl] = v.(string)
	}

	if v, ok := contactPoint.Settings[strcase.LowerCamelCase(grafanaContactPointOpsgenieApiKey)]; ok && v != nil {
		notifier[grafanaContactPointOpsgenieApiKey] = v.(string)
	}

	if v, ok := contactPoint.Settings[strcase.LowerCamelCase(grafanaContactPointOpsgenieAutoClose)]; ok && v != nil {
		notifier[grafanaContactPointOpsgenieAutoClose] = v.(bool)
	}
	if v, ok := contactPoint.Settings[strcase.LowerCamelCase(grafanaContactPointOpsgenieOverridePriority)]; ok && v != nil {
		notifier[grafanaContactPointOpsgenieOverridePriority] = v.(bool)
	}
	if v, ok := contactPoint.Settings[strcase.LowerCamelCase(grafanaContactPointOpsgenieSendTagsAs)]; ok && v != nil {
		notifier[grafanaContactPointOpsgenieSendTagsAs] = v.(string)
	}

	getSecuredFieldsFromSchema(notifier, o.meta().secureFields, o.meta().field, d)

	return notifier, nil
}

func (o opsGenieNotifier) getGrafanaContactPointFromSchema(raw []interface{}, name string, disableResolveMessage bool, uid string) grafana_contact_points.GrafanaContactPoint {
	json := raw[0].(map[string]interface{})
	settings := make(map[string]interface{})

	if v, ok := json[grafanaContactPointOpsgenieApiUrl]; ok && v != nil {
		settings[strcase.LowerCamelCase(grafanaContactPointOpsgenieApiUrl)] = v.(string)
	}
	if v, ok := json[grafanaContactPointOpsgenieApiKey]; ok && v != nil {
		settings[strcase.LowerCamelCase(grafanaContactPointOpsgenieApiKey)] = v.(string)
	}
	if v, ok := json[grafanaContactPointOpsgenieAutoClose]; ok && v != nil {
		settings[strcase.LowerCamelCase(grafanaContactPointOpsgenieAutoClose)] = v.(bool)
	}
	if v, ok := json[grafanaContactPointOpsgenieOverridePriority]; ok && v != nil {
		settings[strcase.LowerCamelCase(grafanaContactPointOpsgenieOverridePriority)] = v.(bool)
	}
	if v, ok := json[grafanaContactPointOpsgenieSendTagsAs]; ok && v != nil {
		settings[strcase.LowerCamelCase(grafanaContactPointOpsgenieSendTagsAs)] = v.(string)
	}

	return grafana_contact_points.GrafanaContactPoint{
		Uid:                   uid,
		Name:                  name,
		Type:                  o.meta().typeStr,
		DisableResolveMessage: disableResolveMessage,
		Settings:              settings,
	}
}

type pagerDutyNotifier struct{}

var _ grafanaContactPointNotifier = (*pagerDutyNotifier)(nil)

func (p pagerDutyNotifier) meta() grafanaContactPointNotifierMeta {
	return grafanaContactPointNotifierMeta{
		field:        grafanaContactPointPagerduty,
		typeStr:      grafanaContactPointPagerduty,
		secureFields: []string{grafanaContactPointPagerdutyIntegrationKey},
	}
}

func (p pagerDutyNotifier) schema() *schema.Resource {
	r := &schema.Resource{
		Schema: map[string]*schema.Schema{},
	}
	r.Schema[grafanaContactPointPagerdutyClass] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	r.Schema[grafanaContactPointPagerdutyComponent] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	r.Schema[grafanaContactPointPagerdutyGroup] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	r.Schema[grafanaContactPointPagerdutyIntegrationKey] = &schema.Schema{
		Type:      schema.TypeString,
		Required:  true,
		Sensitive: true,
	}
	r.Schema[grafanaContactPointPagerdutySeverity] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		ValidateFunc: validation.StringInSlice(
			[]string{grafanaContactPointPagerdutySeverityInfo,
				grafanaContactPointPagerdutySeverityWarning,
				grafanaContactPointPagerdutySeverityError,
				grafanaContactPointPagerdutySeverityCritical},
			false),
	}
	r.Schema[grafanaContactPointPagerdutySummary] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}

	return r
}

func (p pagerDutyNotifier) getGrafanaContactPointFromObject(d *schema.ResourceData, contactPoint grafana_contact_points.GrafanaContactPoint) (interface{}, error) {
	notifier := make(map[string]interface{}, 0)

	if v, ok := contactPoint.Settings[strcase.LowerCamelCase(grafanaContactPointPagerdutyIntegrationKey)]; ok && v != nil {
		notifier[grafanaContactPointPagerdutyIntegrationKey] = v.(string)
	}
	if v, ok := contactPoint.Settings[grafanaContactPointPagerdutySeverity]; ok && v != nil {
		notifier[grafanaContactPointPagerdutySeverity] = v.(string)
	}
	if v, ok := contactPoint.Settings[grafanaContactPointPagerdutyClass]; ok && v != nil {
		notifier[grafanaContactPointPagerdutyClass] = v.(string)
	}
	if v, ok := contactPoint.Settings[grafanaContactPointPagerdutyComponent]; ok && v != nil {
		notifier[grafanaContactPointPagerdutyComponent] = v.(string)
	}
	if v, ok := contactPoint.Settings[grafanaContactPointPagerdutyGroup]; ok && v != nil {
		notifier[grafanaContactPointPagerdutyGroup] = v.(string)
	}
	if v, ok := contactPoint.Settings[grafanaContactPointPagerdutySummary]; ok && v != nil {
		notifier[grafanaContactPointPagerdutySummary] = v.(string)
	}

	getSecuredFieldsFromSchema(notifier, p.meta().secureFields, p.meta().field, d)

	return notifier, nil
}

func (p pagerDutyNotifier) getGrafanaContactPointFromSchema(raw []interface{}, name string, disableResolveMessage bool, uid string) grafana_contact_points.GrafanaContactPoint {
	json := raw[0].(map[string]interface{})
	settings := make(map[string]interface{})

	if v, ok := json[grafanaContactPointPagerdutyIntegrationKey]; ok && v != nil {
		settings[strcase.LowerCamelCase(grafanaContactPointPagerdutyIntegrationKey)] = v.(string)
	}
	if v, ok := json[grafanaContactPointPagerdutySeverity]; ok && v != nil {
		settings[grafanaContactPointPagerdutySeverity] = v.(string)
	}
	if v, ok := json[grafanaContactPointPagerdutyClass]; ok && v != nil {
		settings[grafanaContactPointPagerdutyClass] = v.(string)
	}
	if v, ok := json[grafanaContactPointPagerdutyComponent]; ok && v != nil {
		settings[grafanaContactPointPagerdutyComponent] = v.(string)
	}
	if v, ok := json[grafanaContactPointPagerdutyGroup]; ok && v != nil {
		settings[grafanaContactPointPagerdutyGroup] = v.(string)
	}
	if v, ok := json[grafanaContactPointPagerdutySummary]; ok && v != nil {
		settings[grafanaContactPointPagerdutySummary] = v.(string)
	}

	return grafana_contact_points.GrafanaContactPoint{
		Uid:                   uid,
		Name:                  name,
		Type:                  p.meta().typeStr,
		DisableResolveMessage: disableResolveMessage,
		Settings:              settings,
	}
}

type slackNotifier struct{}

var _ grafanaContactPointNotifier = (*slackNotifier)(nil)

func (s slackNotifier) meta() grafanaContactPointNotifierMeta {
	return grafanaContactPointNotifierMeta{
		field:        grafanaContactPointSlack,
		typeStr:      grafanaContactPointSlack,
		secureFields: []string{grafanaContactPointSlackUrl, grafanaContactPointSlackToken},
	}
}

func (s slackNotifier) schema() *schema.Resource {
	r := &schema.Resource{
		Schema: map[string]*schema.Schema{},
	}
	r.Schema[grafanaContactPointSlackEndpointUrl] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	r.Schema[grafanaContactPointSlackUrl] = &schema.Schema{
		Type:      schema.TypeString,
		Optional:  true,
		Sensitive: true,
	}
	r.Schema[grafanaContactPointSlackToken] = &schema.Schema{
		Type:      schema.TypeString,
		Optional:  true,
		Sensitive: true,
	}
	r.Schema[grafanaContactPointSlackRecipient] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	r.Schema[grafanaContactPointSlackText] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	r.Schema[grafanaContactPointSlackTitle] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	r.Schema[grafanaContactPointSlackUsername] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	r.Schema[grafanaContactPointSlackMentionChannel] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		ValidateFunc: validation.StringInSlice(
			[]string{grafanaContactPointSlackMentionChannelHere,
				grafanaContactPointSlackMentionChannelChannel,
				grafanaContactPointSlackMentionChannelDisable},
			false),
	}
	r.Schema[grafanaContactPointSlackMentionUsers] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	r.Schema[grafanaContactPointSlackMentionGroups] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}

	return r
}

func (s slackNotifier) getGrafanaContactPointFromObject(d *schema.ResourceData, contactPoint grafana_contact_points.GrafanaContactPoint) (interface{}, error) {
	notifier := make(map[string]interface{}, 0)

	if v, ok := contactPoint.Settings[strcase.LowerCamelCase(grafanaContactPointSlackEndpointUrl)]; ok && v != nil {
		notifier[grafanaContactPointSlackEndpointUrl] = v.(string)
	}

	if v, ok := contactPoint.Settings[grafanaContactPointSlackRecipient]; ok && v != nil {
		notifier[grafanaContactPointSlackRecipient] = v.(string)
	}

	if v, ok := contactPoint.Settings[grafanaContactPointSlackText]; ok && v != nil {
		notifier[grafanaContactPointSlackText] = v.(string)
	}

	if v, ok := contactPoint.Settings[grafanaContactPointSlackTitle]; ok && v != nil {
		notifier[grafanaContactPointSlackTitle] = v.(string)
	}

	if v, ok := contactPoint.Settings[grafanaContactPointSlackUsername]; ok && v != nil {
		notifier[grafanaContactPointSlackUsername] = v.(string)
	}

	if v, ok := contactPoint.Settings[strcase.LowerCamelCase(grafanaContactPointSlackMentionChannel)]; ok && v != nil {
		notifier[grafanaContactPointSlackMentionChannel] = v.(string)
	}

	if v, ok := contactPoint.Settings[strcase.LowerCamelCase(grafanaContactPointSlackMentionUsers)]; ok && v != nil {
		notifier[grafanaContactPointSlackMentionUsers] = v.(string)
	}

	if v, ok := contactPoint.Settings[strcase.LowerCamelCase(grafanaContactPointSlackMentionGroups)]; ok && v != nil {
		notifier[grafanaContactPointSlackMentionGroups] = v.(string)
	}

	getSecuredFieldsFromSchema(notifier, s.meta().secureFields, s.meta().field, d)

	return notifier, nil
}

func (s slackNotifier) getGrafanaContactPointFromSchema(raw []interface{}, name string, disableResolveMessage bool, uid string) grafana_contact_points.GrafanaContactPoint {
	json := raw[0].(map[string]interface{})
	settings := make(map[string]interface{})

	if v, ok := json[grafanaContactPointSlackEndpointUrl]; ok && v != nil {
		settings[strcase.LowerCamelCase(grafanaContactPointSlackEndpointUrl)] = v.(string)
	}

	if v, ok := json[grafanaContactPointSlackUrl]; ok && v != nil {
		settings[grafanaContactPointSlackUrl] = v.(string)
	}

	if v, ok := json[grafanaContactPointSlackToken]; ok && v != nil {
		settings[grafanaContactPointSlackToken] = v.(string)
	}

	if v, ok := json[grafanaContactPointSlackRecipient]; ok && v != nil {
		settings[grafanaContactPointSlackRecipient] = v.(string)
	}

	if v, ok := json[grafanaContactPointSlackText]; ok && v != nil {
		settings[grafanaContactPointSlackText] = v.(string)
	}

	if v, ok := json[grafanaContactPointSlackTitle]; ok && v != nil {
		settings[grafanaContactPointSlackTitle] = v.(string)
	}

	if v, ok := json[grafanaContactPointSlackUsername]; ok && v != nil {
		settings[grafanaContactPointSlackUsername] = v.(string)
	}

	if v, ok := json[grafanaContactPointSlackMentionChannel]; ok && v != nil {
		settings[strcase.LowerCamelCase(grafanaContactPointSlackMentionChannel)] = v.(string)
	}

	if v, ok := json[grafanaContactPointSlackMentionUsers]; ok && v != nil {
		settings[strcase.LowerCamelCase(grafanaContactPointSlackMentionUsers)] = v.(string)
	}

	if v, ok := json[grafanaContactPointSlackMentionGroups]; ok && v != nil {
		settings[strcase.LowerCamelCase(grafanaContactPointSlackMentionGroups)] = v.(string)
	}

	return grafana_contact_points.GrafanaContactPoint{
		Uid:                   uid,
		Name:                  name,
		Type:                  s.meta().typeStr,
		DisableResolveMessage: disableResolveMessage,
		Settings:              settings,
	}
}

type teamsNotifier struct{}

var _ grafanaContactPointNotifier = (*teamsNotifier)(nil)

func (t teamsNotifier) meta() grafanaContactPointNotifierMeta {
	return grafanaContactPointNotifierMeta{
		field:        grafanaContactPointMicrosoftTeams,
		typeStr:      grafanaContactPointMicrosoftTeams,
		secureFields: []string{grafanaContactPointMicrosoftTeamsUrl},
	}
}

func (t teamsNotifier) schema() *schema.Resource {
	r := &schema.Resource{
		Schema: map[string]*schema.Schema{},
	}
	r.Schema[grafanaContactPointMicrosoftTeamsUrl] = &schema.Schema{
		Type:      schema.TypeString,
		Required:  true,
		Sensitive: true,
	}
	r.Schema[grafanaContactPointMicrosoftTeamsMessage] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	return r
}

func (t teamsNotifier) getGrafanaContactPointFromObject(d *schema.ResourceData, contactPoint grafana_contact_points.GrafanaContactPoint) (interface{}, error) {
	notifier := make(map[string]interface{}, 0)

	if v, ok := contactPoint.Settings[grafanaContactPointMicrosoftTeamsUrl]; ok && v != nil {
		notifier[grafanaContactPointMicrosoftTeamsUrl] = v.(string)
	}

	if v, ok := contactPoint.Settings[grafanaContactPointMicrosoftTeamsMessage]; ok && v != nil {
		notifier[grafanaContactPointMicrosoftTeamsMessage] = v.(string)
	}

	getSecuredFieldsFromSchema(notifier, t.meta().secureFields, t.meta().field, d)

	return notifier, nil
}

func (t teamsNotifier) getGrafanaContactPointFromSchema(raw []interface{}, name string, disableResolveMessage bool, uid string) grafana_contact_points.GrafanaContactPoint {
	json := raw[0].(map[string]interface{})
	settings := make(map[string]interface{})

	if v, ok := json[grafanaContactPointMicrosoftTeamsUrl]; ok && v != nil {
		settings[grafanaContactPointMicrosoftTeamsUrl] = v.(string)
	}

	if v, ok := json[grafanaContactPointMicrosoftTeamsMessage]; ok && v != nil {
		settings[grafanaContactPointMicrosoftTeamsMessage] = v.(string)
	}

	return grafana_contact_points.GrafanaContactPoint{
		Uid:                   uid,
		Name:                  name,
		Type:                  t.meta().typeStr,
		DisableResolveMessage: disableResolveMessage,
		Settings:              settings,
	}
}
