package logzio

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/logzio/logzio_terraform_client/authentication_groups"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
	"regexp"
	"testing"
)

func TestAccLogzioAuthenticationGroups_AuthenticationGroups(t *testing.T) {
	resourceName := "tf_create_test"
	fullResourceName := resourceAuthenticationGroupsType + "." + resourceName

	userRolesInConfigCreate := []string{
		authentication_groups.AuthGroupsUserRoleAdmin,
		authentication_groups.AuthGroupsUserRoleRegular,
		authentication_groups.AuthGroupsUserRoleReadonly,
	}

	userRolesInConfigUpdate := []string{
		authentication_groups.AuthGroupsUserRoleAdmin,
		authentication_groups.AuthGroupsUserRoleRegular,
		authentication_groups.AuthGroupsUserRoleRegular,
	}
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				// Create
				Config: getAuthGroupsConfig(resourceName, userRolesInConfigCreate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(fullResourceName, authGroupsId),
					resource.TestCheckResourceAttr(fullResourceName, authGroupsAuthGroup+".#", "3"),
					testAccCheckAuthGroups(fullResourceName, userRolesInConfigCreate),
				),
			},
			{
				// Update - change role
				Config: getAuthGroupsConfig(resourceName, userRolesInConfigUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(fullResourceName, authGroupsId),
					resource.TestCheckResourceAttr(fullResourceName, authGroupsAuthGroup+".#", "3"),
					testAccCheckAuthGroups(fullResourceName, userRolesInConfigUpdate),
				),
			},
			{
				// Update - remove a group
				Config: getAuthGroupsConfigUpdate(resourceName, userRolesInConfigUpdate[0:2]),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(fullResourceName, authGroupsId),
					resource.TestCheckResourceAttr(fullResourceName, authGroupsAuthGroup+".#", "2"),
					testAccCheckAuthGroups(fullResourceName, userRolesInConfigUpdate[0:2]),
				),
			},
			{
				// Import
				Config:            getAuthGroupsConfigUpdate(resourceName, userRolesInConfigUpdate[0:2]),
				ResourceName:      fullResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func getAuthGroupsConfig(resourceName string, roles []string) string {
	return fmt.Sprintf(`resource "logzio_authentication_groups" "%s" {
	authentication_group {
		group = "group_testing_1"
		user_role = "%s"
	}
	authentication_group {
		group = "group_testing_2"
		user_role = "%s"
	}
	authentication_group {
		group = "group_testing_3"
		user_role = "%s"
	}
}
`, resourceName, roles[0], roles[1], roles[2])
}

func getAuthGroupsConfigUpdate(resourceName string, roles []string) string {
	return fmt.Sprintf(`resource "logzio_authentication_groups" "%s" {
	authentication_group {
		group = "group_testing_1"
		user_role = "%s"
	}
	authentication_group {
		group = "group_testing_2"
		user_role = "%s"
	}
}
`, resourceName, roles[0], roles[1])
}

func testAccCheckAuthGroups(resourceName string, userRoles []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no authentication groups id")
		}

		rolesFromSchema := make([]string, 0)
		for key, val := range rs.Primary.Attributes {
			matched, err := regexp.MatchString(`^authentication_group\.\d+\.user_role`, key)
			if err != nil {
				return err
			}

			if matched {
				rolesFromSchema = append(rolesFromSchema, val)
			}
		}

		if !compareRoles(rolesFromSchema, userRoles) {
			return fmt.Errorf("user roles are not as set in the configuration")
		}

		return nil
	}
}

// comparing the list based on the number of occurrences of the roles
func compareRoles(roleList1 []string, roleList2 []string) bool {
	occurr1 := make(map[string]int)
	occurr2 := make(map[string]int)

	for _, role1 := range roleList1 {
		occurr1[role1]++
	}

	for _, role2 := range roleList2 {
		occurr2[role2]++
	}

	for key, val := range occurr1 {
		if occurr2[key] != val {
			return false
		}
	}

	return true
}
