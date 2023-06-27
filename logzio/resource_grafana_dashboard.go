package logzio

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/logzio/logzio_terraform_client/grafana_dashboards"
	"reflect"
	"strings"
)

const (
	grafanaDashboardId        = "dashboard_id"
	grafanaDashboardUid       = "dashboard_uid"
	grafanaDashboardUrl       = "url"
	grafanaDashboardFolderUid = "folder_uid"
	grafanaDashboardJson      = "dashboard_json"
	grafanaDashboardMessage   = "message"
	grafanaDashboardVersion   = "version"
	grafanaDashboardOverwrite = "overwrite"

	grafanaDashboardRetryAttempts = 8
)

var (
	grafanaDashboardsFieldsToDelete = []string{"id", "version"}
)

/**
 * the endpoint resource schema, what terraform uses to parse and read the template
 */
func resourceGrafanaDashboard() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGrafanaDashboardCreate,
		ReadContext:   resourceGrafanaDashboardRead,
		UpdateContext: resourceGrafanaDashboardUpdate,
		DeleteContext: resourceGrafanaDashboardDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			grafanaDashboardJson: {
				Type:         schema.TypeString,
				Required:     true,
				StateFunc:    handleGrafanaDashboardConfig,
				ValidateFunc: validateGrafanaDashboardJson,
			},
			grafanaDashboardFolderUid: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			grafanaDashboardMessage: {
				Type:     schema.TypeString,
				Optional: true,
			},
			grafanaDashboardUrl: {
				Type:     schema.TypeString,
				Computed: true,
			},
			grafanaDashboardId: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			grafanaDashboardUid: {
				Type:     schema.TypeString,
				Computed: true,
			},
			grafanaDashboardVersion: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			grafanaDashboardOverwrite: {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func dashboardClient(m interface{}) (*grafana_dashboards.GrafanaObjectsClient, error) {
	client, err := grafana_dashboards.New(m.(Config).apiToken, m.(Config).baseUrl)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func resourceGrafanaDashboardCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := dashboardClient(m)
	if err != nil {
		return diag.FromErr(err)
	}

	req, err := getCreateUpdateGrafanaDashboardFromSchema(d)
	if err != nil {
		return diag.FromErr(err)
	}

	result, err := client.CreateUpdateGrafanaDashboard(req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(result.Uid)

	return resourceGrafanaDashboardRead(ctx, d, m)
}

func resourceGrafanaDashboardRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := dashboardClient(m)
	if err != nil {
		return diag.FromErr(err)
	}

	grafanaDashboard, err := client.GetGrafanaDashboard(d.Id())
	if err != nil {
		if err != nil {
			tflog.Error(ctx, err.Error())
			if strings.Contains(err.Error(), "missing grafana dashboard") {
				// If we were not able to find the resource - delete from state
				d.SetId("")
				return diag.Diagnostics{}
			} else {
				return diag.FromErr(err)
			}
		}
	}

	err = setGrafanaDashboard(d, grafanaDashboard)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceGrafanaDashboardUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if uidChanged(d) {
		return diag.Errorf("Updating uid is not allowed")
	}

	client, err := dashboardClient(m)
	if err != nil {
		return diag.FromErr(err)
	}

	req, err := getCreateUpdateGrafanaDashboardFromSchema(d)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.CreateUpdateGrafanaDashboard(req)
	if err != nil {
		return diag.FromErr(err)
	}

	var diagRet diag.Diagnostics
	readErr := retry.Do(func() error {
		diagRet = resourceGrafanaDashboardRead(ctx, d, m)
		if diagRet.HasError() {
			return fmt.Errorf("received error from read grafana dashboard")
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
					grafanaDashboardFromSchema, _ := getCreateUpdateGrafanaDashboardFromSchema(d)
					return !reflect.DeepEqual(grafanaDashboardFromSchema, req)
				}
			}),
		retry.DelayType(retry.BackOffDelay),
		retry.Attempts(grafanaDashboardRetryAttempts),
	)

	if readErr != nil {
		tflog.Error(ctx, "could not update schema")
		return diagRet
	}

	return nil
}

func resourceGrafanaDashboardDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := dashboardClient(m)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.DeleteGrafanaDashboard(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func getCreateUpdateGrafanaDashboardFromSchema(d *schema.ResourceData) (grafana_dashboards.CreateUpdatePayload, error) {
	var dashboardObject map[string]interface{}
	var payload grafana_dashboards.CreateUpdatePayload

	err := json.Unmarshal([]byte(d.Get(grafanaDashboardJson).(string)), &dashboardObject)
	if err != nil {
		return payload, err
	}

	for _, key := range grafanaDashboardsFieldsToDelete {
		delete(dashboardObject, key)
	}

	payload = grafana_dashboards.CreateUpdatePayload{
		FolderUid: d.Get(grafanaDashboardFolderUid).(string),
		Message:   d.Get(grafanaDashboardMessage).(string),
		Overwrite: d.Get(grafanaDashboardOverwrite).(bool),
		Dashboard: dashboardObject,
	}

	return payload, nil
}

func setGrafanaDashboard(d *schema.ResourceData, result *grafana_dashboards.GetResults) error {
	d.Set(grafanaDashboardUrl, result.Meta["url"].(string))
	d.Set(grafanaDashboardVersion, int(result.Dashboard["version"].(float64)))
	d.Set(grafanaDashboardId, int(result.Dashboard["id"].(float64)))
	d.Set(grafanaDashboardUid, result.Dashboard["uid"].(string))
	d.Set(grafanaDashboardFolderUid, result.Meta["folderUid"].(string))
	dashboardConfigStr := handleGrafanaDashboardConfig(result.Dashboard)
	d.Set(grafanaDashboardJson, dashboardConfigStr)

	return nil
}

func handleGrafanaDashboardConfig(config interface{}) string {
	var dashboardJson map[string]interface{}
	switch c := config.(type) {
	case map[string]interface{}:
		dashboardJson = c
	case string:
		err := json.Unmarshal([]byte(c), &dashboardJson)
		if err != nil {
			return c
		}
	}

	for _, key := range grafanaDashboardsFieldsToDelete {
		delete(dashboardJson, key)
	}

	if panels, ok := dashboardJson["panels"]; ok {
		for _, panel := range panels.([]interface{}) {
			panelObj := panel.(map[string]interface{})
			delete(panelObj, "id")
			if libraryPanel, exists := panelObj["libraryPanel"].(map[string]interface{}); exists {
				for key := range libraryPanel {
					if key != "uid" && key != "name" {
						delete(libraryPanel, key)
					}
				}
			}
		}
	}

	newDashboard, _ := json.Marshal(dashboardJson)
	return string(newDashboard)
}

func validateGrafanaDashboardJson(config interface{}, k string) ([]string, []error) {
	var configMap map[string]interface{}
	err := json.Unmarshal([]byte(config.(string)), &configMap)
	if err != nil {
		return nil, []error{err}
	}

	return nil, nil
}

func uidChanged(d *schema.ResourceData) bool {
	currUid := d.Id()
	dashboardConfigStr := d.Get(grafanaDashboardJson).(string)
	var dashboardMap map[string]interface{}
	_ = json.Unmarshal([]byte(dashboardConfigStr), &dashboardMap)
	if uid, exists := dashboardMap["uid"]; exists {
		return currUid != uid.(string)
	}

	return true
}
