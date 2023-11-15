package logzio

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
	"testing"
)

func TestAccLogzioGrafanaContactPoint_GrafanaPointEmail(t *testing.T) {
	defer utils.SleepAfterTest()
	resourceFullName := "logzio_grafana_contact_point.test_cp_email"
	emailsCreate := "[\"example@example.com\", \"example2@example.com\"]"
	emailsUpdate := "[\"example@example.com\", \"example2@example.com\", \"example3@example.com\"]"
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				// Create
				Config: getGrafanaContactPointConfigEmail(emailsCreate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, grafanaContactPointUid),
					resource.TestCheckResourceAttr(resourceFullName, grafanaContactPointName, "my-email-cp"),
					resource.TestCheckResourceAttr(resourceFullName, grafanaContactPointDisableResolveMessage, "false"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointEmail, grafanaContactPointEmailSingleEmail), "true"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s.#", grafanaContactPointEmail, grafanaContactPointEmailAddresses), "2"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointEmail, grafanaContactPointEmailMessage), "{{ len .Alerts.Firing }} firing."),
				),
			},
			{
				// Update
				Config: getGrafanaContactPointConfigEmail(emailsUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, grafanaContactPointUid),
					resource.TestCheckResourceAttr(resourceFullName, grafanaContactPointName, "my-email-cp"),
					resource.TestCheckResourceAttr(resourceFullName, grafanaContactPointDisableResolveMessage, "false"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointEmail, grafanaContactPointEmailSingleEmail), "true"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s.#", grafanaContactPointEmail, grafanaContactPointEmailAddresses), "3"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointEmail, grafanaContactPointEmailMessage), "{{ len .Alerts.Firing }} firing."),
				),
			},
			{
				ResourceName:      resourceFullName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})

}

func getGrafanaContactPointConfigEmail(emails string) string {
	return fmt.Sprintf(`
resource "logzio_grafana_contact_point" "test_cp_email" {
	name = "my-email-cp"
	disable_resolve_message = false
	email {
		addresses = %s
		single_email = true
		message = "{{ len .Alerts.Firing }} firing."
	}
}
`, emails)
}
