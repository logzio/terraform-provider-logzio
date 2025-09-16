package logzio

import (
	"context"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/logzio/logzio_terraform_client/sub_accounts"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
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
	subAccountUtilizationSettingsFrequencyMinutes   string = "frequency_minutes"
	subAccountUtilizationSettingsUtilizationEnabled string = "utilization_enabled"
	subAccountsSnapSearchRetentionDays              string = "snap_search_retention_days"
	subAccountsIsCapped                             string = "is_capped"
	subAccountsSharedGb                             string = "shared_gb"
	subAccountsTotalTimeBasedDailyGb                string = "total_time_based_daily_gb"
	subAccountIsOwner                               string = "is_owner"
	subAccountSoftLimitGB                           string = "soft_limit_gb"

	delayGetSubAccount      = 2 * time.Second
	subAccountRetryAttempts = 8
)

// The endpoint resource schema, what terraform uses to parse and read the template
func resourceSubAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSubAccountCreate,
		ReadContext:   resourceSubAccountRead,
		UpdateContext: resourceSubAccountUpdate,
		DeleteContext: resourceSubAccountDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			subAccountId: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			subAccountToken: {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			subAccountEmail: {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
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
				Default:  true,
			},
			subAccountUtilizationSettingsFrequencyMinutes: {
				Type:     schema.TypeInt,
				Optional: true,
			},
			subAccountUtilizationSettingsUtilizationEnabled: {
				Type:     schema.TypeBool,
				Optional: true,
			},
			subAccountsSnapSearchRetentionDays: {
				Type:     schema.TypeInt,
				Optional: true,
			},
			subAccountsIsCapped: {
				Type:     schema.TypeBool,
				Computed: true,
			},
			subAccountsSharedGb: {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			subAccountsTotalTimeBasedDailyGb: {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			subAccountIsOwner: {
				Type:     schema.TypeBool,
				Computed: true,
			},
			subAccountSoftLimitGB: {
				Type:     schema.TypeFloat,
				Optional: true,
			},
		},
	}
}

func subAccountClient(m interface{}) *sub_accounts.SubAccountClient {
	var client *sub_accounts.SubAccountClient
	client, _ = sub_accounts.New(m.(Config).apiToken, m.(Config).baseUrl)
	return client
}

func resourceSubAccountCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	createSubAccount := getCreateSubAccountFromSchema(d)
	subAccount, err := subAccountClient(m).CreateSubAccount(createSubAccount)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(int64(subAccount.AccountId), 10))
	d.Set(subAccountToken, subAccount.AccountToken)
	d.Set(subAccountId, subAccount.AccountId)
	return resourceSubAccountRead(ctx, d, m)
}

func resourceSubAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id, err := utils.IdFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}
	subAccount, err := subAccountClient(m).GetSubAccount(id)
	if err != nil {
		tflog.Error(ctx, err.Error())
		if strings.Contains(err.Error(), "missing sub account") {
			// If we were not able to find the resource - delete from state
			d.SetId("")
			return diag.Diagnostics{}
		} else {
			return diag.FromErr(err)
		}

	}

	setSubAccount(d, subAccount)
	// Sub accounts created before v1.2.4 had no account_id, account_token attributes.
	// These lines add those attributes to already existing resources on Read
	err = setTokenAndId(d, m, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSubAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id, err := utils.IdFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	updateSubAccount := getCreateSubAccountFromSchema(d)
	err = subAccountClient(m).UpdateSubAccount(id, updateSubAccount)
	if err != nil {
		return diag.FromErr(err)
	}

	var diagRet diag.Diagnostics
	readErr := retry.Do(
		func() error {
			diagRet = resourceSubAccountRead(ctx, d, m)
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
					subAccountFromSchema := getCreateSubAccountFromSchema(d)
					return !reflect.DeepEqual(subAccountFromSchema, updateSubAccount)
				}
			}),
		retry.DelayType(retry.BackOffDelay),
		retry.Attempts(subAccountRetryAttempts),
	)

	if readErr != nil {
		tflog.Error(ctx, "could not update schema")
		return diagRet
	}

	return nil
}

func resourceSubAccountDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id, err := utils.IdFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	err = subAccountClient(m).DeleteSubAccount(id)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
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
	d.Set(subAccountsSnapSearchRetentionDays, subAccount.SnapSearchRetentionDays)
	d.Set(subAccountsIsCapped, subAccount.IsCapped)
	d.Set(subAccountsSharedGb, subAccount.SharedGB)
	d.Set(subAccountsTotalTimeBasedDailyGb, subAccount.TotalTimeBasedDailyGB)
	d.Set(subAccountIsOwner, subAccount.IsOwner)
	d.Set(subAccountSoftLimitGB, subAccount.SoftLimitGB)

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

	flexible := d.Get(subAccountFlexible).(bool)
	maxDailyGbVal := float32(d.Get(subAccountMaxDailyGB).(float64))
	reservedDailyGbVal := float32(d.Get(subAccountReservedDailyGb).(float64))
	softLimitGBVal := float32(d.Get(subAccountSoftLimitGB).(float64))
	var maxDailyGb, reservedDailyGb, softLimitGB *float32
	if !flexible {
		maxDailyGb = new(float32)
		*maxDailyGb = maxDailyGbVal
		softLimitGB = new(float32)
		*softLimitGB = softLimitGBVal
	}

	if flexible {
		reservedDailyGb = new(float32)
		*reservedDailyGb = reservedDailyGbVal
		if maxDailyGbVal > 0 {
			maxDailyGb = new(float32)
			*maxDailyGb = maxDailyGbVal
		}
	}

	createSubAccount := sub_accounts.CreateOrUpdateSubAccount{
		Email:                  d.Get(subAccountEmail).(string),
		AccountName:            d.Get(subAccountName).(string),
		Flexible:               strconv.FormatBool(d.Get(subAccountFlexible).(bool)),
		ReservedDailyGB:        reservedDailyGb,
		MaxDailyGB:             maxDailyGb,
		RetentionDays:          int32(d.Get(subAccountRetentionDays).(int)),
		Searchable:             strconv.FormatBool(d.Get(subAccountSearchable).(bool)),
		Accessible:             strconv.FormatBool(d.Get(subAccountAccessible).(bool)),
		SharingObjectsAccounts: sharingObjectAccounts,
		DocSizeSetting:         strconv.FormatBool(d.Get(subAccountDocSizeSetting).(bool)),
		UtilizationSettings: sub_accounts.AccountUtilizationSettingsCreateOrUpdate{
			FrequencyMinutes:   int32(d.Get(subAccountUtilizationSettingsFrequencyMinutes).(int)),
			UtilizationEnabled: strconv.FormatBool(d.Get(subAccountUtilizationSettingsUtilizationEnabled).(bool)),
		},
		SnapSearchRetentionDays: getOptionalInt32Pointer(d, subAccountsSnapSearchRetentionDays),
		SoftLimitGB:             softLimitGB,
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

// getOptionalInt32Pointer returns a pointer to the numeric value from the config, or nil if it was not set.
// We don't use d.Get since it returns 0 for nil, when the field is unset.
// And we don't use d.GetOk because it returns false for 0, even if it's explicitly set.
func getOptionalInt32Pointer(d *schema.ResourceData, key string) *int32 {
	val, diags := d.GetRawConfigAt(cty.GetAttrPath(key))
	if diags.HasError() || val.IsNull() {
		return nil
	} else {
		int64Val, _ := val.AsBigFloat().Int64()
		int32Val := int32(int64Val)
		return &int32Val
	}
}
