package logzio

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/logzio/logzio_terraform_client/archive_logs"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
	"os"
	"regexp"
	"testing"
	"time"
)

const (
	envLogzioAwsAccessKeyUpdate = "AWS_ACCESS_KEY_UPDATE" // for update test
	envLogzioAwsSecretKeyUpdate = "AWS_SECRET_KEY_UPDATE" // for update test
)

func TestAccLogzioArchiveLogs_SetupArchiveS3Keys(t *testing.T) {
	path := os.Getenv(envLogzioS3Path)
	accessKey := os.Getenv(envLogzioAwsAccessKey)
	secretKey := os.Getenv(envLogzioAwsSecretKey)
	resourceName := "setup_test_s3_keys"
	fullResourceName := "logzio_archive_logs." + resourceName
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getConfigTestArchiveS3Keys(resourceName, path, accessKey, secretKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(fullResourceName, archiveLogsStorageType, archive_logs.StorageTypeS3),
					resource.TestCheckResourceAttr(fullResourceName, archiveLogsEnabled, "false"),
					resource.TestCheckResourceAttr(fullResourceName, archiveLogsS3CredentialsType, archive_logs.CredentialsTypeKeys),
					resource.TestCheckResourceAttrSet(fullResourceName, archiveLogsS3Path),
					resource.TestCheckResourceAttrSet(fullResourceName, archiveLogsS3AccessKey),
				),
			},
			{
				ResourceName:            fullResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{archiveLogsS3SecretKey},
			},
		},
	})
}

func TestAccLogzioArchiveLogs_SetupArchiveS3Iam(t *testing.T) {
	path := os.Getenv(envLogzioS3Path)
	arn := os.Getenv(envLogzioAwsArn)
	resourceName := "setup_test_s3_iam"
	fullResourceName := "logzio_archive_logs." + resourceName
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getConfigTestArchiveS3Iam(resourceName, path, arn),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(fullResourceName, archiveLogsStorageType, archive_logs.StorageTypeS3),
					resource.TestCheckResourceAttr(fullResourceName, archiveLogsEnabled, "true"),
					resource.TestCheckResourceAttr(fullResourceName, archiveLogsCompressed, "false"),
					resource.TestCheckResourceAttrSet(fullResourceName, archiveLogsS3CredentialsType),
					resource.TestCheckResourceAttrSet(fullResourceName, archiveLogsS3Path),
					resource.TestCheckResourceAttrSet(fullResourceName, archiveLogsS3IamCredentialsArn),
					//resource.TestCheckResourceAttrSet(fullResourceName, archiveLogsS3ExternalId),
				),
			},
			{
				ResourceName:      fullResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLogzioArchiveLogs_SetupArchiveBlob(t *testing.T) {
	tenantId := os.Getenv(envLogzioAzureTenantId)
	clientId := os.Getenv(envLogzioAzureClientId)
	clientSecret := os.Getenv(envLogzioAzureClientSecret)
	accountName := os.Getenv(envLogzioAzureAccountName)
	containerName := os.Getenv(envLogzioAzureContainerName)

	resourceName := "setup_test_blob"
	fullResourceName := "logzio_archive_logs." + resourceName
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getConfigTestArchiveBlob(resourceName, tenantId, clientId, clientSecret, accountName, containerName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(fullResourceName, archiveLogsStorageType, archive_logs.StorageTypeBlob),
					resource.TestCheckResourceAttr(fullResourceName, archiveLogsEnabled, "true"),
					resource.TestCheckResourceAttr(fullResourceName, archiveLogsCompressed, "true"),
					resource.TestCheckResourceAttrSet(fullResourceName, archiveLogsBlobTenantId),
					resource.TestCheckResourceAttrSet(fullResourceName, archiveLogsBlobClientId),
					resource.TestCheckResourceAttrSet(fullResourceName, archiveLogsBlobClientSecret),
					resource.TestCheckResourceAttrSet(fullResourceName, archiveLogsBlobAccountName),
					resource.TestCheckResourceAttrSet(fullResourceName, archiveLogsBlobContainerName),
				)},
			{
				ResourceName:            fullResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{archiveLogsBlobClientSecret},
			},
		},
	})
}

func TestAccLogzioArchiveLogs_SetupArchiveEmptyStorageType(t *testing.T) {
	path := os.Getenv(envLogzioS3Path)
	accessKey := os.Getenv(envLogzioAwsAccessKey)
	secretKey := os.Getenv(envLogzioAwsSecretKey)
	resourceName := "setup_test_empty_storage_type"
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      getConfigTestArchiveEmptyStorageType(resourceName, path, accessKey, secretKey),
				ExpectError: regexp.MustCompile("value for storage type is unknown"),
			},
		},
	})
}

