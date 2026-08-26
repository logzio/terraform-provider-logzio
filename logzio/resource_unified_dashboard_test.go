package logzio

import (
	"fmt"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
}

func TestUnifiedDashboardResource_Schema(t *testing.T) {
	r := resourceUnifiedDashboard()
	for _, field := range []string{unifiedDashboardJson, unifiedDashboardFolderId, unifiedDashboardUid, unifiedDashboardName, unifiedDashboardVersion} {
		if _, ok := r.Schema[field]; !ok {
			t.Fatalf("expected schema field %q", field)
		}
	}
	assert.True(t, r.Schema[unifiedDashboardFolderId].ForceNew)
	assert.True(t, r.Schema[unifiedDashboardJson].Required)
	assert.True(t, r.Schema[unifiedDashboardUid].Computed)
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
