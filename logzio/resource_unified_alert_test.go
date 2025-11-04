package logzio

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccLogzioUnifiedAlert_LogAlert(t *testing.T) {
	email := "test@logz.io"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getUnifiedLogAlertConfig(email),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("logzio_unified_alert.test_log_alert", "alert_id"),
					resource.TestCheckResourceAttr("logzio_unified_alert.test_log_alert", "title", "Test Log Alert"),
					resource.TestCheckResourceAttr("logzio_unified_alert.test_log_alert", "type", "LOG_ALERT"),
					resource.TestCheckResourceAttr("logzio_unified_alert.test_log_alert", "enabled", "true"),
				),
			},
			{
				Config: getUnifiedLogAlertConfigUpdated(email),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("logzio_unified_alert.test_log_alert", "alert_id"),
					resource.TestCheckResourceAttr("logzio_unified_alert.test_log_alert", "title", "Test Log Alert Updated"),
					resource.TestCheckResourceAttr("logzio_unified_alert.test_log_alert", "enabled", "false"),
				),
			},
		},
	})
}

func TestAccLogzioUnifiedAlert_MetricAlert(t *testing.T) {
	email := "test@logz.io"
	datasourceUid := os.Getenv("GRAFANA_DATASOURCE_UID")
	if datasourceUid == "" {
		t.Skip("GRAFANA_DATASOURCE_UID must be set for metric alert tests")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getUnifiedMetricAlertConfig(email, datasourceUid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("logzio_unified_alert.test_metric_alert", "alert_id"),
					resource.TestCheckResourceAttr("logzio_unified_alert.test_metric_alert", "title", "Test Metric Alert"),
					resource.TestCheckResourceAttr("logzio_unified_alert.test_metric_alert", "type", "METRIC_ALERT"),
					resource.TestCheckResourceAttr("logzio_unified_alert.test_metric_alert", "enabled", "true"),
				),
			},
		},
	})
}

func TestAccLogzioUnifiedAlert_MetricAlertMathExpression(t *testing.T) {
	email := "test@logz.io"
	datasourceUid := os.Getenv("GRAFANA_DATASOURCE_UID")
	if datasourceUid == "" {
		t.Skip("GRAFANA_DATASOURCE_UID must be set for metric alert tests")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getUnifiedMetricAlertMathConfig(email, datasourceUid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("logzio_unified_alert.test_metric_math_alert", "alert_id"),
					resource.TestCheckResourceAttr("logzio_unified_alert.test_metric_math_alert", "title", "Test Math Expression Alert"),
					resource.TestCheckResourceAttr("logzio_unified_alert.test_metric_math_alert", "type", "METRIC_ALERT"),
				),
			},
		},
	})
}

func testCheckUnifiedAlertDestroy(s *terraform.State) error {
	client := unifiedAlertClient(testAccProvider.Meta())

	for _, r := range s.RootModule().Resources {
		if r.Type != "logzio_unified_alert" {
			continue
		}

		alertId := r.Primary.ID
		alertType := r.Primary.Attributes["type"]
		urlType := getUrlTypeFromAlertType(alertType)

		_, err := client.GetUnifiedAlert(urlType, alertId)
		if err == nil {
			return fmt.Errorf("alert %s still exists", alertId)
		}
	}

	return nil
}

func getUnifiedLogAlertConfig(email string) string {
	return fmt.Sprintf(`
resource "logzio_unified_alert" "test_log_alert" {
  title       = "Test Log Alert"
  type        = "LOG_ALERT"
  description = "Test log alert description"
  tags        = ["test", "terraform"]
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
`, email)
}

func getUnifiedLogAlertConfigUpdated(email string) string {
	return fmt.Sprintf(`
resource "logzio_unified_alert" "test_log_alert" {
  title       = "Test Log Alert Updated"
  type        = "LOG_ALERT"
  description = "Test log alert description updated"
  tags        = ["test", "terraform", "updated"]
  enabled     = false

  log_alert {
    search_timeframe_minutes = 10

    output {
      type = "TABLE"

      recipients {
        emails = ["%s"]
      }
    }

    sub_components {
      query_definition {
        query                        = "level:ERROR OR level:WARN"
        should_query_on_all_accounts = true

        aggregation {
          aggregation_type = "COUNT"
        }
      }

      trigger {
        operator = "GREATER_THAN"

        severity_threshold_tiers {
          severity  = "MEDIUM"
          threshold = 20
        }
      }

      output {
        should_use_all_fields = true
      }
    }
  }
}
`, email)
}

func getUnifiedMetricAlertConfig(email, datasourceUid string) string {
	return fmt.Sprintf(`
resource "logzio_unified_alert" "test_metric_alert" {
  title       = "Test Metric Alert"
  type        = "METRIC_ALERT"
  description = "Test metric alert description"
  tags        = ["test", "terraform", "metrics"]
  enabled     = true

  metric_alert {
    severity = "HIGH"

    trigger {
      trigger_type             = "THRESHOLD"
      metric_operator          = "ABOVE"
      min_threshold            = 80.0
      search_timeframe_minutes = 5
    }

    queries {
      ref_id = "A"

      query_definition {
        datasource_uid = "%s"
        promql_query   = "avg(cpu_usage)"
      }
    }

    recipients {
      emails = ["%s"]
    }
  }
}
`, datasourceUid, email)
}

func getUnifiedMetricAlertMathConfig(email, datasourceUid string) string {
	return fmt.Sprintf(`
resource "logzio_unified_alert" "test_metric_math_alert" {
  title       = "Test Math Expression Alert"
  type        = "METRIC_ALERT"
  description = "Test metric alert with math expression"
  tags        = ["test", "terraform", "math"]
  enabled     = true

  metric_alert {
    severity = "MEDIUM"

    trigger {
      trigger_type             = "MATH"
      math_expression          = "($A / $B) * 100"
      metric_operator          = "ABOVE"
      min_threshold            = 5.0
      search_timeframe_minutes = 5
    }

    queries {
      ref_id = "A"

      query_definition {
        datasource_uid = "%s"
        promql_query   = "sum(errors)"
      }
    }

    queries {
      ref_id = "B"

      query_definition {
        datasource_uid = "%s"
        promql_query   = "sum(requests)"
      }
    }

    recipients {
      emails = ["%s"]
    }
  }
}
`, datasourceUid, datasourceUid, email)
}
