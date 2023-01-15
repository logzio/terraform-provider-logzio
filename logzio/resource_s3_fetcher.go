package logzio

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/logzio/logzio_terraform_client/s3_fetcher"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
	"google.golang.org/appengine/log"
	"strconv"
)

const (
	s3FetcherId                       = "id"
	s3FetcherAccessKey                = "aws_access_key"
	s3FetcherSecretKey                = "aws_secret_key"
	s3FetcherArn                      = "aws_arn"
	s3FetcherBucket                   = "bucket_name"
	s3FetcherPrefix                   = "prefix"
	s3FetcherActive                   = "active"
	s3FetcherAddS3ObjectKeyAsLogField = "add_s3_object_key_as_log_field"
	s3FetcherRegion                   = "aws_region"
	s3FetcherLogsType                 = "logs_type"
)

func resourceS3Fetcher() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceS3FetcherCreate,
		ReadContext:   resourceS3FetcherRead,
		//UpdateContext: resourceSubAccountUpdate,
		//DeleteContext: resourceSubAccountDelete,
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
	createFetcher, err := getCreateUpdateFromSchema(ctx, d)
	if err != nil {
		return diag.FromErr(err)
	}

	fetcher, err := s3FetcherClient(m).CreateS3Fetcher(createFetcher)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(fetcher.Id, 10))
	return resourceS3FetcherRead(ctx, d, m)
}

func resourceS3FetcherRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
}

func getCreateUpdateFromSchema(ctx context.Context, d *schema.ResourceData) (s3_fetcher.S3FetcherRequest, error) {
	request := s3_fetcher.S3FetcherRequest{
		Bucket:   d.Get(s3FetcherBucket).(string),
		Region:   d.Get(s3FetcherRegion).(s3_fetcher.AwsRegion),
		LogsType: d.Get(s3FetcherLogsType).(s3_fetcher.AwsLogsType),
	}

	arn := d.Get(s3FetcherArn).(string)
	accessKey := d.Get(s3FetcherAccessKey).(string)
	secretKey := d.Get(s3FetcherSecretKey).(string)
	if arn == "" && accessKey == "" && secretKey == "" {
		return request, fmt.Errorf("either %s or %s & %s must be set", s3FetcherArn, s3FetcherAccessKey, s3FetcherSecretKey)
	}

	if arn != "" {
		log.Debugf(ctx, "aws authentication with arn detected")
		request.Arn = arn
	} else {
		if (accessKey == "" && secretKey != "") || (accessKey != "" && secretKey == "") {
			return request, fmt.Errorf("when using keys authentication, both %s and %s must be set", s3FetcherAccessKey, s3FetcherSecretKey)
		}

		log.Debugf(ctx, "aws authentication with keys detected")
		request.AccessKey = accessKey
		request.SecretKey = secretKey
	}

	active := d.Get(s3FetcherActive).(bool)
	request.Active = &active
	addS3ObjectKeyAsLogField := d.Get(s3FetcherAddS3ObjectKeyAsLogField).(bool)
	request.AddS3ObjectKeyAsLogField = &addS3ObjectKeyAsLogField

	return request, nil
}
