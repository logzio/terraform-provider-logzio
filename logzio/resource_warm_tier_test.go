package logzio

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/logzio/logzio_terraform_client/sub_accounts"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
)

// Runs against the warm account: sets the warm retention of a new sub account, checks that updating the sub account
// keeps it, and that values above the main account's are rejected.
func TestAccLogzioWarmTier_SubAccount(t *testing.T) {
	accountId := os.Getenv(envLogzioWarmAccountId)
	email := os.Getenv(envLogzioEmail)
	suffix := getRandomId()
	accountName := "test_warm_tier_" + suffix
	accountNameUpdate := "test_warm_tier_update_" + suffix
	resourceName := "logzio_warm_tier.test_warm_tier"
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckApiTokenWarm(t)
			testAccPreCheckWarmAccountId(t)
			testAccPreCheckEmail(t)
		},
		ProviderFactories: testAccWarmProviderFactories,
		CheckDestroy:      testAccCheckSubaccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLogzioWarmTierSubAccountConfig(email, accountName, accountId, 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, warmTierAccountId, "logzio_subaccount.test_subaccount", subAccountId),
					resource.TestCheckResourceAttr(resourceName, warmTierSnapSearchRetentionDays, "1"),
					testAccCheckWarmRetention(resourceName, 1),
				),
			},
			{
				Config: testAccCheckLogzioWarmTierSubAccountConfig(email, accountName, accountId, 2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, warmTierSnapSearchRetentionDays, "2"),
					testAccCheckWarmRetention(resourceName, 2),
				),
			},
			{
				// updating the sub account replaces the whole account, and must keep the warm retention it has
				Config: testAccCheckLogzioWarmTierSubAccountConfig(email, accountNameUpdate, accountId, 2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("logzio_subaccount.test_subaccount", subAccountName, accountNameUpdate),
					testAccCheckWarmRetention(resourceName, 2),
				),
			},
			{
				Config:      testAccCheckLogzioWarmTierSubAccountConfig(email, accountNameUpdate, accountId, 1000),
				ExpectError: regexp.MustCompile("INVALID_WARM_RETENTION"),
			},
			{
				// 0 turns warm tier off for a sub account
				Config: testAccCheckLogzioWarmTierSubAccountConfig(email, accountNameUpdate, accountId, 0),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, warmTierSnapSearchRetentionDays, "0"),
					testAccCheckWarmRetention(resourceName, 0),
				),
			},
			{
				Config:            testAccCheckLogzioWarmTierSubAccountConfig(email, accountNameUpdate, accountId, 0),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Sets the main account's warm retention to the value it already has, so the account doesn't change. Destroy must
// leave it as it is.
func TestAccLogzioWarmTier_MainAccount(t *testing.T) {
	accountId := os.Getenv(envLogzioWarmAccountId)
	resourceName := "logzio_warm_tier.test_warm_tier"
	dataSourceName := "data.logzio_subaccount.main_account"
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckApiTokenWarm(t)
			testAccPreCheckWarmAccountId(t)
		},
		ProviderFactories: testAccWarmProviderFactories,
		CheckDestroy:      testAccCheckWarmTierKept,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLogzioWarmTierMainAccountConfig(accountId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, subAccountIsOwner, "true"),
					resource.TestCheckResourceAttr(resourceName, warmTierAccountId, accountId),
					resource.TestCheckResourceAttrPair(resourceName, warmTierSnapSearchRetentionDays, dataSourceName, subAccountsSnapSearchRetentionDays),
				),
			},
			{
				Config:            testAccCheckLogzioWarmTierMainAccountConfig(accountId),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Fails at plan time, before any API call.
func TestAccLogzioWarmTier_NegativeRetention(t *testing.T) {
	terraformPlan := `
resource "logzio_warm_tier" "test_warm_tier" {
  account_id = 1
  snap_search_retention_days = -1
}
`

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckApiToken(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      terraformPlan,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("expected snap_search_retention_days to be at least"),
			},
		},
	})
}

func warmTierTestClient() (*sub_accounts.SubAccountClient, error) {
	url := os.Getenv(envLogzioCustomApiUrl)
	if url == "" {
		url = fmt.Sprintf(baseUrl, "")
	}
	return sub_accounts.New(os.Getenv(envLogzioApiTokenWarm), url)
}

// testAccCheckWarmRetention reads the account behind the warm tier resource from the API.
func testAccCheckWarmRetention(resourceName string, expected int32) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("%s not found in the state", resourceName)
		}
		id, err := strconv.ParseInt(r.Primary.ID, 10, 64)
		if err != nil {
			return err
		}
		client, err := warmTierTestClient()
		if err != nil {
			return err
		}
		account, err := client.GetSubAccount(id)
		if err != nil {
			return err
		}
		if account.SnapSearchRetentionDays != expected {
			return fmt.Errorf("account %d has warm retention %d, expected %d", id, account.SnapSearchRetentionDays, expected)
		}
		return nil
	}
}

// testAccCheckWarmTierKept checks that destroy left the main account's warm tier on.
func testAccCheckWarmTierKept(s *terraform.State) error {
	client, err := warmTierTestClient()
	if err != nil {
		return err
	}
	for _, r := range s.RootModule().Resources {
		if r.Type != resourceWarmTierType {
			continue
		}
		id, err := strconv.ParseInt(r.Primary.ID, 10, 64)
		if err != nil {
			return err
		}
		account, err := client.GetSubAccount(id)
		if err != nil {
			return err
		}
		if account.SnapSearchRetentionDays <= 0 {
			return fmt.Errorf("account %d has no warm tier after destroy", id)
		}
	}
	return nil
}

// The sub account takes the main account's volume mode, so the test fits a flexible or a fixed warm test plan. Its hot
// retention is 4, which the sub account update needs when warm tier is set.
func testAccCheckLogzioWarmTierSubAccountConfig(email string, accountName string, accountId string, snapRetention int) string {
	return fmt.Sprintf(`
data "logzio_subaccount" "main_account" {
  account_id = %s
}

resource "logzio_subaccount" "test_subaccount" {
  email = "%s"
  account_name = "%s"
  retention_days = 4
  utilization_enabled = "true"
  flexible = data.logzio_subaccount.main_account.flexible
  max_daily_gb = 1
  reserved_daily_gb = data.logzio_subaccount.main_account.flexible ? 0.5 : null
  sharing_objects_accounts = [
    %s
  ]
}

resource "logzio_warm_tier" "test_warm_tier" {
  account_id = logzio_subaccount.test_subaccount.account_id
  snap_search_retention_days = %d
}
`, accountId, email, accountName, accountId, snapRetention)
}

func testAccCheckLogzioWarmTierMainAccountConfig(accountId string) string {
	return fmt.Sprintf(`
data "logzio_subaccount" "main_account" {
  account_id = %s
}

resource "logzio_warm_tier" "test_warm_tier" {
  account_id = %s
  snap_search_retention_days = data.logzio_subaccount.main_account.snap_search_retention_days
}
`, accountId, accountId)
}
