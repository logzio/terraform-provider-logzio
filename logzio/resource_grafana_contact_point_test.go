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

func TestAccLogzioGrafanaContactPoint_GrafanaPointPagerDuty(t *testing.T) {
	defer utils.SleepAfterTest()
	resourceFullName := "logzio_grafana_contact_point.test_cp_pagerduty"
	apiTokenCreate := "some_api"
	apiTokenUpdate := "other"
	severityCreate := "info"
	severityUpdate := "warning"
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				// Create
				Config: getGrafanaContactPointConfigPagerduty(apiTokenCreate, severityCreate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, grafanaContactPointUid),
					resource.TestCheckResourceAttr(resourceFullName, grafanaContactPointName, "my-pagerduty-cp"),
					resource.TestCheckResourceAttr(resourceFullName, grafanaContactPointDisableResolveMessage, "false"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointPagerduty, grafanaContactPointPagerdutySeverity), severityCreate),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointPagerduty, grafanaContactPointPagerdutyIntegrationKey), apiTokenCreate),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointPagerduty, grafanaContactPointPagerdutyClass), "some-class"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointPagerduty, grafanaContactPointPagerdutyComponent), "some-component"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointPagerduty, grafanaContactPointPagerdutyGroup), "some-group"),
				),
			},
			{
				// Update
				Config: getGrafanaContactPointConfigPagerduty(apiTokenCreate, severityUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, grafanaContactPointUid),
					resource.TestCheckResourceAttr(resourceFullName, grafanaContactPointName, "my-pagerduty-cp"),
					resource.TestCheckResourceAttr(resourceFullName, grafanaContactPointDisableResolveMessage, "false"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointPagerduty, grafanaContactPointPagerdutySeverity), severityUpdate),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointPagerduty, grafanaContactPointPagerdutyIntegrationKey), apiTokenCreate),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointPagerduty, grafanaContactPointPagerdutyClass), "some-class"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointPagerduty, grafanaContactPointPagerdutyComponent), "some-component"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointPagerduty, grafanaContactPointPagerdutyGroup), "some-group"),
				),
			},
			{
				// Update sensitive
				Config: getGrafanaContactPointConfigPagerduty(apiTokenUpdate, severityUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, grafanaContactPointUid),
					resource.TestCheckResourceAttr(resourceFullName, grafanaContactPointName, "my-pagerduty-cp"),
					resource.TestCheckResourceAttr(resourceFullName, grafanaContactPointDisableResolveMessage, "false"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointPagerduty, grafanaContactPointPagerdutySeverity), severityUpdate),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointPagerduty, grafanaContactPointPagerdutyIntegrationKey), apiTokenUpdate),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointPagerduty, grafanaContactPointPagerdutyClass), "some-class"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointPagerduty, grafanaContactPointPagerdutyComponent), "some-component"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointPagerduty, grafanaContactPointPagerdutyGroup), "some-group"),
				),
			},
			{
				ResourceName:            resourceFullName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{fmt.Sprintf("%s.0.%s", grafanaContactPointPagerduty, grafanaContactPointPagerdutyIntegrationKey)},
			},
		},
	})
}

func TestAccLogzioGrafanaContactPoint_GrafanaPointSlack(t *testing.T) {
	defer utils.SleepAfterTest()
	resourceFullName := "logzio_grafana_contact_point.test_cp_slack"
	tokenCreate := "some_api"
	tokenUpdate := "other"
	mentionChannelCreate := "here"
	mentionChannelUpdate := ""
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				// Create
				Config: getGrafanaContactPointConfigSlack(tokenCreate, mentionChannelCreate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, grafanaContactPointUid),
					resource.TestCheckResourceAttr(resourceFullName, grafanaContactPointName, "my-slack-cp"),
					resource.TestCheckResourceAttr(resourceFullName, grafanaContactPointDisableResolveMessage, "false"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointSlack, grafanaContactPointSlackMentionChannel), mentionChannelCreate),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointSlack, grafanaContactPointSlackToken), tokenCreate),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointSlack, grafanaContactPointSlackTitle), "some-title"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointSlack, grafanaContactPointSlackText), "{{ len .Alerts.Firing }} firing."),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointSlack, grafanaContactPointSlackRecipient), "me"),
				),
			},
			{
				// Update
				Config: getGrafanaContactPointConfigSlack(tokenCreate, mentionChannelUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, grafanaContactPointUid),
					resource.TestCheckResourceAttr(resourceFullName, grafanaContactPointName, "my-slack-cp"),
					resource.TestCheckResourceAttr(resourceFullName, grafanaContactPointDisableResolveMessage, "false"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointSlack, grafanaContactPointSlackMentionChannel), mentionChannelUpdate),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointSlack, grafanaContactPointSlackToken), tokenCreate),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointSlack, grafanaContactPointSlackTitle), "some-title"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointSlack, grafanaContactPointSlackText), "{{ len .Alerts.Firing }} firing."),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointSlack, grafanaContactPointSlackRecipient), "me"),
				),
			},
			{
				// Update sensitive
				Config: getGrafanaContactPointConfigSlack(tokenUpdate, mentionChannelUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, grafanaContactPointUid),
					resource.TestCheckResourceAttr(resourceFullName, grafanaContactPointName, "my-slack-cp"),
					resource.TestCheckResourceAttr(resourceFullName, grafanaContactPointDisableResolveMessage, "false"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointSlack, grafanaContactPointSlackMentionChannel), mentionChannelUpdate),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointSlack, grafanaContactPointSlackToken), tokenUpdate),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointSlack, grafanaContactPointSlackTitle), "some-title"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointSlack, grafanaContactPointSlackText), "{{ len .Alerts.Firing }} firing."),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointSlack, grafanaContactPointSlackRecipient), "me"),
				),
			},
			{
				ResourceName:            resourceFullName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{fmt.Sprintf("%s.0.%s", grafanaContactPointSlack, grafanaContactPointSlackToken)},
			},
		},
	})
}

