package logzio

import (
	"context"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/logzio/logzio_terraform_client/metrics_accounts"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
	"reflect"
	"strconv"
	"strings"
)

const (
	metricsAccountId                 string = "account_id"
	metricsAccountEmail              string = "email"
	metricsAccountName               string = "account_name"
	metricsAccountToken              string = "account_token"
	metricsAccountPlanUts            string = "plan_uts"
	metricsAccountAuthorizedAccounts string = "authorized_accounts"
	metricsAccountSoftLimit          string = "soft_limit_unique_metrics"

	metricsAccountRetryAttempts = 8

	metricsAccountNotConsumptionErrorCode = "NOT_CONSUMPTION_ACCOUNT"
)

// The endpoint resource schema, what terraform uses to parse and read the template
func resourceMetricsAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMetricsAccountCreate,
		ReadContext:   resourceMetricsAccountRead,
		UpdateContext: resourceMetricsAccountUpdate,
		DeleteContext: resourceMetricsAccountDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			metricsAccountId: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			metricsAccountToken: {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			metricsAccountEmail: {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			metricsAccountName: {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			metricsAccountPlanUts: {
				Type:     schema.TypeInt,
				Optional: true,
			},
			metricsAccountAuthorizedAccounts: {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Optional: true,
				Computed: true,
			},
			metricsAccountSoftLimit: {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},
		},
	}
}

func MetricsAccountClient(m interface{}) (*metrics_accounts.MetricsAccountClient, error) {
	var client *metrics_accounts.MetricsAccountClient
	var clientError error
	client, clientError = metrics_accounts.New(m.(Config).apiToken, m.(Config).baseUrl)

	if clientError != nil {
		return nil, clientError
	}
	return client, nil
}

func resourceMetricsAccountCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	createSubAccount := getCreateMetricsAccountFromSchema(d)
	MetricsClient, err := MetricsAccountClient(m)
	if err != nil {
		return diag.FromErr(err)
	}

	metricsAccount, err := MetricsClient.CreateMetricsAccount(createSubAccount)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(int64(metricsAccount.Id), 10))
	d.Set(metricsAccountToken, metricsAccount.Token)
	d.Set(metricsAccountId, metricsAccount.Id)

	// the soft limit lives behind its own endpoint, so it can only be applied once the
	// account exists. The raw config is checked rather than GetOk, which reports an explicit 0 as unset.
	if isMetricsAccountSoftLimitConfigured(d) {
		_, err = MetricsClient.UpdateMetricsAccountSoftLimit(int64(metricsAccount.Id),
			metrics_accounts.UpdateMetricsAccountSoftLimit{SoftLimitUniqueMetrics: int32(d.Get(metricsAccountSoftLimit).(int))})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceMetricsAccountRead(ctx, d, m)
}

func resourceMetricsAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id, err := utils.IdFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}
	MetricsClient, err := MetricsAccountClient(m)
	if err != nil {
		return diag.FromErr(err)
	}

	metricsAccount, err := MetricsClient.GetMetricsAccount(id)
	if err != nil {
		tflog.Error(ctx, err.Error())
		if strings.Contains(err.Error(), "missing metrics account") {
			// If we were not able to find the resource - delete from state
			d.SetId("")
			return diag.Diagnostics{}
		} else {
			return diag.FromErr(err)
		}

	}

	setMetricsAccount(d, metricsAccount)

	// The soft limit endpoint is consumption-only and rejects a subscription owner with 400.
	// That is not a failure of this resource, so the field is simply left unset. Any other error is.
	softLimit, err := MetricsClient.GetMetricsAccountSoftLimit(id)
	if err != nil {
		if !strings.Contains(err.Error(), metricsAccountNotConsumptionErrorCode) {
			return diag.FromErr(err)
		}
		tflog.Debug(ctx, fmt.Sprintf("not reading soft limit for metrics account %d: %s", id, err))
	} else if softLimit.SoftLimitUniqueMetrics != nil {
		d.Set(metricsAccountSoftLimit, *softLimit.SoftLimitUniqueMetrics)
	}

	return nil
}

func resourceMetricsAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id, err := utils.IdFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}
	MetricsClient, err := MetricsAccountClient(m)
	if err != nil {
		return diag.FromErr(err)
	}

	updateMetricsAccount := getCreateMetricsAccountFromSchema(d)
	err = MetricsClient.UpdateMetricsAccount(id, updateMetricsAccount)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange(metricsAccountSoftLimit) {
		_, err = MetricsClient.UpdateMetricsAccountSoftLimit(id,
			metrics_accounts.UpdateMetricsAccountSoftLimit{SoftLimitUniqueMetrics: int32(d.Get(metricsAccountSoftLimit).(int))})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	var diagRet diag.Diagnostics
	readErr := retry.Do(
		func() error {
			diagRet = resourceMetricsAccountRead(ctx, d, m)
			if diagRet.HasError() {
				return fmt.Errorf("received error from read subaccount")
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
					MetricsAccountFromSchema := getCreateMetricsAccountFromSchema(d)
					return !reflect.DeepEqual(MetricsAccountFromSchema, updateMetricsAccount)
				}
			}),
		retry.DelayType(retry.BackOffDelay),
		retry.Attempts(metricsAccountRetryAttempts),
	)

	if readErr != nil {
		tflog.Error(ctx, "could not update schema")
		return diagRet
	}

	return nil
}

func resourceMetricsAccountDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id, err := utils.IdFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}
	MetricsClient, err := MetricsAccountClient(m)
	if err != nil {
		return diag.FromErr(err)
	}

	err = MetricsClient.DeleteMetricsAccount(id)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func setMetricsAccount(d *schema.ResourceData, metricsAccount *metrics_accounts.MetricsAccount) {
	d.Set(metricsAccountId, metricsAccount.Id)
	d.Set(metricsAccountName, metricsAccount.AccountName)
	d.Set(metricsAccountPlanUts, metricsAccount.PlanUts)
	d.Set(metricsAccountToken, metricsAccount.AccountToken)

	sharingObjectAccounts := make([]int32, 0)
	for _, accountId := range metricsAccount.AuthorizedAccountsIds {
		sharingObjectAccounts = append(sharingObjectAccounts, accountId)
	}

	d.Set(metricsAccountAuthorizedAccounts, sharingObjectAccounts)
}

func getCreateMetricsAccountFromSchema(d *schema.ResourceData) metrics_accounts.CreateOrUpdateMetricsAccount {
	accounts := d.Get(metricsAccountAuthorizedAccounts).([]interface{})
	// Allows users to insert empty array of authorizedAccounts, but avoiding `nil`
	authorizedAccounts := make([]int32, 0)
	for _, accountId := range accounts {
		authorizedAccounts = append(authorizedAccounts, int32(accountId.(int)))
	}

	PlanUtsVal := int32(d.Get(metricsAccountPlanUts).(int))
	var planUts *int32
	planUts = new(int32)
	*planUts = PlanUtsVal

	createMetricsAccount := metrics_accounts.CreateOrUpdateMetricsAccount{
		Email:                 d.Get(metricsAccountEmail).(string),
		AccountName:           d.Get(metricsAccountName).(string),
		PlanUts:               planUts,
		AuthorizedAccountsIds: authorizedAccounts,
	}

	return createMetricsAccount
}

func isMetricsAccountSoftLimitConfigured(d *schema.ResourceData) bool {
	value, diags := d.GetRawConfigAt(cty.GetAttrPath(metricsAccountSoftLimit))
	return !diags.HasError() && !value.IsNull()
}
