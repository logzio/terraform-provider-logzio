package logzio

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/logzio/logzio_terraform_client/unified_dashboards"
	"github.com/logzio/logzio_terraform_client/unified_projects"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComposeAndParseUnifiedDashboardId(t *testing.T) {
	id := composeUnifiedDashboardId("folder-1", "dash-1")
	assert.Equal(t, "folder-1/dash-1", id)

	folderId, uid, err := parseUnifiedDashboardId(id)
	require.NoError(t, err)
	assert.Equal(t, "folder-1", folderId)
	assert.Equal(t, "dash-1", uid)
}

func TestParseUnifiedDashboardId_Invalid(t *testing.T) {
	cases := []string{"", "only-folder", "/no-folder", "folder/", "folder-only-no-slash"}
	for _, id := range cases {
		_, _, err := parseUnifiedDashboardId(id)
		assert.Error(t, err, "id %q should be rejected", id)
	}
}

func TestValidateUnifiedDashboardJson(t *testing.T) {
	_, errs := validateUnifiedDashboardJson(`{"kind":"Dashboard","metadata":{"name":"cpu"},"spec":{"duration":"1h"}}`, unifiedDashboardJson)
	assert.Empty(t, errs)

	_, errs = validateUnifiedDashboardJson(`not-json`, unifiedDashboardJson)
	assert.NotEmpty(t, errs)

	_, errs = validateUnifiedDashboardJson(`{"metadata":{"name":"cpu"},"spec":{}}`, unifiedDashboardJson)
	assert.NotEmpty(t, errs)

	_, errs = validateUnifiedDashboardJson(`{"kind":"Dashboard","metadata":{},"spec":{}}`, unifiedDashboardJson)
	assert.NotEmpty(t, errs)

	_, errs = validateUnifiedDashboardJson(`{"kind":"Dashboard","metadata":{"name":"cpu"}}`, unifiedDashboardJson)
	assert.NotEmpty(t, errs)
}

func TestHandleUnifiedDashboardConfig_StripsServerMetadata(t *testing.T) {
	input := `{"kind":"Dashboard","metadata":{"name":"cpu","project":"proj-1","createdAt":"2024-01-01"},"spec":{"duration":"1h","layouts":[],"panels":{}}}`
	got := handleUnifiedDashboardConfig(input)
	assert.Equal(t, `{"kind":"Dashboard","metadata":{"name":"cpu"},"spec":{"duration":"1h","layouts":[],"panels":{}}}`, got)
}

func TestHandleUnifiedDashboardConfig_InputForms(t *testing.T) {
	doc := map[string]any{
		"kind":     "Dashboard",
		"metadata": map[string]any{"name": "cpu", "project": "project-1"},
		"spec":     map[string]any{},
	}
	assert.JSONEq(t, `{"kind":"Dashboard","metadata":{"name":"cpu"},"spec":{}}`, handleUnifiedDashboardConfig(doc))
	assert.Equal(t, "not-json", handleUnifiedDashboardConfig("not-json"))
	assert.Empty(t, handleUnifiedDashboardConfig(42))
}

func TestUnifiedDashboardNameFromDoc(t *testing.T) {
	assert.Equal(t, "cpu", unifiedDashboardNameFromDoc(map[string]any{
		"metadata": map[string]any{"name": "cpu"},
	}))
	assert.Equal(t, "", unifiedDashboardNameFromDoc(map[string]any{}))
}

func TestUnifiedDashboardNameChanged(t *testing.T) {
	same := `{"kind":"Dashboard","metadata":{"name":"cpu"},"spec":{}}`
	assert.False(t, unifiedDashboardNameChanged("cpu", same))
	assert.True(t, unifiedDashboardNameChanged("cpu", `{"kind":"Dashboard","metadata":{"name":"other"},"spec":{}}`))
	assert.True(t, unifiedDashboardNameChanged("cpu", `not-json`))
	assert.False(t, unifiedDashboardNameChanged("", same))
}

func TestUnifiedDashboardResource_Schema(t *testing.T) {
	r := resourceUnifiedDashboard()
	require.NoError(t, r.InternalValidate(nil, true))
	for _, field := range []string{unifiedDashboardJson, unifiedDashboardFolderId, unifiedDashboardUid, unifiedDashboardName, unifiedDashboardVersion} {
		if _, ok := r.Schema[field]; !ok {
			t.Fatalf("expected schema field %q", field)
		}
	}
	assert.True(t, r.Schema[unifiedDashboardFolderId].ForceNew)
	assert.True(t, r.Schema[unifiedDashboardJson].Required)
	assert.True(t, r.Schema[unifiedDashboardUid].Computed)
}

