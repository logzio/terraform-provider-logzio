package logzio

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
	"os"
	"regexp"
	"strconv"
	"testing"
)

func TestAccLogzioMetricsAccount_CreateMetricsAccount(t *testing.T) {
	t.Parallel()
	defer utils.SleepAfterTest()
	accountId := os.Getenv(envLogzioAccountId)
	email := os.Getenv(envLogzioEmail)
	accountName := "test_create_metrics_" + getRandomId()
	terraformPlan := testAccCheckLogzioMetricsAccountConfig(email, accountName, accountId)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckApiToken(t)
			testAccPreCheckAccountId(t)
			testAccPreCheckEmail(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMetricsAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: terraformPlan,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"logzio_metrics_account.test_subaccount", "email", email),
				),
			},
			{
				Config:                  terraformPlan,
				ResourceName:            "logzio_metrics_account.test_subaccount",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{metricsAccountEmail},
			},
		},
	})
}

func TestAccLogzioMetricsAccount_CreateMetricsAccountEmptyAuthorizedAccounts(t *testing.T) {
	t.Parallel()
	defer utils.SleepAfterTest()
	email := os.Getenv(envLogzioEmail)
	accountName := "test_metrics_empty_sharing_" + getRandomId()
	terraformPlan := testAccCheckLogzioMetricsAccountConfig(email, accountName, "")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckApiToken(t)
			testAccPreCheckEmail(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMetricsAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: terraformPlan,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"logzio_metrics_account.test_subaccount", "email", email),
				),
			},
			{
				Config:                  terraformPlan,
				ResourceName:            "logzio_metrics_account.test_subaccount",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{metricsAccountEmail},
			},
		},
	})
}

func TestAccLogzioMetricsAccount_CreateMetricsAccountNoEmail(t *testing.T) {
	t.Parallel()
	defer utils.SleepAfterTest()
	email := ""
	accountName := "test_no_email"
	terraformPlan := testAccCheckLogzioMetricsAccountConfig(email, accountName, "")

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

func TestAccLogzioMetricsAccount_CreateMetricsAccountInvalidEmail(t *testing.T) {
	t.Parallel()
	defer utils.SleepAfterTest()
	email := "some@invalid.mail"
	accountName := "test_invalid_email"
	terraformPlan := testAccCheckLogzioMetricsAccountConfig(email, accountName, "")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckApiToken(t)
			testAccPreCheckEmail(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      terraformPlan,
				ExpectError: regexp.MustCompile("There is no registered user for the passed email"),
			},
		},
	})
}

func TestAccLogzioMetricsAccount_CreateMetricsAccountNoName(t *testing.T) {
	t.Parallel()
	defer utils.SleepAfterTest()
	accountId := os.Getenv(envLogzioAccountId)
	email := os.Getenv(envLogzioEmail)
	resourceName := "logzio_metrics_account.test_subaccount"
	terraformPlan := testAccCheckLogzioMetricsAccountConfigNoName(email, accountId)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckApiToken(t)
			testAccPreCheckAccountId(t)
			testAccPreCheckEmail(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMetricsAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: terraformPlan,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, metricsAccountPlanUts, "160"),
					resource.TestCheckResourceAttr(resourceName, metricsAccountEmail, email),
					resource.TestMatchResourceAttr(resourceName, metricsAccountName, regexp.MustCompile(".*_metrics")),
				),
			},
		},
	})
}

func TestAccLogzioMetricsAccount_UpdateMetricsAccount(t *testing.T) {
	t.Parallel()
	defer utils.SleepAfterTest()
	accountId := os.Getenv(envLogzioAccountId)
	email := os.Getenv(envLogzioEmail)
	suffix := getRandomId()
	accountName := "test_metrics_update_before_" + suffix
	accountNameUpdate := "test_metrics_update_after_" + suffix
	resourceName := "logzio_metrics_account.test_subaccount"
	terraformPlan := testAccCheckLogzioMetricsAccountConfig(email, accountName, accountId)
	terraformPlanUpdate := testAccCheckLogzioMetricsAccountConfig(email, accountNameUpdate, accountId)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckApiToken(t)
			testAccPreCheckAccountId(t)
			testAccPreCheckEmail(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMetricsAccountDestroy,
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

func testAccCheckMetricsAccountDestroy(s *terraform.State) error {
	if testAccProvider.Meta() == nil {
		return nil
	}
	client, err := MetricsAccountClient(testAccProvider.Meta())
	if err != nil {
		return err
	}
	for _, r := range s.RootModule().Resources {
		if r.Type != "logzio_metrics_account" {
			continue
		}
		id, err := strconv.ParseInt(r.Primary.ID, 10, 64)
		if err != nil {
			return err
		}
		_, err = client.GetMetricsAccount(id)
		if err == nil {
			return fmt.Errorf("metrics account %s still exists", r.Primary.ID)
		}
	}
	return nil
}

func testAccCheckLogzioMetricsAccountConfig(email string, accountName string, accountId string) string {
	return fmt.Sprintf(`
resource "logzio_metrics_account" "test_subaccount" {
  email = "%s"
  account_name = "%s"
  plan_uts = 100
  authorized_accounts = [
    %s
  ]
}
`, email, accountName, accountId)
}

func testAccCheckLogzioMetricsAccountConfigNoName(email string, accountId string) string {
	return fmt.Sprintf(`
resource "logzio_metrics_account" "test_subaccount" {
  email = "%s"
  plan_uts = 160
  authorized_accounts = [
    %s
  ]
}
`, email, accountId)
}

func testAccCheckLogzioMetricsAccountFromDatasource(email string, accountId string) string {
	return fmt.Sprintf(`
resource "logzio_metrics_account" "test_subaccount" {
  email = "%s"
  account_name = data.logzio_metrics_account.metrics_account_datasource_by_id.account_name
  plan_uts = 100
  authorized_accounts = [
    %s
  ]
}
`, email, accountId)
}
