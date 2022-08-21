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
				ExpectNonEmptyPlan:        true,
				Config:                    utils.ReadFixtureFromFile("kibana_object_datasource.tf"),
				PreventPostDestroyRefresh: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.logzio_kibana_object.ds_kb", kibanaObjectKibanaVersionField, "7.2.1"),
					resource.TestMatchOutput("output_id", regexp.MustCompile("search:tf-provider-test-search")),
				),
			},
		},
	})
}
