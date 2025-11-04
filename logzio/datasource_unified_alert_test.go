package logzio

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceUnifiedAlert_LogAlert(t *testing.T) {
	email := "test@logz.io"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getUnifiedLogAlertDataSourceConfig(email),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.logzio_unified_alert.test_log_alert_ds", "alert_id"),
					resource.TestCheckResourceAttr("data.logzio_unified_alert.test_log_alert_ds", "title", "Test Log Alert for DS"),
					resource.TestCheckResourceAttr("data.logzio_unified_alert.test_log_alert_ds", "type", "LOG_ALERT"),
					resource.TestCheckResourceAttr("data.logzio_unified_alert.test_log_alert_ds", "enabled", "true"),
				),
			},
		},
	})
}

func getUnifiedLogAlertDataSourceConfig(email string) string {
	return fmt.Sprintf(`
resource "logzio_unified_alert" "test_log_alert_resource" {
  title       = "Test Log Alert for DS"
  type        = "LOG_ALERT"
  description = "Test log alert for data source"
  enabled     = true

  log_alert {
    search_timeframe_minutes = 5

    output {
      type = "JSON"

      recipients {
        emails = ["%s"]
      }
    }

    sub_components {
      query_definition {
        query                        = "level:ERROR"
        should_query_on_all_accounts = true

        aggregation {
          aggregation_type = "COUNT"
        }
      }

      trigger {
        operator = "GREATER_THAN"

        severity_threshold_tiers {
          severity  = "HIGH"
          threshold = 10
        }
      }

      output {
        should_use_all_fields = true
      }
    }
  }
}

data "logzio_unified_alert" "test_log_alert_ds" {
  type     = "LOG_ALERT"
  alert_id = logzio_unified_alert.test_log_alert_resource.alert_id

  depends_on = [logzio_unified_alert.test_log_alert_resource]
}
`, email)
}

