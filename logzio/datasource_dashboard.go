package logzio

import (
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceDashboard() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDashboardRead,
		Schema: map[string]*schema.Schema{
			dashboardUid: {
				Type:     schema.TypeString,
				Required: true,
			},
			dashboardStarred: {
				Type:     schema.TypeBool,
				Computed: true,
			},
			dashboardUrl: {
				Type:     schema.TypeString,
				Computed: true,
			},
			dashboardFolderId: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			dashboardFolderUid: {
				Type:     schema.TypeString,
				Computed: true,
			},
			dashboardSlug: {
				Type:     schema.TypeString,
				Computed: true,
			},
			dashboardJson: {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDashboardRead(d *schema.ResourceData, m interface{}) error {
	client, err := dashboardClient(m)
	if err != nil {
		return err
	}
	uid := d.Get(dashboardUid).(string)

	result, err := client.Get(uid)
	if err != nil {
		return err
	}
	d.SetId(uid)

	d.Set(dashboardStarred, result.Meta.IsStarred)
	d.Set(dashboardUrl, result.Meta.Url)
	d.Set(dashboardFolderId, result.Meta.FolderId)
	d.Set(dashboardFolderUid, result.Meta.FolderUid)
	d.Set(dashboardSlug, result.Meta.Slug)

	dashboardObject, err := json.Marshal(result.Dashboard)
	if err != nil {
		return err
	}

	d.Set(dashboardJson, string(dashboardObject))

	return nil
}
