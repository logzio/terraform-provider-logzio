package logzio

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
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

func TestAccLogzioSubaccount_CreateSubaccountWarmRetention(t *testing.T) {
	accountId := os.Getenv(envLogzioWarmAccountId)
	email := os.Getenv(envLogzioEmail)
	accountName := "test_create_subaccountwarm"
	retention := 4
	snapRetention := 2
	terraformPlan := testAccCheckLogzioWarmSubaccountConfig(email, accountName, accountId, retention, snapRetention)
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckApiTokenWarm(t)
			testAccPreCheckWarmAccountId(t)
			testAccPreCheckEmail(t)
		},
		ProviderFactories: testAccWarmProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: terraformPlan,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"logzio_subaccount.test_subaccount", "email", email),
					resource.TestCheckResourceAttr("logzio_subaccount.test_subaccount", "snap_search_retention_days", fmt.Sprintf("%d", snapRetention)),
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

func TestAccLogzioSubaccount_CreateSubaccountWarmRetentionIssues(t *testing.T) {
	accountId := os.Getenv(envLogzioWarmAccountId)
	email := os.Getenv(envLogzioEmail)
	accountName := "test_invalid_snap_retention"
	retention := 4
	snapRetention := 2
	terraformPlanLowRetention := testAccCheckLogzioWarmSubaccountConfig(email, accountName, accountId, 2, snapRetention)
	terraformPlanInvalidSnapRetention := testAccCheckLogzioWarmSubaccountConfig(email, accountName, accountId, retention, 0)
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckApiToken(t)
			testAccPreCheckWarmAccountId(t)
			testAccPreCheckEmail(t)
		},
		ProviderFactories: testAccWarmProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      terraformPlanLowRetention,
				ExpectError: regexp.MustCompile("SnapSearchRetentionDays cannot be set if retentionDays is less than 4"),
			},
			{
				Config:      terraformPlanInvalidSnapRetention,
				ExpectError: regexp.MustCompile("snapSearchRetentionDays should be >= 1"),
			},
		},
	})
}

func TestAccLogzioSubaccount_CreateSubaccountConsumption(t *testing.T) {
	accountId := os.Getenv(envLogzioConsumptionAccountId)
	email := os.Getenv(envLogzioEmail)
	accountName := "test_create_subaccount_consumption"
	isFlexible := "false"
	softLimit := float32(1)
	terraformPlan := testAccCheckLogzioConsumptionSubaccountConfig(email, accountName, accountId, isFlexible, softLimit)
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckApiTokenConsumption(t)
			testAccPreCheckConsumptionAccountId(t)
			testAccPreCheckEmail(t)
		},
		ProviderFactories: testAccConsumptionProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: terraformPlan,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"logzio_subaccount.test_subaccount_consumption", "email", email),
					resource.TestCheckResourceAttr("logzio_subaccount.test_subaccount_consumption", "soft_limit_gb", fmt.Sprintf("%.0f", softLimit)),
				),
			},
			{
				Config:                  terraformPlan,
				ResourceName:            "logzio_subaccount.test_subaccount_consumption",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{subAccountEmail},
			},
		},
	})

}

func TestAccLogzioSubaccount_CreateSubaccountConsumptionIssues(t *testing.T) {
	accountId := os.Getenv(envLogzioConsumptionAccountId)
	email := os.Getenv(envLogzioEmail)
	accountName := "test_invalid_snap_retention"
	softLimit := float32(1)
	softLimitInvalid := float32(0)
	terraformPlanSoftLimitOnFlexible := testAccCheckLogzioConsumptionSubaccountConfig(email, accountName, accountId, "true", softLimit)
	terraformPlanInvalidSoftLimit := testAccCheckLogzioConsumptionSubaccountConfig(email, accountName, accountId, "false", softLimitInvalid)
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckApiTokenConsumption(t)
			testAccPreCheckConsumptionAccountId(t)
			testAccPreCheckEmail(t)
		},
		ProviderFactories: testAccConsumptionProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      terraformPlanSoftLimitOnFlexible,
				ExpectError: regexp.MustCompile("when isFlexible=true SoftLimitGB should be empty or omitted"),
			},
			{
				Config:      terraformPlanInvalidSoftLimit,
				ExpectError: regexp.MustCompile("SoftLimitGB should be > 0 when set"),
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
  frequency_minutes = 3
  utilization_enabled = "true"
  max_daily_gb = 1
  sharing_objects_accounts = [
    %s
  ]
}
`, email, accountName, accountId)
}

func testAccCheckLogzioWarmSubaccountConfig(email string, accountName string, accountId string, retention int, snapRetention int) string {
	return fmt.Sprintf(`
resource "logzio_subaccount" "test_subaccount" {
  email = "%s"
  account_name = "%s"
  retention_days = %d
  utilization_enabled = "true"
  max_daily_gb = 1
  reserved_daily_gb = 0.5
  sharing_objects_accounts = [
    %s
  ]
  flexible = "true"
  snap_search_retention_days = %d
}
`, email, accountName, retention, accountId, snapRetention)
}

func testAccCheckLogzioConsumptionSubaccountConfig(email, accountName, accountId, isFlexible string, softLimitGb float32) string {
	return fmt.Sprintf(`
resource "logzio_subaccount" "test_subaccount_consumption" {
  email = "%s"
  account_name = "%s"
  retention_days = 2
  frequency_minutes = 3
  utilization_enabled = "true"
  max_daily_gb = 1
  sharing_objects_accounts = [
    %s
  ]
  flexible = %s
  soft_limit_gb = %f
}
`, email, accountName, accountId, isFlexible, softLimitGb)
}
