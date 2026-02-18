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
	// Datasource-only lookup field (URL type for API call)
	unifiedAlertDsAlertType = "alert_type"
)

func dataSourceUnifiedAlert() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUnifiedAlertRead,
		Schema: map[string]*schema.Schema{
			unifiedAlertId: {
				Type:     schema.TypeString,
				Required: true,
			},
			unifiedAlertDsAlertType: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{unified_alerts.UrlTypeLogs, unified_alerts.UrlTypeMetrics}, false),
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
	alertType := d.Get(unifiedAlertDsAlertType).(string)
	alertId := d.Get(unifiedAlertId).(string)

	alert, err := client.GetUnifiedAlert(alertType, alertId)
	if err != nil {
		return diag.Errorf("failed to get unified alert by ID: %v", err)
	}

	d.SetId(fmt.Sprintf("%s:%s", alertType, alert.Id))
	return setUnifiedAlert(d, alert)
}
