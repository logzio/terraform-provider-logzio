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

// unifiedDashboardServerOwnedMetadataFields are the metadata keys the API owns.
// They are stripped from both the configured and the returned document so they
// cannot produce a perpetual diff. Every other metadata key — tags among them —
// is author-owned and round-trips untouched.
//
// Probed live against api.logz.io on 2026-09-02: the gateway stamps only
// "project" into metadata and keeps version/createdAt/updatedAt beside the doc
// rather than inside it, at v1 and v2 alike. The other three are listed because
// upstream Perses carries them in metadata, so a gateway that stops flattening
// them would otherwise reintroduce the perpetual diff.
var unifiedDashboardServerOwnedMetadataFields = []string{
	"project",
	"version",
	"createdAt",
	"updatedAt",
}

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

	if err = validateUnifiedDashboardIdentity(dashboard, folderId, uid); err != nil {
		return diag.FromErr(err)
	}
	if err = setUnifiedDashboard(d, dashboard, folderId, uid); err != nil {
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
		// diagRet only carries diagnostics when the read itself failed; when the
		// retries ran out because the dashboard never converged it is empty, and
		// returning it would report the update as successful.
		if diagRet.HasError() {
			return diagRet
		}
		return diag.FromErr(readErr)
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

func validateUnifiedDashboardIdentity(result *unified_dashboards.Dashboard, folderId, uid string) error {
	if result == nil {
		return fmt.Errorf("unified dashboard response was empty")
	}
	if result.Uid != "" && result.Uid != uid {
		return fmt.Errorf("unified dashboard response uid %q does not match requested dashboard_uid %q", result.Uid, uid)
	}
	if result.ProjectId != "" && result.ProjectId != folderId {
		return fmt.Errorf("unified dashboard response project id %q does not match requested folder_id %q", result.ProjectId, folderId)
	}
	return nil
}

func setUnifiedDashboard(d *schema.ResourceData, result *unified_dashboards.Dashboard, folderId, uid string) error {
	if result == nil {
		return fmt.Errorf("cannot set unified dashboard state from an empty response")
	}
	if err := d.Set(unifiedDashboardFolderId, folderId); err != nil {
		return fmt.Errorf("failed to set %s: %w", unifiedDashboardFolderId, err)
	}
	if err := d.Set(unifiedDashboardUid, uid); err != nil {
		return fmt.Errorf("failed to set %s: %w", unifiedDashboardUid, err)
	}
	if err := d.Set(unifiedDashboardVersion, result.Version); err != nil {
		return fmt.Errorf("failed to set %s: %w", unifiedDashboardVersion, err)
	}
	if err := d.Set(unifiedDashboardName, unifiedDashboardNameFromDoc(result.Doc)); err != nil {
		return fmt.Errorf("failed to set %s: %w", unifiedDashboardName, err)
	}
	if err := d.Set(unifiedDashboardJson, handleUnifiedDashboardConfig(result.Doc)); err != nil {
		return fmt.Errorf("failed to set %s: %w", unifiedDashboardJson, err)
	}
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
		for _, key := range unifiedDashboardServerOwnedMetadataFields {
			delete(metadata, key)
		}
		dashboardJson["metadata"] = metadata
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
