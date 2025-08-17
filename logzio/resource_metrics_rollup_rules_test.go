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
	metricsRollupRulesResourceCreateSimple                   = "create_metrics_rollup_rules_simple"
	metricsRollupRulesResourceCreateComplex                  = "create_metrics_rollup_rules_complex"
	metricsRollupRulesResourceCreateInvalidMetricType        = "create_metrics_rollup_rules_invalid_metric_type"
	metricsRollupRulesResourceCreateInvalidRollupFunction    = "create_metrics_rollup_rules_invalid_rollup_function"
	metricsRollupRulesResourceCreateInvalidEliminationMethod = "create_metrics_rollup_rules_invalid_elimination_method"
	metricsRollupRulesResourceCreateEmptyMetricName          = "create_metrics_rollup_rules_empty_metric_name"
	metricsRollupRulesResourceCreateEmptyLabels              = "create_metrics_rollup_rules_empty_labels"
	metricsRollupRulesResourceUpdate                         = "update_metrics_rollup_rules"
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
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesMetricType, "gauge"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesRollupFunction, "last"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesLabelsEliminationMethod, "exclude_by"),
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
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesMetricType, "counter"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesRollupFunction, "sum"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesLabelsEliminationMethod, "exclude_by"),
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
				ExpectError: regexp.MustCompile("expected metric_type to be one of \\[\"gauge\" \"counter\"\\]"),
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
				ExpectError: regexp.MustCompile("expected rollup_function to be one of \\[\"sum\" \"min\" \"max\" \"count\" \"last\"\\]"),
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
				ExpectError: regexp.MustCompile("expected labels_elimination_method to be one of \\[\"exclude_by\"\\]"),
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
				ExpectError: regexp.MustCompile("metric_name cannot be empty"),
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
				ExpectError: regexp.MustCompile("At least 1 \"labels\" must be configured"),
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
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesMetricType, "gauge"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesRollupFunction, "last"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesLabelsEliminationMethod, "exclude_by"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesLabels+".#", "1"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesLabels+".0", "instance_id"),
				),
			},
			{
				Config: resourceTestMetricsRollupRules(resourceName, metricsRollupRulesResourceUpdate, accountId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesAccountId, accountId),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesMetricName, "cpu_usage"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesMetricType, "gauge"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesRollupFunction, "max"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesLabelsEliminationMethod, "exclude_by"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesLabels+".#", "2"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesLabels+".0", "instance_id"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesLabels+".1", "region"),
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
