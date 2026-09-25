package logzio

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/logzio/logzio_terraform_client/tracing_accounts"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
	"os"
	"regexp"
	"strconv"
	"testing"
)

// Runs against the consumption account: the tracing accounts API is consumption-only and gated
// per account by the create-tracing-consumption-account feature flag.
func TestAccLogzioTracingAccount_CreateUpdateTracingAccount(t *testing.T) {
	accountId := os.Getenv(envLogzioConsumptionAccountId)
	email := os.Getenv(envLogzioEmail)
	suffix := getRandomId()
	accountName := "test_tracing_create_" + suffix
	accountNameUpdate := "test_tracing_update_" + suffix
	resourceName := "logzio_tracing_account.test_tracing_account"
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckApiTokenConsumption(t)
			testAccPreCheckConsumptionAccountId(t)
			testAccPreCheckEmail(t)
		},
		ProviderFactories: testAccConsumptionProviderFactories,
		CheckDestroy:      testAccCheckTracingAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLogzioTracingAccountConfig(email, accountName, "1", accountId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, tracingAccountName, accountName),
					resource.TestCheckResourceAttr(resourceName, tracingAccountMaxDailyGB, "1"),
					resource.TestCheckResourceAttr(resourceName, tracingAccountAuthorizedAccounts+".0", accountId),
					resource.TestCheckResourceAttrSet(resourceName, tracingAccountId),
					resource.TestCheckResourceAttrSet(resourceName, tracingAccountRetention),
				),
			},
			{
				Config: testAccCheckLogzioTracingAccountConfig(email, accountNameUpdate, "2", accountId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, tracingAccountName, accountNameUpdate),
					resource.TestCheckResourceAttr(resourceName, tracingAccountMaxDailyGB, "2"),
				),
			},
			{
				Config:                  testAccCheckLogzioTracingAccountConfig(email, accountNameUpdate, "2", accountId),
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{tracingAccountEmail},
			},
		},
	})
}

// Both cases fail at plan time, before any API call, so they do not depend on the tracing
// feature flag being enabled for the test account.

func TestAccLogzioTracingAccount_CreateTracingAccountNoMaxDailyGB(t *testing.T) {
	email := os.Getenv(envLogzioEmail)
	terraformPlan := fmt.Sprintf(`
resource "logzio_tracing_account" "test_tracing_account" {
  email = "%s"
  account_name = "test_tracing_no_max_daily_gb"
  authorized_accounts = []
}
`, email)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckApiToken(t)
			testAccPreCheckEmail(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      terraformPlan,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`The argument "max_daily_gb" is required`),
			},
		},
	})
}

func TestAccLogzioTracingAccount_CreateTracingAccountNegativeMaxDailyGB(t *testing.T) {
	email := os.Getenv(envLogzioEmail)
	terraformPlan := fmt.Sprintf(`
resource "logzio_tracing_account" "test_tracing_account" {
  email = "%s"
  account_name = "test_tracing_negative_max_daily_gb"
  max_daily_gb = -1
  authorized_accounts = []
}
`, email)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckApiToken(t)
			testAccPreCheckEmail(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      terraformPlan,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("expected max_daily_gb to be at least"),
			},
		},
	})
}

// Uses the consumption token directly: the default test provider's token belongs to a
// subscription account, which the tracing accounts API rejects outright - so a lookup through
// it would "fail" for every id and the check would prove nothing.
func testAccCheckTracingAccountDestroy(s *terraform.State) error {
	client, err := tracing_accounts.New(os.Getenv(envLogzioApiTokenConsumption), fmt.Sprintf(baseUrl, ""))
	if err != nil {
		return err
	}
	for _, r := range s.RootModule().Resources {
		if r.Type != resourceTracingAccountType {
			continue
		}
		id, err := strconv.ParseInt(r.Primary.ID, 10, 64)
		if err != nil {
			return err
		}
		_, err = client.GetTracingAccount(id)
		if err == nil {
			return fmt.Errorf("tracing account %s still exists", r.Primary.ID)
		}
	}
	return nil
}

func testAccCheckLogzioTracingAccountConfig(email, accountName, maxDailyGB, accountId string) string {
	return fmt.Sprintf(`
resource "logzio_tracing_account" "test_tracing_account" {
  email = "%s"
  account_name = "%s"
  max_daily_gb = %s
  authorized_accounts = [
    %s
  ]
}
`, email, accountName, maxDailyGB, accountId)
}
