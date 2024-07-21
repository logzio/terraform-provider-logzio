package logzio

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
	"regexp"
	"testing"
)

func TestAccDataSourceKibanaObject(t *testing.T) {
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:                    getResourceConfigKibanaObject(),
				PreventPostDestroyRefresh: true,
			},
			{
				Config:                    getResourceConfigKibanaObject() + utils.ReadFixtureFromFile("kibana_object_datasource.tf"),
				PreventPostDestroyRefresh: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.logzio_kibana_object.ds_kb", kibanaObjectKibanaVersionField, "7.2.1"),
					resource.TestMatchOutput("output_id", regexp.MustCompile("search:tf-provider-datasource-test-search")),
				),
			},
		},
	})
}

func getResourceConfigKibanaObject() string {
	return `resource "logzio_kibana_object" "test_kb_for_datasource" {
  kibana_version = "7.2.1"
  data = file("./testdata/fixtures/kibana_objects/create_search_for_datasource.json")
}
`
}
