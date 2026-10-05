package logzio

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/avast/retry-go"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/logzio/logzio_terraform_client/sub_accounts"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
)

const (
	warmTierAccountId               string = "account_id"
	warmTierSnapSearchRetentionDays string = "snap_search_retention_days"
)

// The warm tier resource manages the warm retention of one existing account: the main account of the API token or
// one of its sub accounts. It doesn't create or delete accounts.
func resourceWarmTier() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWarmTierCreate,
		ReadContext:   resourceWarmTierRead,
		UpdateContext: resourceWarmTierUpdate,
		DeleteContext: resourceWarmTierDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			warmTierAccountId: {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			// The API enforces the rest: the main account must already have warm tier and keep at least 1 day, and a
			// sub account can be set from 0 (no warm tier) up to the main account's value.
			warmTierSnapSearchRetentionDays: {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},
		},
	}
}

func resourceWarmTierCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	accountId := int64(d.Get(warmTierAccountId).(int))
	d.SetId(strconv.FormatInt(accountId, 10))
	return updateWarmTier(ctx, d, m, accountId)
}

func resourceWarmTierRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id, err := utils.IdFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	account, err := subAccountClient(m).GetSubAccount(id)
	if err != nil {
		tflog.Error(ctx, err.Error())
		if strings.Contains(err.Error(), "missing sub account") {
			// If we were not able to find the account - delete from state
			d.SetId("")
			return diag.Diagnostics{}
		}
		return diag.FromErr(err)
	}

	d.Set(warmTierAccountId, account.AccountId)
	d.Set(warmTierSnapSearchRetentionDays, account.SnapSearchRetentionDays)
	return nil
}

func resourceWarmTierUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id, err := utils.IdFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}
	return updateWarmTier(ctx, d, m, id)
}

// Destroy only removes the resource from the state. The main account's warm tier can't be turned off through the API,
// and setting a sub account to 0 deletes its warm data, so the account keeps the warm retention it has.
func resourceWarmTierDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Warn(ctx, fmt.Sprintf("removing warm tier of account %s from the state only, its warm retention is unchanged", d.Id()))
	d.SetId("")
	return nil
}

// updateWarmTier sets the state from the update response, which holds the value the API stored. It then waits until
// reading the account returns that value, so the refresh after apply shows no diff.
func updateWarmTier(ctx context.Context, d *schema.ResourceData, m interface{}, accountId int64) diag.Diagnostics {
	client := subAccountClient(m)
	retentionDetails, err := client.UpdateWarmRetention(accountId, sub_accounts.UpdateWarmRetention{
		SnapSearchRetentionDays: int32(d.Get(warmTierSnapSearchRetentionDays).(int)),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	var updated *sub_accounts.AccountRetentionDetails
	for i := range retentionDetails {
		if int64(retentionDetails[i].AccountId) == accountId {
			updated = &retentionDetails[i]
		}
	}
	if updated == nil {
		return diag.Errorf("account %d is missing from the warm retention update response", accountId)
	}

	var snapSearchRetentionDays int32
	if updated.SnapSearchRetentionDays != nil {
		snapSearchRetentionDays = *updated.SnapSearchRetentionDays
	}
	d.Set(warmTierSnapSearchRetentionDays, snapSearchRetentionDays)

	readErr := retry.Do(
		func() error {
			account, err := client.GetSubAccount(accountId)
			if err != nil {
				return err
			}
			if account.SnapSearchRetentionDays != snapSearchRetentionDays {
				return fmt.Errorf("account %d still reads warm retention %d", accountId, account.SnapSearchRetentionDays)
			}
			return nil
		},
		retry.DelayType(retry.BackOffDelay),
		retry.Attempts(subAccountRetryAttempts),
	)
	if readErr != nil {
		tflog.Warn(ctx, fmt.Sprintf("warm retention of account %d was updated, but reading it back failed: %s", accountId, readErr.Error()))
	}

	return nil
}
