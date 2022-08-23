package logzio

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/logzio/logzio_terraform_client/archive_logs"
	"time"
)

const (
	archiveDatasourceRetries = 3
)

func dataSourceArchiveLogs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceArchiveLogsRead,
		Schema: map[string]*schema.Schema{
			archiveLogsIdField: {
				Type:     schema.TypeInt,
				Optional: true,
			},
			archiveLogsStorageType: {
				Type:     schema.TypeString,
				Computed: true,
			},
			archiveLogsEnabled: {
				Type:     schema.TypeBool,
				Computed: true,
			},
			archiveLogsCompressed: {
				Type:     schema.TypeBool,
				Computed: true,
			},
			archiveLogsS3CredentialsType: {
				Type:     schema.TypeString,
				Computed: true,
			},
			archiveLogsS3Path: {
				Type:     schema.TypeString,
				Computed: true,
			},
			archiveLogsS3AccessKey: {
				Type:     schema.TypeString,
				Computed: true,
			},
			archiveLogsS3SecretKey: {
				Type:     schema.TypeString,
				Computed: true,
			},
			archiveLogsS3IamCredentialsArn: {
				Type:     schema.TypeString,
				Computed: true,
			},
			//archiveLogsS3ExternalId: {
			//	Type:     schema.TypeString,
			//	Computed: true,
			//},
			archiveLogsBlobTenantId: {
				Type:     schema.TypeString,
				Computed: true,
			},
			archiveLogsBlobClientId: {
				Type:     schema.TypeString,
				Computed: true,
			},
			archiveLogsBlobClientSecret: {
				Type:     schema.TypeString,
				Computed: true,
			},
			archiveLogsBlobAccountName: {
				Type:     schema.TypeString,
				Computed: true,
			},
			archiveLogsBlobContainerName: {
				Type:     schema.TypeString,
				Computed: true,
			},
			archiveLogsBlobPath: {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceArchiveLogsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	archiveId, ok := d.GetOk(archiveLogsIdField)

	if ok {
		archive, err := getArchive(int64(archiveId.(int)), archiveDatasourceRetries, m)
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(fmt.Sprintf("%d", archive.Id))
		setArchive(d, archive)
		return nil
	}

	return diag.Errorf("couldn't find archive with specified id")
}

func getArchive(archiveId int64, retries int, m interface{}) (*archive_logs.ArchiveLogs, error) {
	archive, err := archiveLogsClient(m).RetrieveArchiveLogsSetting(int32(archiveId))
	if err != nil && retries > 0 {
		time.Sleep(time.Second * 2)
		archive, err = getArchive(archiveId, retries-1, m)
	}
	return archive, err
}
