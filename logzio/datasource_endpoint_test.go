package logzio

import (
	"fmt"
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
				Config:                    endpointDatasourceConfig(),
				PreventPostDestroyRefresh: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.logzio_endpoint.by_title", "title", "ATestStaticEndpoint"),
					resource.TestMatchResourceAttr("data.logzio_endpoint.by_title", "id", regexp.MustCompile("\\d")),
					resource.TestMatchOutput("test", regexp.MustCompile("\\d")),
				),
			},
		},
	})
}

func endpointDatasourceConfig() string {
	return fmt.Sprintf(`
resource "logzio_endpoint" "slack" {
  endpoint_type = "slack"
  title = "ATestStaticEndpoint"
  description = "this_is_my_description"
  slack {
	url = "https://www.test.com"
  }
}

data "logzio_endpoint" "by_title" {
  title = "ATestStaticEndpoint"
  depends_on = ["logzio_endpoint.slack"]
}

output "test" {
  value = "${data.logzio_endpoint.by_title.id}"
}
`)
}
