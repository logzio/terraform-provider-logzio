package logzio

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"testing"
)

const (
	testsUrl       = "https://jsonplaceholder.typicode.com/todos/1"
	testsUrlUpdate = "https://jsonplaceholder.typicode.com/todos/2"
)

func TestAccLogzioEndpoint_SlackCreateEndpoint(t *testing.T) {
	resourceName := "logzio_endpoint.valid_slack_endpoint"
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: utils.ReadFixtureFromFile("valid_slack_endpoint.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"logzio_endpoint.valid_slack_endpoint", "title", "valid_slack_endpoint"),
					testAccCheckOutputExists("logzio_endpoint.valid_slack_endpoint", "test_id"),
					resource.TestMatchOutput("test_id", regexp.MustCompile("\\d")),
					resource.TestCheckResourceAttr(resourceName, "slack.4281379687.url", testsUrl),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLogzioEndpoint_SlackCreateInvalidEndpoint(t *testing.T) {
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      utils.ReadFixtureFromFile("invalid_slack_endpoint.tf"),
				ExpectError: regexp.MustCompile("Bad URL provided. no protocol"),
			},
		},
	})
}

func TestAccLogzioEndpoint_SlackUpdateEndpoint(t *testing.T) {
	endpointName := "test_create_slack_endpoint"
	resourceName := "logzio_endpoint." + endpointName
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: utils.ReadResourceFromFile(endpointName, "create_slack_endpoint.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "title", "slack_endpoint"),
					resource.TestCheckResourceAttr(resourceName, "slack.4281379687.url", testsUrl),
				),
			},
			{
				Config: utils.ReadResourceFromFile(endpointName, "update_slack_endpoint.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "title", "updated_slack_endpoint"),
					resource.TestCheckResourceAttr(resourceName, "slack.3558733988.url", testsUrlUpdate),
				),
			},
		},
	})
}

func TestAccLogzioEndpoint_CustomCreateEndpoint(t *testing.T) {
	resourceName := "logzio_endpoint.custom"
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: createCustomEndpoint("custom"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "title", "my_custom_title"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLogzioEndpoint_CustomCreateEndpointNoHeaders(t *testing.T) {
	config := `resource "logzio_endpoint" "custom" {
 title = "my_custom_title_no_headers"
 endpoint_type = "custom"
 description = "this_is_my_description"
 custom {
   url = "https://www.test.com"
   method = "POST"
   body_template = jsonencode({
      this: "is"
      my: "template"
    })
 }
}`
	resourceName := "logzio_endpoint.custom"
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "title", "my_custom_title_no_headers"),
					resource.TestCheckResourceAttr(resourceName, "custom.2688962798.headers", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLogzioEndpoint_CustomCreateEndpointEmptyBodyTemplate(t *testing.T) {
	config := `resource "logzio_endpoint" "custom" {
 title = "my_custom_title_empty_body_template"
 endpoint_type = "custom"
 description = "this_is_my_description"
 custom {
   url = "https://www.test.com"
   method = "POST"
   body_template = "{}"
 }
}`
	resourceName := "logzio_endpoint.custom"
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "title", "my_custom_title_empty_body_template"),
					resource.TestCheckResourceAttr(resourceName, "custom.1831890258.body_template", "{}"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLogzioEndpoint_CustomCreateEndpointInvalidMethod(t *testing.T) {
	config := `resource "logzio_endpoint" "custom" {
 title = "my_custom_title_invalid_method"
 endpoint_type = "custom"
 description = "this_is_my_description"
 custom {
   url = "https://www.test.com"
   method = "PATCH"
   body_template = jsonencode({
      this: "is"
      my: "template"
    })
 }
}`
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("invalid HTTP method specified"),
			},
		},
	})
}

func TestAccLogzioEndpoint_CustomUpdateEndpoint(t *testing.T) {
	endpointName := "test_update_custom_endpoint"
	resourceName := "logzio_endpoint." + endpointName
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: utils.ReadResourceFromFile(endpointName, "create_custom_endpoint.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "title", "my_custom_title"),
					resource.TestCheckResourceAttr(resourceName, "custom.552452041.url", testsUrl),
					resource.TestCheckResourceAttr(resourceName, "custom.552452041.method", http.MethodPost),
				),
			},
			{
				Config: utils.ReadResourceFromFile(endpointName, "update_custom_endpoint.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "title", "updated_custom_endpoint"),
					resource.TestCheckResourceAttr(resourceName, "custom.1435031376.url", testsUrlUpdate),
					resource.TestCheckResourceAttr(resourceName, "custom.1435031376.method", http.MethodPut),
					resource.TestCheckResourceAttr(resourceName, "custom.1435031376.headers", ""),
					resource.TestCheckResourceAttr(resourceName, "custom.1435031376.body_template", "{}"),
				),
			},
		},
	})
}

