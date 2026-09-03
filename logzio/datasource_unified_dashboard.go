package logzio

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceUnifiedDashboard() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUnifiedDashboardRead,

		Schema: map[string]*schema.Schema{
			unifiedDashboardFolderId: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique identifier of the unified project containing the dashboard.",
			},
			unifiedDashboardUid: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The stable unique identifier of the unified dashboard.",
			},
			unifiedDashboardJson: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The normalized Perses dashboard document in JSON format.",
			},
			unifiedDashboardName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The dashboard's Perses metadata name.",
			},
			unifiedDashboardVersion: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The dashboard version.",
			},
		},
	}
}

func dataSourceUnifiedDashboardRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := unifiedDashboardClient(m)
	if err != nil {
		return diag.FromErr(err)
	}

	folderId := d.Get(unifiedDashboardFolderId).(string)
	uid := d.Get(unifiedDashboardUid).(string)
	dashboard, err := client.GetDashboard(folderId, uid)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = validateUnifiedDashboardIdentity(dashboard, folderId, uid); err != nil {
		return diag.FromErr(err)
	}
	if err = setUnifiedDashboard(d, dashboard, folderId, uid); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(composeUnifiedDashboardId(folderId, uid))

	return nil
}