func TestAccLogzioGrafanaContactPoint_GrafanaPointTeams(t *testing.T) {
	defer utils.SleepAfterTest()
	resourceFullName := "logzio_grafana_contact_point.test_cp_teams"
	urlCreate := "some.url"
	urlUpdate := "another.url"
	messageCreate := "some message"
	messageUpdate := "another message"
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				// Create
				Config: getGrafanaContactPointConfigTeams(urlCreate, messageCreate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, grafanaContactPointUid),
					resource.TestCheckResourceAttr(resourceFullName, grafanaContactPointName, "my-teams-cp"),
					resource.TestCheckResourceAttr(resourceFullName, grafanaContactPointDisableResolveMessage, "false"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointMicrosoftTeams, grafanaContactPointMicrosoftTeamsUrl), urlCreate),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointMicrosoftTeams, grafanaContactPointMicrosoftTeamsMessage), messageCreate),
				),
			},
			{
				// Update
				Config: getGrafanaContactPointConfigTeams(urlUpdate, messageCreate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, grafanaContactPointUid),
					resource.TestCheckResourceAttr(resourceFullName, grafanaContactPointName, "my-teams-cp"),
					resource.TestCheckResourceAttr(resourceFullName, grafanaContactPointDisableResolveMessage, "false"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointMicrosoftTeams, grafanaContactPointMicrosoftTeamsUrl), urlUpdate),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointMicrosoftTeams, grafanaContactPointMicrosoftTeamsMessage), messageCreate),
				),
			},
			{
				// Update sensitive
				Config: getGrafanaContactPointConfigTeams(urlUpdate, messageUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceFullName, grafanaContactPointUid),
					resource.TestCheckResourceAttr(resourceFullName, grafanaContactPointName, "my-teams-cp"),
					resource.TestCheckResourceAttr(resourceFullName, grafanaContactPointDisableResolveMessage, "false"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointMicrosoftTeams, grafanaContactPointMicrosoftTeamsUrl), urlUpdate),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s", grafanaContactPointMicrosoftTeams, grafanaContactPointMicrosoftTeamsMessage), messageUpdate),
				),
			},
			{
				ResourceName:            resourceFullName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{fmt.Sprintf("%s.0.%s", grafanaContactPointMicrosoftTeams, grafanaContactPointMicrosoftTeamsUrl)},
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

func getGrafanaContactPointConfigPagerduty(token, severity string) string {
	return fmt.Sprintf(`
resource "logzio_grafana_contact_point" "test_cp_pagerduty" {
	name = "my-pagerduty-cp"
	disable_resolve_message = false
	pagerduty {
		integration_key = "%s"
		class = "some-class"
		component = "some-component"
		group = "some-group"
		severity = "%s"
	}
}
`, token, severity)
}

func getGrafanaContactPointConfigSlack(token, mentionChannel string) string {
	return fmt.Sprintf(`
resource "logzio_grafana_contact_point" "test_cp_slack" {
	name = "my-slack-cp"
	disable_resolve_message = false
	slack {
		token = "%s"
		title = "some-title"
		text = "{{ len .Alerts.Firing }} firing."
		mention_channel = "%s"
		recipient = "me"
	}
}
`, token, mentionChannel)
}

func getGrafanaContactPointConfigTeams(url, message string) string {
	return fmt.Sprintf(`
resource "logzio_grafana_contact_point" "test_cp_teams" {
	name = "my-teams-cp"
	disable_resolve_message = false
	teams {
		url = "%s"
		message = "%s"
	}
}
`, url, message)
}