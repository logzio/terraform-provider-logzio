package logzio

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"io/ioutil"
	"log"
	"testing"
)

func TestAccDataSourceLogzIoAlert(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				ExpectNonEmptyPlan: true,
				Config:             testAccDataSourceLogzioAlertConfig("by_title"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.logzio_alert.by_title", "title", "hello"),
					resource.TestCheckResourceAttr("data.logzio_alert.by_title", "query_string", "loglevel:ERROR"),
					resource.TestCheckResourceAttr("data.logzio_alert.by_title", "operation", "GREATER_THAN"),
				),
			},
		},
	})
}

func testAccDataSourceLogzioAlertBase(name string) string {
	content, err := ioutil.ReadFile("testdata/fixtures/create_alert.tf")
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf(fmt.Sprintf("%s", content), name)
}

func testAccDataSourceLogzioAlertConfig(name string) string {
	return testAccDataSourceLogzioAlertBase(name) + `

data "logzio_alert" "by_title" {
  title = "hello"
  depends_on = ["logzio_alert.by_title"]
}
`
}
