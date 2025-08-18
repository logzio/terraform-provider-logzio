package logzio

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/logzio/logzio_terraform_client/kibana_objects"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
)

func TestAccLogzioKibanaObject_CreateUpdateSearch(t *testing.T) {
	resourceName := "logzio_kibana_object.test_kb_obj_search"
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceTestKibanaObject(kibana_objects.ExportTypeSearch, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, kibanaObjectKibanaVersionField, "7.2.1"),
					resource.TestCheckResourceAttr(resourceName, "id", "search:tf-provider-test-search"),
					testCheckKibanaObjectTitle(resourceName, "tf provider test create search"),
				),
			},
			{
				Config: resourceTestKibanaObject(kibana_objects.ExportTypeSearch, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, kibanaObjectKibanaVersionField, "7.2.1"),
					resource.TestCheckResourceAttr(resourceName, "id", "search:tf-provider-test-search"),
					testCheckKibanaObjectTitle(resourceName, "tf provider test update search"),
				),
			},
		},
	})
}

func TestAccLogzioKibanaObject_CreateUpdateVisualization(t *testing.T) {
	resourceName := "logzio_kibana_object.test_kb_obj_visualization"
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceTestKibanaObject(kibana_objects.ExportTypeVisualization, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, kibanaObjectKibanaVersionField, "7.2.1"),
					resource.TestCheckResourceAttr(resourceName, "id", "visualization:tf-provider-test-visualization"),
					testCheckKibanaObjectTitle(resourceName, "tf provider test create visualization"),
				),
			},
			{
				Config: resourceTestKibanaObject(kibana_objects.ExportTypeVisualization, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, kibanaObjectKibanaVersionField, "7.2.1"),
					resource.TestCheckResourceAttr(resourceName, "id", "visualization:tf-provider-test-visualization"),
					testCheckKibanaObjectTitle(resourceName, "tf provider test update visualization"),
				),
			},
		},
	})
}

func TestAccLogzioKibanaObject_CreateUpdateDashboard(t *testing.T) {
	resourceName := "logzio_kibana_object.test_kb_obj_dashboard"
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceTestKibanaObject(kibana_objects.ExportTypeDashboard, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, kibanaObjectKibanaVersionField, "7.2.1"),
					resource.TestCheckResourceAttr(resourceName, "id", "dashboard:tf-provider-test-dashboard"),
					testCheckKibanaObjectTitle(resourceName, "tf provider test create dashboard"),
				),
			},
			{
				Config: resourceTestKibanaObject(kibana_objects.ExportTypeDashboard, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, kibanaObjectKibanaVersionField, "7.2.1"),
					resource.TestCheckResourceAttr(resourceName, "id", "dashboard:tf-provider-test-dashboard"),
					testCheckKibanaObjectTitle(resourceName, "tf provider test update dashboard"),
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
	case kibana_objects.ExportTypeVisualization:
		if update {
			return getKibanaObjectResourceConfig("update_kibana_object_visualization")
		}
		return getKibanaObjectResourceConfig("create_kibana_object_visualization")
	case kibana_objects.ExportTypeDashboard:
		if update {
			return getKibanaObjectResourceConfig("update_kibana_object_dashboard")
		}
		return getKibanaObjectResourceConfig("create_kibana_object_dashboard")
	default:
		// we should never get to this part
		panic("invalid kibana object type")
	}
}

func getKibanaObjectResourceConfig(path string) string {
	content, err := os.ReadFile(fmt.Sprintf("testdata/fixtures/%s.tf", path))
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%s", content)
}

func testCheckKibanaObjectTitle(name, expectedTitle string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No object ID is set")
		}

		var dataObj map[string]interface{}
		json.Unmarshal([]byte(rs.Primary.Attributes[kibanaObjectDataField]), &dataObj)
		objType := dataObj["_source"].(map[string]interface{})["type"].(string)
		title := dataObj["_source"].(map[string]interface{})[objType].(map[string]interface{})["title"].(string)
		if title != expectedTitle {
			return fmt.Errorf("expected %s but got %s", expectedTitle, title)
		}
		return nil
	}
}
