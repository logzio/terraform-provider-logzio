package logzio

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"os"
	"strconv"
	"testing"
)

const (
	envLogzioEmail = "LOGZIO_EMAIL"
)

func TestAccDataSourceSubaccount(t *testing.T) {
	resourceName := "data.logzio_subaccount.subaccount_datasource_by_id"
	if v := os.Getenv(envLogzioAccountId); v == "" {
		t.Log(v)
		t.Fatalf("%s must be set for acceptance tests", envLogzioAccountId)
	}
	accountId, _ := strconv.ParseInt(os.Getenv(envLogzioAccountId), BASE_10, BITSIZE_64)

	if v := os.Getenv(envLogzioEmail); v == "" {
		t.Log(v)
		t.Fatalf("%s must be set for acceptance tests", envLogzioEmail)
	}
	email := os.Getenv(envLogzioEmail)
	terraformPlan := testAccCheckLogzioSubaccountDatasourceConfig(email, accountId)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				ExpectNonEmptyPlan: true,
				Config: terraformPlan,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_name", "test"),
					resource.TestCheckResourceAttr(resourceName, "retention_days", "2"),
					resource.TestCheckResourceAttr(resourceName, "max_daily_gb", "1"),
				),
			},
		},
	})
}

func testAccCheckLogzioSubaccountDatasourceConfig(email string, accountId int64) string {
	return fmt.Sprintf(`
resource "logzio_subaccount" "subaccount_datasource" {
  email = "%s"
  account_name = "test"
  retention_days = 2
  max_daily_gb = 1
  sharing_objects_accounts = [
    %d
  ]
}

data "logzio_subaccount" "subaccount_datasource_by_id" {
  account_id = logzio_subaccount.subaccount_datasource.id
  depends_on = ["logzio_subaccount.subaccount_datasource"]
}
`, email, accountId)
}
