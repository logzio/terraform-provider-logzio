package logzio

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
)

const (
	metricsRollupRulesResourceCreateSimple                       = "create_metrics_rollup_rules_simple"
	metricsRollupRulesResourceCreateComplex                      = "create_metrics_rollup_rules_complex"
	metricsRollupRulesResourceCreateInvalidMetricType            = "create_metrics_rollup_rules_invalid_metric_type"
	metricsRollupRulesResourceCreateInvalidRollupFunction        = "create_metrics_rollup_rules_invalid_rollup_function"
	metricsRollupRulesResourceCreateInvalidEliminationMethod     = "create_metrics_rollup_rules_invalid_elimination_method"
	metricsRollupRulesResourceCreateEmptyMetricName              = "create_metrics_rollup_rules_empty_metric_name"
	metricsRollupRulesResourceCreateEmptyLabels                  = "create_metrics_rollup_rules_empty_labels"
	metricsRollupRulesResourceUpdate                             = "update_metrics_rollup_rules"
	metricsRollupRulesResourceCreateWithFilter                   = "create_metrics_rollup_rules_with_filter"
	metricsRollupRulesResourceCreateMeasurement                  = "create_metrics_rollup_rules_measurement"
	metricsRollupRulesResourceCreateMeasurementP99               = "create_metrics_rollup_rules_measurement_p99"
	metricsRollupRulesResourceCreateCounterWithRollupFunction    = "create_metrics_rollup_rules_counter_with_rollup_function"
	metricsRollupRulesResourceCreateCounterMissingRollupFunction = "create_metrics_rollup_rules_counter_missing_rollup_function"
)

func TestAccLogzioMetricsRollupRules_CreateSimple(t *testing.T) {
	resourceName := "test_create_metrics_rollup_rules_simple"
	resourceFullName := "logzio_metrics_rollup_rules." + resourceName
	accountId := os.Getenv(envLogzioMetricsAccountId)

	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckApiToken(t)
			testAccPreCheckMetricsAccountId(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceTestMetricsRollupRules(resourceName, metricsRollupRulesResourceCreateSimple, accountId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesAccountId, accountId),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesMetricName, "cpu_usage"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesMetricType, "GAUGE"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesRollupFunction, "LAST"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesLabelsEliminationMethod, "EXCLUDE_BY"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesLabels+".#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesLabels+".0", "instance_id"),
				),
			},
			{
				Config:            resourceTestMetricsRollupRules(resourceName, metricsRollupRulesResourceCreateSimple, accountId),
				ResourceName:      resourceFullName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLogzioMetricsRollupRules_CreateComplex(t *testing.T) {
	resourceName := "test_create_metrics_rollup_rules_complex"
	resourceFullName := "logzio_metrics_rollup_rules." + resourceName
	accountId := os.Getenv(envLogzioMetricsAccountId)

	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckApiToken(t)
			testAccPreCheckMetricsAccountId(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceTestMetricsRollupRules(resourceName, metricsRollupRulesResourceCreateComplex, accountId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesAccountId, accountId),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesMetricName, "http_requests_total"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesMetricType, "COUNTER"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesRollupFunction, "SUM"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesLabelsEliminationMethod, "EXCLUDE_BY"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesLabels+".#", "3"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesLabels+".0", "path"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesLabels+".1", "method"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesLabels+".2", "user_agent"),
				),
			},
			{
				Config:            resourceTestMetricsRollupRules(resourceName, metricsRollupRulesResourceCreateComplex, accountId),
				ResourceName:      resourceFullName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLogzioMetricsRollupRules_CreateInvalidMetricType(t *testing.T) {
	resourceName := "test_create_metrics_rollup_rules_invalid_metric_type"
	accountId := os.Getenv(envLogzioMetricsAccountId)

	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckApiToken(t)
			testAccPreCheckMetricsAccountId(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      resourceTestMetricsRollupRules(resourceName, metricsRollupRulesResourceCreateInvalidMetricType, accountId),
				ExpectError: regexp.MustCompile("expected metric_type to be one of"),
			},
		},
	})
}

func TestAccLogzioMetricsRollupRules_CreateInvalidRollupFunction(t *testing.T) {
	resourceName := "test_create_metrics_rollup_rules_invalid_rollup_function"
	accountId := os.Getenv(envLogzioMetricsAccountId)

	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckApiToken(t)
			testAccPreCheckMetricsAccountId(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      resourceTestMetricsRollupRules(resourceName, metricsRollupRulesResourceCreateInvalidRollupFunction, accountId),
				ExpectError: regexp.MustCompile("expected rollup_function to be one of"),
			},
		},
	})
}

func TestAccLogzioMetricsRollupRules_CreateInvalidEliminationMethod(t *testing.T) {
	resourceName := "test_create_metrics_rollup_rules_invalid_elimination_method"
	accountId := os.Getenv(envLogzioMetricsAccountId)

	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckApiToken(t)
			testAccPreCheckMetricsAccountId(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      resourceTestMetricsRollupRules(resourceName, metricsRollupRulesResourceCreateInvalidEliminationMethod, accountId),
				ExpectError: regexp.MustCompile("expected labels_elimination_method to be one of"),
			},
		},
	})
}

func TestAccLogzioMetricsRollupRules_CreateEmptyMetricName(t *testing.T) {
	resourceName := "test_create_metrics_rollup_rules_empty_metric_name"
	accountId := os.Getenv(envLogzioMetricsAccountId)

	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckApiToken(t)
			testAccPreCheckMetricsAccountId(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      resourceTestMetricsRollupRules(resourceName, metricsRollupRulesResourceCreateEmptyMetricName, accountId),
				ExpectError: regexp.MustCompile("one of `filter,metric_name` must be specified"),
			},
		},
	})
}

