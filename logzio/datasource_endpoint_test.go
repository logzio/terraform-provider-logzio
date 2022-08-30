package logzio

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
	"regexp"
	"testing"
)

func TestAccDataSourceEndpoint(t *testing.T) {
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:                    utils.ReadFixtureFromFile("valid_slack_endpoint_datasource.tf"),
				PreventPostDestroyRefresh: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.logzio_endpoint.by_title", "title", "valid_slack_endpoint_datasource"),
					resource.TestMatchResourceAttr("data.logzio_endpoint.by_title", "id", regexp.MustCompile("\\d")),
					resource.TestMatchOutput("valid_slack_endpoint_datasource_id", regexp.MustCompile("\\d")),
				),
			},
		},
	})
}
