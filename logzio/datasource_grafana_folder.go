package logzio

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGrafanaFolder() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGrafanaFolderRead,
		Schema: map[string]*schema.Schema{
			grafanaFolderUid: {
				Type:     schema.TypeString,
				Computed: true,
			},
			grafanaFolderTitle: {
				Type:     schema.TypeString,
				Required: true,
			},
			grafanaFolderId: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			grafanaFolderUrl: {
				Type:     schema.TypeString,
				Computed: true,
			},
			grafanaFolderVersion: {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceGrafanaFolderRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := grafanaFolderClient(m)
	if err != nil {
		return diag.FromErr(err)
	}

	title := d.Get(grafanaFolderTitle).(string)
	folders, err := client.ListGrafanaFolders()
	if err != nil {
		return diag.FromErr(err)
	}

	for _, folder := range folders {
		if folder.Title == title {
			getFolder, err := client.GetGrafanaFolder(folder.Uid)
			if err != nil {
				return diag.FromErr(err)
			}

			d.SetId(getFolder.Uid)
			setGrafanaFolder(d, getFolder)
			return nil
		}
	}

	return diag.Errorf("Could not find Grafana folder with title %s", title)
}
