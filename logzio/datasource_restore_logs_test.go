package logzio

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"os"
	"testing"
)

func TestAccDataSourceRestoreLogs(t *testing.T) {
	archiveName := "archive_for_restore_datasource"
	restoreName := "my_restore_resource"
	restoreDataSourceName := "my_restore_datasource"
	fullDataSourceName := "data.logzio_restore_logs." + restoreDataSourceName
	path := os.Getenv(envLogzioS3Path)
	arn := os.Getenv(envLogzioAwsArn)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckApiToken(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:  getConfigTestArchiveS3Iam(archiveName, path, arn),
				Destroy: false,
			},
			{
				ExpectNonEmptyPlan: true,
				Config:             getConfigTestArchiveS3Iam(archiveName, path, arn) + getConfigTestRestore(restoreName) + getConfigTestRestoreDatasource(restoreName, restoreDataSourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(fullDataSourceName, restoreLogsAccountName),
					resource.TestCheckResourceAttrSet(fullDataSourceName, restoreLogsStartTime),
					resource.TestCheckResourceAttrSet(fullDataSourceName, restoreLogsEndTime),
				),
			},
		},
	})
}

func getConfigTestRestoreDatasource(restoreResourceName string, datasourceName string) string {
	return fmt.Sprintf(`data "logzio_restore_logs" "%s" {
  restore_operation_id = "${logzio_restore_logs.%s.id}"
  depends_on = ["logzio_restore_logs.%s"]
}
`, datasourceName, restoreResourceName, restoreResourceName)
}
