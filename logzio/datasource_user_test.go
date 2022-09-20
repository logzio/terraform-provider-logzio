package logzio

import (
	"fmt"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceUser(t *testing.T) {
	username := "test_datasource_user@tfacctest.com"
	fullname := "test test"
	accountId, _ := strconv.ParseInt(os.Getenv(envLogzioAccountId), utils.BASE_10, utils.BITSIZE_64)
	terraformPlan := testAccCheckLogzioUserDatasourceConfig(username, fullname, accountId)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckApiToken(t)
			testAccPreCheckAccountId(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
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
	str := fmt.Sprintf(`
resource "logzio_user" "test_user" {
  username = "%s"
  fullname = "%s"
  account_id = %d
  role = "USER_ROLE_READONLY"
}

data "logzio_user" "by_username" {
  username = "%s"
  depends_on = ["logzio_user.test_user"]
}

output "test" {
  value = "${data.logzio_user.by_username.fullname}"
}
`, username, fullname, accountId, username)

	return str
}
