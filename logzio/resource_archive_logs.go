package logzio

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/logzio/logzio_terraform_client/archive_logs"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	archiveLogsIdField                  = "archive_id"
	archiveLogsStorageType              = "storage_type"
	archiveLogsEnabled                  = "enabled"
	archiveLogsCompressed               = "compressed"
	archiveLogsAmazonS3StorageSettings  = "amazon_s3_storage_settings"
	archiveLogsS3CredentialsType        = "credentials_type"
	archiveLogsS3Path                   = "s3_path"
	archiveLogsS3SecretCredentials      = "s3_secret_credentials"
	archiveLogsS3AccessKey              = "access_key"
	archiveLogsS3SecretKey              = "secret_key"
	archiveLogsS3IamCredentialsArn      = "s3_iam_credentials_arn"
	archiveLogsS3ExternalId             = "s3_external_id"
	archiveLogsAzureBlobStorageSettings = "azure_blob_storage_settings"
	archiveLogsBlobTenantId             = "tenant_id"
	archiveLogsBlobClientId             = "client_id"
	archiveLogsBlobClientSecret         = "client_secret"
	archiveLogsBlobAccountName          = "account_name"
	archiveLogsBlobContainerName        = "container_name"
	archiveLogsBlobPath                 = "blob_path"
)

// archiveLogsClient returns the archive logs client with the api token from the provider
func archiveLogsClient(m interface{}) *archive_logs.ArchiveLogsClient {
	var client *archive_logs.ArchiveLogsClient
	client, _ = archive_logs.New(m.(Config).apiToken, m.(Config).baseUrl)
	return client
}

func resourceArchiveLogs() *schema.Resource {
	return &schema.Resource{
		Create: resourceArchiveLogsCreate,
		Read:   resourceArchiveLogsRead,
		Update: resourceArchiveLogsUpdate,
		Delete: resourceArchiveLogsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
			archiveLogsAmazonS3StorageSettings: {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						archiveLogsS3CredentialsType: {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: utils.ValidateArchiveLogsAwsCredentialsType,
						},
						archiveLogsS3Path: {
							Type:     schema.TypeString,
							Required: true,
						},
						archiveLogsS3SecretCredentials: {
							Type:     schema.TypeList,
							Optional: true,
							MinItems: 1,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									archiveLogsS3AccessKey: {
										Type:     schema.TypeString,
										Required: true,
									},
									archiveLogsS3SecretKey: {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						archiveLogsS3IamCredentialsArn: {
							Type:     schema.TypeString,
							Optional: true,
						},
						archiveLogsS3ExternalId: {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			archiveLogsAzureBlobStorageSettings: {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						archiveLogsBlobTenantId: {
							Type:     schema.TypeString,
							Required: true,
						},
						archiveLogsBlobClientId: {
							Type:     schema.TypeString,
							Required: true,
						},
						archiveLogsBlobClientSecret: {
							Type:     schema.TypeString,
							Required: true,
						},
						archiveLogsBlobAccountName: {
							Type:     schema.TypeString,
							Required: true,
						},
						archiveLogsBlobContainerName: {
							Type:     schema.TypeString,
							Required: true,
						},
						archiveLogsBlobPath: {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Second),
			Update: schema.DefaultTimeout(5 * time.Second),
			Delete: schema.DefaultTimeout(5 * time.Second),
		},
	}
}

func resourceArchiveLogsCreate(d *schema.ResourceData, m interface{}) error {
	createArchive := getCreateOrUpdateArchiveFromSchema(d)
	archive, err := archiveLogsClient(m).SetupArchive(createArchive)
	if err != nil {
		return err
	}

	d.SetId(strconv.FormatInt(int64(archive.Id), 10))

	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		err = resourceArchiveLogsRead(d, m)
		if err != nil {
			if strings.Contains(err.Error(), "failed with missing archive") {
				return resource.RetryableError(err)
			}
		}

		return resource.NonRetryableError(err)
	})
}

func resourceArchiveLogsRead(d *schema.ResourceData, m interface{}) error {
	id, err := utils.IdFromResourceData(d)
	if err != nil {
		return nil
	}

	archive, err := archiveLogsClient(m).RetrieveArchiveLogsSetting(int32(id))
	setArchive(d, archive)
	return nil
}

func resourceArchiveLogsUpdate(d *schema.ResourceData, m interface{}) error {
	id, err := utils.IdFromResourceData(d)
	if err != nil {
		return err
	}

	updateArchive := getCreateOrUpdateArchiveFromSchema(d)
	_, err = archiveLogsClient(m).UpdateArchiveLogs(int32(id), updateArchive)
	if err != nil {
		return err
	}

	return resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		err = resourceArchiveLogsRead(d, m)
		if err != nil {
			archiveFromSchema := getCreateOrUpdateArchiveFromSchema(d)
			if strings.Contains(err.Error(), "failed with missing archive") &&
				!reflect.DeepEqual(updateArchive, archiveFromSchema) {
				return resource.RetryableError(fmt.Errorf("archive is not updated yet: %s", err.Error()))
			}
		}

		return resource.NonRetryableError(err)
	})
}

