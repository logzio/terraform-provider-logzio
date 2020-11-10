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
			subAccountUtilizationSettings: {
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
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
		setSubAccount(d, subAccount)
		return nil
	}

	return fmt.Errorf("couldn't find sub-account with specified attributes")
}

func setSubAccount(data *schema.ResourceData, subAccount *sub_accounts.SubAccount) {
	data.SetId(fmt.Sprintf("%d", subAccount.Id))
	data.Set(subAccountName, subAccount.AccountName)
	data.Set(subAccountDocSizeSetting, subAccount.DocSizeSetting)
	data.Set(subAccountUtilizationSettings, subAccount.UtilizationSettings)
	data.Set(subAccountAccessible, subAccount.Accessible)
	data.Set(subAccountSearchable, subAccount.Searchable)
	data.Set(subAccountRetentionDays, subAccount.RetentionDays)
	data.Set(subAccountMaxDailyGB, subAccount.MaxDailyGB)
	var sharingObjectAccounts []int32
	for _, account := range subAccount.SharingObjectAccounts {
		sharingObjectAccounts = append(sharingObjectAccounts, int32((account.(map[string]interface{}))["accountId"].(float64)))
	}
	data.Set(subAccountSharingObjectsAccounts, sharingObjectAccounts)
}
