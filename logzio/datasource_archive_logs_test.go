package logzio

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/logzio/logzio_terraform_client/archive_logs"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
	"os"
	"testing"
)

func TestAccDataSourceArchiveLogs(t *testing.T) {
	resourceName := "data.logzio_archive_logs.my_archive_datasource"
	path := os.Getenv(envLogzioS3Path)
	arn := os.Getenv(envLogzioAwsArn)
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
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