func TestAccLogzioMetricsRollupRules_CreateEmptyLabels(t *testing.T) {
	resourceName := "test_create_metrics_rollup_rules_empty_labels"
	accountId := os.Getenv(envLogzioMetricsAccountId)

	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckApiToken(t)
			testAccPreCheckMetricsAccountId(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      resourceTestMetricsRollupRules(resourceName, metricsRollupRulesResourceCreateEmptyLabels, accountId),
				ExpectError: regexp.MustCompile("Attribute labels requires 1 item minimum"),
			},
		},
	})
}

func TestAccLogzioMetricsRollupRules_Update(t *testing.T) {
	resourceName := "test_update_metrics_rollup_rules"
	resourceFullName := "logzio_metrics_rollup_rules." + resourceName
	accountId := os.Getenv(envLogzioMetricsAccountId)

	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckApiToken(t)
			testAccPreCheckMetricsAccountId(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceTestMetricsRollupRules(resourceName, metricsRollupRulesResourceCreateSimple, accountId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesAccountId, accountId),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesMetricName, "cpu_usage"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesMetricType, "GAUGE"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesRollupFunction, "LAST"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesLabelsEliminationMethod, "EXCLUDE_BY"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesLabels+".#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesLabels+".0", "instance_id"),
				),
			},
			{
				Config: resourceTestMetricsRollupRules(resourceName, metricsRollupRulesResourceUpdate, accountId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesAccountId, accountId),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesMetricName, "cpu_usage"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesMetricType, "GAUGE"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesRollupFunction, "MAX"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesLabelsEliminationMethod, "EXCLUDE_BY"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesLabels+".#", "2"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesLabels+".0", "instance_id"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesLabels+".1", "region"),
				),
			},
		},
	})
}

func TestAccLogzioMetricsRollupRules_CreateWithFilter(t *testing.T) {
	resourceName := "test_create_metrics_rollup_rules_with_filter"
	resourceFullName := fmt.Sprintf("logzio_metrics_rollup_rules.%s", resourceName)
	accountId := os.Getenv(envLogzioMetricsAccountId)

	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckApiToken(t)
			testAccPreCheckMetricsAccountId(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceTestMetricsRollupRules(resourceName, metricsRollupRulesResourceCreateWithFilter, accountId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesAccountId, accountId),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesName, "frontend_metrics_rollup"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesMetricType, "COUNTER"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesRollupFunction, "SUM"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesLabelsEliminationMethod, "GROUP_BY"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesNewMetricNameTemplate, "rollup.frontend.{{metricName}}"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesDropOriginalMetric, "true"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesFilter+".#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesFilter+".0.expression.#", "2"),
				),
			},
		},
	})
}

func TestAccLogzioMetricsRollupRules_CreateMeasurement(t *testing.T) {
	resourceName := "test_create_metrics_rollup_rules_measurement"
	resourceFullName := fmt.Sprintf("logzio_metrics_rollup_rules.%s", resourceName)
	accountId := os.Getenv(envLogzioMetricsAccountId)

	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckApiToken(t)
			testAccPreCheckMetricsAccountId(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceTestMetricsRollupRules(resourceName, metricsRollupRulesResourceCreateMeasurement, accountId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesAccountId, accountId),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesMetricName, "response_time"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesMetricType, "MEASUREMENT"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesRollupFunction, "MEAN"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesLabelsEliminationMethod, "EXCLUDE_BY"),
				),
			},
		},
	})
}

func TestAccLogzioMetricsRollupRules_CreateCounterWithRollupFunction(t *testing.T) {
	resourceName := "test_create_metrics_rollup_rules_counter_with_rollup_function"
	accountId := os.Getenv(envLogzioMetricsAccountId)

	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckApiToken(t)
			testAccPreCheckMetricsAccountId(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      resourceTestMetricsRollupRules(resourceName, metricsRollupRulesResourceCreateCounterWithRollupFunction, accountId),
				ExpectError: regexp.MustCompile("for COUNTER metrics, rollup_function must be SUM"),
			},
		},
	})
}

func TestAccLogzioMetricsRollupRules_CreateCounterMissingRollupFunction(t *testing.T) {
	resourceName := "test_create_metrics_rollup_rules_counter_missing_rollup_function"
	accountId := os.Getenv(envLogzioMetricsAccountId)

	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckApiToken(t)
			testAccPreCheckMetricsAccountId(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      resourceTestMetricsRollupRules(resourceName, metricsRollupRulesResourceCreateCounterMissingRollupFunction, accountId),
				ExpectError: regexp.MustCompile("rollup_function must be set for COUNTER metrics"),
			},
		},
	})
}

func TestAccLogzioMetricsRollupRules_CreateMeasurementWithP99(t *testing.T) {
	resourceName := "test_create_metrics_rollup_rules_measurement_p99"
	accountId := os.Getenv(envLogzioMetricsAccountId)

	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckApiToken(t)
			testAccPreCheckMetricsAccountId(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceTestMetricsRollupRules(resourceName, metricsRollupRulesResourceCreateMeasurementP99, accountId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("logzio_metrics_rollup_rules."+resourceName, metricsRollupRulesMetricType, "MEASUREMENT"),
					resource.TestCheckResourceAttr("logzio_metrics_rollup_rules."+resourceName, metricsRollupRulesRollupFunction, "P99"),
					resource.TestCheckResourceAttr("logzio_metrics_rollup_rules."+resourceName, metricsRollupRulesMetricName, "response_time"),
				),
			},
		},
	})
}

func resourceTestMetricsRollupRules(name, path, accountId string) string {
	content, err := os.ReadFile(fmt.Sprintf("testdata/fixtures/%s.tf", path))
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf(string(content), name, accountId)
}
