package logzio

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/logzio/logzio_terraform_client/restore_logs"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
	"strconv"
	"strings"
	"time"
)

const (
	restoreLogsId               = "restore_operation_id"
	restoreLogsAccountName      = "account_name"
	restoreLogsStartTime        = "start_time"
	restoreLogsEndTime          = "end_time"
	restoreLogsAccountId        = "account_id"
	restoreLogsRestoredVolumeGb = "restored_volume_gb"
	restoreLogsStatus           = "status"
	restoreLogsCreatedAt        = "created_at"
	restoreLogsStartedAt        = "started_at"
	restoreLogsFinishedAt       = "finished_at"
	restoreLogsExpiresAt        = "expires_at"
)

// restoreLogsClient returns the restore logs client with the api token from the provider
func restoreLogsClient(m interface{}) *restore_logs.RestoreClient {
	var client *restore_logs.RestoreClient
	client, _ = restore_logs.New(m.(Config).apiToken, m.(Config).baseUrl)
	return client
}

func resourceRestoreLogs() *schema.Resource {
	return &schema.Resource{
		Create: resourceRestoreLogsCreate,
		Read:   resourceRestoreLogsRead,
		Delete: resourceRestoreLogsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			restoreLogsId: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			restoreLogsAccountName: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			restoreLogsStartTime: {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			restoreLogsEndTime: {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			restoreLogsAccountId: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			restoreLogsRestoredVolumeGb: {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			restoreLogsStatus: {
				Type:     schema.TypeString,
				Computed: true,
			},
			restoreLogsCreatedAt: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			restoreLogsStartedAt: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			restoreLogsFinishedAt: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			restoreLogsExpiresAt: {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Second),
			Update: schema.DefaultTimeout(5 * time.Second),
			Delete: schema.DefaultTimeout(5 * time.Second),
		},
	}
}

func resourceRestoreLogsCreate(d *schema.ResourceData, m interface{}) error {
	initiateRestore := getCreateRestoreFromSchema(d)
	restore, err := restoreLogsClient(m).InitiateRestoreOperation(initiateRestore)
	if err != nil {
		return err
	}

	d.SetId(strconv.FormatInt(int64(restore.Id), 10))
	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		err = resourceRestoreLogsRead(d, m)
		if err != nil {
			if strings.Contains(err.Error(), "failed with missing restore") {
				return resource.RetryableError(err)
			}
		}

		return resource.NonRetryableError(err)
	})
}

func resourceRestoreLogsRead(d *schema.ResourceData, m interface{}) error {
	id, err := utils.IdFromResourceData(d)
	if err != nil {
		return nil
	}

	restore, err := restoreLogsClient(m).GetRestoreOperation(int32(id))
	if err != nil {
		return err
	}

	setRestore(d, restore)
	return nil
}

func resourceRestoreLogsDelete(d *schema.ResourceData, m interface{}) error {
	id, err := utils.IdFromResourceData(d)
	if err != nil {
		return err
	}

	return resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		_, err = restoreLogsClient(m).DeleteRestoreOperation(int32(id))
		if err != nil {
			return resource.RetryableError(err)
		}

		return resource.NonRetryableError(err)
	})
}

func getCreateRestoreFromSchema(d *schema.ResourceData) restore_logs.InitiateRestore {
	return restore_logs.InitiateRestore{
		AccountName: d.Get(restoreLogsAccountName).(string),
		StartTime:   int64(d.Get(restoreLogsStartTime).(int)),
		EndTime:     int64(d.Get(restoreLogsEndTime).(int)),
	}
}

func setRestore(d *schema.ResourceData, restore *restore_logs.RestoreOperation) {
	d.Set(restoreLogsId, restore.Id)
	d.Set(restoreLogsAccountName, restore.AccountName)
	d.Set(restoreLogsStartTime, restore.StartTime)
	d.Set(restoreLogsEndTime, restore.EndTime)
	d.Set(restoreLogsAccountId, restore.AccountId)
	d.Set(restoreLogsRestoredVolumeGb, restore.RestoredVolumeGb)
	d.Set(restoreLogsStatus, restore.Status)
	d.Set(restoreLogsCreatedAt, restore.CreatedAt)
	d.Set(restoreLogsStartedAt, restore.StartedAt)
	d.Set(restoreLogsFinishedAt, restore.FinishedAt)
	d.Set(restoreLogsExpiresAt, restore.ExpiresAt)
}