func resourceArchiveLogsDelete(d *schema.ResourceData, m interface{}) error {
	archiveId, err := utils.IdFromResourceData(d)
	if err != nil {
		return err
	}

	return resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		err := archiveLogsClient(m).DeleteArchiveLogs(int32(archiveId))
		if err != nil {
			return resource.RetryableError(err)
		}

		return resource.NonRetryableError(err)
	})
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
		settings := d.Get(archiveLogsAmazonS3StorageSettings).([]interface{})
		createArchive.AmazonS3StorageSettings = getS3StorageSettingsFromSchema(settings[0].(map[string]interface{}))
	case archive_logs.StorageTypeBlob:
		settings := d.Get(archiveLogsAzureBlobStorageSettings).([]interface{})
		createArchive.AzureBlobStorageSettings = getBlobStorageSettingsFromSchema(settings[0].(map[string]interface{}))
	default:
		panic("unknown storage type")
	}

	return createArchive
}

func getS3StorageSettingsFromSchema(settings map[string]interface{}) *archive_logs.S3StorageSettings {
	var s3Settings archive_logs.S3StorageSettings
	s3Settings.Path = settings[archiveLogsS3Path].(string)
	s3Settings.CredentialsType = settings[archiveLogsS3CredentialsType].(string)
	switch s3Settings.CredentialsType {
	case archive_logs.CredentialsTypeKeys:
		keysSettings := settings[archiveLogsS3SecretCredentials].([]interface{})[0].(map[string]interface{})
		keys := archive_logs.S3SecretCredentialsObject{
			AccessKey: keysSettings[archiveLogsS3AccessKey].(string),
			SecretKey: keysSettings[archiveLogsS3SecretKey].(string),
		}
		s3Settings.S3SecretCredentials = &keys
	case archive_logs.CredentialsTypeIam:
		iam := archive_logs.S3IamCredentials{Arn: settings[archiveLogsS3IamCredentialsArn].(string)}
		s3Settings.S3IamCredentials = &iam
	default:
		panic("unknown s3 credentials type")
	}

	return &s3Settings
}

func getBlobStorageSettingsFromSchema(settings map[string]interface{}) *archive_logs.BlobSettings {
	var blobSettings archive_logs.BlobSettings

	blobSettings.TenantId = settings[archiveLogsBlobTenantId].(string)
	blobSettings.ClientId = settings[archiveLogsBlobClientId].(string)
	blobSettings.ClientSecret = settings[archiveLogsBlobClientSecret].(string)
	blobSettings.AccountName = settings[archiveLogsBlobAccountName].(string)
	blobSettings.ContainerName = settings[archiveLogsBlobContainerName].(string)

	if path, ok := settings[archiveLogsBlobPath]; ok {
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
	settingsMap := make(map[string]interface{}, 0)
	settingsMap[archiveLogsS3CredentialsType] = s3Settings.CredentialsType
	settingsMap[archiveLogsS3Path] = s3Settings.Path
	switch s3Settings.CredentialsType {
	case archive_logs.CredentialsTypeKeys:
		keys := map[string]interface{}{
			archiveLogsS3AccessKey: s3Settings.S3SecretCredentials.AccessKey,
			archiveLogsS3SecretKey: s3Settings.S3SecretCredentials.SecretKey,
		}
		settingsMap[archiveLogsS3SecretCredentials] = []interface{}{keys}
	case archive_logs.CredentialsTypeIam:
		settingsMap[archiveLogsS3IamCredentialsArn] = s3Settings.S3IamCredentials.Arn
		settingsMap[archiveLogsS3ExternalId] = s3Settings.S3IamCredentials.ExternalId
	default:
		panic("unknown s3 credentials type while setting archive")
	}

	d.Set(archiveLogsAmazonS3StorageSettings, []interface{}{settingsMap})
}

func setBlobSettings(d *schema.ResourceData, blobSettings archive_logs.BlobSettings) {
	settings := map[string]interface{}{
		archiveLogsBlobTenantId:      blobSettings.TenantId,
		archiveLogsBlobClientId:      blobSettings.ClientId,
		archiveLogsBlobClientSecret:  blobSettings.ClientSecret,
		archiveLogsBlobAccountName:   blobSettings.AccountName,
		archiveLogsBlobContainerName: blobSettings.ContainerName,
	}

	if len(blobSettings.Path) > 0 {
		settings[archiveLogsBlobPath] = blobSettings.Path
	}

	d.Set(archiveLogsAzureBlobStorageSettings, []interface{}{settings})
}