func TestAccLogzioEndpoint_PagerDutyCreateEndpoint(t *testing.T) {
	resourceName := "logzio_endpoint.pagerduty"
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: utils.ReadFixtureFromFile("create_pagerduty_endpoint.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "title", "my_pagerduty_title"),
					resource.TestCheckResourceAttr(resourceName, "pagerduty.1955626064.service_key", "my_service_key"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLogzioEndpoint_PagerDutyCreateEndpointEmptyServiceKey(t *testing.T) {
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      utils.ReadFixtureFromFile("create_pagerduty_endpoint_empty_service_key.tf"),
				ExpectError: regexp.MustCompile("service key must be set for type pagerduty"),
			},
		},
	})
}

func TestAccLogzioEndpoint_PagerDutyUpdateEndpoint(t *testing.T) {
	resourceName := "logzio_endpoint.pagerduty"
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: utils.ReadFixtureFromFile("create_pagerduty_endpoint.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "title", "my_pagerduty_title"),
					resource.TestCheckResourceAttr(resourceName, "pagerduty.1955626064.service_key", "my_service_key"),
				),
			},
			{
				Config: utils.ReadFixtureFromFile("update_pagerduty_endpoint.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "title", "pagerduty_title_updated"),
					resource.TestCheckResourceAttr(resourceName, "pagerduty.3330485350.service_key", "another_service_key"),
				),
			},
		},
	})
}

func TestAccLogzioEndpoint_BigPandaCreateEndpoint(t *testing.T) {
	resourceName := "logzio_endpoint.bigpanda"
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: utils.ReadFixtureFromFile("create_bigpanda_endpoint.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "title", "my_bigpanda_title"),
					resource.TestCheckResourceAttr(resourceName, "bigpanda.1922960384.api_token", "my_api_token"),
					resource.TestCheckResourceAttr(resourceName, "bigpanda.1922960384.app_key", "my_app_key"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLogzioEndpoint_BigPandaCreateEndpointEmptyApiToken(t *testing.T) {
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      utils.ReadFixtureFromFile("create_bigpanda_endpoint_empty_api_token.tf"),
				ExpectError: regexp.MustCompile("api token must be set for type bigpanda"),
			},
		},
	})
}

func TestAccLogzioEndpoint_BigPandaCreateEndpointEmptyAppKey(t *testing.T) {
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      utils.ReadFixtureFromFile("create_bigpanda_endpoint_empty_app_key.tf"),
				ExpectError: regexp.MustCompile("app key must be set for type bigpanda"),
			},
		},
	})
}

func TestAccLogzioEndpoint_BigPandaUpdateEndpoint(t *testing.T) {
	resourceName := "logzio_endpoint.bigpanda"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: utils.ReadFixtureFromFile("create_bigpanda_endpoint.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "title", "my_bigpanda_title"),
					resource.TestCheckResourceAttr(resourceName, "bigpanda.1922960384.api_token", "my_api_token"),
					resource.TestCheckResourceAttr(resourceName, "bigpanda.1922960384.app_key", "my_app_key"),
				),
			},
			{
				Config: utils.ReadFixtureFromFile("update_bigpanda_endpoint.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "title", "bigpanda_title_updated"),
					resource.TestCheckResourceAttr(resourceName, "bigpanda.1493627637.api_token", "updated_api_token"),
					resource.TestCheckResourceAttr(resourceName, "bigpanda.1493627637.app_key", "updated_app_key"),
				),
			},
		},
	})
}

func TestAccLogzioEndpoint_DataDogCreateEndpoint(t *testing.T) {
	resourceName := "logzio_endpoint.datadog"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: utils.ReadFixtureFromFile("create_datadog_endpoint.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "title", "my_datadog_title"),
					resource.TestCheckResourceAttr(resourceName, "datadog.411979392.api_key", "my_api_key"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLogzioEndpoint_DataDogCreateEndpointEmptyApiKey(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      utils.ReadFixtureFromFile("create_datadog_endpoint_empty_api_key.tf"),
				ExpectError: regexp.MustCompile("api key must be set for type datadog"),
			},
		},
	})
}

