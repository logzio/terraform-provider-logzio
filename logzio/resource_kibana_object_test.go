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

const (
	kibanaObjectCreateDashboardJsonFileName     = "create_dashboard"
	kibanaObjectCreateVisualizationJsonFileName = "create_visualization"
	kibanaObjectCreateSearchJsonFileName        = "create_search"
	kibanaObjectUpdateDashboardJsonFileName     = "update_dashboard"
	kibanaObjectUpdateVisualizationJsonFileName = "update_visualization"
	kibanaObjectUpdateSearchJsonFileName        = "update_search"
)

func TestAccLogzioKibanaObject_CreateUpdateSearch(t *testing.T) {
	resourceName := "logzio_kibana_object.test_kb_obj_search"
	defer utils.SleepAfterTest()

	suffix := getRandomId()
	resId := fmt.Sprintf("search:tf-provider-test-search-%s", suffix)
	titleCreate := fmt.Sprintf("tf provider test create search-%s", suffix)
	titleUpdate := fmt.Sprintf("tf provider test update search-%s", suffix)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceTestKibanaObject(kibana_objects.ExportTypeSearch, false, suffix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, kibanaObjectKibanaVersionField, "7.2.1"),
					resource.TestCheckResourceAttr(resourceName, "id", resId),
					testCheckKibanaObjectTitle(resourceName, titleCreate),
				),
			},
			{
				Config: resourceTestKibanaObject(kibana_objects.ExportTypeSearch, true, suffix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, kibanaObjectKibanaVersionField, "7.2.1"),
					resource.TestCheckResourceAttr(resourceName, "id", resId),
					testCheckKibanaObjectTitle(resourceName, titleUpdate),
				),
			},
		},
	})
}

func TestAccLogzioKibanaObject_CreateUpdateVisualization(t *testing.T) {
	resourceName := "logzio_kibana_object.test_kb_obj_visualization"
	defer utils.SleepAfterTest()

	suffix := getRandomId()
	resId := fmt.Sprintf("visualization:tf-provider-test-visualization-%s", suffix)
	titleCreate := fmt.Sprintf("tf provider test create visualization-%s", suffix)
	titleUpdate := fmt.Sprintf("tf provider test update visualization-%s", suffix)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceTestKibanaObject(kibana_objects.ExportTypeVisualization, false, suffix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, kibanaObjectKibanaVersionField, "7.2.1"),
					resource.TestCheckResourceAttr(resourceName, "id", resId),
					testCheckKibanaObjectTitle(resourceName, titleCreate),
				),
			},
			{
				Config: resourceTestKibanaObject(kibana_objects.ExportTypeVisualization, true, suffix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, kibanaObjectKibanaVersionField, "7.2.1"),
					resource.TestCheckResourceAttr(resourceName, "id", resId),
					testCheckKibanaObjectTitle(resourceName, titleUpdate),
				),
			},
		},
	})
}

func TestAccLogzioKibanaObject_CreateUpdateDashboard(t *testing.T) {
	resourceName := "logzio_kibana_object.test_kb_obj_dashboard"
	defer utils.SleepAfterTest()

	suffix := getRandomId()
	resId := fmt.Sprintf("dashboard:tf-provider-test-dashboard-%s", suffix)
	titleCreate := fmt.Sprintf("tf provider test create dashboard-%s", suffix)
	titleUpdate := fmt.Sprintf("tf provider test update dashboard-%s", suffix)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceTestKibanaObject(kibana_objects.ExportTypeDashboard, false, suffix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, kibanaObjectKibanaVersionField, "7.2.1"),
					resource.TestCheckResourceAttr(resourceName, "id", resId),
					testCheckKibanaObjectTitle(resourceName, titleCreate),
				),
			},
			{
				Config: resourceTestKibanaObject(kibana_objects.ExportTypeDashboard, true, suffix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, kibanaObjectKibanaVersionField, "7.2.1"),
					resource.TestCheckResourceAttr(resourceName, "id", resId),
					testCheckKibanaObjectTitle(resourceName, titleUpdate),
				),
			},
		},
	})
}

func resourceTestKibanaObject(objType kibana_objects.ExportType, update bool, suffix string) string {
	switch objType {
	case kibana_objects.ExportTypeSearch:
		testName := "test_kb_obj_search"
		if update {
			return getKibanaObjectResourceConfig(testName, kibanaObjectUpdateSearchJsonFileName, suffix)
		}
		return getKibanaObjectResourceConfig(testName, kibanaObjectCreateSearchJsonFileName, suffix)
	case kibana_objects.ExportTypeVisualization:
		testName := "test_kb_obj_visualization"
		if update {
			return getKibanaObjectResourceConfig(testName, kibanaObjectUpdateVisualizationJsonFileName, suffix)
		}
		return getKibanaObjectResourceConfig(testName, kibanaObjectCreateVisualizationJsonFileName, suffix)
	case kibana_objects.ExportTypeDashboard:
		testName := "test_kb_obj_dashboard"
		if update {
			return getKibanaObjectResourceConfig(testName, kibanaObjectUpdateDashboardJsonFileName, suffix)
		}
		return getKibanaObjectResourceConfig(testName, kibanaObjectCreateDashboardJsonFileName, suffix)
	default:
		// we should never get to this part
		panic("invalid kibana object type")
	}
}

func getKibanaObjectResourceConfig(resourceName, path, randomSuffix string) string {
	content, err := os.ReadFile(fmt.Sprintf("./testdata/fixtures/kibana_objects/%s.json", path))
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

	return fmt.Sprintf(`resource "logzio_kibana_object" "%s" {
		kibana_version = "7.2.1"
  		data = <<EOF
%s
EOF
}`, resourceName, string(updated))
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

func updateFieldsRecursive(obj interface{}, randomSuffix string) {
	switch v := obj.(type) {
	case map[string]interface{}:
		for key, val := range v {
			if key == "id" || key == "_id" || key == "title" {
				if strVal, ok := val.(string); ok {
					v[key] = fmt.Sprintf("%s-%s", strVal, randomSuffix)
				}
			} else {
				updateFieldsRecursive(val, randomSuffix)
			}
		}
	case []interface{}:
		for _, item := range v {
			updateFieldsRecursive(item, randomSuffix)
		}
	}
}
