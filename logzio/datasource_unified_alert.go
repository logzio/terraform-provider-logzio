package logzio

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceUnifiedAlert() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUnifiedAlertRead,
		Schema: map[string]*schema.Schema{
			unifiedAlertId: {
				Type:     schema.TypeString,
				Optional: true,
			},
			unifiedAlertTitle: {
				Type:     schema.TypeString,
				Optional: true,
			},
			unifiedAlertType: {
				Type:     schema.TypeString,
				Required: true,
			},
			unifiedAlertDescription: {
				Type:     schema.TypeString,
				Computed: true,
			},
			unifiedAlertTags: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			unifiedAlertFolderId: {
				Type:     schema.TypeString,
				Computed: true,
			},
			unifiedAlertDashboardId: {
				Type:     schema.TypeString,
				Computed: true,
			},
			unifiedAlertPanelId: {
				Type:     schema.TypeString,
				Computed: true,
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
			unifiedAlertCreatedAt: {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			unifiedAlertUpdatedAt: {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			unifiedAlertLogAlert: {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     resourceLogAlertConfig(),
			},
			unifiedAlertMetricAlert: {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     resourceMetricAlertConfig(),
			},
		},
	}
}

func dataSourceUnifiedAlertRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := unifiedAlertClient(m)
	alertType := d.Get(unifiedAlertType).(string)
	urlType := getUrlTypeFromAlertType(alertType)

	alertId, alertIdOk := d.GetOk(unifiedAlertId)
	_, alertTitleOk := d.GetOk(unifiedAlertTitle)

	if !alertIdOk && !alertTitleOk {
		return diag.Errorf("either alert_id or title must be specified")
	}

	if alertIdOk {
		// Lookup by ID
		alert, err := client.GetUnifiedAlert(urlType, alertId.(string))
		if err != nil {
			return diag.Errorf("failed to get unified alert by ID: %v", err)
		}

		d.SetId(alert.Id)
		return setUnifiedAlert(d, alert)
	}

	// Lookup by title not directly supported by API
	// This would require listing all alerts and filtering, which is not ideal
	// For now, return error directing users to use ID
	return diag.Errorf("lookup by title is not supported for unified alerts, please use alert_id instead")
}
