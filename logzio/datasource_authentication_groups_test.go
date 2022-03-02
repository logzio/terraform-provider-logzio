package logzio

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/logzio/logzio_terraform_client/authentication_groups"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
	"testing"
)

func TestAccDataSourceAuthenticationGroups(t *testing.T) {
	resourceName := "test_resource"
	fullResourceName := resourceAuthenticationGroupsType + "." + resourceName
	datasourceName := "my_auth_group_datasource"
	fullName := "data." + resourceAuthenticationGroupsType + "." + datasourceName
	userRolesInConfigCreate := []string{
		authentication_groups.AuthGroupsUserRoleAdmin,
		authentication_groups.AuthGroupsUserRoleRegular,
		authentication_groups.AuthGroupsUserRoleReadonly,
	}
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckApiToken(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: getAuthGroupsConfig(resourceName, userRolesInConfigCreate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(fullResourceName, authGroupsId),
					resource.TestCheckResourceAttr(fullResourceName, authGroupsAuthGroup+".#", "3"),
					testAccCheckAuthGroups(fullResourceName, userRolesInConfigCreate),
				),
			},
			{
				Config: getAuthGroupsConfig(resourceName, userRolesInConfigCreate) +
					getDatasourceConfig(datasourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(fullName, authGroupsAuthGroup+".#", "3"),
				),
			},
		},
	})
}

func getDatasourceConfig(resourceName string) string {
	return fmt.Sprintf(`data "logzio_authentication_groups" "%s" {
}
`, resourceName)
}
