package logzio

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/avast/retry-go"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/logzio/logzio_terraform_client/users"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
)

const (
	userId        string = "id"
	userUsername  string = "username"
	userFullName  string = "fullname"
	userAccountId string = "account_id"
	userRole      string = "role"
	userActive    string = "active"

	userRetryAttempts = 16
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
	user, err := usersClient(m).GetUser(int32(id))
	if err != nil {
		tflog.Error(ctx, err.Error())
		if strings.Contains(err.Error(), "missing user") {
			// If we were not able to find the resource - delete from state
			d.SetId("")
			return diag.Diagnostics{}
		} else {
			return diag.FromErr(err)
		}

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
	currUser, err := usersClient(m).GetUser(int32(id))
	if err != nil {
		return diag.FromErr(err)
	}

	// check if we need to activate or deactivate user
	activeFromSchema := d.Get(userActive).(bool)
	if activeFromSchema != currUser.Active {
		tflog.Debug(ctx, fmt.Sprintf("detected activation change, from %t to %t", currUser.Active, activeFromSchema))
		err = changeUserActivation(int32(id), activeFromSchema, m)
		if err != nil {
			tflog.Error(ctx, "error occurred while trying to change user activation. If other changes were planned, they will not be applied")
			return diag.FromErr(err)
		}
	}

	// check if there are more updates to be applied
	if updateUser.UserName != currUser.UserName ||
		updateUser.FullName != currUser.FullName ||
		updateUser.Role != currUser.Role ||
		updateUser.AccountId != currUser.AccountId {
		_, err = usersClient(m).UpdateUser(int32(id), updateUser)

		var diagRet diag.Diagnostics
		readErr := retry.Do(
			func() error {
				diagRet = resourceUserRead(ctx, d, m)
				if diagRet.HasError() {
					return fmt.Errorf("received error from read user")
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
						userAfterUpdate := getCreateOrUpdateUserFromSchema(d)
						return !reflect.DeepEqual(userAfterUpdate, updateUser)
					}
				}),
			retry.DelayType(retry.BackOffDelay),
			retry.Attempts(userRetryAttempts),
		)

		if readErr != nil {
			tflog.Error(ctx, "could not update schema")
			return diagRet
		}

		return nil
	}

	return nil
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id, err := utils.IdFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	err = usersClient(m).DeleteUser(int32(id))

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func getCreateOrUpdateUserFromSchema(d *schema.ResourceData) users.CreateUpdateUser {
	return users.CreateUpdateUser{
		UserName:  d.Get(userUsername).(string),
		FullName:  d.Get(userFullName).(string),
		AccountId: int32(d.Get(userAccountId).(int)),
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

func changeUserActivation(id int32, activate bool, m interface{}) error {
	var err error
	if activate {
		err = usersClient(m).UnSuspendUser(id)
	} else {
		err = usersClient(m).SuspendUser(id)
	}

	return err
}
