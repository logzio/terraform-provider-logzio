package logzio

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"regexp"
	"testing"
)

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
