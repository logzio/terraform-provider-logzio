package logzio

import (
	"github.com/yyyogev/logzio_terraform_provider/logzio"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jonboydell/logzio_client/sub_accounts"
)

const (
	subAccountId  						string = "accountId"
	subAccountEmail						string = "email"
	subAccountName  					string = "accountName"
	subAccountToken 					string = "accountToken"
	subAccountMaxDailyGB				string = "maxDailyGB"
	subAccountRetentionDays				string = "retentionDays"
	subAccountSearchable				string = "searchable"
	subAccountAccessible				string = "accessible"
	subAccountDocSizeSetting			string = "docSizeSetting"
	subAccountSharingObjectsAccounts	string = "sharingObjectsAccounts"
	subAccountUtilizationSettings		string = "utilizationSettings"
)

/**
 * the endpoint resource schema, what terraform uses to parse and read the template
 */
func resourceSubAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceSubAccountCreate,
		Read:   resourceSubAccountRead,
		Update: resourceSubAccountUpdate,
		Delete: resourceSubAccountDelete,

		Schema: map[string]*schema.Schema{
			subAccountEmail: {
				Type:	schema.TypeString,
			},
			subAccountToken: {
				Type:	schema.TypeString,
			},
			subAccountMaxDailyGB: {
				Type:	schema.TypeFloat,
			},
			subAccountRetentionDays: {
				Type:	schema.TypeInt,
			},
			subAccountSearchable: {
				Type:	schema.TypeBool,
			},
			subAccountDocSizeSetting: {
				Type:	schema.TypeBool,
			},
			subAccountSharingObjectsAccounts: {
				Type:	schema.TypeList,
			},
			subAccountUtilizationSettings: {
				Type:	schema.TypeMap,

			},
		},
	}
}

func getSubAccountFromResource(data *schema.ResourceData, id int64) sub_accounts.SubAccount {
	return sub_accounts.SubAccount{
		Id:                    id,
		AccountName:           data.Get(subAccountName).(string),
		AccountToken:          data.Get(subAccountToken).(string),
		Email:                 data.Get(subAccountEmail).(string),
		MaxDailyGB:            data.Get(subAccountMaxDailyGB).(float32),
		RetentionDays:         data.Get(subAccountRetentionDays).(int32),
		Searchable:            data.Get(subAccountSearchable).(bool),
		Accessible:            data.Get(subAccountAccessible).(bool),
		SharingObjectAccounts: data.Get(subAccountSharingObjectsAccounts).([]interface{}),
		UtilizationSettings:   data.Get(subAccountUtilizationSettings).(map[string]interface{}),
		DocSizeSetting:        data.Get(subAccountDocSizeSetting).(bool),
	}
}

func subAccountClient(m interface{}) *sub_accounts.SubAccountClient {
	var client *sub_accounts.SubAccountClient
	client, _ = sub_accounts.New(m.(Config).apiToken, m.(Config).baseUrl)
	return client
}

func resourceSubAccountCreate(d *schema.ResourceData, m interface{}) error {
	subAccount := getSubAccountFromResource(d, int64(d.Get(subAccountId).(int)))

	u, err := subAccountClient(m).CreateSubAccount(subAccount)
	if err != nil {
		return err
	}
	subAccountId := strconv.FormatInt(u.Id, logzio.BASE_10)
	d.SetId(subAccountId)

	return nil
}

func resourceSubAccountRead(d *schema.ResourceData, m interface{}) error {
	id, err := idFromResourceData(d)
	if err != nil {
		return err
	}

	subAccount, err := subAccountClient(m).GetSubAccount(id)
	if err != nil {
		return err
	}

	setSubAccount(d, subAccount)
	return nil
}

func resourceSubAccountUpdate(d *schema.ResourceData, m interface{}) error {
	id, err := idFromResourceData(d)
	if err != nil {
		return err
	}

	subAccount := getSubAccountFromResource(d, id)

	err = subAccountClient(m).UpdateSubAccount(id, subAccount)
	if err != nil {
		return err
	}

	return nil
}

func resourceSubAccountDelete(d *schema.ResourceData, m interface{}) error {
	id, err := idFromResourceData(d)
	if err != nil {
		return err
	}

	err = subAccountClient(m).DeleteSubAccount(id)
	if err != nil {
		return err
	}

	return nil
}
