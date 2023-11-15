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
