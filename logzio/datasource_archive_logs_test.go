package logzio

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/logzio/logzio_terraform_client/archive_logs"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
	"os"
	"testing"
)

func TestAccDataSourceArchiveLogs(t *testing.T) {
	resourceName := "data.logzio_archive_logs.my_archive_datasource"
	path := os.Getenv(envLogzioS3Path)
	arn := os.Getenv(envLogzioAwsArn)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckApiToken(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				ExpectNonEmptyPlan: true,
				Config: fmt.Sprintf(utils.ReadFixtureFromFile("create_archive_logs_datasource.tf"),
					path, arn),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, archiveLogsStorageType, archive_logs.StorageTypeS3),
				),
			},
		},
	})
}
