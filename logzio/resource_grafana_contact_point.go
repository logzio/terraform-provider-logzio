package logzio

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

const (
	grafanaContactPointName                  = "name"
	grafanaContactPointUid                   = "uid"
	grafanaContactPointType                  = "type"
	grafanaContactPointSettings              = "settings"
	grafanaContactPointDisableResolveMessage = "disable_resolve_message"
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
			grafanaContactPointType: {
				Type:     schema.TypeString,
				Required: true,
				// TODO : validation of type
			},
			grafanaContactPointSettings: {
				Type:      schema.TypeMap,
				Optional:  true,
				Sensitive: true,
				Default:   map[string]interface{}{},
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			grafanaContactPointDisableResolveMessage: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}
