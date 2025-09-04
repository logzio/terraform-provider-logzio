package logzio

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
)

func TestAccDataSourceMetricsRollupRules_Basic(t *testing.T) {
	resourceName := "datasource_test_resource_metrics_rollup_rules_find"
	datasourceName := "datasource_test_find_metrics_rollup_rules_by_id"
	resourceFullName := "logzio_metrics_rollup_rules." + resourceName
	datasourceFullName := "data.logzio_metrics_rollup_rules." + datasourceName
	accountId := os.Getenv(envLogzioMetricsAccountId)

	resourceConfig := datasourceResourceTestMetricsRollupRules(resourceName, accountId)

	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckApiToken(t)
			testAccPreCheckMetricsAccountId(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesAccountId, accountId),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesMetricName, "memory_usage"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesMetricType, "GAUGE"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesRollupFunction, "LAST"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesLabelsEliminationMethod, "EXCLUDE_BY"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesLabels+".#", "2"),
				),
			},
			{
				Config: resourceConfig +
					datasourceMetricsRollupRulesById(datasourceName, fmt.Sprintf("${%s.id}", resourceFullName)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceFullName, metricsRollupRulesAccountId, accountId),
					resource.TestCheckResourceAttr(datasourceFullName, metricsRollupRulesMetricName, "memory_usage"),
					resource.TestCheckResourceAttr(datasourceFullName, metricsRollupRulesMetricType, "GAUGE"),
					resource.TestCheckResourceAttr(datasourceFullName, metricsRollupRulesRollupFunction, "LAST"),
					resource.TestCheckResourceAttr(datasourceFullName, metricsRollupRulesLabelsEliminationMethod, "EXCLUDE_BY"),
					resource.TestCheckResourceAttr(datasourceFullName, metricsRollupRulesLabels+".#", "2"),
					resource.TestCheckResourceAttr(datasourceFullName, metricsRollupRulesLabels+".0", "instance_id"),
					resource.TestCheckResourceAttr(datasourceFullName, metricsRollupRulesLabels+".1", "region"),
				),
			},
		},
	})
}

func TestAccDataSourceMetricsRollupRules_NotFound(t *testing.T) {
	datasourceName := "datasource_test_metrics_rollup_rules_not_existing"

	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckApiToken(t)
			testAccPreCheckMetricsAccountId(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      datasourceMetricsRollupRulesById(datasourceName, "nonexistent-id-123"),
				ExpectError: regexp.MustCompile("Error"),
			},
		},
	})
}

func TestAccDataSourceMetricsRollupRules_SearchByAttributes(t *testing.T) {
	resourceName := "datasource_test_resource_metrics_rollup_rules_search"
	datasourceName := "datasource_test_find_metrics_rollup_rules_by_attrs"
	resourceFullName := "logzio_metrics_rollup_rules." + resourceName
	datasourceFullName := "data.logzio_metrics_rollup_rules." + datasourceName
	accountId := os.Getenv(envLogzioMetricsAccountId)

	resourceConfig := datasourceResourceTestMetricsRollupRulesWithName(resourceName, accountId)

	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckApiToken(t)
			testAccPreCheckMetricsAccountId(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesAccountId, accountId),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesName, "Test Search Rollup"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesMetricName, "disk_usage"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesMetricType, "GAUGE"),
				),
			},
			{
				Config: resourceConfig +
					datasourceMetricsRollupRulesByAttributes(datasourceName, accountId, "disk_usage", "GAUGE"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceFullName, metricsRollupRulesAccountId, accountId),
					resource.TestCheckResourceAttr(datasourceFullName, metricsRollupRulesName, "Test Search Rollup"),
					resource.TestCheckResourceAttr(datasourceFullName, metricsRollupRulesMetricName, "disk_usage"),
					resource.TestCheckResourceAttr(datasourceFullName, metricsRollupRulesMetricType, "GAUGE"),
					resource.TestCheckResourceAttr(datasourceFullName, metricsRollupRulesRollupFunction, "MEAN"),
					resource.TestCheckResourceAttr(datasourceFullName, metricsRollupRulesLabelsEliminationMethod, "GROUP_BY"),
				),
			},
		},
	})
}

func datasourceResourceTestMetricsRollupRules(resourceName, accountId string) string {
	return fmt.Sprintf(`resource "logzio_metrics_rollup_rules" "%s" {
  account_id = %s
  metric_name = "memory_usage"
  metric_type = "GAUGE"
  rollup_function = "LAST"
  labels_elimination_method = "EXCLUDE_BY"
  labels = ["instance_id", "region"]
}
`, resourceName, accountId)
}

func datasourceResourceTestMetricsRollupRulesWithName(resourceName, accountId string) string {
	return fmt.Sprintf(`resource "logzio_metrics_rollup_rules" "%s" {
  account_id = %s
  name = "Test Search Rollup"
  metric_name = "disk_usage"
  metric_type = "GAUGE"
  rollup_function = "MEAN"
  labels_elimination_method = "GROUP_BY"
  labels = ["host", "environment"]
}
`, resourceName, accountId)
}

func datasourceMetricsRollupRulesById(datasourceName, rollupRuleId string) string {
	return fmt.Sprintf(`
data "logzio_metrics_rollup_rules" "%s" {
  id = "%s"
}
`, datasourceName, rollupRuleId)
}

func datasourceMetricsRollupRulesByAttributes(datasourceName, accountId, metricName, metricType string) string {
	return fmt.Sprintf(`
data "logzio_metrics_rollup_rules" "%s" {
  account_id  = %s
  metric_name = "%s"
  metric_type = "%s"
}
`, datasourceName, accountId, metricName, metricType)
}
