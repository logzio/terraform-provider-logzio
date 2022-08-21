package logzio

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/logzio/logzio_terraform_client/users"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUserRead,
		Schema: map[string]*schema.Schema{
			userId: {
				Type:     schema.TypeInt,
				Optional: true,
			},
			userUsername: {
				Type:     schema.TypeString,
				Optional: true,
			},
			userFullName: {
				Type:     schema.TypeString,
				Computed: true,
			},
			userAccountId: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			userRole: {
				Type:     schema.TypeString,
				Computed: true,
			},
			userActive: {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var client *users.UsersClient
	client, _ = users.New(m.(Config).apiToken, m.(Config).baseUrl)

	id, ok := d.GetOk(userId)
	if ok {
		user, err := client.GetUser(int32(id.(int64)))
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(fmt.Sprintf("%d", user.Id))
		d.Set(userUsername, user.UserName)
		d.Set(userFullName, user.FullName)
		d.Set(userAccountId, user.AccountId)
		d.Set(userRole, user.Role)
		d.Set(userActive, user.Active)
		return nil
	}

	username, ok := d.GetOk(userUsername)
	if ok {
		list, err := client.ListUsers()
		if err != nil {
			return diag.FromErr(err)
		}
		for i := 0; i < len(list); i++ {
			user := list[i]
			if user.UserName == username {
				d.SetId(fmt.Sprintf("%d", user.Id))
				d.Set(userUsername, user.UserName)
				d.Set(userFullName, user.UserName)
				d.Set(userAccountId, user.AccountId)
				d.Set(userRole, user.Role)
				d.Set(userActive, user.Active)
				return nil
			}
		}
	}

	return diag.Errorf("couldn't find user with specified attributes")
}
