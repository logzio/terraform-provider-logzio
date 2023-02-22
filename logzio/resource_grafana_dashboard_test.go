package logzio

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
	"os"
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
				Config: getGrafanaDashboardConfig(folderUid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, "dashboard_id"),
					resource.TestCheckResourceAttr(resourceFullName, "dashboard_uid", "_terraform_provider_test"),
				),
			},
		},
	})
}

func getGrafanaDashboardConfig(folderUid string) string {
	return fmt.Sprintf(`
resource "logzio_grafana_dashboard" "test_dashboard" {
  dashboard_json = file("./testdata/fixtures/grafana_dashboard/create.json")
  folder_uid = "%s"
}
`, folderUid)
}
