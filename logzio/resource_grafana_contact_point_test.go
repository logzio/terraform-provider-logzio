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

func TestAccLogzioGrafanaContactPoint_GrafanaPointGoogleChat(t *testing.T) {
	defer utils.SleepAfterTest()
	resourceFullName := "logzio_grafana_contact_point.test_cp_googlechat"
	urlCreate := "some.url"
	urlUpdate := "other.url"
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				// Create
				Config: getGrafanaContactPointConfigGoogleChat(urlCreate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, grafanaContactPointUid),
					resource.TestCheckResourceAttr(resourceFullName, grafanaContactPointName, "my-googlechat-cp"),
					resource.TestCheckResourceAttr(resourceFullName, grafanaContactPointDisableResolveMessage, "false"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointGoogleChat, grafanaContactPointGoogleChatUrl), urlCreate),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointGoogleChat, grafanaContactPointGoogleChatMessage), "{{ len .Alerts.Firing }} firing."),
				),
			},
			{
				// Update
				Config: getGrafanaContactPointConfigGoogleChat(urlUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, grafanaContactPointUid),
					resource.TestCheckResourceAttr(resourceFullName, grafanaContactPointName, "my-googlechat-cp"),
					resource.TestCheckResourceAttr(resourceFullName, grafanaContactPointDisableResolveMessage, "false"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointGoogleChat, grafanaContactPointGoogleChatUrl), urlUpdate),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointGoogleChat, grafanaContactPointGoogleChatMessage), "{{ len .Alerts.Firing }} firing."),
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

func TestAccLogzioGrafanaContactPoint_GrafanaPointOpsgenie(t *testing.T) {
	defer utils.SleepAfterTest()
	resourceFullName := "logzio_grafana_contact_point.test_cp_opsgenie"
	urlCreate := "some.url"
	urlUpdate := "other.url"
	apiTokenCreate := "some_api"
	apiTokenUpdate := "other"
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				// Create
				Config: getGrafanaContactPointConfigOpsgenie(urlCreate, apiTokenCreate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, grafanaContactPointUid),
					resource.TestCheckResourceAttr(resourceFullName, grafanaContactPointName, "my-opsgenie-cp"),
					resource.TestCheckResourceAttr(resourceFullName, grafanaContactPointDisableResolveMessage, "false"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointOpsgenie, grafanaContactPointOpsgenieApiUrl), urlCreate),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointOpsgenie, grafanaContactPointOpsgenieApiKey), apiTokenCreate),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointOpsgenie, grafanaContactPointOpsgenieAutoClose), "false"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointOpsgenie, grafanaContactPointOpsgenieOverridePriority), "true"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointOpsgenie, grafanaContactPointOpsgenieSendTagsAs), grafanaContactPointOpsgenieSendTagsBoth),
				),
			},
			{
				// Update
				Config: getGrafanaContactPointConfigOpsgenie(urlUpdate, apiTokenCreate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, grafanaContactPointUid),
					resource.TestCheckResourceAttr(resourceFullName, grafanaContactPointName, "my-opsgenie-cp"),
					resource.TestCheckResourceAttr(resourceFullName, grafanaContactPointDisableResolveMessage, "false"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointOpsgenie, grafanaContactPointOpsgenieApiUrl), urlUpdate),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointOpsgenie, grafanaContactPointOpsgenieApiKey), apiTokenCreate),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointOpsgenie, grafanaContactPointOpsgenieAutoClose), "false"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointOpsgenie, grafanaContactPointOpsgenieOverridePriority), "true"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointOpsgenie, grafanaContactPointOpsgenieSendTagsAs), grafanaContactPointOpsgenieSendTagsBoth),
				),
			},
			{
				// Update sensitive
				Config: getGrafanaContactPointConfigOpsgenie(urlUpdate, apiTokenUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, grafanaContactPointUid),
					resource.TestCheckResourceAttr(resourceFullName, grafanaContactPointName, "my-opsgenie-cp"),
					resource.TestCheckResourceAttr(resourceFullName, grafanaContactPointDisableResolveMessage, "false"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointOpsgenie, grafanaContactPointOpsgenieApiUrl), urlUpdate),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointOpsgenie, grafanaContactPointOpsgenieApiKey), apiTokenUpdate),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointOpsgenie, grafanaContactPointOpsgenieAutoClose), "false"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointOpsgenie, grafanaContactPointOpsgenieOverridePriority), "true"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointOpsgenie, grafanaContactPointOpsgenieSendTagsAs), grafanaContactPointOpsgenieSendTagsBoth),
				),
			},
			{
				ResourceName:            resourceFullName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{fmt.Sprintf("%s.0.%s", grafanaContactPointOpsgenie, grafanaContactPointOpsgenieApiKey)},
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

func getGrafanaContactPointConfigGoogleChat(url string) string {
	return fmt.Sprintf(`
resource "logzio_grafana_contact_point" "test_cp_googlechat" {
	name = "my-googlechat-cp"
	disable_resolve_message = false
	googlechat {
		url = "%s"
		message = "{{ len .Alerts.Firing }} firing."
	}
}
`, url)
}

func getGrafanaContactPointConfigOpsgenie(url, apiToken string) string {
	return fmt.Sprintf(`
resource "logzio_grafana_contact_point" "test_cp_opsgenie" {
	name = "my-opsgenie-cp"
	disable_resolve_message = false
	opsgenie {
		api_url = "%s"
		api_key = "%s"
		auto_close = false
		override_priority = true
		send_tags_as = "both"
	}
}
`, url, apiToken)
}
