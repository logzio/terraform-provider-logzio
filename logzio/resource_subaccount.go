package logzio

import (
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jonboydell/logzio_client/sub_accounts"
)

const (
	subAccountId  						string = "account_id"
	subAccountEmail						string = "email"
	subAccountName  					string = "account_name"
	subAccountToken 					string = "account_token"
	subAccountMaxDailyGB				string = "max_daily_gb"
	subAccountRetentionDays				string = "retention_days"
	subAccountSearchable				string = "searchable"
	subAccountAccessible				string = "accessible"
	subAccountDocSizeSetting			string = "doc_size_setting"
	subAccountSharingObjectsAccounts	string = "sharing_objects_accounts"
	subAccountUtilizationSettings		string = "utilization_settings"
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
				Required:	true,
			},
			subAccountName: {
				Type:	schema.TypeString,
				Required:	true,
			},
			subAccountMaxDailyGB: {
				Type:	schema.TypeFloat,
				Optional:	true,
			},
			subAccountRetentionDays: {
				Type:	schema.TypeInt,
				Required:	true,
			},
			subAccountSearchable: {
				Type:	schema.TypeBool,
				Optional:	true,
			},
			subAccountAccessible: {
				Type:	schema.TypeBool,
				Optional:	true,
			},
			subAccountSharingObjectsAccounts: {
				Type:	schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Required:	true,
			},
			subAccountDocSizeSetting: {
				Type:	schema.TypeBool,
				Optional:	true,

			},
			subAccountUtilizationSettings: {
				Type:	schema.TypeMap,
				Optional:	true,
			},
		},
	}
}

func getMinimalSubAccountFromResource(data *schema.ResourceData) sub_accounts.SubAccount {
	return sub_accounts.SubAccount{
		AccountName:           data.Get(subAccountName).(string),
		Email:                 data.Get(subAccountEmail).(string),
		RetentionDays:         int32(data.Get(subAccountRetentionDays).(int)),
		SharingObjectAccounts: data.Get(subAccountSharingObjectsAccounts).([]interface{}),

	}
}

func subAccountClient(m interface{}) *sub_accounts.SubAccountClient {
	var client *sub_accounts.SubAccountClient
	client, _ = sub_accounts.New(m.(Config).apiToken, m.(Config).baseUrl)
	return client
}

func resourceSubAccountCreate(d *schema.ResourceData, m interface{}) error {
	subAccount := getMinimalSubAccountFromResource(d)

	u, err := subAccountClient(m).CreateSubAccount(subAccount)
	if err != nil {
		return err
	}
	subAccountId := strconv.FormatInt(u.Id, BASE_10)
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

	subAccount := getMinimalSubAccountFromResource(d)
	subAccount.Id = id

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
