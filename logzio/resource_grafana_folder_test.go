package logzio

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
	"testing"
)

func TestAccLogzioGrafanaFolder_GrafanaFolder(t *testing.T) {
	defer utils.SleepAfterTest()

	title := "tf_provider_test"
	resourceType := "logzio_grafana_folder"
	resourceName := "test_folder"
	fullResourceName := fmt.Sprintf("%s.%s", resourceType, resourceName)
	newTitle := "tf_provider_updated"
	resource.Test(t, resource.TestCase{
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
				Config: getGrafanaFolderConfig(newTitle),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fullResourceName, grafanaFolderTitle, newTitle),
					resource.TestCheckResourceAttrSet(fullResourceName, grafanaFolderUid),
					resource.TestCheckResourceAttrSet(fullResourceName, grafanaFolderId),
					resource.TestCheckResourceAttrSet(fullResourceName, grafanaFolderVersion),
					resource.TestCheckResourceAttrSet(fullResourceName, grafanaFolderUrl),
				),
			},
			{
				Config:            getGrafanaFolderConfig(newTitle),
				ResourceName:      fullResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func getGrafanaFolderConfig(title string) string {
	return fmt.Sprintf(`
resource "logzio_grafana_folder" "test_folder" {
  title = "%s"
}
`, title)
}