func TestUnifiedDashboardExplicitIdentityStateMapping(t *testing.T) {
	r := resourceUnifiedDashboard()
	d := schema.TestResourceDataRaw(t, r.Schema, map[string]any{
		unifiedDashboardFolderId: "project-1",
		unifiedDashboardJson:     `{"kind":"Dashboard","metadata":{"name":"cpu"},"spec":{"duration":"1h"}}`,
	})

	doc, err := unifiedDashboardDocFromSchema(d)
	require.NoError(t, err)
	assert.Equal(t, "cpu", unifiedDashboardNameFromDoc(doc))

	d.SetId("state-folder/state-dashboard")
	err = setUnifiedDashboard(d, &unified_dashboards.Dashboard{
		Uid:     "dashboard-1",
		Version: 3,
		Doc: map[string]any{
			"kind":     "Dashboard",
			"metadata": map[string]any{"name": "cpu", "project": "project-1", "updatedAt": "server-value"},
			"spec":     map[string]any{"duration": "6h"},
		},
	}, "project-1", "dashboard-1")
	require.NoError(t, err)
	assert.Equal(t, "state-folder/state-dashboard", d.Id())
	assert.Equal(t, "project-1", d.Get(unifiedDashboardFolderId))
	assert.Equal(t, "dashboard-1", d.Get(unifiedDashboardUid))
	assert.Equal(t, "cpu", d.Get(unifiedDashboardName))
	assert.Equal(t, 3, d.Get(unifiedDashboardVersion))
	assert.JSONEq(t, `{"kind":"Dashboard","metadata":{"name":"cpu"},"spec":{"duration":"6h"}}`, d.Get(unifiedDashboardJson).(string))
}

func TestSetUnifiedDashboardPropagatesStateErrors(t *testing.T) {
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		unifiedDashboardFolderId: {
			Type:     schema.TypeInt,
			Optional: true,
		},
	}, map[string]any{})

	err := setUnifiedDashboard(d, &unified_dashboards.Dashboard{}, "project-1", "dashboard-1")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to set folder_id")
}

