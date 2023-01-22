package logzio

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/logzio/logzio_terraform_client/s3_fetcher"
	"strconv"
)

func dataSourceS3Fetcher() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceS3FetcherRead,
		Schema: map[string]*schema.Schema{
			s3FetcherId: {
				Type:     schema.TypeInt,
				Required: true,
			},
			s3FetcherAccessKey: {
				Type:      schema.TypeString,
				Computed:  true,
				Optional:  true,
				Sensitive: true,
			},
			s3FetcherArn: {
				Type:      schema.TypeString,
				Computed:  true,
				Optional:  true,
				Sensitive: true,
			},
			s3FetcherBucket: {
				Type:     schema.TypeString,
				Computed: true,
			},
			s3FetcherPrefix: {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			s3FetcherActive: {
				Type:     schema.TypeBool,
				Computed: true,
			},
			s3FetcherAddS3ObjectKeyAsLogField: {
				Type:     schema.TypeBool,
				Computed: true,
			},
			s3FetcherRegion: {
				Type:     schema.TypeString,
				Computed: true,
			},
			s3FetcherLogsType: {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceS3FetcherRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id, ok := d.GetOk(s3FetcherId)
	if ok {
		client, err := s3_fetcher.New(m.(Config).apiToken, m.(Config).baseUrl)
		if err != nil {
			return diag.FromErr(err)
		}

		fetcher, err := client.GetS3Fetcher(int64(id.(int)))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(strconv.FormatInt(fetcher.Id, 10))
		setS3Fetcher(d, *fetcher)

		return nil
	}

	return diag.Errorf("could not get fetcher id")
}