func TestAccLogzioArchiveLogs_SetupArchiveEmptyS3CredentialsType(t *testing.T) {
	path := os.Getenv(envLogzioS3Path)
	accessKey := os.Getenv(envLogzioAwsAccessKey)
	secretKey := os.Getenv(envLogzioAwsSecretKey)
	resourceName := "setup_test_empty_credentials_type"
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      getConfigTestArchiveEmptyCredentialsType(resourceName, path, accessKey, secretKey),
				ExpectError: regexp.MustCompile("value for credentials type is unknown"),
			},
		},
	})
}

func TestAccLogzioArchiveLogs_UpdateArchive(t *testing.T) {
	path := os.Getenv(envLogzioS3Path)
	accessKey := os.Getenv(envLogzioAwsAccessKey)
	secretKey := os.Getenv(envLogzioAwsSecretKey)
	resourceName := "update_test"
	fullResourceName := "logzio_archive_logs." + resourceName
	newAccessKey := os.Getenv(envLogzioAwsAccessKeyUpdate)
	newSecretKey := os.Getenv(envLogzioAwsSecretKeyUpdate)
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getConfigTestArchiveS3Keys(resourceName, path, accessKey, secretKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(fullResourceName, archiveLogsStorageType, archive_logs.StorageTypeS3),
					resource.TestCheckResourceAttr(fullResourceName, archiveLogsEnabled, "false"),
					resource.TestCheckResourceAttrSet(fullResourceName, archiveLogsS3CredentialsType),
					resource.TestCheckResourceAttrSet(fullResourceName, archiveLogsS3Path),
					resource.TestCheckResourceAttrSet(fullResourceName, archiveLogsS3AccessKey),
					resource.TestCheckResourceAttrSet(fullResourceName, archiveLogsS3SecretKey),
				),
			},
			{
				PreConfig: func() {
					time.Sleep(3 * time.Second)
				},
				Config: getConfigTestArchiveS3Keys(resourceName, path, newAccessKey, newSecretKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(fullResourceName, archiveLogsStorageType, archive_logs.StorageTypeS3),
					resource.TestCheckResourceAttr(fullResourceName, archiveLogsEnabled, "false"),
					resource.TestCheckResourceAttrSet(fullResourceName, archiveLogsS3CredentialsType),
					resource.TestCheckResourceAttrSet(fullResourceName, archiveLogsS3Path),
					resource.TestCheckResourceAttrSet(fullResourceName, archiveLogsS3AccessKey),
					resource.TestCheckResourceAttrSet(fullResourceName, archiveLogsS3SecretKey),
				),
			},
		},
	})
}

func getConfigTestArchiveS3Keys(name string, path string, accessKey string, secretKey string) string {
	return fmt.Sprintf(`resource "logzio_archive_logs" "%s" {
 storage_type = "S3"
 enabled = false
 credentials_type = "KEYS"
 s3_path = "%s"
 aws_access_key = "%s"
 aws_secret_key = "%s"
}
`, name, path, accessKey, secretKey)
}

func getConfigTestArchiveS3Iam(name string, path string, arn string) string {
	return fmt.Sprintf(`resource "logzio_archive_logs" "%s" {
 storage_type = "S3"
 compressed = false
 credentials_type = "IAM"
 s3_path = "%s"
 s3_iam_credentials_arn = "%s"
}
`, name, path, arn)
}

func getConfigTestArchiveBlob(name string, tenantId string, clientId string,
	clientSecret string, accountName string, containerName string) string {
	return fmt.Sprintf(`resource "logzio_archive_logs" "%s" {
 storage_type = "BLOB"
 tenant_id = "%s"
 client_id = "%s"
 client_secret = "%s"
 account_name = "%s"
 container_name = "%s" 
}
`, name, tenantId, clientId, clientSecret, accountName, containerName)
}

func getConfigTestArchiveEmptyStorageType(name string, path string, accessKey string, secretKey string) string {
	return fmt.Sprintf(`resource "logzio_archive_logs" "%s" {
 storage_type = ""
 enabled = false
 credentials_type = "KEYS"
 s3_path = "%s"
 aws_access_key = "%s"
 aws_secret_key = "%s"
}
`, name, path, accessKey, secretKey)
}

func getConfigTestArchiveEmptyCredentialsType(name string, path string, accessKey string, secretKey string) string {
	return fmt.Sprintf(`resource "logzio_archive_logs" "%s" {
 storage_type = "S3"
 enabled = false
 credentials_type = ""
 s3_path = "%s"
 aws_access_key = "%s"
 aws_secret_key = "%s"
}
`, name, path, accessKey, secretKey)
}
