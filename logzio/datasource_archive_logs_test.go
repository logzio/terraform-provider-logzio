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
				Config:  getConfigResourceArchiveLogs(path, arn),
				Destroy: false,
			},
			{
				Config: getConfigResourceArchiveLogs(path, arn) + getConfigDatasourceArchiveLogs(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, archiveLogsStorageType, archive_logs.StorageTypeS3),
				),
			},
		},
	})
}

func getConfigResourceArchiveLogs(path, arn string) string {
	return fmt.Sprintf(`resource "logzio_archive_logs" "test_to_datasource" {
  storage_type = "S3"
  compressed = false
  aws_credentials_type = "IAM"
  aws_s3_path = "%s"
  aws_s3_iam_credentials_arn = "%s"
}
`, path, arn)
}

func getConfigDatasourceArchiveLogs() string {
	return `data "logzio_archive_logs" "my_archive_datasource" {
  archive_id = "${logzio_archive_logs.test_to_datasource.id}"
  depends_on = ["logzio_archive_logs.test_to_datasource"]
}
`
}
