package logzio

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"regexp"
	"testing"
)

func TestAccDataSourceEndpoint(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckApiToken(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				ExpectNonEmptyPlan:        true,
				Config:                    ReadFixtureFromFile("valid_slack_endpoint_datasource.tf"),
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
