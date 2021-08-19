package logzio

import (
	"fmt"
	"github.com/avast/retry-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/logzio/logzio_terraform_client/sub_accounts"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	subAccountId                                    string = "account_id"
	subAccountEmail                                 string = "email"
	subAccountName                                  string = "account_name"
	subAccountToken                                 string = "account_token"
	subAccountFlexible                              string = "flexible"
	subAccountReservedDailyGb                       string = "reserved_daily_gb"
	subAccountMaxDailyGB                            string = "max_daily_gb"
	subAccountRetentionDays                         string = "retention_days"
	subAccountSearchable                            string = "searchable"
	subAccountAccessible                            string = "accessible"
	subAccountDocSizeSetting                        string = "doc_size_setting"
	subAccountSharingObjectsAccounts                string = "sharing_objects_accounts"
	subAccountUtilizationSettings                   string = "utilization_settings"
	subAccountUtilizationSettingsFrequencyMinutes   string = "frequency_minutes"
	subAccountUtilizationSettingsUtilizationEnabled string = "utilization_enabled"

	delayGetSubAccount = 2 * time.Second
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
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			subAccountId: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			subAccountToken: {
				Type:     schema.TypeString,
				Computed: true,
			},
			subAccountEmail: {
				Type:     schema.TypeString,
				Required: true,
			},
			subAccountName: {
				Type:     schema.TypeString,
				Required: true,
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
				Required: true,
			},
			subAccountSearchable: {
				Type:     schema.TypeBool,
				Optional: true,
			},
			subAccountAccessible: {
				Type:     schema.TypeBool,
				Optional: true,
			},
			subAccountSharingObjectsAccounts: {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Optional: true,
				Computed: true,
			},
			subAccountDocSizeSetting: {
				Type:     schema.TypeBool,
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
			subAccountUtilizationSettings: {
				Type:     schema.TypeMap,
				Optional: true,
				Deprecated: fmt.Sprintf(
					"this attribute is deprecated, please use attributes %s and %s instead",
					subAccountUtilizationSettingsFrequencyMinutes,
					subAccountUtilizationSettingsUtilizationEnabled),
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Second),
			Update: schema.DefaultTimeout(5 * time.Second),
			Delete: schema.DefaultTimeout(5 * time.Second),
		},
	}
}

func subAccountClient(m interface{}) *sub_accounts.SubAccountClient {
	var client *sub_accounts.SubAccountClient
	client, _ = sub_accounts.New(m.(Config).apiToken, m.(Config).baseUrl)
	return client
}

func resourceSubAccountCreate(d *schema.ResourceData, m interface{}) error {
	createSubAccount := getCreateSubAccountFromSchema(d)
	subAccount, err := subAccountClient(m).CreateSubAccount(createSubAccount)
	if err != nil {
		return err
	}

	d.SetId(strconv.FormatInt(int64(subAccount.AccountId), 10))
	d.Set(subAccountToken, subAccount.AccountToken)
	d.Set(subAccountId, subAccount.AccountId)

	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		err = resourceSubAccountRead(d, m)
		if err != nil {
			if strings.Contains(err.Error(), "failed with missing sub account") {
				return resource.RetryableError(err)
			}

			if strings.Contains(err.Error(), "failed with status code 500") {
				return resource.RetryableError(err)
			}
		}

		return resource.NonRetryableError(err)
	})

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
	// Sub accounts created before v1.2.4 had no account_id, account_token attributes.
	// These lines add those attributes to already existing resources on Read
	err = setTokenAndId(d, m, id)
	if err != nil {
		return err
	}

	return nil
}

func resourceSubAccountUpdate(d *schema.ResourceData, m interface{}) error {
	id, err := idFromResourceData(d)
	if err != nil {
		return err
	}

	updateSubAccount := getCreateSubAccountFromSchema(d)
	err = subAccountClient(m).UpdateSubAccount(id, updateSubAccount)
	if err != nil {
		return err
	}

	return resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		err = resourceSubAccountRead(d, m)
		if err != nil {
			if strings.Contains(err.Error(), "failed with status code 500") {
				return resource.RetryableError(err)
			}

			subAccountFromSchema := getCreateSubAccountFromSchema(d)
			if strings.Contains(err.Error(), "failed with missing sub account") &&
				!reflect.DeepEqual(subAccountFromSchema, updateSubAccount) {
				return resource.RetryableError(fmt.Errorf("sub account is not updated yet: %s", err.Error()))
			}
		}

		return resource.NonRetryableError(err)
	})
}

