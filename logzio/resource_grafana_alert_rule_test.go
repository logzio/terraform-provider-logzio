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
					resource.TestCheckResourceAttr(resourceFullName, grafanaAlertRuleTitle, "my_grafana_alert"),
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
      intervalMs    = 1000
      maxDataPoints = 43200
      refId         = "B"
    })
    relative_time_range {
      from = 700
      to   = 0
    }
  }
  labels = {
    "hey" = "ho"
    "lets" = "go"
  }
  is_paused = false
  exec_err_state = "Alerting"
  folder_uid = "%s"
  for = 3
  no_data_state = "OK"
  org_id = 1
  rule_group = "rule_group_1"
  title = "my_grafana_alert"
}
`, folderUid)
}
