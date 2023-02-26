package logzio

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
	"os"
	"testing"
	"time"
)

func TestAccDataSourceGrafanaDashboard(t *testing.T) {
	defer utils.SleepAfterTest()
	folderUid := os.Getenv(grafanaFolderIdEnv)
	resourceFullName := "logzio_grafana_dashboard.test_dashboard"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
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
				PreConfig: func() {
					time.Sleep(time.Second * 3)
				},
				Config: getGrafanaDashboardConfig(folderUid, "create") + getGrafanaDashboardDatasourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.logzio_grafana_dashboard.my_grafana_dashboard_ds", grafanaDashboardUrl),
					resource.TestCheckResourceAttrSet("data.logzio_grafana_dashboard.my_grafana_dashboard_ds", grafanaDashboardFolderUid),
					resource.TestCheckResourceAttrSet("data.logzio_grafana_dashboard.my_grafana_dashboard_ds", grafanaDashboardJson),
				),
			},
		},
	})
}

func getGrafanaDashboardDatasourceConfig() string {
	return `
data "logzio_grafana_dashboard" "my_grafana_dashboard_ds" {
  dashboard_uid = "_terraform_provider_test"
  depends_on = ["logzio_grafana_dashboard.test_dashboard"]
}
`
}
