package logzio

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
	"os"
	"regexp"
	"testing"
)

const (
	folderIdEnv = "GRAFANA_FOLDER_UID"
)

func TestAccLogzioGrafanaDashboard_CreateUpdateDashboard(t *testing.T) {
	defer utils.SleepAfterTest()
	folderUid := os.Getenv(folderIdEnv)
	resourceFullName := "logzio_grafana_dashboard.test_dashboard"
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				// Create
				Config: getGrafanaDashboardConfig(folderUid, "create"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "dashboard_id"),
					resource.TestCheckResourceAttr(resourceFullName, "dashboard_uid", "_terraform_provider_test"),
					resource.TestCheckResourceAttr(resourceFullName, "folder_uid", folderUid),
					resource.TestCheckResourceAttr(resourceFullName, "dashboard_json", expectedCreate()),
					resource.TestCheckResourceAttrSet(resourceFullName, "url"),
					resource.TestCheckResourceAttrSet(resourceFullName, "version"),
				),
			},
			{
				// Update
				Config: getGrafanaDashboardConfig(folderUid, "update"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "dashboard_id"),
					resource.TestCheckResourceAttr(resourceFullName, "dashboard_uid", "_terraform_provider_test"),
					resource.TestCheckResourceAttr(resourceFullName, "folder_uid", folderUid),
					resource.TestCheckResourceAttr(resourceFullName, "dashboard_json", expectedUpdate()),
					resource.TestCheckResourceAttrSet(resourceFullName, "url"),
					resource.TestCheckResourceAttrSet(resourceFullName, "version"),
				),
			},
			{
				ResourceName:            resourceFullName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"overwrite"},
			},
		},
	})
}

func TestAccLogzioGrafanaDashboard_CreateUpdateDashboardChangeUid(t *testing.T) {
	defer utils.SleepAfterTest()
	folderUid := os.Getenv(folderIdEnv)
	resourceFullName := "logzio_grafana_dashboard.test_dashboard"
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				// Create
				Config: getGrafanaDashboardConfig(folderUid, "create"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "dashboard_id"),
					resource.TestCheckResourceAttr(resourceFullName, "dashboard_uid", "_terraform_provider_test"),
					resource.TestCheckResourceAttr(resourceFullName, "folder_uid", folderUid),
					resource.TestCheckResourceAttr(resourceFullName, "dashboard_json", expectedCreate()),
					resource.TestCheckResourceAttrSet(resourceFullName, "url"),
					resource.TestCheckResourceAttrSet(resourceFullName, "version"),
				),
			},
			{
				// Update uid
				Config:      getGrafanaDashboardConfig(folderUid, "update_uid"),
				ExpectError: regexp.MustCompile("Updating uid is not allowed"),
			},
		},
	})
}

func getGrafanaDashboardConfig(folderUid, operation string) string {
	return fmt.Sprintf(`
resource "logzio_grafana_dashboard" "test_dashboard" {
  dashboard_json = file("./testdata/fixtures/grafana_dashboard/%s.json")
  folder_uid = "%s"
  overwrite = true
}
`, operation, folderUid)
}

func expectedCreate() string {
	return "{\"panels\":[],\"title\":\"_terraform_provider_test\",\"uid\":\"_terraform_provider_test\"}"
}

func expectedUpdate() string {
	return "{\"message\":\"this is an update\",\"panels\":[],\"tags\":[\"some\",\"tags\",\"blah\"],\"title\":\"terraform test update\",\"uid\":\"_terraform_provider_test\"}"
}

func expectedUpdateUid() string {
	return "{\"message\":\"this is an update\",\"panels\":[],\"tags\":[\"some\",\"tags\",\"blah\"],\"title\":\"terraform test update\"}"
}
