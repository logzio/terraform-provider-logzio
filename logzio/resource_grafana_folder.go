package logzio

import (
	"context"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/logzio/logzio_terraform_client/grafana_folders"
	"reflect"
	"strings"
)

const (
	grafanaFolderUid     = "uid"
	grafanaFolderTitle   = "title"
	grafanaFolderId      = "folder_id"
	grafanaFolderUrl     = "url"
	grafanaFolderVersion = "version"

	grafanaFolderRetryAttempts = 8
)

func resourceGrafanaFolder() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGrafanaFolderCreate,
		ReadContext:   resourceGrafanaFolderRead,
		UpdateContext: resourceGrafanaFolderUpdate,
		DeleteContext: resourceGrafanaFolderDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
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

func grafanaFolderClient(m interface{}) (*grafana_folders.GrafanaFolderClient, error) {
	client, err := grafana_folders.New(m.(Config).apiToken, m.(Config).baseUrl)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func resourceGrafanaFolderCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := grafanaFolderClient(m)
	if err != nil {
		return diag.FromErr(err)
	}

	req := getCreateGrafanaFolderFromSchema(d)
	result, err := client.CreateGrafanaFolder(req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(result.Uid)
	return resourceGrafanaFolderRead(ctx, d, m)
}

func resourceGrafanaFolderRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := grafanaFolderClient(m)
	if err != nil {
		return diag.FromErr(err)
	}

	folder, err := client.GetGrafanaFolder(d.Id())
	if err != nil {
		tflog.Error(ctx, err.Error())
		if strings.Contains(err.Error(), "missing grafana folder") {
			// If we were not able to find the resource - delete from state
			d.SetId("")
			return diag.Diagnostics{}
		} else {
			return diag.FromErr(err)
		}
	}

	setGrafanaFolder(d, folder)
	return nil
}

func resourceGrafanaFolderUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := grafanaFolderClient(m)
	if err != nil {
		return diag.FromErr(err)
	}

	updateFolder := getUpdateGrafanaFolderFromSchema(d)
	err = client.UpdateGrafanaFolder(updateFolder)
	if err != nil {
		return diag.FromErr(err)
	}

	var diagRet diag.Diagnostics
	readErr := retry.Do(func() error {
		diagRet = resourceGrafanaFolderRead(ctx, d, m)
		if diagRet.HasError() {
			return fmt.Errorf("received error from read grafana folder")
		}

		return nil
	},
		retry.RetryIf(
			// Retry ONLY if the resource was not updated yet
			func(err error) bool {
				if err != nil {
					return false
				} else {
					// Check if the update shows on read
					// if not updated yet - retry
					grafanaFolderFromSchema := getUpdateGrafanaFolderFromSchema(d)
					return !reflect.DeepEqual(grafanaFolderFromSchema, updateFolder)
				}
			}),
		retry.DelayType(retry.BackOffDelay),
		retry.Attempts(grafanaFolderRetryAttempts),
	)

	if readErr != nil {
		tflog.Error(ctx, "could not update schema")
		return diagRet
	}

	return nil
}

func resourceGrafanaFolderDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := grafanaFolderClient(m)
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteGrafanaFolder(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func setGrafanaFolder(d *schema.ResourceData, folder *grafana_folders.GrafanaFolder) {
	d.Set(grafanaFolderUid, folder.Uid)
	d.Set(grafanaFolderId, folder.Id)
	d.Set(grafanaFolderTitle, folder.Title)
	d.Set(grafanaFolderUrl, folder.Url)
	d.Set(grafanaFolderVersion, folder.Version)
}

func getCreateGrafanaFolderFromSchema(d *schema.ResourceData) grafana_folders.CreateUpdateFolder {
	return grafana_folders.CreateUpdateFolder{
		Uid:   d.Get(grafanaFolderUid).(string),
		Title: d.Get(grafanaFolderTitle).(string),
	}
}

func getUpdateGrafanaFolderFromSchema(d *schema.ResourceData) grafana_folders.CreateUpdateFolder {
	folder := getCreateGrafanaFolderFromSchema(d)
	folder.Overwrite = true
	return folder
}
