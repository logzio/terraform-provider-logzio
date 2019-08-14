package logzio

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jonboydell/logzio_client/sub_accounts"
	"strconv"
)

const (
	accountId              string = "account_id"   //required
	email                  string = "email"       //required
	accountName            string = "account_name" //required
	maxDailyGB             string = "max_daily_gb"
	retentionDays          string = "retention_days" //required
	accessible             string = "accessible"
	searchable             string = "searchable"
	sharingObjectsAccounts string = "sharing_objects_account" //required
	docSizeSetting         string = "doc_size_setting"
	utilizationSettings    string = "utilization_settings"
	frequencyMinutes       string = "frequency_minutes"
	utilizationEnabled     string = "utilization_enabled"
	accountToken           string = "account_token"
	dailyUsagesList        string = "daily_usages_list"
)

func resourceSubAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceSubAccountCreate,
		Read:   resourceSubAccountRead,
		Update: resourceSubAccountUpdate,
		Delete: resourceSubAccountDelete,

		Schema: map[string]*schema.Schema{
			email: {
				Type:     schema.TypeString,
				Required: true,
			},
			accountName: {
				Type:     schema.TypeString,
				Required: true,
			},
			maxDailyGB: {
				Type:     schema.TypeFloat,
				Optional: true,
				Default:  0.0,
			},
			retentionDays: {
				Type:     schema.TypeInt,
				Required: true,
			},
			accessible: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			searchable: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			sharingObjectsAccounts: {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Required: true,
			},
			docSizeSetting: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			utilizationSettings: {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						frequencyMinutes: {
							Type:     schema.TypeInt,
							Required: true,
						},
						utilizationEnabled: {
							Type:     schema.TypeBool,
							Required: true,
						},
					},
				},
				Optional: true,
			},
		},
	}
}

func subaccountClient(m interface{}) *sub_accounts.SubAccountClient {
	var client *sub_accounts.SubAccountClient
	client, _ = sub_accounts.New(m.(Config).apiToken, m.(Config).baseUrl)
	return client
}

func resourceSubAccountCreate(d *schema.ResourceData, m interface{}) error {
	subAccount := sub_accounts.SubAccount{
		Email:                 d.Get(email).(string),
		AccountName:           d.Get(accountName).(string),
		MaxDailyGB:            float32(d.Get(maxDailyGB).(float64)),
		RetentionDays:         int32(d.Get(retentionDays).(int)),
		Searchable:            d.Get(searchable).(bool),
		Accessible:            d.Get(accessible).(bool),
		SharingObjectAccounts: []interface{}{},
		DocSizeSetting:        d.Get(docSizeSetting).(bool),
		UtilizationSettings: make(map[string]interface{}),
	}

	opts, _ := MappingsFromResourceData(d, utilizationSettings)
	if opts != nil {
		subAccount.UtilizationSettings[frequencyMinutes] = opts[frequencyMinutes]
		subAccount.UtilizationSettings[utilizationEnabled] = opts[utilizationEnabled]
	} else {
		subAccount.UtilizationSettings = map[string]interface{}{"frequencyMinutes":nil, "utilizationEnabled":false}
	}

	client := subaccountClient(m)
	s, err := client.CreateSubAccount(subAccount)
	if err != nil {
		return err
	}

	subAccountId := strconv.FormatInt(s.Id, BASE_10)
	d.SetId(subAccountId)

	return nil
}

func resourceSubAccountRead(d *schema.ResourceData, m interface{}) error {
	client := subaccountClient(m)
	subAccountId, _ := IdFromResourceData(d)

	var subAccount *sub_accounts.SubAccount
	subAccount, err := client.GetSubAccount(subAccountId)
	if err != nil {
		return err
	}

	if len(subAccount.Email) > 0 {
		d.Set(email, subAccount.Email)
	}

	d.Set(accountName, subAccount.AccountName)
	d.Set(maxDailyGB, subAccount.MaxDailyGB)
	d.Set(retentionDays, subAccount.RetentionDays)
	d.Set(searchable, subAccount.Searchable)
	d.Set(accessible, subAccount.Accessible)
	d.Set(docSizeSetting, subAccount.DocSizeSetting)

	set := make([]map[string]interface{}, 1)
	set[0] = map[string]interface{}{
		frequencyMinutes: subAccount.UtilizationSettings["frequencyMinutes"],
		utilizationEnabled: subAccount.UtilizationSettings["utilizationEnabled"],
	}
	d.Set(utilizationSettings, set)

	return nil
}

func resourceSubAccountUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceSubAccountDelete(d *schema.ResourceData, m interface{}) error {
	subAccountId, _ := IdFromResourceData(d)
	client := subaccountClient(m)
	err := client.DeleteSubAccount(subAccountId)
	if err != nil {
		return err
	}

	return nil
}