func TestAccLogzioEndpoint_DataDogUpdateEndpoint(t *testing.T) {
	resourceName := "logzio_endpoint.datadog"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: utils.ReadFixtureFromFile("create_datadog_endpoint.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "title", "my_datadog_title"),
					resource.TestCheckResourceAttr(resourceName, "datadog.411979392.api_key", "my_api_key"),
				),
			},
			{
				Config: utils.ReadFixtureFromFile("update_datadog_endpoint.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "title", "datadog_title_updated"),
					resource.TestCheckResourceAttr(resourceName, "datadog.2413799041.api_key", "updated_api_key"),
				),
			},
		},
	})
}

func TestAccLogzioEndpoint_VictorOpsCreateEndpoint(t *testing.T) {
	resourceName := "logzio_endpoint.victorops"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: utils.ReadFixtureFromFile("create_victorops_endpoint.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "title", "my_victorops_title"),
					resource.TestCheckResourceAttr(resourceName, "victorops.3725242508.routing_key", "my_routing_key"),
					resource.TestCheckResourceAttr(resourceName, "victorops.3725242508.message_type", "my_message_type"),
					resource.TestCheckResourceAttr(resourceName, "victorops.3725242508.service_api_key", "my_service_api_key"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLogzioEndpoint_VictorOpsCreateEndpointEmptyRoutingKey(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      utils.ReadFixtureFromFile("create_victorops_endpoint_empty_routing_key.tf"),
				ExpectError: regexp.MustCompile("routing key must be set for type victorops"),
			},
		},
	})
}

func TestAccLogzioEndpoint_VictorOpsCreateEndpointEmptyMessageType(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      utils.ReadFixtureFromFile("create_victorops_endpoint_empty_message_type.tf"),
				ExpectError: regexp.MustCompile("message type must be set for type victorops"),
			},
		},
	})
}

func TestAccLogzioEndpoint_VictorOpsCreateEndpointEmptyServiceApiKey(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      utils.ReadFixtureFromFile("create_victorops_endpoint_empty_service_api_key.tf"),
				ExpectError: regexp.MustCompile("service api key must be set for type victorops"),
			},
		},
	})
}

func TestAccLogzioEndpoint_VictorOpsUpdateEndpoint(t *testing.T) {
	resourceName := "logzio_endpoint.victorops"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: utils.ReadFixtureFromFile("create_victorops_endpoint.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "title", "my_victorops_title"),
					resource.TestCheckResourceAttr(resourceName, "victorops.3725242508.routing_key", "my_routing_key"),
					resource.TestCheckResourceAttr(resourceName, "victorops.3725242508.message_type", "my_message_type"),
					resource.TestCheckResourceAttr(resourceName, "victorops.3725242508.service_api_key", "my_service_api_key"),
				),
			},
			{
				Config: utils.ReadFixtureFromFile("update_victorops_endpoint.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "title", "updated_victorops_title"),
					resource.TestCheckResourceAttr(resourceName, "victorops.1896695128.routing_key", "updated_routing_key"),
					resource.TestCheckResourceAttr(resourceName, "victorops.1896695128.message_type", "updated_message_type"),
					resource.TestCheckResourceAttr(resourceName, "victorops.1896695128.service_api_key", "updated_service_api_key"),
				),
			},
		},
	})
}

func TestAccLogzioEndpoint_OpsGenieCreateEndpoint(t *testing.T) {
	resourceName := "logzio_endpoint.opsgenie"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: utils.ReadFixtureFromFile("create_opsgenie_endpoint.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "title", "my_opsgenie_title"),
					resource.TestCheckResourceAttr(resourceName, "opsgenie.411979392.api_key", "my_api_key"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLogzioEndpoint_OpsGenieCreateEndpointEmptyApiKey(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      utils.ReadFixtureFromFile("create_opsgenie_endpoint_empty_api_key.tf"),
				ExpectError: regexp.MustCompile("api key must be set for type opsgenie"),
			},
		},
	})
}

func TestAccLogzioEndpoint_OpsGenieUpdateEndpoint(t *testing.T) {
	resourceName := "logzio_endpoint.opsgenie"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: utils.ReadFixtureFromFile("create_opsgenie_endpoint.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "title", "my_opsgenie_title"),
					resource.TestCheckResourceAttr(resourceName, "opsgenie.411979392.api_key", "my_api_key"),
				),
			},
			{
				Config: utils.ReadFixtureFromFile("update_opsgenie_endpoint.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "title", "updated_opsgenie_title"),
					resource.TestCheckResourceAttr(resourceName, "opsgenie.2413799041.api_key", "updated_api_key"),
				),
			},
		},
	})
}

