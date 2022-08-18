package logzio

import (
	"context"
	"github.com/avast/retry-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/logzio/logzio_terraform_client/users"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
	"strconv"
	"strings"
)

const (
	userId        string = "id"
	userUsername  string = "username"
	userFullName  string = "full_name"
	userAccountId string = "account_id"
	userRole      string = "role"
	userActive    string = "active"
)

/**
 * the endpoint resource schema, what terraform uses to parse and read the template
 */
func resourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			userUsername: {
				Type:     schema.TypeString,
				Required: true,
			},
			userFullName: {
				Type:     schema.TypeString,
				Required: true,
			},
			userAccountId: {
				Type:     schema.TypeInt,
				Required: true,
			},
			userRole: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: utils.ValidateUserRoleUser,
			},
			userActive: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func usersClient(m interface{}) *users.UsersClient {
	var client *users.UsersClient
	client, _ = users.New(m.(Config).apiToken, m.(Config).baseUrl)
	return client
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	createUser := getCreateOrUpdateUserFromSchema(d)
	user, err := usersClient(m).CreateUser(createUser)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(int64(user.Id), 10))
	return resourceUserRead(ctx, d, m)
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id, err := utils.IdFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	var user *users.User
	readErr := retry.Do(
		func() error {
			user, err = usersClient(m).GetUser(int32(id))
			if err != nil {
				return err
			}

			return nil
		},
		retry.RetryIf(
			func(err error) bool {
				if err != nil {
					if strings.Contains(err.Error(), "failed with missing user") {
						return true
					}
				}
				return false
			}),
		retry.DelayType(retry.BackOffDelay),
		retry.Attempts(15),
	)

	if readErr != nil {
		// If we were not able to find the resource - delete from state
		d.SetId("")
		return diag.FromErr(err)
	}

	setUser(d, user)
	return nil
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id, err := utils.IdFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	updateUser := getCreateOrUpdateUserFromSchema(d)

	//accountId, err := strconv.ParseInt(d.Get(userAccountId).(string), utils.BASE_10, utils.BITSIZE_64)
	//if err != nil {
	//	return err
	//}
	//
	//user := users.User{
	//	Id:        id,
	//	AccountId: accountId,
	//	Username:  d.Get(userUsername).(string),
	//	Fullname:  d.Get(userFullName).(string),
	//	Roles:     d.Get(userRoles).([]int32),
	//	Active:    d.Get(userActive).(bool),
	//}
	//
	//_, err = usersClient(m).UpdateUser(user)
	//if err != nil {
	//	return err
	//}
	//
	//return nil
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id, err := utils.IdFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	deleteErr := retry.Do(
		func() error {
			return usersClient(m).DeleteUser(int32(id))
		},
		retry.RetryIf(
			func(err error) bool {
				return err != nil
			}),
		retry.DelayType(retry.BackOffDelay),
		retry.Attempts(15),
	)

	if deleteErr != nil {
		return diag.FromErr(deleteErr)
	}

	d.SetId("")
	return nil
}

func getCreateOrUpdateUserFromSchema(d *schema.ResourceData) users.CreateUpdateUser {
	return users.CreateUpdateUser{
		UserName:  d.Get(userUsername).(string),
		FullName:  d.Get(userFullName).(string),
		AccountId: int32(d.Get(userAccountId).(int64)),
		Role:      d.Get(userRole).(string),
	}
}

func setUser(d *schema.ResourceData, user *users.User) {
	d.Set(userUsername, user.UserName)
	d.Set(userFullName, user.FullName)
	d.Set(userAccountId, user.AccountId)
	d.Set(userRole, user.Role)
	d.Set(userActive, user.Active)
}
