package logzio

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/logzio/logzio_terraform_client/authentication_groups"
	"math/rand"
	"strconv"
	"time"
)

const (
	authGroupsDatasourceRetries = 3
)

func dataSourceAuthenticationGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAuthenticationGroupsRead,
		Schema: map[string]*schema.Schema{
			authGroupsAuthGroup: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						authGroupGroup: {
							Type:     schema.TypeString,
							Computed: true,
						},
						authGroupUserRole: {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAuthenticationGroupsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	groups, err := getAuthGroups(authGroupsDatasourceRetries, m)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(groups) > 0 {
		// Logz.io authentication groups API doesn't return id, we need to create a random id for TF.
		rand.Seed(time.Now().UnixNano())
		id := rand.Int()
		d.SetId(strconv.FormatInt(int64(id), 10))
		setAuthenticationGroupsDatasource(id, groups, d)
	}

	return nil
}

func getAuthGroups(retries int, m interface{}) ([]authentication_groups.AuthenticationGroup, error) {
	groups, err := authenticationGroupsClient(m).GetAuthenticationGroups()
	if err != nil && retries > 0 {
		time.Sleep(time.Second * 2)
		groups, err = getAuthGroups(retries-1, m)
	}

	return groups, err
}

func setAuthenticationGroupsDatasource(id int, groups []authentication_groups.AuthenticationGroup, d *schema.ResourceData) {
	var groupsToSchema []map[string]interface{}

	if id != 0 {
		d.Set(authGroupsId, id)
	}

	for _, group := range groups {
		groupMap := map[string]interface{}{
			authGroupGroup:    group.Group,
			authGroupUserRole: group.UserRole,
		}

		groupsToSchema = append(groupsToSchema, groupMap)
	}

	d.Set(authGroupsAuthGroup, groupsToSchema)
}
