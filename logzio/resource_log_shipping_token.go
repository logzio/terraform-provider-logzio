package logzio

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/logzio/logzio_terraform_client/log_shipping_tokens"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	logShippingTokenId        = "id"
	logShippingTokenName      = "name"
	logShippingTokenEnabled   = "enabled"
	logShippingTokenToken     = "token"
	logShippingTokenUpdatedAt = "updated_at"
	logShippingTokenUpdatedBy = "updated_by"
	logShippingTokenCreatedAt = "created_at"
	logShippingTokenCreatedBy = "created_by"
	logShippingTokenTokenId   = "token_id"
)

func resourceLogShippingToken() *schema.Resource {
	return &schema.Resource{
		Create: resourceLogShippingTokenCreate,
		Read:   resourceLogShippingTokenRead,
		Update: resourceLogShippingTokenUpdate,
		Delete: resourceLogShippingTokenDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			logShippingTokenTokenId: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			logShippingTokenName: {
				Type:     schema.TypeString,
				Required: true,
			},
			logShippingTokenEnabled: {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			logShippingTokenToken: {
				Type:     schema.TypeString,
				Computed: true,
			},
			logShippingTokenUpdatedAt: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			logShippingTokenUpdatedBy: {
				Type:     schema.TypeString,
				Computed: true,
			},
			logShippingTokenCreatedAt: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			logShippingTokenCreatedBy: {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Second),
			Update: schema.DefaultTimeout(10 * time.Second),
			Delete: schema.DefaultTimeout(10 * time.Second),
		},
	}
}

// Creates a new log shipping token in logz.io
func resourceLogShippingTokenCreate(d *schema.ResourceData, m interface{}) error {
	createToken := log_shipping_tokens.CreateLogShippingToken{Name: d.Get(logShippingTokenName).(string)}
	tokenLimits, err := logShippingTokenClient(m).GetLogShippingLimitsToken()
	if err != nil {
		return err
	}

	if tokenLimits.MaxAllowedTokens > tokenLimits.NumOfEnabledTokens {
		token, err := logShippingTokenClient(m).CreateLogShippingToken(createToken)
		if err != nil {
			return err
		}

		d.SetId(strconv.FormatInt(int64(token.Id), 10))

		return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			err = resourceLogShippingTokenRead(d, m)
			if err != nil {
				if strings.Contains(err.Error(), "failed with missing log shipping") {
					return resource.RetryableError(err)
				}
			}

			return resource.NonRetryableError(err)
		})
	}

	return fmt.Errorf("cannot create new log shipping token. max allowed tokens for account: %d. number of enabled tokens: :%d",
		tokenLimits.MaxAllowedTokens, tokenLimits.NumOfEnabledTokens)

}

func resourceLogShippingTokenRead(d *schema.ResourceData, m interface{}) error {
	id, err := idFromResourceData(d)
	if err != nil {
		return err
	}

	token, err := logShippingTokenClient(m).GetLogShippingToken(int32(id))
	if err != nil {
		return err
	}

	setLogShippingToken(d, token)
	return nil
}

func resourceLogShippingTokenUpdate(d *schema.ResourceData, m interface{}) error {
	id, err := idFromResourceData(d)
	if err != nil {
		return err
	}

	updateToken := log_shipping_tokens.UpdateLogShippingToken{
		Name:    d.Get(logShippingTokenName).(string),
		Enabled: strconv.FormatBool(d.Get(logShippingTokenEnabled).(bool)),
	}

	_, err = logShippingTokenClient(m).UpdateLogShippingToken(int32(id), updateToken)
	if err != nil {
		return err
	}

	return resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		err = resourceLogShippingTokenRead(d, m)
		if err != nil {
			tokenFromSchema := log_shipping_tokens.UpdateLogShippingToken{
				Name:    d.Get(logShippingTokenName).(string),
				Enabled: strconv.FormatBool(d.Get(logShippingTokenEnabled).(bool)),
			}

			if strings.Contains(err.Error(), "failed with missing log shipping") &&
				!reflect.DeepEqual(updateToken, tokenFromSchema) {
				return resource.RetryableError(fmt.Errorf("log shipping token is not updated yet: %s", err.Error()))
			}
		}

		return resource.NonRetryableError(err)
	})
}

func resourceLogShippingTokenDelete(d *schema.ResourceData, m interface{}) error {
	id, err := idFromResourceData(d)
	if err != nil {
		return err
	}

	return resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		err = logShippingTokenClient(m).DeleteLogShippingToken(int32(id))
		if err != nil {
			return resource.RetryableError(err)
		}

		return resource.NonRetryableError(err)
	})
}

func logShippingTokenClient(m interface{}) *log_shipping_tokens.LogShippingTokensClient {
	var client *log_shipping_tokens.LogShippingTokensClient
	client, _ = log_shipping_tokens.New(m.(Config).apiToken, m.(Config).baseUrl)
	return client
}

func setLogShippingToken(d *schema.ResourceData, token *log_shipping_tokens.LogShippingToken) {
	d.Set(logShippingTokenTokenId, token.Id)
	d.Set(logShippingTokenName, token.Name)
	d.Set(logShippingTokenEnabled, token.Enabled)
	d.Set(logShippingTokenToken, token.Token)
	d.Set(logShippingTokenUpdatedAt, token.UpdatedAt)
	d.Set(logShippingTokenUpdatedBy, token.UpdatedBy)
	d.Set(logShippingTokenCreatedAt, token.CreatedAt)
	d.Set(logShippingTokenCreatedBy, token.CreatedBy)
}
