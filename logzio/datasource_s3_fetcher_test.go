package logzio

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/logzio/logzio_terraform_client/s3_fetcher"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
	"testing"
	"time"
)

func TestAccDataSourceS3Fetcher(t *testing.T) {
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getS3FetcherConfigKeys(false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("logzio_s3_fetcher.test_fetcher", s3FetcherBucket, bucketName),
					resource.TestCheckResourceAttrSet("logzio_s3_fetcher.test_fetcher", s3FetcherAccessKey),
					resource.TestCheckResourceAttrSet("logzio_s3_fetcher.test_fetcher", s3FetcherSecretKey),
					resource.TestCheckResourceAttrSet("logzio_s3_fetcher.test_fetcher", s3FetcherId),
					resource.TestCheckResourceAttr("logzio_s3_fetcher.test_fetcher", s3FetcherArn, ""),
					resource.TestCheckResourceAttr("logzio_s3_fetcher.test_fetcher", s3FetcherActive, "false"),
					resource.TestCheckResourceAttr("logzio_s3_fetcher.test_fetcher", s3FetcherAddS3ObjectKeyAsLogField, "false"),
					resource.TestCheckResourceAttr("logzio_s3_fetcher.test_fetcher", s3FetcherRegion, s3_fetcher.RegionUsEast1.String()),
					resource.TestCheckResourceAttr("logzio_s3_fetcher.test_fetcher", s3FetcherLogsType, s3AccessType),
				),
			},
			{
				PreConfig: func() {
					time.Sleep(time.Second * 3)
				},
				Config: getS3FetcherConfigKeys(false) + getDataSourceS3Fetcher(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.logzio_s3_fetcher.my_ds_fetcher", s3FetcherBucket, bucketName),
					resource.TestCheckResourceAttrSet("data.logzio_s3_fetcher.my_ds_fetcher", s3FetcherAccessKey),
					resource.TestCheckResourceAttrSet("data.logzio_s3_fetcher.my_ds_fetcher", s3FetcherId),
					resource.TestCheckResourceAttr("data.logzio_s3_fetcher.my_ds_fetcher", s3FetcherArn, ""),
					resource.TestCheckResourceAttr("data.logzio_s3_fetcher.my_ds_fetcher", s3FetcherActive, "false"),
					resource.TestCheckResourceAttr("data.logzio_s3_fetcher.my_ds_fetcher", s3FetcherAddS3ObjectKeyAsLogField, "false"),
					resource.TestCheckResourceAttr("data.logzio_s3_fetcher.my_ds_fetcher", s3FetcherRegion, s3_fetcher.RegionUsEast1.String()),
					resource.TestCheckResourceAttr("data.logzio_s3_fetcher.my_ds_fetcher", s3FetcherLogsType, s3AccessType),
				),
			},
		},
	})
}

func getDataSourceS3Fetcher() string {
	return fmt.Sprintf(`
data "logzio_s3_fetcher" "my_ds_fetcher" {
  fetcher_id = "${logzio_s3_fetcher.test_fetcher.fetcher_id}"
  depends_on = ["logzio_s3_fetcher.test_fetcher"]
}`)
}
