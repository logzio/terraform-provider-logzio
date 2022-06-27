package logzio

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/logzio/logzio_terraform_client/kibana_objects"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
	"io/ioutil"
	"log"
	"testing"
)

func TestAccLogzioKibanaObject_CreateUpdateSearch(t *testing.T) {
	resourceName := "logzio_kibana_object.test_kb_obj_search"
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckApiToken(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: resourceTestKibanaObject(kibana_objects.ExportTypeSearch, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, kibanaObjectKibanaVersionField, "7.2.1"),
					resource.TestCheckResourceAttr(resourceName, "id", "tf-provider-test-search"),
				),
			},
			{
				Config: resourceTestKibanaObject(kibana_objects.ExportTypeSearch, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, kibanaObjectKibanaVersionField, "7.2.1"),
					resource.TestCheckResourceAttr(resourceName, "id", "tf-provider-test-search"),
				),
			},
		},
	})
}

func resourceTestKibanaObject(objType kibana_objects.ExportType, update bool) string {
	switch objType {
	case kibana_objects.ExportTypeSearch:
		if update {
			return getKibanaObjectResourceConfig("update_kibana_object_search")
		}
		return getKibanaObjectResourceConfig("create_kibana_object_search")
	default:
		// we should never get to this part
		panic("invalid kibana object type")
	}
}

func getKibanaObjectResourceConfig(path string) string {
	content, err := ioutil.ReadFile(fmt.Sprintf("testdata/fixtures/%s.tf", path))
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%s", content)
}
