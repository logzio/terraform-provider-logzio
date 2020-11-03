package logzio

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"os"
	"strconv"
	"testing"
)

func TestAccLogzioSubaccount_CreateSubaccount(t *testing.T) {

	accountId, _ := strconv.ParseInt(os.Getenv(envLogzioAccountId), BASE_10, BITSIZE_64)
	email := os.Getenv(envLogzioEmail)
	terraformPlan := testAccCheckLogzioSubaccountConfig(email, accountId)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckApiToken(t)
			testAccPreCheckAccountId(t)
			testAccPreCheckEmail(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: terraformPlan,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"logzio_subaccount.test_subaccount", "email", email),
				),
			},
		},
	})
}

func testAccCheckLogzioSubaccountConfig(email string, accountId int64) string {
	return fmt.Sprintf(`
resource "logzio_subaccount" "test_subaccount" {
  email = "%s"
  account_name = "test"
  retention_days = 2
  max_daily_gb = 1
  sharing_objects_accounts = [
    %d
  ]
}
`, email, accountId)
}
