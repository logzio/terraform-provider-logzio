package logzio

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"os"
	"testing"
)

func TestAccLogzioSubaccount_CreateSubaccount(t *testing.T) {
	accountId := os.Getenv(envLogzioAccountId)
	email := os.Getenv(envLogzioEmail)
	accountName := "test"
	terraformPlan := testAccCheckLogzioSubaccountConfig(email, accountName, accountId)

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
			{
				Config:                  terraformPlan,
				ResourceName:            "logzio_subaccount.test_subaccount",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"email"},
			},
		},
	})
}

func TestAccLogzioSubaccount_CreateSubaccountEmptySharingObject(t *testing.T) {
	email := os.Getenv(envLogzioEmail)
	accountName := "test_empty_sharing_object"
	terraformPlan := testAccCheckLogzioSubaccountConfig(email, accountName, "")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckApiToken(t)
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
			{
				Config:                  terraformPlan,
				ResourceName:            "logzio_subaccount.test_subaccount",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"email"},
			},
		},
	})
}

func TestAccLogzioSubaccount_UpdateSubaccount(t *testing.T) {
	accountId := os.Getenv(envLogzioAccountId)
	email := os.Getenv(envLogzioEmail)
	accountName := "test_update_before"
	accountNameUpdate := "test_update_after"
	resourceName := "logzio_subaccount.test_subaccount"
	terraformPlan := testAccCheckLogzioSubaccountConfig(email, accountName, accountId)
	terraformPlanUpdate := testAccCheckLogzioSubaccountConfig(email, accountNameUpdate, accountId)

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
						resourceName, "email", email),
				),
			},
			{
				Config: terraformPlanUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, subAccountName, accountNameUpdate),
				),
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{subAccountEmail},
			},
		},
	})
}

func testAccCheckLogzioSubaccountConfig(email string, accountName string, accountId string) string {
	return fmt.Sprintf(`
resource "logzio_subaccount" "test_subaccount" {
  email = "%s"
  account_name = "%s"
  retention_days = 2
  max_daily_gb = 1
  sharing_objects_accounts = [
    %s
  ]
}
`, email, accountName, accountId)
}
