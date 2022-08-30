package logzio

import (
	"context"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/logzio/logzio_terraform_client/log_shipping_tokens"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
	"reflect"
	"strconv"
	"strings"
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
		CreateContext: resourceLogShippingTokenCreate,
		ReadContext:   resourceLogShippingTokenRead,
		UpdateContext: resourceLogShippingTokenUpdate,
		DeleteContext: resourceLogShippingTokenDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
	}
}

// resourceLogShippingTokenCreate creates a new log shipping token in logz.io
func resourceLogShippingTokenCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	createToken := log_shipping_tokens.CreateLogShippingToken{Name: d.Get(logShippingTokenName).(string)}
	tokenLimits, err := logShippingTokenClient(m).GetLogShippingLimitsToken()
	if err != nil {
		return diag.FromErr(err)
	}

	// Check if we exceeded the number of max allowed tokens
	if tokenLimits.NumOfEnabledTokens < tokenLimits.MaxAllowedTokens {
		token, err := logShippingTokenClient(m).CreateLogShippingToken(createToken)
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(strconv.FormatInt(int64(token.Id), 10))
		return resourceLogShippingTokenRead(ctx, d, m)
	}

	return diag.Errorf("cannot create new log shipping token. max allowed tokens for account: %d. number of enabled tokens: :%d",
		tokenLimits.MaxAllowedTokens, tokenLimits.NumOfEnabledTokens)
}

// resourceLogShippingTokenRead gets log shipping token by id
func resourceLogShippingTokenRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id, err := utils.IdFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	var token *log_shipping_tokens.LogShippingToken
	readErr := retry.Do(
		func() error {
			token, err = logShippingTokenClient(m).GetLogShippingToken(int32(id))
			if err != nil {
				return err
			}

			return nil
		},
		retry.RetryIf(
			func(err error) bool {
				if err != nil {
					if strings.Contains(err.Error(), "failed with missing log shipping") {
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

	setLogShippingToken(d, token)
	return nil
}

// resourceLogShippingTokenUpdate updates log shipping token by id
func resourceLogShippingTokenUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id, err := utils.IdFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	updateToken := log_shipping_tokens.UpdateLogShippingToken{
		Name:    d.Get(logShippingTokenName).(string),
		Enabled: strconv.FormatBool(d.Get(logShippingTokenEnabled).(bool)),
	}

	_, err = logShippingTokenClient(m).UpdateLogShippingToken(int32(id), updateToken)
	if err != nil {
		return diag.FromErr(err)
	}

	var diagRet diag.Diagnostics
	readErr := retry.Do(
		func() error {
			diagRet = resourceLogShippingTokenRead(ctx, d, m)
			if diagRet.HasError() {
				return fmt.Errorf("received error from read log shipping token")
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
					tokenFromSchema := log_shipping_tokens.UpdateLogShippingToken{
						Name:    d.Get(logShippingTokenName).(string),
						Enabled: strconv.FormatBool(d.Get(logShippingTokenEnabled).(bool)),
					}

					return !reflect.DeepEqual(updateToken, tokenFromSchema)
				}
			}),
		retry.DelayType(retry.BackOffDelay),
		retry.Attempts(15),
	)

	if readErr != nil {
		tflog.Error(ctx, "could not update schema")
		return diagRet
	}

	return nil
}

// resourceLogShippingTokenDelete deletes log shipping token by id
func resourceLogShippingTokenDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id, err := utils.IdFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	deleteErr := retry.Do(
		func() error {
			return logShippingTokenClient(m).DeleteLogShippingToken(int32(id))
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
