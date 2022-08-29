package logzio

import (
	"context"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/logzio/logzio_terraform_client/archive_logs"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
	"reflect"
	"strconv"
	"strings"
)

const (
	archiveLogsIdField             = "archive_id"
	archiveLogsStorageType         = "storage_type"
	archiveLogsEnabled             = "enabled"
	archiveLogsCompressed          = "compressed"
	archiveLogsS3CredentialsType   = "credentials_type"
	archiveLogsS3Path              = "s3_path"
	archiveLogsS3AccessKey         = "aws_access_key"
	archiveLogsS3SecretKey         = "aws_secret_key"
	archiveLogsS3IamCredentialsArn = "s3_iam_credentials_arn"
	archiveLogsBlobTenantId        = "tenant_id"
	archiveLogsBlobClientId        = "client_id"
	archiveLogsBlobClientSecret    = "client_secret"
	archiveLogsBlobAccountName     = "account_name"
	archiveLogsBlobContainerName   = "container_name"
	archiveLogsBlobPath            = "blob_path"
)

// archiveLogsClient returns the archive logs client with the api token from the provider
func archiveLogsClient(m interface{}) *archive_logs.ArchiveLogsClient {
	var client *archive_logs.ArchiveLogsClient
	client, _ = archive_logs.New(m.(Config).apiToken, m.(Config).baseUrl)
	return client
}

