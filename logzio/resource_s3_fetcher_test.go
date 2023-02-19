package logzio

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/logzio/logzio_terraform_client/s3_fetcher"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
	"os"
	"regexp"
	"testing"
)

const (
	bucketName = "terraform-auto-tests"
)

func TestAccLogzioS3Fetcher_S3FetcherKeys(t *testing.T) {
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getS3FetcherConfigKeys(true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("logzio_s3_fetcher.test_fetcher", s3FetcherBucket, bucketName),
					resource.TestCheckResourceAttrSet("logzio_s3_fetcher.test_fetcher", s3FetcherAccessKey),
					resource.TestCheckResourceAttrSet("logzio_s3_fetcher.test_fetcher", s3FetcherSecretKey),
					resource.TestCheckResourceAttrSet("logzio_s3_fetcher.test_fetcher", s3FetcherId),
					resource.TestCheckResourceAttr("logzio_s3_fetcher.test_fetcher", s3FetcherArn, ""),
					resource.TestCheckResourceAttr("logzio_s3_fetcher.test_fetcher", s3FetcherActive, "true"),
					resource.TestCheckResourceAttr("logzio_s3_fetcher.test_fetcher", s3FetcherAddS3ObjectKeyAsLogField, "false"),
					resource.TestCheckResourceAttr("logzio_s3_fetcher.test_fetcher", s3FetcherRegion, s3_fetcher.RegionUsEast1.String()),
					resource.TestCheckResourceAttr("logzio_s3_fetcher.test_fetcher", s3FetcherLogsType, s3_fetcher.LogsTypeElb.String()),
				),
			},
			{
				Config: getS3FetcherConfigKeys(false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("logzio_s3_fetcher.test_fetcher", s3FetcherBucket, bucketName),
					resource.TestCheckResourceAttrSet("logzio_s3_fetcher.test_fetcher", s3FetcherAccessKey),
					resource.TestCheckResourceAttrSet("logzio_s3_fetcher.test_fetcher", s3FetcherSecretKey),
					resource.TestCheckResourceAttr("logzio_s3_fetcher.test_fetcher", s3FetcherArn, ""),
					resource.TestCheckResourceAttr("logzio_s3_fetcher.test_fetcher", s3FetcherActive, "false"),
					resource.TestCheckResourceAttr("logzio_s3_fetcher.test_fetcher", s3FetcherAddS3ObjectKeyAsLogField, "false"),
					resource.TestCheckResourceAttr("logzio_s3_fetcher.test_fetcher", s3FetcherRegion, s3_fetcher.RegionUsEast1.String()),
					resource.TestCheckResourceAttr("logzio_s3_fetcher.test_fetcher", s3FetcherLogsType, s3_fetcher.LogsTypeElb.String()),
				),
			},
			{
				Config:                  getS3FetcherConfigKeys(false),
				ResourceName:            "logzio_s3_fetcher.test_fetcher",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{s3FetcherSecretKey},
			},
		},
	})
}

