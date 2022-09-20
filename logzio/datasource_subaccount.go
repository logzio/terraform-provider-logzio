package logzio

import (
	"context"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/logzio/logzio_terraform_client/sub_accounts"
	"strconv"
	"strings"
	"time"
)

func dataSourceSubAccount() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSubaccountReadWrapper,
		Schema: map[string]*schema.Schema{
			subAccountId: {
				Type:     schema.TypeInt,
				Optional: true,
			},
			subAccountName: {
				Type:     schema.TypeString,
				Optional: true,
			},
			subAccountEmail: {
				Type:     schema.TypeString,
				Optional: true,
			},
			subAccountToken: {
				Type:     schema.TypeString,
				Optional: true,
			},
			subAccountFlexible: {
				Type:     schema.TypeBool,
				Optional: true,
			},
			subAccountReservedDailyGb: {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			subAccountMaxDailyGB: {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			subAccountRetentionDays: {
				Type:     schema.TypeInt,
				Optional: true,
			},
			subAccountSearchable: {
				Type:     schema.TypeBool,
				Optional: true,
			},
			subAccountAccessible: {
				Type:     schema.TypeBool,
				Optional: true,
			},
			subAccountDocSizeSetting: {
				Type:     schema.TypeBool,
				Optional: true,
			},
			subAccountSharingObjectsAccounts: {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Optional: true,
			},
			subAccountUtilizationSettingsFrequencyMinutes: {
				Type:     schema.TypeInt,
				Optional: true,
			},
			subAccountUtilizationSettingsUtilizationEnabled: {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Second),
		},
	}
}

func dataSourceSubaccountReadWrapper(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var err error
	readErr := retry.Do(
		func() error {
			err = dataSourceSubaccountRead(d, m)
			if err != nil {
				return err
			}

			return nil
		},
		retry.RetryIf(
			func(err error) bool {
				if err != nil {
					if strings.Contains(err.Error(), "failed with missing sub account") ||
						strings.Contains(err.Error(), "failed with status code 500") {
						return true
					}
				}
				return false
			}),
		retry.DelayType(retry.BackOffDelay),
		retry.Attempts(15),
	)

	if readErr != nil {
		return diag.FromErr(err)
	}

	return nil
}

func dataSourceSubaccountRead(d *schema.ResourceData, m interface{}) error {
	var client *sub_accounts.SubAccountClient
	client, _ = sub_accounts.New(m.(Config).apiToken, m.(Config).baseUrl)

	accountId, ok := d.GetOk(subAccountId)
	if ok {
		subAccount, err := client.GetSubAccount(int64(accountId.(int)))
		if err != nil {
			return err
		}
		d.SetId(strconv.FormatInt(int64(accountId.(int)), 10))
		setSubAccount(d, subAccount)
		err = setTokenAndId(d, m, int64(accountId.(int)))
		if err != nil {
			return err
		}
		return nil
	}

	accountName, ok := d.GetOk(subAccountName)
	if ok {
		subAccounts, err := client.ListSubAccounts()
		if err != nil {
			return err
		}

		for _, account := range subAccounts {
			if account.AccountName == accountName.(string) {
				d.SetId(strconv.FormatInt(int64(account.AccountId), 10))
				setSubAccount(d, &account)
				err = setTokenAndId(d, m, int64(account.AccountId))
				if err != nil {
					return err
				}
				return nil
			}
		}
	}

	return fmt.Errorf("couldn't find sub-account with specified attributes")
}