func TestAccLogzioEndpoint_ServiceNowCreateEndpoint(t *testing.T) {
	resourceName := "logzio_endpoint.servicenow"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: utils.ReadFixtureFromFile("create_servicenow_endpoint.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "title", "my_servicenow_title"),
					resource.TestCheckResourceAttr(resourceName, "servicenow.500131967.username", "my_username"),
					resource.TestCheckResourceAttr(resourceName, "servicenow.500131967.password", "my_password"),
					resource.TestCheckResourceAttr(resourceName, "servicenow.500131967.url", testsUrl),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLogzioEndpoint_ServiceNowCreateEndpointEmptyUsername(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      utils.ReadFixtureFromFile("create_servicenow_endpoint_empty_username.tf"),
				ExpectError: regexp.MustCompile("username must be set for type servicenow"),
			},
		},
	})
}

func TestAccLogzioEndpoint_ServiceNowCreateEndpointEmptyPassword(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      utils.ReadFixtureFromFile("create_servicenow_endpoint_empty_password.tf"),
				ExpectError: regexp.MustCompile("password must be set for type servicenow"),
			},
		},
	})
}

func TestAccLogzioEndpoint_ServiceNowCreateEndpointEmptyUrl(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      utils.ReadFixtureFromFile("create_servicenow_endpoint_empty_url.tf"),
				ExpectError: regexp.MustCompile("url must be set for type servicenow"),
			},
		},
	})
}

func TestAccLogzioEndpoint_ServiceNowUpdateEndpoint(t *testing.T) {
	resourceName := "logzio_endpoint.servicenow"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: utils.ReadFixtureFromFile("create_servicenow_endpoint.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "title", "my_servicenow_title"),
					resource.TestCheckResourceAttr(resourceName, "servicenow.500131967.username", "my_username"),
					resource.TestCheckResourceAttr(resourceName, "servicenow.500131967.password", "my_password"),
					resource.TestCheckResourceAttr(resourceName, "servicenow.500131967.url", testsUrl),
				),
			},
			{
				Config: utils.ReadFixtureFromFile("update_servicenow_endpoint.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "title", "updated_servicenow_title"),
					resource.TestCheckResourceAttr(resourceName, "servicenow.1337040240.username", "updated_username"),
					resource.TestCheckResourceAttr(resourceName, "servicenow.1337040240.password", "updated_password"),
					resource.TestCheckResourceAttr(resourceName, "servicenow.1337040240.url", testsUrlUpdate),
				),
			},
		},
	})
}

func TestAccLogzioEndpoint_MicrosoftTeamsCreateEndpoint(t *testing.T) {
	resourceName := "logzio_endpoint.microsoftteams"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: utils.ReadFixtureFromFile("create_microsoftteams_endpoint.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "title", "my_microsoftteams_title"),
					resource.TestCheckResourceAttr(resourceName, "microsoftteams.4281379687.url", testsUrl),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLogzioEndpoint_MicrosoftTeamsCreateEndpointEmptyUrl(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      utils.ReadFixtureFromFile("create_microsoftteams_endpoint_empty_url.tf"),
				ExpectError: regexp.MustCompile("url must be set for type microsoftteams"),
			},
		},
	})
}

func TestAccLogzioEndpoint_MicrosoftTeamsUpdateEndpoint(t *testing.T) {
	resourceName := "logzio_endpoint.microsoftteams"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: utils.ReadFixtureFromFile("create_microsoftteams_endpoint.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "title", "my_microsoftteams_title"),
					resource.TestCheckResourceAttr(resourceName, "microsoftteams.4281379687.url", testsUrl),
				),
			},
			{
				Config: utils.ReadFixtureFromFile("update_microsoftteams_endpoint.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "title", "updated_microsoftteams_title"),
					resource.TestCheckResourceAttr(resourceName, "microsoftteams.3558733988.url", testsUrlUpdate),
				),
			},
		},
	})
}

func testAccCheckOutputExists(n string, o string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		id := rs.Primary.ID
		os, ok := s.RootModule().Outputs[o]

		if rs.Primary.ID == "" {
			return errors.New("no endpoint ID is set")
		}

		if os.Value != id {
			return fmt.Errorf("can't find resource that matches output ID")
		}

		return nil
	}
}

func createCustomEndpoint(name string) string {
	content, err := ioutil.ReadFile("testdata/fixtures/create_custom_endpoint.tf")
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf(fmt.Sprintf("%s", content), name)
}
