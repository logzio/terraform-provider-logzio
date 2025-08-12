package logzio

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
	"os"
	"regexp"
	"testing"
)

func TestAccDataSourceDropMetric(t *testing.T) {
	resourceFilterName := "datasource_test_resource_drop_metric_find"
	datasourceFilterName := "datasource_test_find_drop_metric_by_id"
	resourceName := "logzio_drop_metrics." + resourceFilterName
	datasourceName := "data.logzio_drop_metrics." + datasourceFilterName
	accountId := os.Getenv(envLogzioMetricsAccountId)

	resourceConfig := datasourceResourceTestDropMetrics(resourceFilterName, accountId)

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
					resource.TestCheckResourceAttr(resourceName, dropMetricsAccountId, accountId),
					resource.TestCheckResourceAttr(resourceName, dropMetricsFilters+".#", "2"),
					resource.TestCheckResourceAttr(resourceName, dropMetricsFilterOperator, "AND"),
					resource.TestCheckResourceAttr(resourceName, dropMetricsActive, "true"),
				),
			},
			{
				Config: resourceConfig +
					datasourceDropMetricMatchingResource(datasourceFilterName, accountId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, dropMetricsAccountId, accountId),
					resource.TestCheckResourceAttr(datasourceName, dropMetricsFilters+".#", "2"),
					resource.TestCheckResourceAttr(datasourceName, dropMetricsFilterOperator, "AND"),
					resource.TestCheckResourceAttr(datasourceName, dropMetricsActive, "true"),
				),
			},
		},
	})
}

func TestAccDataSourceDropMetricById(t *testing.T) {
	resourceFilterName := "datasource_test_resource_drop_metric_find_by_id"
	datasourceFilterName := "datasource_test_drop_metric_find_by_id"
	resourceName := "logzio_drop_metrics." + resourceFilterName
	datasourceName := "data.logzio_drop_metrics." + datasourceFilterName
	dropMetricsId := fmt.Sprintf("${%s.id}", resourceName)
	accountId := os.Getenv(envLogzioMetricsAccountId)

	resourceConfig := datasourceResourceTestDropMetrics(resourceFilterName, accountId)

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
					resource.TestCheckResourceAttr(resourceName, dropMetricsAccountId, accountId),
					resource.TestCheckResourceAttr(resourceName, dropMetricsFilters+".#", "2"),
					resource.TestCheckResourceAttr(resourceName, dropMetricsFilterOperator, "AND"),
					resource.TestCheckResourceAttr(resourceName, dropMetricsActive, "true"),
				),
			},
			{
				Config: resourceConfig +
					datasourceDropMetricById(datasourceFilterName, dropMetricsId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, dropMetricsAccountId, accountId),
					resource.TestCheckResourceAttr(datasourceName, dropMetricsFilters+".#", "2"),
					resource.TestCheckResourceAttr(datasourceName, dropMetricsFilterOperator, "AND"),
					resource.TestCheckResourceAttr(datasourceName, dropMetricsActive, "true"),
				),
			},
		},
	})
}

func TestAccDataSourceDropMetricNotFoundAndNotEnoughCriteriaToSearch(t *testing.T) {
	datasourceFilterName := "datasource_test_drop_metric_not_existing_and_no_search_criteria"

	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckApiToken(t)
			testAccPreCheckMetricsAccountId(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      datasourceDropMetricById(datasourceFilterName, "123"),
				ExpectError: regexp.MustCompile("could not find drop metrics filter with id 123"),
			},
		},
	})
}

func TestAccDataSourceDropMetricNotFoundSearch(t *testing.T) {
	datasourceFilterName := "datasource_test_drop_metric_not_existing"
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
				Config:      datasourceDropMetricWithCriteria(datasourceFilterName, accountId),
				ExpectError: regexp.MustCompile("no metrics drop filters matched the criteria:"),
			},
		},
	})
}

