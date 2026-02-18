package logzio

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
)

func TestAccLogzioUnifiedAlert_LogAlert(t *testing.T) {
	defer utils.SleepAfterTest()
	email := "test@logz.io"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testCheckUnifiedAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: getUnifiedLogAlertConfig(email),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("logzio_unified_alert.test_log_alert", "alert_id"),
					resource.TestCheckResourceAttr("logzio_unified_alert.test_log_alert", "title", "Test Log Alert"),
					resource.TestCheckResourceAttr("logzio_unified_alert.test_log_alert", "enabled", "true"),
					resource.TestCheckResourceAttr("logzio_unified_alert.test_log_alert", "alert_configuration.0.type", "LOG_ALERT"),
					resource.TestCheckResourceAttr("logzio_unified_alert.test_log_alert", "alert_configuration.0.search_timeframe_minutes", "5"),
				),
			},
			{
				Config: getUnifiedLogAlertConfigUpdated(email),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("logzio_unified_alert.test_log_alert", "alert_id"),
					resource.TestCheckResourceAttr("logzio_unified_alert.test_log_alert", "title", "Test Log Alert Updated"),
					resource.TestCheckResourceAttr("logzio_unified_alert.test_log_alert", "enabled", "false"),
					resource.TestCheckResourceAttr("logzio_unified_alert.test_log_alert", "alert_configuration.0.search_timeframe_minutes", "10"),
				),
			},
		},
	})
}

func TestAccLogzioUnifiedAlert_MetricAlert(t *testing.T) {
	defer utils.SleepAfterTest()
	email := "test@logz.io"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testCheckUnifiedAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: getUnifiedMetricAlertConfig(email),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("logzio_unified_alert.test_metric_alert", "alert_id"),
					resource.TestCheckResourceAttr("logzio_unified_alert.test_metric_alert", "title", "Test Metric Alert"),
					resource.TestCheckResourceAttr("logzio_unified_alert.test_metric_alert", "enabled", "true"),
					resource.TestCheckResourceAttr("logzio_unified_alert.test_metric_alert", "alert_configuration.0.type", "METRIC_ALERT"),
					resource.TestCheckResourceAttr("logzio_unified_alert.test_metric_alert", "alert_configuration.0.severity", "HIGH"),
				),
			},
		},
	})
}

func TestAccLogzioUnifiedAlert_MetricAlertMathExpression(t *testing.T) {
	defer utils.SleepAfterTest()
	email := "test@logz.io"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testCheckUnifiedAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: getUnifiedMetricAlertMathConfig(email),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("logzio_unified_alert.test_metric_math_alert", "alert_id"),
					resource.TestCheckResourceAttr("logzio_unified_alert.test_metric_math_alert", "title", "Test Math Expression Alert"),
					resource.TestCheckResourceAttr("logzio_unified_alert.test_metric_math_alert", "alert_configuration.0.type", "METRIC_ALERT"),
					resource.TestCheckResourceAttr("logzio_unified_alert.test_metric_math_alert", "alert_configuration.0.trigger.0.type", "math"),
				),
			},
		},
	})
}

func TestAccLogzioUnifiedAlert_Import(t *testing.T) {
	defer utils.SleepAfterTest()
	email := "test@logz.io"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testCheckUnifiedAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: getUnifiedLogAlertConfig(email),
			},
			{
				ResourceName:      "logzio_unified_alert.test_log_alert",
				ImportState:       true,
				ImportStateVerify: true,
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

		compositeId := r.Primary.ID
		urlType, alertId, err := parseUnifiedAlertCompositeId(compositeId)
		if err != nil {
			return fmt.Errorf("failed to parse composite ID %s: %v", compositeId, err)
		}

		_, err = client.GetUnifiedAlert(urlType, alertId)
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
  description = "Test log alert description"
  tags        = ["test", "terraform"]
  enabled     = true

  recipients {
    emails = ["%s"]
  }

  alert_configuration {
    type                           = "LOG_ALERT"
    search_timeframe_minutes       = 5
    alert_output_template_type     = "JSON"

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
  description = "Test log alert description updated"
  tags        = ["test", "terraform", "updated"]
  enabled     = false

  recipients {
    emails = ["%s"]
  }

  alert_configuration {
    type                           = "LOG_ALERT"
    search_timeframe_minutes       = 10
    alert_output_template_type     = "TEXT"

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

func getUnifiedMetricAlertConfig(email string) string {
	return fmt.Sprintf(`
resource "logzio_unified_alert" "test_metric_alert" {
  title       = "Test Metric Alert"
  description = "Test metric alert description"
  tags        = ["test", "terraform", "metrics"]
  enabled     = true

  recipients {
    emails = ["%s"]
  }

  alert_configuration {
    type     = "METRIC_ALERT"
    severity = "HIGH"

    trigger {
      type = "threshold"

      condition {
        operator_type = "above"
        threshold     = 80.0
      }
    }

    queries {
      ref_id = "A"

      query_definition {
        promql_query = "avg(cpu_usage)"
      }
    }
  }
}
`, email)
}

func getUnifiedMetricAlertMathConfig(email string) string {
	return fmt.Sprintf(`
resource "logzio_unified_alert" "test_metric_math_alert" {
  title       = "Test Math Expression Alert"
  description = "Test metric alert with math expression"
  tags        = ["test", "terraform", "math"]
  enabled     = true

  recipients {
    emails = ["%s"]
  }

  alert_configuration {
    type     = "METRIC_ALERT"
    severity = "MEDIUM"

    trigger {
      type       = "math"
      expression = "($A / $B) * 100"
    }

    queries {
      ref_id = "A"

      query_definition {
        promql_query = "sum(errors)"
      }
    }

    queries {
      ref_id = "B"

      query_definition {
        promql_query = "sum(requests)"
      }
    }
  }
}
`, email)
}
