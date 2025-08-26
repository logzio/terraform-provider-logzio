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
	dropMetricResourceCreateSimple           = "create_drop_metrics_simple"
	dropMetricResourceCreateComplex          = "create_drop_metrics_complex"
	dropMetricResourceCreateWithName         = "create_drop_metrics_with_name"
	dropMetricResourceCreateNameTooLong      = "create_drop_metrics_name_too_long"
	dropMetricResourceCreateNoLabelName      = "create_drop_metrics_missing_label_name"
	dropMetricResourceCreateNoValue          = "create_drop_metrics_missing_value"
	dropMetricResourceCreateEmptyCondition   = "create_drop_metrics_empty_condition"
	dropMetricResourceCreateInvalidCondition = "create_drop_metrics_invalid_condition"
	dropMetricResourceCreateInvalidOperator  = "create_drop_metrics_invalid_operator"
	dropMetricResourceUpdate                 = "update_drop_metrics"
	dropMetricResourceUpdateWithName         = "update_drop_metrics_with_name"
)

func TestAccLogzioDropMetric_CreateDropMetricSimple(t *testing.T) {
	filterName := "test_create_drop_metrics_simple"
	resourceName := "logzio_drop_metrics." + filterName
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
				Config: resourceTestDropMetrics(filterName, dropMetricResourceCreateSimple, accountId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, dropMetricsAccountId, accountId),
					resource.TestCheckResourceAttr(resourceName, dropMetricsName, ""),
					resource.TestCheckResourceAttr(resourceName, dropMetricsFilters+".#", "1"),
					resource.TestCheckResourceAttr(resourceName, dropMetricsFilterOperator, "AND"),
					resource.TestCheckResourceAttr(resourceName, dropMetricsActive, "true"),
				),
			},
			{
				Config:            resourceTestDropMetrics(filterName, dropMetricResourceCreateSimple, accountId),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLogzioDropMetric_CreateDropMetricComplex(t *testing.T) {
	filterName := "test_create_drop_metrics_complex"
	resourceName := "logzio_drop_metrics." + filterName
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
				Config: resourceTestDropMetrics(filterName, dropMetricResourceCreateComplex, accountId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, dropMetricsAccountId, accountId),
					resource.TestCheckResourceAttr(resourceName, dropMetricsFilters+".#", "2"),
					resource.TestCheckResourceAttr(resourceName, dropMetricsFilterOperator, "AND"),
					resource.TestCheckResourceAttr(resourceName, dropMetricsActive, "true"),
				),
			},
			{
				Config:            resourceTestDropMetrics(filterName, dropMetricResourceCreateComplex, accountId),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLogzioDropMetric_CreateDropMetricNoLabelName(t *testing.T) {
	filterName := "test_create_drop_metrics_no_label_name"
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
				Config:      resourceTestDropMetrics(filterName, dropMetricResourceCreateNoLabelName, accountId),
				ExpectError: regexp.MustCompile("The argument \"name\" is required, but no definition was found"),
			},
		},
	})
}

func TestAccLogzioDropMetric_CreateDropMetricNoLabelValue(t *testing.T) {
	filterName := "test_create_drop_metrics_no_label_value"
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
				Config:      resourceTestDropMetrics(filterName, dropMetricResourceCreateNoValue, accountId),
				ExpectError: regexp.MustCompile("The argument \"value\" is required, but no definition was found"),
			},
		},
	})
}

func TestAccLogzioDropMetric_CreateDropMetricNoCondition(t *testing.T) {
	filterName := "test_create_drop_metrics_no_filter_condition"
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
				Config:      resourceTestDropMetrics(filterName, dropMetricResourceCreateEmptyCondition, accountId),
				ExpectError: regexp.MustCompile("The argument \"condition\" is required, but no definition was found"),
			},
		},
	})
}

func TestAccLogzioDropMetric_CreateDropMetricInvalidCondition(t *testing.T) {
	filterName := "test_create_drop_metrics_invalid_filter_condition"
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
				Config:      resourceTestDropMetrics(filterName, dropMetricResourceCreateInvalidCondition, accountId),
				ExpectError: regexp.MustCompile("expected filters\\.0\\.condition to be one of \\[\"EQ\" \"NOT_EQ\" \"REGEX_MATCH\" \"REGEX_NO_MATCH\"]"),
			},
		},
	})
}

func TestAccLogzioDropMetric_CreateDropMetricInvalidOperator(t *testing.T) {
	filterName := "test_create_drop_metrics_invalid_operator"
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
				Config:      resourceTestDropMetrics(filterName, dropMetricResourceCreateInvalidOperator, accountId),
				ExpectError: regexp.MustCompile("expected operator to be one of \\[\"AND\"]"),
			},
		},
	})
}

func TestAccLogzioDropMetric_UpdateDropMetrics(t *testing.T) {
	filterName := "test_update_drop_metrics"
	resourceName := "logzio_drop_metrics." + filterName
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
				Config: resourceTestDropMetrics(filterName, dropMetricResourceCreateSimple, accountId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, dropMetricsAccountId, accountId),
					resource.TestCheckResourceAttr(resourceName, dropMetricsFilters+".#", "1"),
					resource.TestCheckResourceAttr(resourceName, dropMetricsFilterOperator, "AND"),
					resource.TestCheckResourceAttr(resourceName, dropMetricsActive, "true"),
				),
			},
			{
				Config: resourceTestDropMetrics(filterName, dropMetricResourceUpdate, accountId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, dropMetricsAccountId, accountId),
					resource.TestCheckResourceAttr(resourceName, dropMetricsFilters+".#", "2"),
					resource.TestCheckResourceAttr(resourceName, dropMetricsFilterOperator, "AND"),
					resource.TestCheckResourceAttr(resourceName, dropMetricsActive, "false"),
				),
			},
		},
	})
}

