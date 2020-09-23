package logzio

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"os"
	"strconv"
	"testing"
)

func TestAccLogzioUser_CreateUser(t *testing.T) {

	username := "test_resource_user@tfacctest.com"
	accountId, _ := strconv.ParseInt(os.Getenv(envLogzioAccountId), BASE_10, BITSIZE_64)
	terraformPlan := testAccCheckLogzioUserConfig(username, "test test", accountId)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckApiToken(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: terraformPlan,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"logzio_user.test_user", "username", username),
				),
			},
			{
				Config: terraformPlan,
				ResourceName:      "logzio_user.test_user",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckLogzioUserConfig(username string, fullname string, accountId int64) string {
	return fmt.Sprintf(`
resource "logzio_user" "test_user" {
  username = "%s"
  fullname = "%s"
  account_id = %d
  roles = [2]
}
`, username, fullname, accountId)
}