func resourceSubAccountDelete(d *schema.ResourceData, m interface{}) error {
	id, err := idFromResourceData(d)
	if err != nil {
		return err
	}

	return resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		err = subAccountClient(m).DeleteSubAccount(id)
		if err != nil {
			return resource.RetryableError(err)
		}

		return resource.NonRetryableError(err)
	})
}

func setSubAccount(d *schema.ResourceData, subAccount *sub_accounts.SubAccount) {
	d.Set(subAccountId, subAccount.AccountId)
	d.Set(subAccountName, subAccount.AccountName)
	d.Set(subAccountFlexible, subAccount.Flexible)
	d.Set(subAccountReservedDailyGb, subAccount.ReservedDailyGB)
	d.Set(subAccountMaxDailyGB, subAccount.MaxDailyGB)
	d.Set(subAccountRetentionDays, subAccount.RetentionDays)
	d.Set(subAccountSearchable, subAccount.Searchable)
	d.Set(subAccountAccessible, subAccount.Accessible)
	d.Set(subAccountDocSizeSetting, subAccount.DocSizeSetting)
	d.Set(subAccountUtilizationSettingsFrequencyMinutes, subAccount.UtilizationSettings.FrequencyMinutes)
	d.Set(subAccountUtilizationSettingsUtilizationEnabled, subAccount.UtilizationSettings.UtilizationEnabled)

	sharingObjectAccounts := make([]int32, 0)
	for _, account := range subAccount.SharingObjectsAccounts {
		sharingObjectAccounts = append(sharingObjectAccounts, account.AccountId)
	}

	d.Set(subAccountSharingObjectsAccounts, sharingObjectAccounts)
}

func setTokenAndId(d *schema.ResourceData, m interface{}, id int64) error {
	accountToken, okToken := d.GetOk(subAccountToken)
	accountId, okId := d.GetOk(subAccountId)

	if !okToken || !okId || accountId.(int) == 0 || len(accountToken.(string)) == 0 {
		err := insertAccountTokenAndId(d, m, id)
		if err != nil {
			return err
		}
	}

	return nil
}

func getCreateSubAccountFromSchema(d *schema.ResourceData) sub_accounts.CreateOrUpdateSubAccount {
	sharingAccounts := d.Get(subAccountSharingObjectsAccounts).([]interface{})
	// Allows users to insert empty array of sharingObjectAccounts, but avoiding `nil`
	sharingObjectAccounts := make([]int32, 0)
	for _, accountId := range sharingAccounts {
		sharingObjectAccounts = append(sharingObjectAccounts, int32(accountId.(int)))
	}

	createSubAccount := sub_accounts.CreateOrUpdateSubAccount{
		Email:                  d.Get(subAccountEmail).(string),
		AccountName:            d.Get(subAccountName).(string),
		Flexible:               strconv.FormatBool(d.Get(subAccountFlexible).(bool)),
		ReservedDailyGB:        float32(d.Get(subAccountReservedDailyGb).(float64)),
		MaxDailyGB:             float32(d.Get(subAccountMaxDailyGB).(float64)),
		RetentionDays:          int32(d.Get(subAccountRetentionDays).(int)),
		Searchable:             strconv.FormatBool(d.Get(subAccountSearchable).(bool)),
		Accessible:             strconv.FormatBool(d.Get(subAccountAccessible).(bool)),
		SharingObjectsAccounts: sharingObjectAccounts,
		DocSizeSetting:         strconv.FormatBool(d.Get(subAccountDocSizeSetting).(bool)),
		UtilizationSettings: sub_accounts.AccountUtilizationSettingsCreateOrUpdate{
			FrequencyMinutes:   int32(d.Get(subAccountUtilizationSettingsFrequencyMinutes).(int)),
			UtilizationEnabled: strconv.FormatBool(d.Get(subAccountUtilizationSettingsUtilizationEnabled).(bool)),
		},
	}

	return createSubAccount
}

func getDetailedSubAccount(m interface{}, id int64) (*sub_accounts.DetailedSubAccount, error) {
	subAccount, err := subAccountClient(m).GetDetailedSubAccount(id)
	if err != nil {
		return nil, err
	}

	return subAccount, nil
}

func insertAccountTokenAndId(d *schema.ResourceData, m interface{}, id int64) error {
	return retry.Do(
		func() error {
			detailed, err := getDetailedSubAccount(m, id)
			if err != nil {
				return err
			}

			d.Set(subAccountId, detailed.Account.AccountId)
			d.Set(subAccountToken, detailed.Account.AccountToken)

			return nil
		},
		retry.RetryIf(
			func(err error) bool {
				if err != nil {
					match, _ := regexp.MatchString("^404.*errorCode", err.Error())
					if match {
						return true
					}
				}
				return false
			}),
		retry.Delay(delayGetSubAccount),
	)
}
