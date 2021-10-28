package logzio

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/logzio/logzio_terraform_client/authentication_groups"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	authGroupsId = "manage_groups_id"
	authGroupsManageGroups = "manage_groups"
	authGroupGroup = "group"
	authGroupUserRole = "user_role"
)

func resourceAuthenticationGroups() *schema.Resource {
	return &schema.Resource{
		Create: resourceAuthenticationGroupsCreate,
		Read:   resourceAuthenticationGroupsRead,
		Update: resourceAuthenticationGroupsUpdate,
		Delete: resourceAuthenticationGroupDelete,
		// TODO: implement my own import function, since the API does not return an id for group.
		//Importer: &schema.ResourceImporter{
		//	State: schema.ImportStatePassthrough,
		//},
		Schema: map[string]*schema.Schema{
			authGroupsId: {
				Type: schema.TypeInt,
				Computed: true,
			},
			authGroupsManageGroups: {
				Type: schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						authGroupGroup: {
							Type: schema.TypeString,
							Optional: true,
							ValidateFunc: utils.ValidateGroupName,
						},
						authGroupUserRole: {
							Type: schema.TypeString,
							Required: true,
							ValidateFunc: utils.ValidateUserRole,
						},
					},
				},
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Second),
			Update: schema.DefaultTimeout(5 * time.Second),
			Delete: schema.DefaultTimeout(5 * time.Second),
		},
	}
}

func authenticationGroupsClient(m interface{}) *authentication_groups.AuthenticationGroupsClient {
	var client *authentication_groups.AuthenticationGroupsClient
	client, _ = authentication_groups.New(m.(Config).apiToken, m.(Config).baseUrl)
	return client
}

func resourceAuthenticationGroupsCreate(d *schema.ResourceData, m interface{}) error {
	createGroups := getAuthenticationGroupsFromSchema(d)
	groups, err := authenticationGroupsClient(m).PostAuthenticationGroups(createGroups)
	if err != nil {
		return err
	}

	// Logz.io authentication groups API doesn't return id, we need to create a random id for TF.
	rand.Seed(time.Now().UnixNano())
	id := rand.Int()
	d.SetId(strconv.FormatInt(int64(id), 10))

	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		err = resourceAuthenticationGroupsRead(d, m)
		if err != nil {
			if strings.Contains(err.Error(), "failed with missing authentication groups") {
				return resource.RetryableError(err)
			}
		}

		return resource.NonRetryableError(err)
	})
}

func resourceAuthenticationGroupsRead(d *schema.ResourceData, m interface{}) error {
	id, err := utils.IdFromResourceData(d)
	if err != nil {
		return nil
	}

	groups, err := authenticationGroupsClient(m).GetAuthenticationGroups()
	if err != nil {
		return nil
	}

	setAuthenticationGroups(id, groups, d)
	return nil
}

func resourceAuthenticationGroupsUpdate(d *schema.ResourceData, m interface{}) error {
	updateAuthGroup := getAuthenticationGroupsFromSchema(d)

	// Prevent deleting auth group by sending empty set.
	// Makes the user use a destroy action instead, to keep with the TF conventions
	if len(updateAuthGroup) == 0 {
		return fmt.Errorf("can't delete by sending an empty set. you need to destroy the resource in order to delete all groups")
	}

	return resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		err := resourceAuthenticationGroupsRead(d, m)
		if err != nil {
			groupsFromSchema := getAuthenticationGroupsFromSchema(d)
			if strings.Contains(err.Error(), "failed with missing authentication groups") &&
				!reflect.DeepEqual(updateAuthGroup, groupsFromSchema) {
				return resource.RetryableError(fmt.Errorf("authentication groups not updated yet: %s", err.Error()))
			}
		}

		return resource.NonRetryableError(err)
	})
}

func resourceAuthenticationGroupDelete(d *schema.ResourceData, m interface{}) error {
	return resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		_, err := authenticationGroupsClient(m).PostAuthenticationGroups([]authentication_groups.AuthenticationGroup{})
		if err != nil {
			return resource.RetryableError(err)
		}

		return resource.NonRetryableError(err)
	})
}

func getAuthenticationGroupsFromSchema(d *schema.ResourceData) []authentication_groups.AuthenticationGroup {
	var groups []authentication_groups.AuthenticationGroup
	var groupToAdd authentication_groups.AuthenticationGroup
	groupsFromSchema := d.Get(authGroupsManageGroups).(*schema.Set).List()

	for _, group := range groupsFromSchema {
		groupToAdd.Group = group.(map[string]interface{})[authGroupGroup].(string)
		groupToAdd.UserRole = group.(map[string]interface{})[authGroupUserRole].(string)
		groups = append(groups, groupToAdd)
	}

	return groups
}

func setAuthenticationGroups(id int64, groups []authentication_groups.AuthenticationGroup, d *schema.ResourceData) {
	var groupsToSchema []interface{}
	d.Set(authGroupsId, id)

	for _, group := range groups {
		groupMap := map[string]interface{} {
			authGroupGroup: group.Group,
			authGroupUserRole: group.UserRole,
		}

		groupsToSchema = append(groupsToSchema, groupMap)
	}

	d.Set(authGroupsManageGroups, groupsToSchema)
}