func TestAccLogzioDropMetric_UpdateDropMetricsEnable(t *testing.T) {
	filterName := "test_update_enable_drop_metrics"
	resourceName := "logzio_drop_metrics." + filterName
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
				Config: resourceTestDropMetrics(filterName, dropMetricResourceUpdate, accountId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, dropMetricsAccountId, accountId),
					resource.TestCheckResourceAttr(resourceName, dropMetricsFilters+".#", "2"),
					resource.TestCheckResourceAttr(resourceName, dropMetricsFilterOperator, "AND"),
					resource.TestCheckResourceAttr(resourceName, dropMetricsActive, "false"),
				),
			},
			{
				Config: resourceTestDropMetrics(filterName, dropMetricResourceCreateComplex, accountId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, dropMetricsAccountId, accountId),
					resource.TestCheckResourceAttr(resourceName, dropMetricsFilters+".#", "2"),
					resource.TestCheckResourceAttr(resourceName, dropMetricsFilterOperator, "AND"),
					resource.TestCheckResourceAttr(resourceName, dropMetricsActive, "true"),
				),
			},
		},
	})
}

func TestAccLogzioDropMetric_UpdateDropMetricsDisable(t *testing.T) {
	filterName := "test_update_disable_drop_metrics"
	resourceName := "logzio_drop_metrics." + filterName
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
				Config: resourceTestDropMetrics(filterName, dropMetricResourceCreateComplex, accountId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, dropMetricsAccountId, accountId),
					resource.TestCheckResourceAttr(resourceName, dropMetricsFilters+".#", "2"),
					resource.TestCheckResourceAttr(resourceName, dropMetricsFilterOperator, "AND"),
					resource.TestCheckResourceAttr(resourceName, dropMetricsActive, "true"),
				),
			},
			{
				Config: resourceTestDropMetrics(filterName, dropMetricResourceUpdate, accountId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, dropMetricsAccountId, accountId),
					resource.TestCheckResourceAttr(resourceName, dropMetricsFilters+".#", "2"),
					resource.TestCheckResourceAttr(resourceName, dropMetricsFilterOperator, "AND"),
					resource.TestCheckResourceAttr(resourceName, dropMetricsActive, "false"),
				),
			},
		},
	})
}

func TestAccLogzioDropMetric_CreateDropMetricWithName(t *testing.T) {
	filterName := "test_create_drop_metrics_with_name"
	resourceName := "logzio_drop_metrics." + filterName
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
				Config: resourceTestDropMetrics(filterName, dropMetricResourceCreateWithName, accountId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, dropMetricsAccountId, accountId),
					resource.TestCheckResourceAttr(resourceName, dropMetricsName, "test-drop-metrics-filter"),
					resource.TestCheckResourceAttr(resourceName, dropMetricsFilters+".#", "1"),
					resource.TestCheckResourceAttr(resourceName, dropMetricsFilterOperator, "AND"),
					resource.TestCheckResourceAttr(resourceName, dropMetricsActive, "true"),
				),
			},
			{
				Config:            resourceTestDropMetrics(filterName, dropMetricResourceCreateWithName, accountId),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLogzioDropMetric_CreateDropMetricNameTooLong(t *testing.T) {
	filterName := "test_create_drop_metrics_name_too_long"
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
				Config:      resourceTestDropMetrics(filterName, dropMetricResourceCreateNameTooLong, accountId),
				ExpectError: regexp.MustCompile("expected length of name to be in the range \\(0 - 256\\)"),
			},
		},
	})
}

func TestAccLogzioDropMetric_UpdateDropMetricsWithName(t *testing.T) {
	filterName := "test_update_drop_metrics_with_name"
	resourceName := "logzio_drop_metrics." + filterName
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
				Config: resourceTestDropMetrics(filterName, dropMetricResourceCreateSimple, accountId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, dropMetricsAccountId, accountId),
					resource.TestCheckResourceAttr(resourceName, dropMetricsName, ""),
					resource.TestCheckResourceAttr(resourceName, dropMetricsFilters+".#", "1"),
					resource.TestCheckResourceAttr(resourceName, dropMetricsFilterOperator, "AND"),
					resource.TestCheckResourceAttr(resourceName, dropMetricsActive, "true"),
				),
			},
			{
				Config: resourceTestDropMetrics(filterName, dropMetricResourceUpdateWithName, accountId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, dropMetricsAccountId, accountId),
					resource.TestCheckResourceAttr(resourceName, dropMetricsName, "updated-drop-metrics-filter"),
					resource.TestCheckResourceAttr(resourceName, dropMetricsFilters+".#", "2"),
					resource.TestCheckResourceAttr(resourceName, dropMetricsFilterOperator, "AND"),
					resource.TestCheckResourceAttr(resourceName, dropMetricsActive, "false"),
				),
			},
		},
	})
}

func resourceTestDropMetrics(name, path, accountId string) string {
	content, err := os.ReadFile(fmt.Sprintf("testdata/fixtures/%s.tf", path))
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf(string(content), name, accountId)
}
