package logzio

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
)

func TestAccLogzioGrafanaFolder_GrafanaFolder(t *testing.T) {
	defer utils.SleepAfterTest()
	randomSuffix := utils.RandomString(6)
	title := fmt.Sprintf("tf_provider_test_%s", randomSuffix)
	resourceType := "logzio_grafana_folder"
	resourceName := "test_folder"
	fullResourceName := fmt.Sprintf("%s.%s", resourceType, resourceName)
	newTitle := fmt.Sprintf("tf_provider_updated_%s", randomSuffix)
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getGrafanaFolderConfig(title),
				Check: resource.ComposeTestCheckFunc(
					awaitApply(30),
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
					awaitApply(30),
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
