package logzio

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/avast/retry-go"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/logzio/logzio_terraform_client/unified_dashboards"
)

const (
	unifiedDashboardFolderId = "folder_id"
	unifiedDashboardJson     = "dashboard_json"
	unifiedDashboardUid      = "dashboard_uid"
	unifiedDashboardName     = "name"
	unifiedDashboardVersion  = "version"

	unifiedDashboardRetryAttempts = 8
)

func resourceUnifiedDashboard() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUnifiedDashboardCreate,
		ReadContext:   resourceUnifiedDashboardRead,
		UpdateContext: resourceUnifiedDashboardUpdate,
		DeleteContext: resourceUnifiedDashboardDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			unifiedDashboardFolderId: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			unifiedDashboardJson: {
				Type:         schema.TypeString,
				Required:     true,
				StateFunc:    handleUnifiedDashboardConfig,
				ValidateFunc: validateUnifiedDashboardJson,
			},
			unifiedDashboardUid: {
				Type:     schema.TypeString,
				Computed: true,
			},
			unifiedDashboardName: {
				Type:     schema.TypeString,
				Computed: true,
			},
			unifiedDashboardVersion: {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func unifiedDashboardClient(m any) (*unified_dashboards.DashboardsClient, error) {
	client, err := unified_dashboards.New(m.(Config).apiToken, m.(Config).baseUrl)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func resourceUnifiedDashboardCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := unifiedDashboardClient(m)
	if err != nil {
		return diag.FromErr(err)
	}

	folderId := d.Get(unifiedDashboardFolderId).(string)
	doc, err := unifiedDashboardDocFromSchema(d)
	if err != nil {
		return diag.FromErr(err)
	}

	result, err := client.CreateDashboard(folderId, unified_dashboards.CreateDashboardRequest{Doc: doc})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(composeUnifiedDashboardId(folderId, result.Uid))
	return resourceUnifiedDashboardRead(ctx, d, m)
}

func resourceUnifiedDashboardRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := unifiedDashboardClient(m)
	if err != nil {
		return diag.FromErr(err)
	}

	folderId, uid, err := parseUnifiedDashboardId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	dashboard, err := client.GetDashboard(folderId, uid)
	if err != nil {
		tflog.Error(ctx, err.Error())
		if strings.Contains(err.Error(), "missing unified dashboard") {
			d.SetId("")
			return diag.Diagnostics{}
		}
		return diag.FromErr(err)
	}

	if err = setUnifiedDashboard(d, dashboard); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceUnifiedDashboardUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	if unifiedDashboardNameChanged(d.Get(unifiedDashboardName).(string), d.Get(unifiedDashboardJson).(string)) {
		return diag.Errorf("Updating metadata.name is not allowed")
	}

	client, err := unifiedDashboardClient(m)
	if err != nil {
		return diag.FromErr(err)
	}

	folderId, uid, err := parseUnifiedDashboardId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	doc, err := unifiedDashboardDocFromSchema(d)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateDashboard(folderId, uid, unified_dashboards.UpdateDashboardRequest{Doc: doc})
	if err != nil {
		return diag.FromErr(err)
	}

	expectedJson := handleUnifiedDashboardConfig(doc)
	var diagRet diag.Diagnostics
	readErr := retry.Do(func() error {
		diagRet = resourceUnifiedDashboardRead(ctx, d, m)
		if diagRet.HasError() {
			return fmt.Errorf("received error from read unified dashboard")
		}
		if d.Get(unifiedDashboardJson).(string) != expectedJson {
			return fmt.Errorf("unified dashboard has not been updated yet")
		}
		return nil
	},
		retry.RetryIf(func(err error) bool {
			return err != nil
		}),
		retry.DelayType(retry.BackOffDelay),
		retry.Attempts(unifiedDashboardRetryAttempts),
	)

	if readErr != nil {
		tflog.Error(ctx, "could not update schema")
		return diagRet
	}

	return nil
}

func resourceUnifiedDashboardDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := unifiedDashboardClient(m)
	if err != nil {
		return diag.FromErr(err)
	}

	folderId, uid, err := parseUnifiedDashboardId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if err = client.DeleteDashboard(folderId, uid); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func composeUnifiedDashboardId(folderId, uid string) string {
	return folderId + "/" + uid
}

func parseUnifiedDashboardId(id string) (string, string, error) {
	parts := strings.SplitN(id, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("unexpected id format %q, expected folder_id/dashboard_uid", id)
	}
	return parts[0], parts[1], nil
}

func unifiedDashboardDocFromSchema(d *schema.ResourceData) (map[string]any, error) {
	var doc map[string]any
	if err := json.Unmarshal([]byte(d.Get(unifiedDashboardJson).(string)), &doc); err != nil {
		return nil, err
	}
	return doc, nil
}

func setUnifiedDashboard(d *schema.ResourceData, result *unified_dashboards.Dashboard) error {
	folderId, uid, err := parseUnifiedDashboardId(d.Id())
	if err != nil {
		return err
	}

	d.Set(unifiedDashboardFolderId, folderId)
	d.Set(unifiedDashboardUid, uid)
	d.Set(unifiedDashboardVersion, result.Version)
	d.Set(unifiedDashboardName, unifiedDashboardNameFromDoc(result.Doc))
	d.Set(unifiedDashboardJson, handleUnifiedDashboardConfig(result.Doc))
	return nil
}

func handleUnifiedDashboardConfig(config any) string {
	var dashboardJson map[string]any
	switch c := config.(type) {
	case map[string]any:
		dashboardJson = c
	case string:
		if err := json.Unmarshal([]byte(c), &dashboardJson); err != nil {
			return c
		}
	default:
		return ""
	}

	if metadata, ok := dashboardJson["metadata"].(map[string]any); ok {
		cleaned := map[string]any{}
		if name, ok := metadata["name"].(string); ok {
			cleaned["name"] = name
		}
		dashboardJson["metadata"] = cleaned
	}

	newDashboard, _ := json.Marshal(dashboardJson)
	return string(newDashboard)
}

func validateUnifiedDashboardJson(config any, k string) ([]string, []error) {
	var doc map[string]any
	if err := json.Unmarshal([]byte(config.(string)), &doc); err != nil {
		return nil, []error{err}
	}
	if kind, _ := doc["kind"].(string); kind != "Dashboard" {
		return nil, []error{fmt.Errorf("%s must be a Perses Dashboard document with kind %q", k, "Dashboard")}
	}
	if unifiedDashboardNameFromDoc(doc) == "" {
		return nil, []error{fmt.Errorf("%s metadata.name is required", k)}
	}
	if _, ok := doc["spec"].(map[string]any); !ok {
		return nil, []error{fmt.Errorf("%s spec is required", k)}
	}
	return nil, nil
}

func unifiedDashboardNameFromDoc(doc map[string]any) string {
	metadata, ok := doc["metadata"].(map[string]any)
	if !ok {
		return ""
	}
	name, _ := metadata["name"].(string)
	return name
}

func unifiedDashboardNameChanged(currentName, dashboardJson string) bool {
	var doc map[string]any
	if err := json.Unmarshal([]byte(dashboardJson), &doc); err != nil {
		return true
	}
	name := unifiedDashboardNameFromDoc(doc)
	return name != "" && currentName != "" && name != currentName
}
