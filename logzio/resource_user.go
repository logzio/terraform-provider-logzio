package logzio

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jonboydell/logzio_client/users"
)

const (
	userId        string = "id"
	userUsername  string = "username"
	userFullname  string = "fullname"
	userAccountId string = "account_id"
	userRoles     string = "roles"
	userActive    string = "active"
)

/**
 * the endpoint resource schema, what terraform uses to parse and read the template
 */
func resourceUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserCreate,
		Read:   resourceUserRead,
		Update: resourceUserUpdate,
		Delete: resourceUserDelete,

		Schema: map[string]*schema.Schema{
			userUsername: {
				Type:     schema.TypeString,
				Required: true,
			},
			userFullname: {
				Type:     schema.TypeString,
				Required: true,
			},
			userAccountId: {
				Type:     schema.TypeInt,
				Required: true,
			},
			userRoles: {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
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

func resourceUserCreate(d *schema.ResourceData, m interface{}) error {
	var roles []int32
	for _, v := range d.Get(userRoles).([]interface{}) {
		roles = append(roles, int32(v.(int)))
	}

	user := users.User{
		AccountId: int64(d.Get(userAccountId).(int)),
		Username:  d.Get(userUsername).(string),
		Fullname:  d.Get(userFullname).(string),
		Roles:     roles,
	}

	u, err := usersClient(m).CreateUser(user)
	if err != nil {
		return err
	}
	userId := strconv.FormatInt(u.Id, BASE_10)
	d.SetId(userId)

	return nil
}

func resourceUserRead(d *schema.ResourceData, m interface{}) error {
	id, err := idFromResourceData(d)
	if err != nil {
		return err
	}

	user, err := usersClient(m).GetUser(id)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", user.Id))
	d.Set(userAccountId, fmt.Sprintf("%d", user.AccountId))
	d.Set(userUsername, user.Username)
	d.Set(userFullname, user.Fullname)

	var roles []interface{}
	for _, v := range user.Roles {
		roles = append(roles, int(v))
	}

	d.Set(userRoles, roles)
	d.Set(userActive, user.Active)

	return nil
}

func resourceUserUpdate(d *schema.ResourceData, m interface{}) error {
	id, err := idFromResourceData(d)
	if err != nil {
		return err
	}

	accountId, err := strconv.ParseInt(d.Get(userAccountId).(string), BASE_10, BITSIZE_64)
	if err != nil {
		return err
	}

	user := users.User{
		Id:        id,
		AccountId: accountId,
		Username:  d.Get(userUsername).(string),
		Fullname:  d.Get(userFullname).(string),
		Roles:     d.Get(userRoles).([]int32),
		Active:    d.Get(userActive).(bool),
	}

	_, err = usersClient(m).UpdateUser(user)
	if err != nil {
		return err
	}

	return nil
}

func resourceUserDelete(d *schema.ResourceData, m interface{}) error {
	id, err := idFromResourceData(d)
	if err != nil {
		return err
	}

	err = usersClient(m).DeleteUser(id)
	if err != nil {
		return err
	}

	return nil
}
