package logzio

import (
	"context"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/logzio/logzio_terraform_client/tracing_accounts"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
)

const (
	tracingAccountId                 string = "account_id"
	tracingAccountEmail              string = "email"
	tracingAccountName               string = "account_name"
	tracingAccountToken              string = "account_token"
	tracingAccountMaxDailyGB         string = "max_daily_gb"
	tracingAccountRetention          string = "retention"
	tracingAccountAuthorizedAccounts string = "authorized_accounts"
)

// The tracing account resource schema, what terraform uses to parse and read the template
func resourceTracingAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTracingAccountCreate,
		ReadContext:   resourceTracingAccountRead,
		UpdateContext: resourceTracingAccountUpdate,
		DeleteContext: resourceTracingAccountDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			tracingAccountId: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			tracingAccountToken: {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			tracingAccountEmail: {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			tracingAccountName: {
				Type:     schema.TypeString,
				Required: true,
			},
			// max_daily_gb IS the tracing account's soft cap: the API reads the soft limit
			// from maxDailyGB and writes maxDailyGB from it, so there is deliberately no
			// separate soft_limit_gb argument here - two fields would write one value.
			tracingAccountMaxDailyGB: {
				Type:     schema.TypeFloat,
				Optional: true,
				Computed: true,
			},
			tracingAccountRetention: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			tracingAccountAuthorizedAccounts: {
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

func tracingAccountClient(m interface{}) (*tracing_accounts.TracingAccountClient, error) {
	client, clientError := tracing_accounts.New(m.(Config).apiToken, m.(Config).baseUrl)
	if clientError != nil {
		return nil, clientError
	}
	return client, nil
}

func resourceTracingAccountCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := tracingAccountClient(m)
	if err != nil {
		return diag.FromErr(err)
	}

	tracingAccount, err := client.CreateTracingAccount(getCreateTracingAccountFromSchema(d))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(int64(tracingAccount.AccountId), 10))
	d.Set(tracingAccountId, tracingAccount.AccountId)
	d.Set(tracingAccountToken, tracingAccount.Token)

	return resourceTracingAccountRead(ctx, d, m)
}

func resourceTracingAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id, err := utils.IdFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}
	client, err := tracingAccountClient(m)
	if err != nil {
		return diag.FromErr(err)
	}

	tracingAccount, err := client.GetTracingAccount(id)
	if err != nil {
		tflog.Error(ctx, err.Error())
		if strings.Contains(err.Error(), "missing tracing account") {
			// If we were not able to find the resource - delete from state
			d.SetId("")
			return diag.Diagnostics{}
		}
		return diag.FromErr(err)
	}

	setTracingAccount(d, tracingAccount)

	return nil
}

func resourceTracingAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id, err := utils.IdFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}
	client, err := tracingAccountClient(m)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateTracingAccount(id, getCreateTracingAccountFromSchema(d))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceTracingAccountRead(ctx, d, m)
}

func resourceTracingAccountDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id, err := utils.IdFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}
	client, err := tracingAccountClient(m)
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteTracingAccount(id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func setTracingAccount(d *schema.ResourceData, tracingAccount *tracing_accounts.TracingAccount) {
	d.Set(tracingAccountId, tracingAccount.AccountId)
	d.Set(tracingAccountName, tracingAccount.AccountName)
	d.Set(tracingAccountRetention, tracingAccount.Retention)
	if tracingAccount.Token != "" {
		d.Set(tracingAccountToken, tracingAccount.Token)
	}
	if tracingAccount.MaxDailyGB != nil {
		d.Set(tracingAccountMaxDailyGB, *tracingAccount.MaxDailyGB)
	}

	authorizedAccounts := make([]int32, 0)
	for _, account := range tracingAccount.AuthorizedAccounts {
		authorizedAccounts = append(authorizedAccounts, account.AccountId)
	}
	d.Set(tracingAccountAuthorizedAccounts, authorizedAccounts)
}

func getCreateTracingAccountFromSchema(d *schema.ResourceData) tracing_accounts.CreateOrUpdateTracingAccount {
	accounts := d.Get(tracingAccountAuthorizedAccounts).([]interface{})
	// Allows users to insert an empty array of authorizedAccounts, but avoiding `nil`
	authorizedAccounts := make([]int32, 0)
	for _, accountId := range accounts {
		authorizedAccounts = append(authorizedAccounts, int32(accountId.(int)))
	}

	return tracing_accounts.CreateOrUpdateTracingAccount{
		AccountName:          d.Get(tracingAccountName).(string),
		AuthorizedAccountIds: authorizedAccounts,
		MaxDailyGB:           float32(d.Get(tracingAccountMaxDailyGB).(float64)),
		Email:                d.Get(tracingAccountEmail).(string),
	}
}
