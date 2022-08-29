package logzio

import (
	"context"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/logzio/logzio_terraform_client/authentication_groups"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	authGroupsId        = "manage_groups_id"
	authGroupsAuthGroup = "authentication_group"
	authGroupGroup      = "group"
	authGroupUserRole   = "user_role"
)

func resourceAuthenticationGroups() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAuthenticationGroupsCreate,
		ReadContext:   resourceAuthenticationGroupsRead,
		UpdateContext: resourceAuthenticationGroupsUpdate,
		DeleteContext: resourceAuthenticationGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			// Id created by TF to keep with conventions, because the Logz.io auth groups API doesn't create one.
			authGroupsId: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			authGroupsAuthGroup: {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						authGroupGroup: {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: utils.ValidateGroupName,
						},
						authGroupUserRole: {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: utils.ValidateUserRole,
						},
					},
				},
			},
		},
	}
}

func authenticationGroupsClient(m interface{}) *authentication_groups.AuthenticationGroupsClient {
	var client *authentication_groups.AuthenticationGroupsClient
	client, _ = authentication_groups.New(m.(Config).apiToken, m.(Config).baseUrl)
	return client
}

func resourceAuthenticationGroupsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	createGroups := getAuthenticationGroupsFromSchema(d)
	_, err := authenticationGroupsClient(m).PostAuthenticationGroups(createGroups)
	if err != nil {
		return diag.FromErr(err)
	}

	// Logz.io authentication groups API doesn't return id, we need to create a random id for TF.
	rand.Seed(time.Now().UnixNano())
	id := rand.Int()
	d.SetId(strconv.FormatInt(int64(id), 10))

	return resourceAuthenticationGroupsRead(ctx, d, m)
}

func resourceAuthenticationGroupsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id, err := utils.IdFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	var groups []authentication_groups.AuthenticationGroup
	readErr := retry.Do(
		func() error {
			groups, err = authenticationGroupsClient(m).GetAuthenticationGroups()
			if err != nil {
				return err
			}

			return nil
		},
		retry.RetryIf(
			func(err error) bool {
				if err != nil {
					if strings.Contains(err.Error(), "failed with missing authentication groups") {
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

	setAuthenticationGroups(id, groups, d)
	return nil
}

func resourceAuthenticationGroupsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	updateAuthGroup := getAuthenticationGroupsFromSchema(d)

	// Prevent deleting auth group by sending empty set.
	// Makes the user use a destroy action instead, to keep with the TF conventions
	if len(updateAuthGroup) == 0 {
		return diag.Errorf("can't delete by sending an empty set. you need to destroy the resource in order to delete all groups")
	}

	_, err := authenticationGroupsClient(m).PostAuthenticationGroups(updateAuthGroup)
	if err != nil {
		return diag.FromErr(err)
	}

	var diagRet diag.Diagnostics
	readErr := retry.Do(
		func() error {
			diagRet = resourceAuthenticationGroupsRead(ctx, d, m)
			if diagRet.HasError() {
				return fmt.Errorf("received error from read authentication groups")
			}

			return nil
		},
		retry.RetryIf(
			func(err error) bool {
				if err != nil {
					return true
				} else {
					// Check if the update shows on read
					// if not updated yet - retry
					groupsFromSchema := getAuthenticationGroupsFromSchema(d)
					return !isSameAuthGroups(updateAuthGroup, groupsFromSchema)
				}
			}),
		retry.DelayType(retry.BackOffDelay),
		retry.Attempts(15),
	)

	if readErr != nil {
		tflog.Error(ctx, "could not update schema")
		return diagRet
	}

	return nil
}

func resourceAuthenticationGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	deleteErr := retry.Do(
		func() error {
			_, err := authenticationGroupsClient(m).PostAuthenticationGroups([]authentication_groups.AuthenticationGroup{})
			return err
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

func getAuthenticationGroupsFromSchema(d *schema.ResourceData) []authentication_groups.AuthenticationGroup {
	var groups []authentication_groups.AuthenticationGroup
	var groupToAdd authentication_groups.AuthenticationGroup
	groupsFromSchema := d.Get(authGroupsAuthGroup).(*schema.Set).List()

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
		groupMap := map[string]interface{}{
			authGroupGroup:    group.Group,
			authGroupUserRole: group.UserRole,
		}

		groupsToSchema = append(groupsToSchema, groupMap)
	}

	d.Set(authGroupsAuthGroup, groupsToSchema)
}

func isSameAuthGroups(authGroups1, authGroups2 []authentication_groups.AuthenticationGroup) bool {
	if len(authGroups1) != len(authGroups2) {
		return false
	}

	diff := make(map[string]int, len(authGroups1))
	for _, group1 := range authGroups1 {
		diff[group1.Group+group1.UserRole]++
	}

	for _, group2 := range authGroups2 {
		if _, ok := diff[group2.Group+group2.UserRole]; !ok {
			return false
		}

		diff[group2.Group+group2.UserRole] -= 1
		if diff[group2.Group+group2.UserRole] == 0 {
			delete(diff, group2.Group+group2.UserRole)
		}
	}

	return len(diff) == 0
}
