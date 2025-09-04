package logzio

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
)

func TestAccDataSourceKibanaObject(t *testing.T) {
	defer utils.SleepAfterTest()

	suffix := getRandomId()
	dataSourceConf := fmt.Sprintf(utils.ReadFixtureFromFile("kibana_object_datasource.tf"), suffix)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:                    getResourceConfigKibanaObject(suffix),
				PreventPostDestroyRefresh: true,
			},
			{
				Config:                    getResourceConfigKibanaObject(suffix) + dataSourceConf,
				PreventPostDestroyRefresh: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.logzio_kibana_object.ds_kb", kibanaObjectKibanaVersionField, "7.2.1"),
					resource.TestMatchOutput("output_id", regexp.MustCompile("search:tf-provider-datasource-test-search")),
				),
			},
		},
	})
}

func getResourceConfigKibanaObject(randomSuffix string) string {
	content, err := os.ReadFile("./testdata/fixtures/kibana_objects/create_search_for_datasource.json")

	if err != nil {
		log.Fatal(err)
	}

	var obj interface{}
	var updated []byte
	if err := json.Unmarshal(content, &obj); err != nil {
		log.Print("could not unmarshal json content, returning original content")
		updated = content
	} else {
		updateFieldsRecursive(obj, randomSuffix)
		updated, err = json.Marshal(obj)
		if err != nil {
			log.Print("could not marshal updated content, returning original content")
			updated = content
		}
	}

	return fmt.Sprintf(`resource "logzio_kibana_object" "test_kb_for_datasource" {
  kibana_version = "7.2.1"
  data = <<EOF
%s
EOF
}
`, string(updated))
}
