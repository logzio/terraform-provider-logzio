package logzio

import (
	"context"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/logzio/logzio_terraform_client/metrics_accounts"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	metricsAccountId                 string = "account_id"
	metricsAccountEmail              string = "email"
	metricsAccountName               string = "account_name"
	metricsAccountToken              string = "account_token"
	metricsAccountPlanUts            string = "plan_uts"
	metricsAccountAuthorizedAccounts string = "authorized_accounts"

	delayGetMetricsAccount      = 2 * time.Second
	metricsAccountRetryAttempts = 8
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
				Required: true,
			},
			metricsAccountAuthorizedAccounts: {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Optional: true,
				Computed: true,
			},
		},
	}
}

func MetricsAccountClient(m interface{}) *metrics_accounts.MetricsAccountClient {
	var client *metrics_accounts.MetricsAccountClient
	client, _ = metrics_accounts.New(m.(Config).apiToken, m.(Config).baseUrl)
	return client
}

func resourceMetricsAccountCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	createSubAccount := getCreateMetricsAccountFromSchema(d)
	subAccount, err := MetricsAccountClient(m).CreateMetricsAccount(createSubAccount)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(int64(subAccount.Id), 10))
	d.Set(metricsAccountToken, subAccount.Token)
	d.Set(metricsAccountId, subAccount.Id)
	return resourceMetricsAccountRead(ctx, d, m)
}

func resourceMetricsAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id, err := utils.IdFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}
	subAccount, err := MetricsAccountClient(m).GetMetricsAccount(id)
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

	setMetricsAccount(d, subAccount)

	return nil
}

func resourceMetricsAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id, err := utils.IdFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	updateSubAccount := getCreateMetricsAccountFromSchema(d)
	err = MetricsAccountClient(m).UpdateMetricsAccount(id, updateSubAccount)
	if err != nil {
		return diag.FromErr(err)
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
					subAccountFromSchema := getCreateMetricsAccountFromSchema(d)
					return !reflect.DeepEqual(subAccountFromSchema, updateSubAccount)
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

	err = MetricsAccountClient(m).DeleteMetricsAccount(id)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func setMetricsAccount(d *schema.ResourceData, subAccount *metrics_accounts.MetricsAccount) {
	d.Set(metricsAccountId, subAccount.Id)
	d.Set(metricsAccountName, subAccount.AccountName)
	d.Set(metricsAccountPlanUts, subAccount.PlanUts)
	d.Set(metricsAccountToken, subAccount.Token)

	sharingObjectAccounts := make([]int32, 0)
	for _, accountId := range subAccount.AuthorizedAccountsIds {
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

	createSubAccount := metrics_accounts.CreateOrUpdateMetricsAccount{
		Email:                 d.Get(metricsAccountEmail).(string),
		AccountName:           d.Get(metricsAccountName).(string),
		PlanUts:               int32(d.Get(metricsAccountPlanUts).(int)),
		AuthorizedAccountsIds: authorizedAccounts,
	}

	return createSubAccount
}
