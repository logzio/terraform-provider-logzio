package logzio

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
)

const (
	envLogzioEmail = "LOGZIO_EMAIL"
)

func TestAccDataSourceSubaccount(t *testing.T) {
	resourceName := "logzio_subaccount.subaccount_datasource"
	dataSourceName := "data.logzio_subaccount.subaccount_datasource_by_id"
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
				Config: testAccSubAccountDataSourceResource(email, accountId, accountName),
				Check: resource.ComposeAggregateTestCheckFunc(
					awaitApply(15),
					resource.TestCheckResourceAttr(resourceName, "account_name", accountName),
					resource.TestCheckResourceAttr(resourceName, "retention_days", "2"),
					resource.TestCheckResourceAttr(resourceName, "max_daily_gb", "1"),
				),
			},
			{
				Config: testAccSubAccountDataSourceResource(email, accountId, accountName) +
					testAccCheckLogzioSubaccountDatasourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					awaitApply(15),
					resource.TestCheckResourceAttr(dataSourceName, "account_name", accountName),
					resource.TestCheckResourceAttr(dataSourceName, "retention_days", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "max_daily_gb", "1"),
				),
			},
		},
	})
}

func TestAccDataSourceSubaccountByAccountName(t *testing.T) {
	resourceName := "logzio_subaccount.subaccount_datasource"
	dataSourceName := "data.logzio_subaccount.subaccount_datasource_by_id"
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
				Config: testAccSubAccountDataSourceResource(email, accountId, accountName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_name", accountName),
					resource.TestCheckResourceAttr(resourceName, "retention_days", "2"),
					resource.TestCheckResourceAttr(resourceName, "max_daily_gb", "1"),
				),
			},
			{
				Config: testAccSubAccountDataSourceResource(email, accountId, accountName) +
					testAccCheckLogzioSubaccountDatasourceConfigAccountName(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "account_name", accountName),
					resource.TestCheckResourceAttr(dataSourceName, "retention_days", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "max_daily_gb", "1"),
				),
			},
		},
	})
}

func TestAccDataSourceSubaccountWarm(t *testing.T) {
	resourceName := "logzio_subaccount.subaccount_datasource"
	dataSourceName := "data.logzio_subaccount.subaccount_datasource_by_id"
	accountId, _ := strconv.ParseInt(os.Getenv(envLogzioWarmAccountId), utils.BASE_10, utils.BITSIZE_64)
	email := os.Getenv(envLogzioEmail)
	accountName := "test_datasource_create"
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckApiTokenWarm(t)
			testAccPreCheckEmail(t)
			testAccPreCheckWarmAccountId(t)
		},
		ProviderFactories: testAccWarmProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSubAccountWarmDataSourceResource(email, accountId, accountName),
				Check: resource.ComposeAggregateTestCheckFunc(
					awaitApply(15),
					resource.TestCheckResourceAttr(resourceName, "account_name", accountName),
					resource.TestCheckResourceAttr(resourceName, "snap_search_retention_days", "2"),
					resource.TestCheckResourceAttr(resourceName, "is_capped", "true"),
					resource.TestCheckResourceAttr(resourceName, "shared_gb", "9"),
					resource.TestCheckResourceAttr(resourceName, "total_time_based_daily_gb", "10"),
					resource.TestCheckResourceAttr(resourceName, "is_owner", "false"),
				),
			},
			{
				Config: testAccSubAccountWarmDataSourceResource(email, accountId, accountName) +
					testAccCheckLogzioSubaccountDatasourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					awaitApply(15),
					resource.TestCheckResourceAttr(dataSourceName, "account_name", accountName),
					resource.TestCheckResourceAttr(dataSourceName, "snap_search_retention_days", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "is_capped", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "shared_gb", "9"),
					resource.TestCheckResourceAttr(dataSourceName, "total_time_based_daily_gb", "10"),
					resource.TestCheckResourceAttr(dataSourceName, "is_owner", "false"),
				),
			},
		},
	})
}

func TestAccDataSourceSubaccountConsumption(t *testing.T) {
	resourceName := "logzio_subaccount.subaccount_datasource"
	dataSourceName := "data.logzio_subaccount.subaccount_datasource_by_id"
	accountId, _ := strconv.ParseInt(os.Getenv(envLogzioConsumptionAccountId), utils.BASE_10, utils.BITSIZE_64)
	email := os.Getenv(envLogzioEmail)
	accountName := "test_datasource_create"
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckApiTokenConsumption(t)
			testAccPreCheckEmail(t)
			testAccPreCheckConsumptionAccountId(t)
		},
		ProviderFactories: testAccConsumptionProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSubAccountConsumptionDataSourceResource(email, accountId, accountName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_name", accountName),
					resource.TestCheckResourceAttr(resourceName, "soft_limit_gb", "1"),
				),
			},
			{
				Config: testAccSubAccountConsumptionDataSourceResource(email, accountId, accountName) +
					testAccCheckLogzioSubaccountDatasourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "account_name", accountName),
					resource.TestCheckResourceAttr(dataSourceName, "soft_limit_gb", "1"),
				),
			},
		},
	})
}

func TestAccDataSourceSubaccountNotExists(t *testing.T) {
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckLogzioSubaccountDatasourceConfigNotExist(),
				ExpectError: regexp.MustCompile("couldn't find sub-account with specified attributes"),
			},
		},
	})
}

func testAccSubAccountDataSourceResource(email string, accountId int64, accountName string) string {
	return fmt.Sprintf(`resource "logzio_subaccount" "subaccount_datasource" {
  email = "%s"
  account_name = "%s"
  retention_days = 2
  max_daily_gb = 1
  sharing_objects_accounts = [
    %d
  ]
}
`, email, accountName, accountId)
}

func testAccCheckLogzioSubaccountDatasourceConfig() string {
	return fmt.Sprint(`data "logzio_subaccount" "subaccount_datasource_by_id" {
  account_id = "${logzio_subaccount.subaccount_datasource.id}"
}
`)
}

func testAccCheckLogzioSubaccountDatasourceConfigAccountName() string {
	return fmt.Sprint(`data "logzio_subaccount" "subaccount_datasource_by_id" {
  account_name = "${logzio_subaccount.subaccount_datasource.account_name}"
}
`)
}

func testAccCheckLogzioSubaccountDatasourceConfigNotExist() string {
	return fmt.Sprint(`data "logzio_subaccount" "subaccount_datasource_by_id" {
  account_name = "name_not_exist"
}
`)
}

func testAccSubAccountWarmDataSourceResource(email string, accountId int64, accountName string) string {
	return fmt.Sprintf(`resource "logzio_subaccount" "subaccount_datasource" {
  email = "%s"
  account_name = "%s"
  retention_days = 4
  max_daily_gb = 1
  reserved_daily_gb = 0.5
  sharing_objects_accounts = [
    %d
  ]
  flexible = "true"
  snap_search_retention_days = 2
}
`, email, accountName, accountId)
}

func testAccSubAccountConsumptionDataSourceResource(email string, accountId int64, accountName string) string {
	return fmt.Sprintf(`resource "logzio_subaccount" "subaccount_datasource" {
  email = "%s"
  account_name = "%s"
  retention_days = 4
  max_daily_gb = 1
  sharing_objects_accounts = [
    %d
  ]
  soft_limit_gb = 1
}
`, email, accountName, accountId)
}
