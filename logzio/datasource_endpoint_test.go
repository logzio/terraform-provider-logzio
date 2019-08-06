package logzio

import (
	"github.com/hashicorp/terraform/helper/resource"
	"regexp"
	"testing"
)

func TestAccDataSourceEndpoint(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				ExpectNonEmptyPlan:        true,
				Config:                    readFixtureFromFile("valid_slack_endpoint_datasource.tf"),
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