package main

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceUser(t *testing.T) {

	username := "test@sometest.com"
	fullname := "test test"
	accountId, _ := strconv.ParseInt(os.Getenv(envLogzioAccountId), BASE_10, BITSIZE_64)
	terraformPlan := testAccCheckLogzioUserDatasourceConfig(username, fullname, accountId)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				ExpectNonEmptyPlan:        true,
				Config:                    terraformPlan,
				PreventPostDestroyRefresh: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.logzio_user.by_username", "username", username),
					resource.TestCheckResourceAttr("data.logzio_user.by_username", "fullname", fullname),
					resource.TestCheckOutput("test", fullname),
				),
			},
		},
	})
}

func testAccCheckLogzioUserDatasourceConfig(username string, fullname string, accountId int64) string {
	return fmt.Sprintf(`
resource "logzio_user" "test_user" {
  username = "%s"
  fullname = "%s"
  account_id = %d
  roles = [2]
}

data "logzio_user" "by_username" {
  username = "%s"
  depends_on = ["logzio_user.test_user"]
}

output "test" {
  value = "${data.logzio_user.by_username.fullname}"
}
`, username, fullname, accountId, username)
}
