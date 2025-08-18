package logzio

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
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
					resource.TestMatchOutput("output_id", regexp.MustCompile("search:tf-provider-datasource-test-search-[0-9]+")),
				),
			},
		},
	})
}

func getResourceConfigKibanaObject() string {
	// Load and make unique the JSON data
	jsonContent, err := os.ReadFile("testdata/fixtures/kibana_objects/create_search_for_datasource.json")
	if err != nil {
		log.Fatal(err)
	}

	uniqueJsonData := utils.MakeKibanaObjectDataUnique(string(jsonContent))
	escapedJson := strings.ReplaceAll(uniqueJsonData, `"`, `\"`)

	return fmt.Sprintf(`resource "logzio_kibana_object" "test_kb_for_datasource" {
  kibana_version = "7.2.1"
  data = "%s"
}
`, escapedJson)
}
