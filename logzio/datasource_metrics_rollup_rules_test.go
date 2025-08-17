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
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesMetricType, "gauge"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesRollupFunction, "last"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesLabelsEliminationMethod, "exclude_by"),
					resource.TestCheckResourceAttr(resourceFullName, metricsRollupRulesLabels+".#", "2"),
				),
			},
			{
				Config: resourceConfig +
					datasourceMetricsRollupRulesById(datasourceName, fmt.Sprintf("${%s.id}", resourceFullName)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceFullName, metricsRollupRulesAccountId, accountId),
					resource.TestCheckResourceAttr(datasourceFullName, metricsRollupRulesMetricName, "memory_usage"),
					resource.TestCheckResourceAttr(datasourceFullName, metricsRollupRulesMetricType, "gauge"),
					resource.TestCheckResourceAttr(datasourceFullName, metricsRollupRulesRollupFunction, "last"),
					resource.TestCheckResourceAttr(datasourceFullName, metricsRollupRulesLabelsEliminationMethod, "exclude_by"),
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

func datasourceResourceTestMetricsRollupRules(resourceName, accountId string) string {
	return fmt.Sprintf(`resource "logzio_metrics_rollup_rules" "%s" {
  account_id = %s
  metric_name = "memory_usage"
  metric_type = "gauge"
  rollup_function = "last"
  labels_elimination_method = "exclude_by"
  labels = ["instance_id", "region"]
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
