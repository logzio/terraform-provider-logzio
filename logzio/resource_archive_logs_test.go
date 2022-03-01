package logzio

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
		PreCheck:  func() { testAccPreCheckApiToken(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: getConfigTestArchiveS3Keys(resourceName, path, accessKey, secretKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(fullResourceName, archiveLogsStorageType, archive_logs.StorageTypeS3),
					resource.TestCheckResourceAttr(fullResourceName, archiveLogsEnabled, "false"),
					resource.TestCheckResourceAttr(fullResourceName,
						archiveLogsAmazonS3StorageSettings+".0."+archiveLogsS3CredentialsType,
						archive_logs.CredentialsTypeKeys),
					resource.TestCheckResourceAttr(fullResourceName,
						archiveLogsAmazonS3StorageSettings+".0."+archiveLogsS3Path, path),
					resource.TestCheckResourceAttr(fullResourceName,
						archiveLogsAmazonS3StorageSettings+".0."+archiveLogsS3SecretCredentials+".0."+archiveLogsS3AccessKey,
						accessKey),
					resource.TestCheckResourceAttr(fullResourceName,
						archiveLogsAmazonS3StorageSettings+".0."+archiveLogsS3SecretCredentials+".0."+archiveLogsS3SecretKey,
						secretKey),
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

func TestAccLogzioArchiveLogs_SetupArchiveS3Iam(t *testing.T) {
	path := os.Getenv(envLogzioS3Path)
	arn := os.Getenv(envLogzioAwsArn)
	resourceName := "setup_test_s3_iam"
	fullResourceName := "logzio_archive_logs." + resourceName
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckApiToken(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: getConfigTestArchiveS3Iam(resourceName, path, arn),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(fullResourceName, archiveLogsStorageType, archive_logs.StorageTypeS3),
					resource.TestCheckResourceAttr(fullResourceName, archiveLogsEnabled, "true"),
					resource.TestCheckResourceAttr(fullResourceName, archiveLogsCompressed, "false"),
					resource.TestCheckResourceAttr(fullResourceName,
						archiveLogsAmazonS3StorageSettings+".0."+archiveLogsS3CredentialsType,
						archive_logs.CredentialsTypeIam),
					resource.TestCheckResourceAttr(fullResourceName,
						archiveLogsAmazonS3StorageSettings+".0."+archiveLogsS3Path, path),
					resource.TestCheckResourceAttr(fullResourceName,
						archiveLogsAmazonS3StorageSettings+".0."+archiveLogsS3IamCredentialsArn,
						arn),
					resource.TestCheckResourceAttrSet(fullResourceName,
						archiveLogsAmazonS3StorageSettings+".0."+archiveLogsS3ExternalId),
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
		PreCheck:  func() { testAccPreCheckApiToken(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: getConfigTestArchiveBlob(resourceName, tenantId, clientId, clientSecret, accountName, containerName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(fullResourceName, archiveLogsStorageType, archive_logs.StorageTypeBlob),
					resource.TestCheckResourceAttr(fullResourceName, archiveLogsEnabled, "true"),
					resource.TestCheckResourceAttr(fullResourceName, archiveLogsCompressed, "true"),
					resource.TestCheckResourceAttr(fullResourceName,
						archiveLogsAzureBlobStorageSettings+".0."+archiveLogsBlobTenantId,
						tenantId),
					resource.TestCheckResourceAttr(fullResourceName,
						archiveLogsAzureBlobStorageSettings+".0."+archiveLogsBlobClientId, clientId),
					resource.TestCheckResourceAttr(fullResourceName,
						archiveLogsAzureBlobStorageSettings+".0."+archiveLogsBlobClientSecret,
						clientSecret),
					resource.TestCheckResourceAttr(fullResourceName,
						archiveLogsAzureBlobStorageSettings+".0."+archiveLogsBlobAccountName, accountName),
					resource.TestCheckResourceAttr(fullResourceName,
						archiveLogsAzureBlobStorageSettings+".0."+archiveLogsBlobContainerName, containerName),
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

func TestAccLogzioArchiveLogs_SetupArchiveEmptyStorageType(t *testing.T) {
	path := os.Getenv(envLogzioS3Path)
	accessKey := os.Getenv(envLogzioAwsAccessKey)
	secretKey := os.Getenv(envLogzioAwsSecretKey)
	resourceName := "setup_test_empty_storage_type"
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckApiToken(t) },
		Providers: testAccProviders,
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
		PreCheck:  func() { testAccPreCheckApiToken(t) },
		Providers: testAccProviders,
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
		PreCheck:  func() { testAccPreCheckApiToken(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: getConfigTestArchiveS3Keys(resourceName, path, accessKey, secretKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(fullResourceName, archiveLogsStorageType, archive_logs.StorageTypeS3),
					resource.TestCheckResourceAttr(fullResourceName, archiveLogsEnabled, "false"),
					resource.TestCheckResourceAttr(fullResourceName,
						archiveLogsAmazonS3StorageSettings+".0."+archiveLogsS3CredentialsType,
						archive_logs.CredentialsTypeKeys),
					resource.TestCheckResourceAttr(fullResourceName,
						archiveLogsAmazonS3StorageSettings+".0."+archiveLogsS3Path, path),
					resource.TestCheckResourceAttr(fullResourceName,
						archiveLogsAmazonS3StorageSettings+".0."+archiveLogsS3SecretCredentials+".0."+archiveLogsS3AccessKey,
						accessKey),
					resource.TestCheckResourceAttr(fullResourceName,
						archiveLogsAmazonS3StorageSettings+".0."+archiveLogsS3SecretCredentials+".0."+archiveLogsS3SecretKey,
						secretKey),
				),
			},
			{
				PreConfig: func() {
					time.Sleep(2 * time.Second)
				},
				Config: getConfigTestArchiveS3Keys(resourceName, path, newAccessKey, newSecretKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(fullResourceName, archiveLogsStorageType, archive_logs.StorageTypeS3),
					resource.TestCheckResourceAttr(fullResourceName, archiveLogsEnabled, "false"),
					resource.TestCheckResourceAttr(fullResourceName,
						archiveLogsAmazonS3StorageSettings+".0."+archiveLogsS3CredentialsType,
						archive_logs.CredentialsTypeKeys),
					resource.TestCheckResourceAttr(fullResourceName,
						archiveLogsAmazonS3StorageSettings+".0."+archiveLogsS3Path, path),
					resource.TestCheckResourceAttr(fullResourceName,
						archiveLogsAmazonS3StorageSettings+".0."+archiveLogsS3SecretCredentials+".0."+archiveLogsS3AccessKey,
						newAccessKey),
					resource.TestCheckResourceAttr(fullResourceName,
						archiveLogsAmazonS3StorageSettings+".0."+archiveLogsS3SecretCredentials+".0."+archiveLogsS3SecretKey,
						newSecretKey),
				),
			},
		},
	})
}

func getConfigTestArchiveS3Keys(name string, path string, accessKey string, secretKey string) string {
	return fmt.Sprintf(`resource "logzio_archive_logs" "%s" {
 storage_type = "S3"
 enabled = false
 amazon_s3_storage_settings {
   credentials_type = "KEYS"
   s3_path = "%s"
   s3_secret_credentials {
		access_key = "%s"
		secret_key = "%s"
	}
 }
}
`, name, path, accessKey, secretKey)
}

func getConfigTestArchiveS3Iam(name string, path string, arn string) string {
	return fmt.Sprintf(`resource "logzio_archive_logs" "%s" {
 storage_type = "S3"
 compressed = false
 amazon_s3_storage_settings {
   credentials_type = "IAM"
   s3_path = "%s"
   s3_iam_credentials_arn = "%s"
 }
}
`, name, path, arn)
}

func getConfigTestArchiveBlob(name string, tenantId string, clientId string,
	clientSecret string, accountName string, containerName string) string {
	return fmt.Sprintf(`resource "logzio_archive_logs" "%s" {
 storage_type = "BLOB"
 azure_blob_storage_settings {
	tenant_id = "%s"
	client_id = "%s"
	client_secret = "%s"
	account_name = "%s"
	container_name = "%s" 
 }
}
`, name, tenantId, clientId, clientSecret, accountName, containerName)
}

func getConfigTestArchiveEmptyStorageType(name string, path string, accessKey string, secretKey string) string {
	return fmt.Sprintf(`resource "logzio_archive_logs" "%s" {
 storage_type = ""
 enabled = false
 amazon_s3_storage_settings {
   credentials_type = "KEYS"
   s3_path = "%s"
   s3_secret_credentials {
		access_key = "%s"
		secret_key = "%s"
	}
 }
}
`, name, path, accessKey, secretKey)
}

func getConfigTestArchiveEmptyCredentialsType(name string, path string, accessKey string, secretKey string) string {
	return fmt.Sprintf(`resource "logzio_archive_logs" "%s" {
 storage_type = "S3"
 enabled = false
 amazon_s3_storage_settings {
   credentials_type = ""
   s3_path = "%s"
   s3_secret_credentials {
		access_key = "%s"
		secret_key = "%s"
	}
 }
}
`, name, path, accessKey, secretKey)
}
