package logzio

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
	"os"
	"regexp"
	"testing"
)

func TestAccLogzioSubaccount_CreateSubaccount(t *testing.T) {
	accountId := os.Getenv(envLogzioAccountId)
	email := os.Getenv(envLogzioEmail)
	accountName := "test_create_subaccount"
	terraformPlan := testAccCheckLogzioSubaccountConfig(email, accountName, accountId)
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckApiToken(t)
			testAccPreCheckAccountId(t)
			testAccPreCheckEmail(t)
		},
		ProviderFactories: testAccProviderFactories,
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
				ImportStateVerifyIgnore: []string{subAccountEmail},
			},
		},
	})
}

func TestAccLogzioSubaccount_CreateSubaccountEmptySharingObject(t *testing.T) {
	email := os.Getenv(envLogzioEmail)
	accountName := "test_empty_sharing_object"
	terraformPlan := testAccCheckLogzioSubaccountConfig(email, accountName, "")
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckApiToken(t)
			testAccPreCheckEmail(t)
		},
		ProviderFactories: testAccProviderFactories,
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
				ImportStateVerifyIgnore: []string{subAccountEmail},
			},
		},
	})
}

func TestAccLogzioSubaccount_CreateSubaccountNoEmail(t *testing.T) {
	email := ""
	accountName := "test_no_email"
	terraformPlan := testAccCheckLogzioSubaccountConfig(email, accountName, "")
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckApiToken(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      terraformPlan,
				ExpectError: regexp.MustCompile("email must be set"),
			},
		},
	})
}

func TestAccLogzioSubaccount_CreateSubaccountInvalidEmail(t *testing.T) {
	email := "some@invalid.mail"
	accountName := "test_invalid_email"
	terraformPlan := testAccCheckLogzioSubaccountConfig(email, accountName, "")
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckApiToken(t)
			testAccPreCheckEmail(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      terraformPlan,
				ExpectError: regexp.MustCompile("Email must belong to an existing user"),
			},
		},
	})
}

func TestAccLogzioSubaccount_CreateSubaccountNoName(t *testing.T) {
	accountId := os.Getenv(envLogzioAccountId)
	email := os.Getenv(envLogzioEmail)
	accountName := ""
	terraformPlan := testAccCheckLogzioSubaccountConfig(email, accountName, accountId)
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckApiToken(t)
			testAccPreCheckAccountId(t)
			testAccPreCheckEmail(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      terraformPlan,
				ExpectError: regexp.MustCompile("account name must be set"),
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
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckApiToken(t)
			testAccPreCheckAccountId(t)
			testAccPreCheckEmail(t)
		},
		ProviderFactories: testAccProviderFactories,
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
