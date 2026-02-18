package logzio

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/logzio/logzio_terraform_client/unified_alerts"
)

const (
	unifiedAlertDsType = "type"
)

func dataSourceUnifiedAlert() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUnifiedAlertRead,
		Schema: map[string]*schema.Schema{
			unifiedAlertId: {
				Type:     schema.TypeString,
				Required: true,
			},
			unifiedAlertDsType: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{unified_alerts.TypeLogAlert, unified_alerts.TypeMetricAlert}, false),
			},
			unifiedAlertTitle: {
				Type:     schema.TypeString,
				Computed: true,
			},
			unifiedAlertDescription: {
				Type:     schema.TypeString,
				Computed: true,
			},
			unifiedAlertTags: {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			unifiedAlertLinkedPanel: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						linkedPanelFolderId: {
							Type:     schema.TypeString,
							Computed: true,
						},
						linkedPanelDashboardId: {
							Type:     schema.TypeString,
							Computed: true,
						},
						linkedPanelPanelId: {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			unifiedAlertRunbook: {
				Type:     schema.TypeString,
				Computed: true,
			},
			unifiedAlertEnabled: {
				Type:     schema.TypeBool,
				Computed: true,
			},
			unifiedAlertRca: {
				Type:     schema.TypeBool,
				Computed: true,
			},
			unifiedAlertRcaNotificationEndpointIds: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			unifiedAlertUseAlertNotificationEndpointsForRca: {
				Type:     schema.TypeBool,
				Computed: true,
			},
			unifiedAlertRecipients: {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     resourceRecipients(),
			},
			unifiedAlertAlertConfiguration: {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     resourceAlertConfiguration(),
			},
			unifiedAlertCreatedAt: {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			unifiedAlertUpdatedAt: {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			unifiedAlertCreatedBy: {
				Type:     schema.TypeString,
				Computed: true,
			},
			unifiedAlertUpdatedBy: {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceUnifiedAlertRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := unifiedAlertClient(m)
	configType := d.Get(unifiedAlertDsType).(string)
	alertId := d.Get(unifiedAlertId).(string)

	urlType := configTypeToUrlType(configType)
	alert, err := client.GetUnifiedAlert(urlType, alertId)
	if err != nil {
		return diag.Errorf("failed to get unified alert by ID: %v", err)
	}

	d.SetId(fmt.Sprintf("%s:%s", urlType, alert.Id))
	return setUnifiedAlert(d, alert)
}