func TestAccLogzioS3Fetcher_S3FetcherArn(t *testing.T) {
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getS3FetcherConfigArn(true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("logzio_s3_fetcher.test_fetcher", s3FetcherBucket, bucketName),
					resource.TestCheckResourceAttrSet("logzio_s3_fetcher.test_fetcher", s3FetcherArn),
					resource.TestCheckResourceAttr("logzio_s3_fetcher.test_fetcher", s3FetcherAccessKey, ""),
					resource.TestCheckResourceAttr("logzio_s3_fetcher.test_fetcher", s3FetcherActive, "true"),
					resource.TestCheckResourceAttr("logzio_s3_fetcher.test_fetcher", s3FetcherAddS3ObjectKeyAsLogField, "false"),
					resource.TestCheckResourceAttr("logzio_s3_fetcher.test_fetcher", s3FetcherRegion, s3_fetcher.RegionUsEast1.String()),
					resource.TestCheckResourceAttr("logzio_s3_fetcher.test_fetcher", s3FetcherLogsType, s3_fetcher.LogsTypeElb.String()),
				),
			},
			{
				Config: getS3FetcherConfigArn(false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("logzio_s3_fetcher.test_fetcher", s3FetcherBucket, bucketName),
					resource.TestCheckResourceAttrSet("logzio_s3_fetcher.test_fetcher", s3FetcherArn),
					resource.TestCheckResourceAttr("logzio_s3_fetcher.test_fetcher", s3FetcherAccessKey, ""),
					resource.TestCheckResourceAttr("logzio_s3_fetcher.test_fetcher", s3FetcherActive, "false"),
					resource.TestCheckResourceAttr("logzio_s3_fetcher.test_fetcher", s3FetcherAddS3ObjectKeyAsLogField, "false"),
					resource.TestCheckResourceAttr("logzio_s3_fetcher.test_fetcher", s3FetcherRegion, s3_fetcher.RegionUsEast1.String()),
					resource.TestCheckResourceAttr("logzio_s3_fetcher.test_fetcher", s3FetcherLogsType, s3_fetcher.LogsTypeElb.String()),
				),
			},
			{
				Config:            getS3FetcherConfigArn(false),
				ResourceName:      "logzio_s3_fetcher.test_fetcher",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLogzioS3Fetcher_S3FetcherInvalidRegion(t *testing.T) {
	defer utils.SleepAfterTest()
	terraformPlan := fmt.Sprintf(`
resource "logzio_s3_fetcher" "test_fetcher" {
  aws_access_key = "%s"
  aws_secret_key = "%s"
  bucket_name = "%s"
  active = false
  aws_region = "SOME_REGION"
  logs_type = "elb"
}
`, os.Getenv(envLogzioAwsAccessKey), os.Getenv(envLogzioAwsSecretKey), bucketName)
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      terraformPlan,
				ExpectError: regexp.MustCompile("is not in the allowed aws regions list"),
			},
		},
	})
}

func TestAccLogzioS3Fetcher_S3FetcherInvalidLogsType(t *testing.T) {
	defer utils.SleepAfterTest()
	terraformPlan := fmt.Sprintf(`
resource "logzio_s3_fetcher" "test_fetcher" {
  aws_access_key = "%s"
  aws_secret_key = "%s"
  bucket_name = "%s"
  active = false
  aws_region = "US_EAST_1"
  logs_type = "my_type"
}
`, os.Getenv(envLogzioAwsAccessKey), os.Getenv(envLogzioAwsSecretKey), bucketName)
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      terraformPlan,
				ExpectError: regexp.MustCompile("is not in the allowed logs types list"),
			},
		},
	})
}

func TestAccLogzioS3Fetcher_S3FetcherNoBucketName(t *testing.T) {
	defer utils.SleepAfterTest()
	terraformPlan := fmt.Sprintf(`
resource "logzio_s3_fetcher" "test_fetcher" {
  aws_access_key = "%s"
  aws_secret_key = "%s"
  active = false
  aws_region = "US_EAST_1"
  logs_type = "elb"
}
`, os.Getenv(envLogzioAwsAccessKey), os.Getenv(envLogzioAwsSecretKey))
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      terraformPlan,
				ExpectError: regexp.MustCompile("The argument \"bucket_name\" is required, but no definition was found"),
			},
		},
	})
}

func TestAccLogzioS3Fetcher_S3FetcherNoActive(t *testing.T) {
	defer utils.SleepAfterTest()
	terraformPlan := fmt.Sprintf(`
resource "logzio_s3_fetcher" "test_fetcher" {
  aws_access_key = "%s"
  aws_secret_key = "%s"
  bucket_name = "%s"
  aws_region = "US_EAST_1"
  logs_type = "elb"
}
`, os.Getenv(envLogzioAwsAccessKey), os.Getenv(envLogzioAwsSecretKey), bucketName)
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      terraformPlan,
				ExpectError: regexp.MustCompile("The argument \"active\" is required, but no definition was found"),
			},
		},
	})
}