func TestUnifiedDashboardResourceLifecycle(t *testing.T) {
	var methods []string
	doc := map[string]any{
		"kind":     "Dashboard",
		"metadata": map[string]any{"name": "cpu"},
		"spec":     map[string]any{"duration": "1h"},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		methods = append(methods, r.Method)
		if r.Method == http.MethodPost {
			assert.Equal(t, "/perses-public/api/v1/projects/project-1/dashboards", r.URL.Path)
		} else {
			assert.Equal(t, "/perses-public/api/v1/projects/project-1/dashboards/dashboard-1", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		case http.MethodPost, http.MethodPut:
			var request struct {
				Doc map[string]any `json:"doc"`
			}
			assert.NoError(t, json.NewDecoder(r.Body).Decode(&request))
			doc = request.Doc
			fmt.Fprintf(w, `{"uid":"dashboard-1","version":2,"doc":%s}`, mustMarshalJSON(t, doc))
		case http.MethodGet:
			fmt.Fprintf(w, `{"uid":"dashboard-1","version":2,"doc":%s}`, mustMarshalJSON(t, doc))
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
	defer server.Close()

	r := resourceUnifiedDashboard()
	d := schema.TestResourceDataRaw(t, r.Schema, map[string]any{
		unifiedDashboardFolderId: "project-1",
		unifiedDashboardJson:     `{"kind":"Dashboard","metadata":{"name":"cpu"},"spec":{"duration":"1h"}}`,
	})
	config := Config{apiToken: "token", baseUrl: server.URL}

	diags := resourceUnifiedDashboardCreate(context.Background(), d, config)
	require.False(t, diags.HasError(), diagnosticsString(diags))
	assert.Equal(t, "project-1/dashboard-1", d.Id())
	assert.Equal(t, "cpu", d.Get(unifiedDashboardName))

	require.NoError(t, d.Set(unifiedDashboardJson, `{"kind":"Dashboard","metadata":{"name":"cpu"},"spec":{"duration":"6h"}}`))
	diags = resourceUnifiedDashboardUpdate(context.Background(), d, config)
	require.False(t, diags.HasError(), diagnosticsString(diags))
	assert.JSONEq(t, `{"kind":"Dashboard","metadata":{"name":"cpu"},"spec":{"duration":"6h"}}`, d.Get(unifiedDashboardJson).(string))

	diags = resourceUnifiedDashboardDelete(context.Background(), d, config)
	require.False(t, diags.HasError(), diagnosticsString(diags))
	assert.Empty(t, d.Id())
	assert.Equal(t, []string{http.MethodPost, http.MethodGet, http.MethodPut, http.MethodGet, http.MethodDelete}, methods)
}

func TestUnifiedDashboardReadClearsMissingResource(t *testing.T) {
	server := httptest.NewServer(http.NotFoundHandler())
	defer server.Close()

	d := schema.TestResourceDataRaw(t, resourceUnifiedDashboard().Schema, map[string]any{})
	d.SetId("project-1/missing")
	diags := resourceUnifiedDashboardRead(context.Background(), d, Config{apiToken: "token", baseUrl: server.URL})

	assert.False(t, diags.HasError(), diagnosticsString(diags))
	assert.Empty(t, d.Id())
}

func TestAccLogzioUnifiedDashboard_CreateUpdateDashboard(t *testing.T) {
	defer utils.SleepAfterTest()
	var folderId string
	defer deleteTestUnifiedProject(t, &folderId)

	resourceFullName := "logzio_unified_dashboard.test_dashboard"
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck: func() {
			testAccPreCheckApiToken(t)
			folderId = createTestUnifiedProject(t)
		},
		Steps: []resource.TestStep{
			{
				Config: getUnifiedDashboardConfig(folderId, "create"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, unifiedDashboardFolderId, folderId),
					resource.TestCheckResourceAttr(resourceFullName, unifiedDashboardName, "_terraform_provider_unidash_test"),
					resource.TestCheckResourceAttr(resourceFullName, unifiedDashboardJson, expectedUnifiedDashboardCreate()),
					resource.TestCheckResourceAttrSet(resourceFullName, unifiedDashboardUid),
					resource.TestCheckResourceAttrSet(resourceFullName, unifiedDashboardVersion),
				),
			},
			{
				Config: getUnifiedDashboardConfig(folderId, "update"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, unifiedDashboardFolderId, folderId),
					resource.TestCheckResourceAttr(resourceFullName, unifiedDashboardName, "_terraform_provider_unidash_test"),
					resource.TestCheckResourceAttr(resourceFullName, unifiedDashboardJson, expectedUnifiedDashboardUpdate()),
					resource.TestCheckResourceAttrSet(resourceFullName, unifiedDashboardUid),
					resource.TestCheckResourceAttrSet(resourceFullName, unifiedDashboardVersion),
				),
			},
			{
				ResourceName:      resourceFullName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLogzioUnifiedDashboard_CreateUpdateDashboardChangeName(t *testing.T) {
	defer utils.SleepAfterTest()
	var folderId string
	defer deleteTestUnifiedProject(t, &folderId)

	resourceFullName := "logzio_unified_dashboard.test_dashboard"
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck: func() {
			testAccPreCheckApiToken(t)
			folderId = createTestUnifiedProject(t)
		},
		Steps: []resource.TestStep{
			{
				Config: getUnifiedDashboardConfig(folderId, "create"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, unifiedDashboardName, "_terraform_provider_unidash_test"),
				),
			},
			{
				Config:      getUnifiedDashboardConfig(folderId, "update_name"),
				ExpectError: regexp.MustCompile("Updating metadata.name is not allowed"),
			},
		},
	})
}

func getUnifiedDashboardConfig(folderId, operation string) string {
	return fmt.Sprintf(`
resource "logzio_unified_dashboard" "test_dashboard" {
  dashboard_json = file("./testdata/fixtures/unified_dashboard/%s.json")
  folder_id      = "%s"
}
`, operation, folderId)
}

func expectedUnifiedDashboardCreate() string {
	return `{"kind":"Dashboard","metadata":{"name":"_terraform_provider_unidash_test"},"spec":{"display":{"name":"terraform provider unidash test"},"duration":"1h","layouts":[],"panels":{}}}`
}

func expectedUnifiedDashboardUpdate() string {
	return `{"kind":"Dashboard","metadata":{"name":"_terraform_provider_unidash_test"},"spec":{"display":{"name":"terraform provider unidash test update"},"duration":"1h","layouts":[],"panels":{}}}`
}

func createTestUnifiedProject(t *testing.T) string {
	t.Helper()
	client, err := unified_projects.New(os.Getenv(envLogzioApiToken), testAccUnifiedDashboardBaseUrl())
	if err != nil {
		t.Fatalf("failed to create unified projects client: %v", err)
	}

	name := "tf-provider-unidash-" + time.Now().Format("20060102150405")
	proj, err := client.CreateProject(unified_projects.CreateProjectRequest{Name: name})
	if err != nil {
		t.Fatalf("failed to create unified project: %v", err)
	}
	time.Sleep(2 * time.Second)
	return proj.Id
}

func deleteTestUnifiedProject(t *testing.T, folderId *string) {
	t.Helper()
	if folderId == nil || *folderId == "" {
		return
	}
	client, err := unified_projects.New(os.Getenv(envLogzioApiToken), testAccUnifiedDashboardBaseUrl())
	if err != nil {
		t.Logf("cleanup: failed to create unified projects client: %v", err)
		return
	}
	if err = client.DeleteProject(*folderId); err != nil {
		t.Logf("cleanup: failed to delete unified project %s: %v", *folderId, err)
	}
}

func testAccUnifiedDashboardBaseUrl() string {
	if custom := os.Getenv(envLogzioCustomApiUrl); custom != "" {
		return custom
	}
	region := os.Getenv(envLogzioRegion)
	regionCode := ""
	if region != "" && region != "us" {
		regionCode = fmt.Sprintf("-%s", region)
	}
	return fmt.Sprintf(baseUrl, regionCode)
}
