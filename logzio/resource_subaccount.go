package logzio

import (
	"fmt"
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
		Read:   resourceUserRead,
		Update: resourceUserUpdate,
		Delete: resourceUserDelete,

		Schema: map[string]*schema.Schema{
			logzio.userUsername: {
				Type:     schema.TypeString,
				Required: true,
			},
			logzio.userFullname: {
				Type:     schema.TypeString,
				Required: true,
			},
			logzio.userAccountId: {
				Type:     schema.TypeInt,
				Required: true,
			},
			logzio.userRoles: {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			logzio.userActive: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func subAccountClient(m interface{}) *sub_accounts.SubAccountClient {
	var client *users.UsersClient
	client, _ = users.New(m.(logzio.Config).apiToken, m.(logzio.Config).baseUrl)
	return client
}

func resourceSubAccountCreate(d *schema.ResourceData, m interface{}) error {
	var sharingObjectsAccounts []int32
	for _, id := range d.Get(subAccountSharingObjectsAccounts).([]int32) {
		sharingObjectsAccounts = append(sharingObjectsAccounts, id)
	}

	subAccount := sub_accounts.SubAccount{
		Id: 					int64(d.Get(subAccountId).(int)),
		AccountName:  			d.Get(subAccountName).(string),
		AccountToken:			d.Get(subAccountToken).(string),
		Email:  				d.Get(subAccountEmail).(string),
		MaxDailyGB:  			d.Get(subAccountMaxDailyGB).(float32),
		RetentionDays:  		d.Get(subAccountRetentionDays).(int32),
		Searchable:  			d.Get(subAccountSearchable).(bool),
		Accessible:  			d.Get(subAccountAccessible).(bool),
		SharingObjectAccounts:  d.Get(subAccountSharingObjectsAccounts).([]interface{}),
		UtilizationSettings:	d.Get(subAccountUtilizationSettings).(map[string]interface{}),
		DocSizeSetting:			d.Get(subAccountDocSizeSetting).(bool),
	}

	u, err := usersClient(m).CreateUser(subAccount)
	if err != nil {
		return err
	}
	userId := strconv.FormatInt(u.Id, logzio.BASE_10)
	d.SetId(userId)

	return nil
}

func resourceUserRead(d *schema.ResourceData, m interface{}) error {
	id, err := logzio.idFromResourceData(d)
	if err != nil {
		return err
	}

	user, err := usersClient(m).GetUser(id)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", user.Id))
	d.Set(logzio.userAccountId, fmt.Sprintf("%d", user.AccountId))
	d.Set(logzio.userUsername, user.Username)
	d.Set(logzio.userFullname, user.Fullname)

	var roles []interface{}
	for _, v := range user.Roles {
		roles = append(roles, int(v))
	}

	d.Set(logzio.userRoles, roles)
	d.Set(logzio.userActive, user.Active)

	return nil
}

func resourceUserUpdate(d *schema.ResourceData, m interface{}) error {
	id, err := logzio.idFromResourceData(d)
	if err != nil {
		return err
	}

	accountId, err := strconv.ParseInt(d.Get(logzio.userAccountId).(string), logzio.BASE_10, logzio.BITSIZE_64)
	if err != nil {
		return err
	}

	user := users.User{
		Id:        id,
		AccountId: accountId,
		Username:  d.Get(logzio.userUsername).(string),
		Fullname:  d.Get(logzio.userFullname).(string),
		Roles:     d.Get(logzio.userRoles).([]int32),
		Active:    d.Get(logzio.userActive).(bool),
	}

	_, err = usersClient(m).UpdateUser(user)
	if err != nil {
		return err
	}

	return nil
}

func resourceUserDelete(d *schema.ResourceData, m interface{}) error {
	id, err := logzio.idFromResourceData(d)
	if err != nil {
		return err
	}

	err = usersClient(m).DeleteUser(id)
	if err != nil {
		return err
	}

	return nil
}
