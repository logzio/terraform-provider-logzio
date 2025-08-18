package logzio

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
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
					testCheckKibanaObjectIdContains(resourceName, "search:tf-provider-test-search"),
					testCheckKibanaObjectTitleContains(resourceName, "tf provider test create search"),
				),
			},
			{
				Config: resourceTestKibanaObject(kibana_objects.ExportTypeSearch, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, kibanaObjectKibanaVersionField, "7.2.1"),
					testCheckKibanaObjectIdContains(resourceName, "search:tf-provider-test-search"),
					testCheckKibanaObjectTitleContains(resourceName, "tf provider test update search"),
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
					testCheckKibanaObjectIdContains(resourceName, "visualization:tf-provider-test-visualization"),
					testCheckKibanaObjectTitleContains(resourceName, "tf provider test create visualization"),
				),
			},
			{
				Config: resourceTestKibanaObject(kibana_objects.ExportTypeVisualization, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, kibanaObjectKibanaVersionField, "7.2.1"),
					testCheckKibanaObjectIdContains(resourceName, "visualization:tf-provider-test-visualization"),
					testCheckKibanaObjectTitleContains(resourceName, "tf provider test update visualization"),
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
					testCheckKibanaObjectIdContains(resourceName, "dashboard:tf-provider-test-dashboard"),
					testCheckKibanaObjectTitleContains(resourceName, "tf provider test create dashboard"),
				),
			},
			{
				Config: resourceTestKibanaObject(kibana_objects.ExportTypeDashboard, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, kibanaObjectKibanaVersionField, "7.2.1"),
					testCheckKibanaObjectIdContains(resourceName, "dashboard:tf-provider-test-dashboard"),
					testCheckKibanaObjectTitleContains(resourceName, "tf provider test update dashboard"),
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

	configStr := fmt.Sprintf("%s", content)

	// For configs that reference JSON files directly, we need to replace them with inline data
	if strings.Contains(configStr, "file(\"./testdata/fixtures/kibana_objects/") {
		var jsonPath string
		if strings.Contains(path, "search") {
			jsonPath = "testdata/fixtures/kibana_objects/create_search.json"
		} else if strings.Contains(path, "visualization") {
			jsonPath = "testdata/fixtures/kibana_objects/create_visualization.json"
		} else if strings.Contains(path, "dashboard") {
			jsonPath = "testdata/fixtures/kibana_objects/create_dashboard.json"
		}

		if jsonPath != "" {
			jsonContent, err := os.ReadFile(jsonPath)
			if err != nil {
				log.Fatal(err)
			}

			uniqueJsonData := utils.MakeKibanaObjectDataUnique(string(jsonContent))
			// Escape the JSON for Terraform string
			escapedJson := strings.ReplaceAll(uniqueJsonData, `"`, `\"`)

			// Replace file() reference with inline data
			configStr = strings.ReplaceAll(configStr,
				"file(\"./testdata/fixtures/kibana_objects/create_search.json\")",
				fmt.Sprintf(`"%s"`, escapedJson))
			configStr = strings.ReplaceAll(configStr,
				"file(\"./testdata/fixtures/kibana_objects/create_visualization.json\")",
				fmt.Sprintf(`"%s"`, escapedJson))
			configStr = strings.ReplaceAll(configStr,
				"file(\"./testdata/fixtures/kibana_objects/create_dashboard.json\")",
				fmt.Sprintf(`"%s"`, escapedJson))
		}
	}

	return configStr
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

// testCheckKibanaObjectIdContains checks if the object ID contains the expected base ID
func testCheckKibanaObjectIdContains(name, expectedBaseId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No object ID is set")
		}

		if !strings.Contains(rs.Primary.ID, expectedBaseId) {
			return fmt.Errorf("expected ID to contain %s but got %s", expectedBaseId, rs.Primary.ID)
		}
		return nil
	}
}

// testCheckKibanaObjectTitleContains checks if the object title contains the expected base title
func testCheckKibanaObjectTitleContains(name, expectedBaseTitle string) resource.TestCheckFunc {
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
		if !strings.Contains(title, expectedBaseTitle) {
			return fmt.Errorf("expected title to contain %s but got %s", expectedBaseTitle, title)
		}
		return nil
	}
}
