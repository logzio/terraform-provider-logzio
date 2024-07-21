package logzio

import (
	"context"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/logzio/logzio_terraform_client/s3_fetcher"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
	"reflect"
	"strconv"
	"strings"
)

const (
	s3FetcherId                       = "fetcher_id"
	s3FetcherAccessKey                = "aws_access_key"
	s3FetcherSecretKey                = "aws_secret_key"
	s3FetcherArn                      = "aws_arn"
	s3FetcherBucket                   = "bucket_name"
	s3FetcherPrefix                   = "prefix"
	s3FetcherActive                   = "active"
	s3FetcherAddS3ObjectKeyAsLogField = "add_s3_object_key_as_log_field"
	s3FetcherRegion                   = "aws_region"
	s3FetcherLogsType                 = "logs_type"

	s3FetcherRetryAttempts = 8
)

func resourceS3Fetcher() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceS3FetcherCreate,
		ReadContext:   resourceS3FetcherRead,
		UpdateContext: resourceS3FetcherUpdate,
		DeleteContext: resourceS3FetcherDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			s3FetcherId: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			s3FetcherAccessKey: {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			s3FetcherSecretKey: {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			s3FetcherArn: {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			s3FetcherBucket: {
				Type:     schema.TypeString,
				Required: true,
			},
			s3FetcherPrefix: {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			s3FetcherActive: {
				Type:     schema.TypeBool,
				Required: true,
			},
			s3FetcherAddS3ObjectKeyAsLogField: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			s3FetcherRegion: {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: utils.ValidateS3FetcherRegion,
			},
			s3FetcherLogsType: {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: utils.ValidateS3FetcherLogsType,
			},
		},
	}
}

func s3FetcherClient(m interface{}) *s3_fetcher.S3FetcherClient {
	var client *s3_fetcher.S3FetcherClient
	client, _ = s3_fetcher.New(m.(Config).apiToken, m.(Config).baseUrl)
	return client
}

func resourceS3FetcherCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	createFetcher, err := getCreateUpdateS3FetcherFromSchema(ctx, d)
	if err != nil {
		return diag.FromErr(err)
	}

	fmt.Println("NAAMA TEST:")
	fmt.Println(createFetcher)
	fetcher, err := s3FetcherClient(m).CreateS3Fetcher(createFetcher)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(fetcher.Id, 10))
	return resourceS3FetcherRead(ctx, d, m)
}

func resourceS3FetcherRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id, err := utils.IdFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	fetcher, err := s3FetcherClient(m).GetS3Fetcher(id)
	if err != nil {
		tflog.Error(ctx, err.Error())
		if strings.Contains(err.Error(), "missing s3 fetcher") {
			// If we were not able to find the resource - delete from state
			d.SetId("")
			return diag.Diagnostics{}
		} else {
			return diag.FromErr(err)
		}
	}

	setS3Fetcher(d, *fetcher)
	return nil

}

func resourceS3FetcherUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id, err := utils.IdFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	updateFetcher, err := getCreateUpdateS3FetcherFromSchema(ctx, d)
	if err != nil {
		return diag.FromErr(err)
	}

	err = s3FetcherClient(m).UpdateS3Fetcher(id, updateFetcher)
	if err != nil {
		return diag.FromErr(err)
	}

	var diagRet diag.Diagnostics
	readErr := retry.Do(func() error {
		diagRet = resourceS3FetcherRead(ctx, d, m)
		if diagRet.HasError() {
			return fmt.Errorf("received error from read s3 fetcher")
		}

		return nil
	},
		retry.RetryIf(
			// Retry ONLY if the resource was not updated yet
			func(err error) bool {
				if err != nil {
					return false
				} else {
					// Check if the update shows on read
					// if not updated yet - retry
					s3FetcherFromSchema, _ := getCreateUpdateS3FetcherFromSchema(ctx, d)
					return !reflect.DeepEqual(s3FetcherFromSchema, updateFetcher)
				}
			}),
		retry.DelayType(retry.BackOffDelay),
		retry.Attempts(s3FetcherRetryAttempts),
	)

	if readErr != nil {
		tflog.Error(ctx, "could not update schema")
		return diagRet
	}

	return nil
}

func resourceS3FetcherDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id, err := utils.IdFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	err = s3FetcherClient(m).DeleteS3Fetcher(id)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func setS3Fetcher(d *schema.ResourceData, response s3_fetcher.S3FetcherResponse) {
	d.Set(s3FetcherId, response.Id)
	d.Set(s3FetcherBucket, response.Bucket)
	d.Set(s3FetcherActive, response.Active)
	d.Set(s3FetcherRegion, response.Region)
	d.Set(s3FetcherLogsType, response.LogsType)
	d.Set(s3FetcherPrefix, response.Prefix)
	d.Set(s3FetcherAddS3ObjectKeyAsLogField, response.AddS3ObjectKeyAsLogField)
	d.Set(s3FetcherAccessKey, response.AccessKey)
	d.Set(s3FetcherArn, response.Arn)
}

func getCreateUpdateS3FetcherFromSchema(ctx context.Context, d *schema.ResourceData) (s3_fetcher.S3FetcherRequest, error) {
	var request s3_fetcher.S3FetcherRequest
	request.Bucket = d.Get(s3FetcherBucket).(string)
	request.Region = s3_fetcher.AwsRegion(d.Get(s3FetcherRegion).(string))
	request.LogsType = s3_fetcher.AwsLogsType(d.Get(s3FetcherLogsType).(string))
	arn := d.Get(s3FetcherArn).(string)
	accessKey := d.Get(s3FetcherAccessKey).(string)
	secretKey := d.Get(s3FetcherSecretKey).(string)
	if arn == "" && accessKey == "" && secretKey == "" {
		return request, fmt.Errorf("either %s or %s & %s must be set", s3FetcherArn, s3FetcherAccessKey, s3FetcherSecretKey)
	}

	if arn != "" {
		if accessKey != "" && secretKey != "" {
			return request, fmt.Errorf("cannot use both authentication methods. Choose authenticating either with keys OR arn")
		}
		tflog.Debug(ctx, "aws authentication with arn detected")
		request.Arn = arn
	} else {
		if (accessKey == "" && secretKey != "") || (accessKey != "" && secretKey == "") {
			return request, fmt.Errorf("when using keys authentication, both %s and %s must be set", s3FetcherAccessKey, s3FetcherSecretKey)
		}

		tflog.Debug(ctx, "aws authentication with keys detected")
		request.AccessKey = accessKey
		request.SecretKey = secretKey
	}

	active := d.Get(s3FetcherActive).(bool)
	request.Active = &active
	addS3ObjectKeyAsLogField := d.Get(s3FetcherAddS3ObjectKeyAsLogField).(bool)
	request.AddS3ObjectKeyAsLogField = &addS3ObjectKeyAsLogField
	request.Prefix = d.Get(s3FetcherPrefix).(string)

	return request, nil
}
