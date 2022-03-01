package logzio

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
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
	defer utils.SleepAfterTest()

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

func testAccDropFilterDataSourceDropFilterById() string {
	return fmt.Sprintf(`
data "logzio_drop_filter" "my_drop_filter_datasource" {
drop_filter_id = "${logzio_drop_filter.test_create_drop_filter_for_ds.id}"
}
`)
}