func resourceArchiveLogs() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceArchiveLogsCreate,
		ReadContext:   resourceArchiveLogsRead,
		UpdateContext: resourceArchiveLogsUpdate,
		DeleteContext: resourceArchiveLogsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			archiveLogsIdField: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			archiveLogsStorageType: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: utils.ValidateArchiveLogsStorageType,
			},
			archiveLogsEnabled: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			archiveLogsCompressed: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			archiveLogsS3CredentialsType: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: utils.ValidateArchiveLogsAwsCredentialsType,
				Sensitive:    true,
			},
			archiveLogsS3Path: {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			archiveLogsS3AccessKey: {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			archiveLogsS3SecretKey: {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			archiveLogsS3IamCredentialsArn: {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			archiveLogsBlobTenantId: {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			archiveLogsBlobClientId: {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			archiveLogsBlobClientSecret: {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			archiveLogsBlobAccountName: {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			archiveLogsBlobContainerName: {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			archiveLogsBlobPath: {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceArchiveLogsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	createArchive := getCreateOrUpdateArchiveFromSchema(d)
	archive, err := archiveLogsClient(m).SetupArchive(createArchive)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(int64(archive.Id), 10))
	if createArchive.StorageType == archive_logs.StorageTypeS3 &&
		createArchive.AmazonS3StorageSettings.CredentialsType == archive_logs.CredentialsTypeKeys {
		setAwsSecretKey(d, createArchive.AmazonS3StorageSettings.S3SecretCredentials.SecretKey)
	}

	if createArchive.StorageType == archive_logs.StorageTypeBlob {
		setBlobClientSecret(d, createArchive.AzureBlobStorageSettings.ClientSecret)
	}

	return resourceArchiveLogsRead(ctx, d, m)
}

func resourceArchiveLogsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id, err := utils.IdFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	var archive *archive_logs.ArchiveLogs
	readErr := retry.Do(
		func() error {
			archive, err = archiveLogsClient(m).RetrieveArchiveLogsSetting(int32(id))
			if err != nil {
				return err
			}

			return nil
		},
		retry.RetryIf(
			func(err error) bool {
				if err != nil {
					if strings.Contains(err.Error(), "failed with missing archive") {
						return true
					}
				}
				return false
			}),
		retry.DelayType(retry.BackOffDelay),
		retry.Attempts(15),
	)

	if readErr != nil {
		// If we were not able to find the resource - delete from state
		d.SetId("")
		return diag.FromErr(err)
	}

	setArchive(d, archive)
	return nil
}

func resourceArchiveLogsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id, err := utils.IdFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	updateArchive := getCreateOrUpdateArchiveFromSchema(d)
	_, err = archiveLogsClient(m).UpdateArchiveLogs(int32(id), updateArchive)
	if err != nil {
		return diag.FromErr(err)
	}

	var diagRet diag.Diagnostics
	readErr := retry.Do(
		func() error {
			diagRet = resourceArchiveLogsRead(ctx, d, m)
			if diagRet.HasError() {
				return fmt.Errorf("received error from read archive")
			}

			return nil
		},
		retry.RetryIf(
			func(err error) bool {
				if err != nil {
					return true
				} else {
					// Check if the update shows on read
					// if not updated yet - retry
					archiveFromSchema := getCreateOrUpdateArchiveFromSchema(d)
					return !reflect.DeepEqual(updateArchive, archiveFromSchema)
				}
			}),
		retry.DelayType(retry.BackOffDelay),
		retry.Attempts(15),
	)

	if readErr != nil {
		tflog.Error(ctx, "could not update schema")
		return diagRet
	}

	if updateArchive.StorageType == archive_logs.StorageTypeS3 &&
		updateArchive.AmazonS3StorageSettings.CredentialsType == archive_logs.CredentialsTypeKeys {
		setAwsSecretKey(d, updateArchive.AmazonS3StorageSettings.S3SecretCredentials.SecretKey)
	}

	if updateArchive.StorageType == archive_logs.StorageTypeBlob {
		setBlobClientSecret(d, updateArchive.AzureBlobStorageSettings.ClientSecret)
	}

	return nil
}

func resourceArchiveLogsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	archiveId, err := utils.IdFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	deleteErr := retry.Do(
		func() error {
			return archiveLogsClient(m).DeleteArchiveLogs(int32(archiveId))
		},
		retry.RetryIf(
			func(err error) bool {
				return err != nil
			}),
		retry.DelayType(retry.BackOffDelay),
		retry.Attempts(15),
	)

	if deleteErr != nil {
		return diag.FromErr(deleteErr)
	}

	d.SetId("")
	return nil
}

func getCreateOrUpdateArchiveFromSchema(d *schema.ResourceData) archive_logs.CreateOrUpdateArchiving {
	var createArchive archive_logs.CreateOrUpdateArchiving

	createArchive.StorageType = d.Get(archiveLogsStorageType).(string)
	createArchive.Enabled = new(bool)
	*createArchive.Enabled = d.Get(archiveLogsEnabled).(bool)
	createArchive.Compressed = new(bool)
	*createArchive.Compressed = d.Get(archiveLogsCompressed).(bool)
	switch createArchive.StorageType {
	case archive_logs.StorageTypeS3:
		createArchive.AmazonS3StorageSettings = getS3StorageSettingsFromSchema(d)
	case archive_logs.StorageTypeBlob:
		createArchive.AzureBlobStorageSettings = getBlobStorageSettingsFromSchema(d)
	default:
		panic("unknown storage type")
	}

	return createArchive
}

func getS3StorageSettingsFromSchema(d *schema.ResourceData) *archive_logs.S3StorageSettings {
	var s3Settings archive_logs.S3StorageSettings
	s3Settings.Path = d.Get(archiveLogsS3Path).(string)
	s3Settings.CredentialsType = d.Get(archiveLogsS3CredentialsType).(string)
	switch s3Settings.CredentialsType {
	case archive_logs.CredentialsTypeKeys:
		s3Settings.S3SecretCredentials = new(archive_logs.S3SecretCredentialsObject)
		s3Settings.S3SecretCredentials.AccessKey = d.Get(archiveLogsS3AccessKey).(string)
		s3Settings.S3SecretCredentials.SecretKey = d.Get(archiveLogsS3SecretKey).(string)
	case archive_logs.CredentialsTypeIam:
		s3Settings.S3IamCredentials = new(archive_logs.S3IamCredentials)
		s3Settings.S3IamCredentials.Arn = d.Get(archiveLogsS3IamCredentialsArn).(string)
	default:
		panic("unknown s3 credentials type")
	}

	return &s3Settings
}

func getBlobStorageSettingsFromSchema(d *schema.ResourceData) *archive_logs.BlobSettings {
	var blobSettings archive_logs.BlobSettings

	blobSettings.TenantId = d.Get(archiveLogsBlobTenantId).(string)
	blobSettings.ClientId = d.Get(archiveLogsBlobClientId).(string)
	blobSettings.ClientSecret = d.Get(archiveLogsBlobClientSecret).(string)
	blobSettings.AccountName = d.Get(archiveLogsBlobAccountName).(string)
	blobSettings.ContainerName = d.Get(archiveLogsBlobContainerName).(string)

	if path, ok := d.GetOk(archiveLogsBlobPath); ok {
		blobSettings.Path = path.(string)
	}

	return &blobSettings
}

func setArchive(d *schema.ResourceData, archive *archive_logs.ArchiveLogs) {
	d.Set(archiveLogsIdField, archive.Id)
	d.Set(archiveLogsStorageType, archive.Settings.StorageType)
	d.Set(archiveLogsEnabled, archive.Settings.Enabled)
	d.Set(archiveLogsCompressed, archive.Settings.Compressed)
	switch archive.Settings.StorageType {
	case archive_logs.StorageTypeS3:
		setS3Settings(d, archive.Settings.AmazonS3StorageSettings)
	case archive_logs.StorageTypeBlob:
		setBlobSettings(d, archive.Settings.AzureBlobStorageSettings)
	default:
		panic("unknown storage type while setting archive")
	}
}

func setS3Settings(d *schema.ResourceData, s3Settings archive_logs.S3StorageSettings) {
	d.Set(archiveLogsS3CredentialsType, s3Settings.CredentialsType)
	d.Set(archiveLogsS3Path, s3Settings.Path)
	switch s3Settings.CredentialsType {
	case archive_logs.CredentialsTypeKeys:
		d.Set(archiveLogsS3AccessKey, s3Settings.S3SecretCredentials.AccessKey)
	case archive_logs.CredentialsTypeIam:
		d.Set(archiveLogsS3IamCredentialsArn, s3Settings.S3IamCredentials.Arn)
	default:
		panic("unknown s3 credentials type while setting archive")
	}
}

func setBlobSettings(d *schema.ResourceData, blobSettings archive_logs.BlobSettings) {
	d.Set(archiveLogsBlobTenantId, blobSettings.TenantId)
	d.Set(archiveLogsBlobClientId, blobSettings.ClientId)
	d.Set(archiveLogsBlobAccountName, blobSettings.AccountName)
	d.Set(archiveLogsBlobContainerName, blobSettings.ContainerName)

	if len(blobSettings.Path) > 0 {
		d.Set(archiveLogsBlobPath, blobSettings.Path)
	}
}

func setAwsSecretKey(d *schema.ResourceData, secretKey string) {
	d.Set(archiveLogsS3SecretKey, secretKey)
}

func setBlobClientSecret(d *schema.ResourceData, clientSecret string) {
	d.Set(archiveLogsBlobClientSecret, clientSecret)
}