func TestAccDataSourceDropMetricTooManyMatches(t *testing.T) {
	resourceFilterName := "datasource_test_resource_drop_metric_too_many_matches"
	datasourceFilterName := "datasource_test_drop_metric_too_many_matches"
	resourceName := "logzio_drop_metrics." + resourceFilterName
	accountId := os.Getenv(envLogzioMetricsAccountId)

	resourceConfig := datasourceResourceTestDropMetrics(resourceFilterName, accountId)
	resourceConfig2 := datasourceResourceTestDropMetrics(fmt.Sprintf("%s_second", resourceFilterName), accountId)

	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckApiToken(t)
			testAccPreCheckMetricsAccountId(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig + resourceConfig2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, dropMetricsAccountId, accountId),
					resource.TestCheckResourceAttr(resourceName, dropMetricsFilters+".#", "2"),
					resource.TestCheckResourceAttr(resourceName, dropMetricsFilterOperator, "AND"),
					resource.TestCheckResourceAttr(resourceName, dropMetricsActive, "true"),
				),
			},
			{
				Config:      datasourceDropMetricWithCriteria(datasourceFilterName, accountId),
				ExpectError: regexp.MustCompile("found multiple \\(2\\) metrics drop filters matching the critiria, add more expressions or specify an id"),
			},
		},
	})
}

func TestAccDataSourceDropMetricNotFoundSearchWithResults(t *testing.T) {
	resourceFilterName := "datasource_test_resource_drop_metric_no_match_found_search"
	datasourceFilterName := "datasource_test_drop_metric_no_match_found_search"
	resourceName := "logzio_drop_metrics." + resourceFilterName
	accountId := os.Getenv(envLogzioMetricsAccountId)

	resourceConfig := datasourceResourceTestDropMetrics(resourceFilterName, accountId)
	resourceConfig2 := datasourceResourceTestDropMetrics(fmt.Sprintf("%s_second", resourceFilterName), accountId)

	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckApiToken(t)
			testAccPreCheckMetricsAccountId(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig + resourceConfig2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, dropMetricsAccountId, accountId),
					resource.TestCheckResourceAttr(resourceName, dropMetricsFilters+".#", "2"),
					resource.TestCheckResourceAttr(resourceName, dropMetricsFilterOperator, "AND"),
					resource.TestCheckResourceAttr(resourceName, dropMetricsActive, "true"),
				),
			},
			{
				Config:      datasourceDropMetricWithCriteriaAndExpressions(datasourceFilterName, accountId),
				ExpectError: regexp.MustCompile("could not find metrics drop filter with the specified attributes:"),
			},
		},
	})

}

func datasourceResourceTestDropMetrics(resourceFilterName, accountId string) string {
	return fmt.Sprintf(`resource "logzio_drop_metrics" "%s" {
  account_id = %s

  filters {
    name = "__name__"
    value = "my_metric"
    condition = "EQ"
  }

  filters {
    name = "some_label"
    value = "some_value"
    condition = "NOT_EQ"
  }
}
`, resourceFilterName, accountId)
}

func datasourceDropMetricById(datasourceFilterName, dropMetricsId string) string {
	return fmt.Sprintf(`
data "logzio_drop_metrics" "%s" {
drop_metric_id = "%s"
}
`, datasourceFilterName, dropMetricsId)
}

func datasourceDropMetricWithCriteria(datasourceFilterName, accountId string) string {
	return fmt.Sprintf(`
data "logzio_drop_metrics" "%s" {
account_id = "%s"
}
`, datasourceFilterName, accountId)
}

func datasourceDropMetricWithCriteriaAndExpressions(datasourceFilterName, accountId string) string {
	return fmt.Sprintf(`
data "logzio_drop_metrics" "%s" {
  account_id = "%s"

  filters {
    name = "__name__"
    value = "my_other_metric"
    condition = "EQ"
  }
}
`, datasourceFilterName, accountId)
}

func datasourceDropMetricMatchingResource(datasourceFilterName, accountId string) string {
	return fmt.Sprintf(`
data "logzio_drop_metrics" "%s" {
  account_id = "%s"

  filters {
    name = "__name__"
    value = "my_metric"
    condition = "EQ"
  }

  filters {
    name = "some_label"
    value = "some_value"
    condition = "NOT_EQ"
  }
}
`, datasourceFilterName, accountId)
}
