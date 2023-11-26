package logzio

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
	"testing"
	"time"
)

func TestAccDataSourceGrafanaFolder(t *testing.T) {
	defer utils.SleepAfterTest()
	title := "my_title"
	resourceType := "logzio_grafana_folder"
	resourceName := "test_folder"
	dsName := "my_ds_folder"
	fullResourceName := fmt.Sprintf("%s.%s", resourceType, resourceName)
	fullDs := fmt.Sprintf("%s.%s.%s", "data", resourceType, dsName)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getGrafanaFolderConfig(title),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fullResourceName, grafanaFolderTitle, title),
					resource.TestCheckResourceAttrSet(fullResourceName, grafanaFolderUid),
					resource.TestCheckResourceAttrSet(fullResourceName, grafanaFolderId),
					resource.TestCheckResourceAttrSet(fullResourceName, grafanaFolderVersion),
					resource.TestCheckResourceAttrSet(fullResourceName, grafanaFolderUrl),
				),
			},
			{
				PreConfig: func() {
					time.Sleep(time.Second * 3)
				},
				Config: getGrafanaFolderConfig(title) + getDataSourceGrafanaFolder(title),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fullDs, grafanaFolderTitle, title),
					resource.TestCheckResourceAttrSet(fullDs, grafanaFolderUid),
					resource.TestCheckResourceAttrSet(fullDs, grafanaFolderId),
					resource.TestCheckResourceAttrSet(fullDs, grafanaFolderUrl),
					resource.TestCheckResourceAttrSet(fullDs, grafanaFolderVersion),
				),
			},
		},
	})
}

func getDataSourceGrafanaFolder(title string) string {
	return fmt.Sprintf(`
data "logzio_grafana_folder" "my_ds_folder" {
  title = "%s"
  depends_on = ["logzio_grafana_folder.test_folder"]
}`, title)
}