func TestAccLogzioS3Fetcher_S3FetcherNoAwsAuth(t *testing.T) {
	defer utils.SleepAfterTest()
	terraformPlan := fmt.Sprintf(`
resource "logzio_s3_fetcher" "test_fetcher" {
  bucket_name = "%s"
  active = false
  aws_region = "US_EAST_1"
  logs_type = "elb"
}
`, bucketName)
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      terraformPlan,
				ExpectError: regexp.MustCompile(fmt.Sprintf("either %s or %s & %s must be set", s3FetcherArn, s3FetcherAccessKey, s3FetcherSecretKey)),
			},
		},
	})
}

func TestAccLogzioS3Fetcher_S3FetcherMissingSecretKey(t *testing.T) {
	defer utils.SleepAfterTest()
	terraformPlan := fmt.Sprintf(`
resource "logzio_s3_fetcher" "test_fetcher" {
  aws_access_key = "%s"
  bucket_name = "%s"
  active = false
  aws_region = "US_EAST_1"
  logs_type = "elb"
}
`, os.Getenv(envLogzioAwsAccessKey), bucketName)
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      terraformPlan,
				ExpectError: regexp.MustCompile(fmt.Sprintf("when using keys authentication, both %s and %s must be set", s3FetcherAccessKey, s3FetcherSecretKey)),
			},
		},
	})
}

func TestAccLogzioS3Fetcher_S3FetcherMissingAccessKey(t *testing.T) {
	defer utils.SleepAfterTest()
	terraformPlan := fmt.Sprintf(`
resource "logzio_s3_fetcher" "test_fetcher" {
  aws_secret_key = "%s"
  bucket_name = "%s"
  active = false
  aws_region = "US_EAST_1"
  logs_type = "elb"
}
`, os.Getenv(envLogzioAwsSecretKey), bucketName)
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      terraformPlan,
				ExpectError: regexp.MustCompile(fmt.Sprintf("when using keys authentication, both %s and %s must be set", s3FetcherAccessKey, s3FetcherSecretKey)),
			},
		},
	})
}

func TestAccLogzioS3Fetcher_S3FetcherAuthNoAccess(t *testing.T) {
	defer utils.SleepAfterTest()
	terraformPlan := fmt.Sprintf(`
resource "logzio_s3_fetcher" "test_fetcher" {
  aws_access_key = "some_access_key"
  aws_secret_key = "some_secret_key"
  bucket_name = "%s"
  active = false
  aws_region = "US_EAST_1"
  logs_type = "elb"
}
`, bucketName)
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      terraformPlan,
				ExpectError: regexp.MustCompile("API call CreateS3Fetcher failed with status code 403"),
			},
		},
	})
}

func TestAccLogzioS3Fetcher_S3FetcherAllAuthMethods(t *testing.T) {
	defer utils.SleepAfterTest()
	terraformPlan := fmt.Sprintf(`
resource "logzio_s3_fetcher" "test_fetcher" {
  aws_access_key = "%s"
  aws_secret_key = "%s"
  aws_arn = "%s"
  bucket_name = "%s"
  active = false
  aws_region = "US_EAST_1"
  logs_type = "elb"
}
`, os.Getenv(envLogzioAwsAccessKey), os.Getenv(envLogzioAwsSecretKey), os.Getenv(envLogzioAwsArnS3Fetcher), bucketName)
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      terraformPlan,
				ExpectError: regexp.MustCompile("cannot use both authentication methods. Choose authenticating either with keys OR arn"),
			},
		},
	})
}

func getS3FetcherConfigKeys(active bool) string {
	return fmt.Sprintf(`
resource "logzio_s3_fetcher" "test_fetcher" {
  aws_access_key = "%s"
  aws_secret_key = "%s"
  bucket_name = "%s"
  active = %t
  aws_region = "US_EAST_1"
  logs_type = "elb"
}
`, os.Getenv(envLogzioAwsAccessKey), os.Getenv(envLogzioAwsSecretKey), bucketName, active)
}

func getS3FetcherConfigArn(active bool) string {
	return fmt.Sprintf(`
resource "logzio_s3_fetcher" "test_fetcher" {
  aws_arn = "%s"
  bucket_name = "%s"
  active = %t
  aws_region = "US_EAST_1"
  logs_type = "elb"
}
`, os.Getenv(envLogzioAwsArnS3Fetcher), bucketName, active)
}
