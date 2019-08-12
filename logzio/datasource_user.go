package logzio

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jonboydell/logzio_client/users"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceUserRead,
		Schema: map[string]*schema.Schema{
			userId: {
				Type:     schema.TypeInt,
				Optional: true,
			},
			userUsername: {
				Type:     schema.TypeString,
				Optional: true,
			},
			userFullname: {
				Type:     schema.TypeString,
				Computed: true,
			},
			userAccountId: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			userRoles: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			userActive: {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceUserRead(d *schema.ResourceData, m interface{}) error {
	var client *users.UsersClient
	client, _ = users.New(m.(Config).apiToken, m.(Config).baseUrl)

	userId, ok := d.GetOk(userId)
	if ok {
		user, err := client.GetUser(userId.(int64))
		if err != nil {
			return err
		}
		d.SetId(fmt.Sprintf("%d", user.Id))
		d.Set(userUsername, user.Username)
		d.Set(userFullname, user.Fullname)
		d.Set(userAccountId, user.AccountId)

		var roles []interface{}
		for _, v := range user.Roles {
			roles = append(roles, int(v))
		}

		d.Set(userRoles, roles)
		d.Set(userActive, user.Active)
		return nil
	}

	username, ok := d.GetOk(userUsername)
	if ok {
		list, err := client.ListUsers()
		if err != nil {
			return err
		}
		for i := 0; i < len(list); i++ {
			user := list[i]
			if user.Username == username {
				d.SetId(fmt.Sprintf("%d", user.Id))
				d.Set(userUsername, user.Username)
				d.Set(userFullname, user.Fullname)
				d.Set(userAccountId, user.AccountId)

				var roles []interface{}
				for _, v := range user.Roles {
					roles = append(roles, int(v))
				}

				d.Set(userRoles, roles)
				d.Set(userActive, user.Active)
				return nil
			}
		}
	}

	return fmt.Errorf("couldn't find user with specified attributes")
}
