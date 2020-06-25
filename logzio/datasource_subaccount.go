package logzio

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/logzio/logzio_terraform_client/sub_accounts"
)

func dataSourceSubAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSubaccountRead,
		Schema: map[string]*schema.Schema{
			subAccountId: {
				Type:	schema.TypeInt,
				Required:	true,
			},
			subAccountEmail: {
				Type:	schema.TypeString,
				Optional:	true,
			},
			subAccountToken: {
				Type:	schema.TypeString,
				Required:	true,
			},
			subAccountMaxDailyGB: {
				Type:	schema.TypeFloat,
				Optional:	true,
			},
			subAccountRetentionDays: {
				Type:	schema.TypeInt,
				Optional:	true,
			},
			subAccountSearchable: {
				Type:	schema.TypeBool,
				Optional:	true,
			},
			subAccountDocSizeSetting: {
				Type:	schema.TypeBool,
				Optional:	true,
			},
			subAccountSharingObjectsAccounts: {
				Type:	schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Optional:	true,
			},
			subAccountUtilizationSettings: {
				Type:	schema.TypeMap,
				Optional:	true,
			},
		},
	}
}

func dataSourceSubaccountRead(d *schema.ResourceData, m interface{}) error {
	var client *sub_accounts.SubAccountClient
	client, _ = sub_accounts.New(m.(Config).apiToken, m.(Config).baseUrl)

	accountId, ok := d.GetOk(subAccountId)
	if ok {
		subAccount, err := client.GetSubAccount(accountId.(int64))
		if err != nil {
			return err
		}
		setSubAccount(d, subAccount)
		return nil
	}

	accountToken, ok := d.GetOk(subAccountToken)
	if ok {
		list, err := client.ListSubAccounts()
		if err != nil {
			return err
		}
		for i := 0; i < len(list); i++ {
			subAccount := list[i]
			if subAccount.AccountToken == accountToken {
				setSubAccount(d, &subAccount)
				return nil
			}
		}
	}

	return fmt.Errorf("couldn't find sub-account with specified attributes")
}

func setSubAccount(data *schema.ResourceData, subAccount *sub_accounts.SubAccount) {
	data.SetId(fmt.Sprintf("%d", subAccount.Id))
	data.Set(subAccountName, subAccount.AccountName)
	data.Set(subAccountEmail, subAccount.Email)
	data.Set(subAccountToken, subAccount.AccountToken)
	data.Set(subAccountDocSizeSetting, subAccount.DocSizeSetting)
	data.Set(subAccountUtilizationSettings, subAccount.UtilizationSettings)
	data.Set(subAccountAccessible, subAccount.Accessible)
	data.Set(subAccountSearchable, subAccount.Searchable)
	data.Set(subAccountRetentionDays, subAccount.RetentionDays)
	data.Set(subAccountMaxDailyGB, subAccount.MaxDailyGB)
	data.Set(subAccountSharingObjectsAccounts, subAccount.SharingObjectAccounts)
}
