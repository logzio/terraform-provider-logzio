package logzio

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/logzio/logzio_terraform_client/restore_logs"
	"time"
)

const (
	restoreDatasourceRetries = 3
)

func dataSourceRestoreLogs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRestoreLogsRead,
		Schema: map[string]*schema.Schema{
			restoreLogsId: {
				Type:     schema.TypeInt,
				Optional: true,
			},
			restoreLogsAccountName: {
				Type:     schema.TypeString,
				Computed: true,
			},
			restoreLogsStartTime: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			restoreLogsEndTime: {
				Type:     schema.TypeInt,
				Computed: true,
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
			Read: schema.DefaultTimeout(5 * time.Second),
		},
	}
}

func dataSourceRestoreLogsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, _ := restore_logs.New(m.(Config).apiToken, m.(Config).baseUrl)
	restoreIdStr, ok := d.GetOk(restoreLogsId)

	if ok {
		id := int32(restoreIdStr.(int))
		restore, err := getRestore(client, id, restoreDatasourceRetries)
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(fmt.Sprintf("%d", restore.Id))
		setRestore(d, restore)
		return nil
	}

	return diag.Errorf("couldn't find restore operation with specified id")

}

func getRestore(client *restore_logs.RestoreClient, restoreId int32, retries int) (*restore_logs.RestoreOperation, error) {
	restore, err := client.GetRestoreOperation(restoreId)
	if err != nil && retries > 0 {
		time.Sleep(time.Second * 2)
		restore, err = getRestore(client, restoreId, retries-1)
	}
	return restore, err
}
