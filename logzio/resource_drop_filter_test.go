package logzio

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"io/ioutil"
	"log"
	"testing"
)

const (
	dropFilterResourceCreateAlert = "create_drop_filter"
)

func TestAccLogzioDropFilter_CreateDropFilter(t *testing.T) {
	filterName := "test_create_drop_filter"
	resourceName := "logzio_drop_filter." + filterName

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckApiToken(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: resourceTestDropFilter(filterName, dropFilterResourceCreateAlert),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, dropFilterLogType, "some_type"),
					resource.TestCheckResourceAttr(resourceName, dropFilterFieldConditions+".#", "2"),
				),
			},
			{
				Config:            resourceTestDropFilter(filterName, dropFilterResourceCreateAlert),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func resourceTestDropFilter(name string, path string) string {
	content, err := ioutil.ReadFile(fmt.Sprintf("testdata/fixtures/%s.tf", path))
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf(fmt.Sprintf("%s", content), name)
}
