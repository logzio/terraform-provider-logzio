package logzio

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"os"
	"regexp"
	"strconv"
	"testing"
)

const (
	envLogzioEmail = "LOGZIO_EMAIL"
)

func TestAccDataSourceSubaccount(t *testing.T) {
	resourceName := "logzio_subaccount.subaccount_datasource"
	dataSourceName := "data.logzio_subaccount.subaccount_datasource_by_id"
	accountId, _ := strconv.ParseInt(os.Getenv(envLogzioAccountId), BASE_10, BITSIZE_64)
	email := os.Getenv(envLogzioEmail)
	accountName := "test_datasource"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckApiToken(t)
			testAccPreCheckEmail(t)
			testAccPreCheckAccountId(t)
		},
		Providers: testAccProviders,
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
					testAccCheckLogzioSubaccountDatasourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
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
	accountId, _ := strconv.ParseInt(os.Getenv(envLogzioAccountId), BASE_10, BITSIZE_64)
	email := os.Getenv(envLogzioEmail)
	accountName := "test_datasource_account_name"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckApiToken(t)
			testAccPreCheckEmail(t)
			testAccPreCheckAccountId(t)
		},
		Providers: testAccProviders,
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

func TestAccDataSourceSubaccountNotExists(t *testing.T) {
	resourceName := "logzio_subaccount.subaccount_datasource"
	accountId, _ := strconv.ParseInt(os.Getenv(envLogzioAccountId), BASE_10, BITSIZE_64)
	email := os.Getenv(envLogzioEmail)
	accountName := "some_account_to_add"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckApiToken(t)
			testAccPreCheckEmail(t)
			testAccPreCheckAccountId(t)
		},
		Providers: testAccProviders,
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
