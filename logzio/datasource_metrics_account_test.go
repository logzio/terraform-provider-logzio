package logzio

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
	"os"
	"regexp"
	"strconv"
	"testing"
)

func TestAccDataSourceMetricsAccount(t *testing.T) {
	resourceName := "logzio_metrics_account.metrics_account_datasource"
	dataSourceName := "data.logzio_metrics_account.metrics_account_datasource_by_id"
	accountId, _ := strconv.ParseInt(os.Getenv(envLogzioAccountId), utils.BASE_10, utils.BITSIZE_64)
	email := os.Getenv(envLogzioEmail)
	accountName := "test_datasource_create"
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckApiToken(t)
			testAccPreCheckEmail(t)
			testAccPreCheckAccountId(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMetricsAccountDataSourceResource(email, accountId, accountName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_name", accountName),
					resource.TestCheckResourceAttr(resourceName, "plan_uts", "100"),
				),
			},
			{
				Config: testAccMetricsAccountDataSourceResource(email, accountId, accountName) +
					testAccCheckLogzioMetricsAccountDatasourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "account_name", accountName),
					resource.TestCheckResourceAttr(dataSourceName, "plan_uts", "100"),
				),
			},
		},
	})
}

func TestAccDataSourceMetricsAccountByAccountName(t *testing.T) {
	resourceName := "logzio_metrics_account.metrics_account_datasource"
	dataSourceName := "data.logzio_metrics_account.metrics_account_datasource_by_id"
	accountId, _ := strconv.ParseInt(os.Getenv(envLogzioAccountId), utils.BASE_10, utils.BITSIZE_64)
	email := os.Getenv(envLogzioEmail)
	accountName := "test_datasource_account_name"
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckApiToken(t)
			testAccPreCheckEmail(t)
			testAccPreCheckAccountId(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMetricsAccountDataSourceResource(email, accountId, accountName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_name", accountName),
					resource.TestCheckResourceAttr(resourceName, "plan_uts", "100"),
				),
			},
			{
				Config: testAccMetricsAccountDataSourceResource(email, accountId, accountName) +
					testAccCheckLogzioMetricsAccountDatasourceConfigAccountName(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "account_name", accountName),
					resource.TestCheckResourceAttr(dataSourceName, "plan_uts", "100"),
				),
			},
		},
	})
}

func TestAccDataSourceMetricsAccountNotExists(t *testing.T) {
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckLogzioMetricsAccountDatasourceConfigNotExist(),
				ExpectError: regexp.MustCompile("couldn't find metrics account with specified attributes"),
			},
		},
	})
}

func testAccMetricsAccountDataSourceResource(email string, accountId int64, accountName string) string {
	return fmt.Sprintf(`resource "logzio_metrics_account" "metrics_account_datasource" {
  email = "%s"
  account_name = "%s"
  plan_uts = 100
  authorized_accounts = [
    %d
  ]
}
`, email, accountName, accountId)
}

func testAccCheckLogzioMetricsAccountDatasourceConfig() string {
	return fmt.Sprint(`data "logzio_metrics_account" "metrics_account_datasource_by_id" {
  account_id = "${logzio_metrics_account.metrics_account_datasource.Id}"
}
`)
}

func testAccCheckLogzioMetricsAccountDatasourceConfigAccountName() string {
	return fmt.Sprint(`data "logzio_metrics_account" "metrics_account_datasource_by_id" {
  account_name = "${logzio_metrics_account.metrics_account_datasource.account_name}"
}
`)
}

func testAccCheckLogzioMetricsAccountDatasourceConfigNotExist() string {
	return fmt.Sprint(`data "logzio_metrics_account" "metrics_account_datasource_by_id" {
  account_name = "name_not_exist"
}
`)
}
