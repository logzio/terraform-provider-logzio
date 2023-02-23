package logzio

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGrafanaDashboard() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDashboardRead,
		Schema: map[string]*schema.Schema{
			grafanaDashboardUid: {
				Type:     schema.TypeString,
				Required: true,
			},
			grafanaDashboardUrl: {
				Type:     schema.TypeString,
				Computed: true,
			},
			grafanaDashboardFolderUid: {
				Type:     schema.TypeString,
				Computed: true,
			},
			grafanaDashboardJson: {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDashboardRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := dashboardClient(m)
	if err != nil {
		return diag.FromErr(err)
	}

	uid := d.Get(grafanaDashboardUid).(string)

	result, err := client.GetGrafanaDashboard(uid)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(uid)

	d.Set(grafanaDashboardUrl, result.Meta["url"])
	d.Set(grafanaDashboardFolderUid, result.Meta["folderUid"])

	dashboardObject, err := json.Marshal(result.Dashboard)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set(grafanaDashboardJson, string(dashboardObject))

	return nil
}
