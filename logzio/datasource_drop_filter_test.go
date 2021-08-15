package logzio

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"regexp"
	"testing"
)

func TestAccDataSourceDropFilter(t *testing.T) {
	resourceName := "logzio_drop_filter.test_create_drop_filter_for_ds"
	dataSourceName := "data.logzio_drop_filter.my_drop_filter_datasource"
	resourceConfig := fmt.Sprintf(`resource "logzio_drop_filter" "test_create_drop_filter_for_ds" {
  log_type = "some_type_create_datadource"

  field_conditions {
    field_name = "some_field"
    value = "some_string_value"
  }
  field_conditions {
    field_name = "another_field"
    value = 200
  }
}
`)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckApiToken(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, dropFilterLogType, "some_type_create_datadource"),
					resource.TestCheckResourceAttr(resourceName, dropFilterFieldConditions+".#", "2"),
				),
			},
			{
				Config: resourceConfig +
					testAccDropFilterDataSourceDropFilterById(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, dropFilterLogType, "some_type_create_datadource"),
					resource.TestCheckResourceAttr(dataSourceName, "field_conditions.#", "2"),
				),
			},
		},
	})
}

func TestAccDataSourceDropFilterByAttributes(t *testing.T) {
	resourceName := "logzio_drop_filter.test_create_drop_filter_for_ds_by_att"
	dataSourceName := "data.logzio_drop_filter.my_drop_filter_datasource"
	resourceConfig := `resource "logzio_drop_filter" "test_create_drop_filter_for_ds_by_att" {
  log_type = "some_type_create_datadource_by_att"

  field_conditions {
    field_name = "some_field"
    value = "some_string_value"
  }
  field_conditions {
    field_name = "another_field"
    value = 200
  }
}
`
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckApiToken(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, dropFilterLogType, "some_type_create_datadource_by_att"),
					resource.TestCheckResourceAttr(resourceName, dropFilterFieldConditions+".#", "2"),
				),
			},
			{
				Config: resourceConfig +
					testAccDropFilterDataSourceDropFilterByAttributes(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, dropFilterLogType, "some_type_create_datadource_by_att"),
					resource.TestCheckResourceAttr(dataSourceName, "field_conditions.#", "2"),
					resource.TestCheckResourceAttrSet(dataSourceName, "drop_filter_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
				),
			},
		},
	})
}

func TestAccDataSourceDropFilterNotExist(t *testing.T) {
	resourceName := "logzio_drop_filter.test_create_drop_filter_for_ds_to_list"
	resourceConfig := fmt.Sprint(`resource "logzio_drop_filter" "test_create_drop_filter_for_ds_to_list" {
  log_type = "some_type_create_to_list"

  field_conditions {
    field_name = "some_field"
    value = "some_string_value"
  }
  field_conditions {
    field_name = "another_field"
    value = 200
  }
}
`)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckApiToken(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, dropFilterLogType, "some_type_create_to_list"),
					resource.TestCheckResourceAttr(resourceName, dropFilterFieldConditions+".#", "2"),
				),
			},
			{
				Config:      testAccDropFilterDataSourceDropFilterByIdNotExists(),
				ExpectError: regexp.MustCompile("couldn't find drop filter with specified attributes"),
			},
		},
	})
}

func testAccDropFilterDataSourceDropFilterById() string {
	return fmt.Sprintf(`
data "logzio_drop_filter" "my_drop_filter_datasource" {
drop_filter_id = "${logzio_drop_filter.test_create_drop_filter_for_ds.id}"
}
`)
}

func testAccDropFilterDataSourceDropFilterByAttributes() string {
	return fmt.Sprintf(`
data "logzio_drop_filter" "my_drop_filter_datasource" {
  log_type = "some_type_create_datadource_by_att"

  field_conditions {
    field_name = "some_field"
    value = "some_string_value"
  }
  field_conditions {
    field_name = "another_field"
    value = 200
  }
}
`)
}

func testAccDropFilterDataSourceDropFilterByIdNotExists() string {
	return fmt.Sprintf(`
data "logzio_drop_filter" "my_drop_filter_datasource" {
drop_filter_id = "1234"
}
`)
}