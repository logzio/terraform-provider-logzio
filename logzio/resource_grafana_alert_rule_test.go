package logzio

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
	"os"
	"testing"
)

func TestAccLogzioGrafanaAlertRule_CreateUpdateDashboard(t *testing.T) {
	defer utils.SleepAfterTest()
	folderUid := os.Getenv(grafanaFolderIdEnv)
	resourceFullName := "logzio_grafana_alert_rule.test_grafana_alert"
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				// Create
				Config: getGrafanaAlertRuleConfigCreate(folderUid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, grafanaAlertRuleUid),
					resource.TestCheckResourceAttrSet(resourceFullName, grafanaAlertRuleId),
					resource.TestCheckResourceAttrSet(resourceFullName, grafanaAlertRuleFolderUid),
					resource.TestCheckResourceAttr(resourceFullName, grafanaAlertRuleTitle, "my_grafana_alert"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.foo", grafanaAlertRuleAnnotations), "bar"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.hello", grafanaAlertRuleAnnotations), "world"),
					resource.TestCheckResourceAttr(resourceFullName, grafanaAlertRuleCondition, "A"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaAlertRuleData, grafanaAlertRuleDataModel), "{\"hide\":false,\"refId\":\"B\"}"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.hey", grafanaAlertRuleLabels), "oh"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.lets", grafanaAlertRuleLabels), "go"),
				),
			},
			{
				// Import
				ResourceName:      resourceFullName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Update title
				Config: getGrafanaAlertRuleConfigUpdateTitle(folderUid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, grafanaAlertRuleUid),
					resource.TestCheckResourceAttrSet(resourceFullName, grafanaAlertRuleId),
					resource.TestCheckResourceAttrSet(resourceFullName, grafanaAlertRuleFolderUid),
					resource.TestCheckResourceAttr(resourceFullName, grafanaAlertRuleTitle, "updated_title"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.foo", grafanaAlertRuleAnnotations), "bar"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.hello", grafanaAlertRuleAnnotations), "world"),
					resource.TestCheckResourceAttr(resourceFullName, grafanaAlertRuleCondition, "A"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaAlertRuleData, grafanaAlertRuleDataModel), "{\"hide\":false,\"refId\":\"B\"}"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.hey", grafanaAlertRuleLabels), "oh"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.lets", grafanaAlertRuleLabels), "go"),
				),
			},
			{
				// Update model
				Config: getGrafanaAlertRuleConfigUpdateModel(folderUid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, grafanaAlertRuleUid),
					resource.TestCheckResourceAttrSet(resourceFullName, grafanaAlertRuleId),
					resource.TestCheckResourceAttrSet(resourceFullName, grafanaAlertRuleFolderUid),
					resource.TestCheckResourceAttr(resourceFullName, grafanaAlertRuleTitle, "updated_title"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.foo", grafanaAlertRuleAnnotations), "bar"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.hello", grafanaAlertRuleAnnotations), "world"),
					resource.TestCheckResourceAttr(resourceFullName, grafanaAlertRuleCondition, "A"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaAlertRuleData, grafanaAlertRuleDataModel), "{\"hide\":false,\"intervalMs\":2000,\"refId\":\"B\"}"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.hey", grafanaAlertRuleLabels), "oh"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.lets", grafanaAlertRuleLabels), "go"),
				),
			},
		},
	})
}

func getGrafanaAlertRuleConfigCreate(folderUid string) string {
	return fmt.Sprintf(`
resource "logzio_grafana_alert_rule" "test_grafana_alert" {
  annotations = {
    "foo" = "bar"
    "hello" = "world"
  }
  condition = "A"
  data {
    ref_id = "B"
    datasource_uid = "AB1C234567D89012E"
    query_type = ""
    model = jsonencode({
      hide          = false
      refId         = "B"
    })
    relative_time_range {
      from = 700
      to   = 0
    }
  }
  labels = {
    "hey" = "oh"
    "lets" = "go"
  }
  is_paused = false
  exec_err_state = "Alerting"
  folder_uid = "%s"
  for = 3
  no_data_state = "OK"
  rule_group = "rule_group_1"
  title = "my_grafana_alert"
}
`, folderUid)
}

func getGrafanaAlertRuleConfigUpdateTitle(folderUid string) string {
	return fmt.Sprintf(`
resource "logzio_grafana_alert_rule" "test_grafana_alert" {
  annotations = {
    "foo" = "bar"
    "hello" = "world"
  }
  condition = "A"
  data {
    ref_id = "B"
    datasource_uid = "AB1C234567D89012E"
    query_type = ""
    model = jsonencode({
      hide          = false
      refId         = "B"
    })
    relative_time_range {
      from = 700
      to   = 0
    }
  }
  labels = {
    "hey" = "oh"
    "lets" = "go"
  }
  is_paused = false
  exec_err_state = "Alerting"
  folder_uid = "%s"
  for = 3
  no_data_state = "OK"
  rule_group = "rule_group_1"
  title = "updated_title"
}
`, folderUid)
}

func getGrafanaAlertRuleConfigUpdateModel(folderUid string) string {
	return fmt.Sprintf(`
resource "logzio_grafana_alert_rule" "test_grafana_alert" {
  annotations = {
    "foo" = "bar"
    "hello" = "world"
  }
  condition = "A"
  data {
    ref_id = "B"
    datasource_uid = "AB1C234567D89012E"
    query_type = ""
    model = jsonencode({
      hide          = false
      intervalMs    = 2000
      refId         = "B"
    })
    relative_time_range {
      from = 700
      to   = 0
    }
  }
  labels = {
    "hey" = "oh"
    "lets" = "go"
  }
  is_paused = false
  exec_err_state = "Alerting"
  folder_uid = "%s"
  for = 3
  no_data_state = "OK"
  rule_group = "rule_group_1"
  title = "updated_title"
}
`, folderUid)
}
