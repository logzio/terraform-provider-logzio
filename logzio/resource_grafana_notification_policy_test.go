package logzio

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
	"regexp"
	"testing"
)

const (
	grafanaDefaultReceiver = "grafana-default-email"
)

func TestAccLogzioGrafanaNotificationPolicy_ManageGrafanaNotificationPolicy(t *testing.T) {
	defer utils.SleepAfterTest()
	resourceFullName := "logzio_grafana_notification_policy.test_np"
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				// Create
				Config: getGrafanaNotificationPolicyConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, grafanaNotificationPolicyContactPoint, grafanaDefaultReceiver),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.#", grafanaNotificationPolicyGroupBy), "1"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0", grafanaNotificationPolicyGroupBy), "p8s_logzio_name"),
					resource.TestCheckResourceAttr(resourceFullName, grafanaNotificationPolicyGroupWait, "50s"),
					resource.TestCheckResourceAttr(resourceFullName, grafanaNotificationPolicyGroupInterval, "7m"),
					resource.TestCheckResourceAttr(resourceFullName, grafanaNotificationPolicyRepeatInterval, "4h"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.#", grafanaNotificationPolicyPolicy), "1"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s.0.%s",
						grafanaNotificationPolicyPolicy,
						grafanaNotificationPolicyMatcher,
						grafanaNotificationPolicyMatcherLabel), "some_label"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s.0.%s",
						grafanaNotificationPolicyPolicy,
						grafanaNotificationPolicyMatcher,
						grafanaNotificationPolicyMatcherMatch), "="),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s.0.%s",
						grafanaNotificationPolicyPolicy,
						grafanaNotificationPolicyMatcher,
						grafanaNotificationPolicyMatcherValue), "some_value"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s",
						grafanaNotificationPolicyPolicy,
						grafanaNotificationPolicyContinue), "true"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s",
						grafanaNotificationPolicyPolicy,
						grafanaNotificationPolicyContactPoint), grafanaDefaultReceiver),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s",
						grafanaNotificationPolicyPolicy,
						grafanaNotificationPolicyGroupWait), "50s"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s.0.%s.0.%s",
						grafanaNotificationPolicyPolicy,
						grafanaNotificationPolicyPolicy,
						grafanaNotificationPolicyMatcher,
						grafanaNotificationPolicyMatcherLabel), "another_label"),
				),
			},
			{
				// Update
				Config: getGrafanaNotificationPolicyConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, grafanaNotificationPolicyContactPoint, grafanaDefaultReceiver),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.#", grafanaNotificationPolicyGroupBy), "2"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0", grafanaNotificationPolicyGroupBy), "p8s_logzio_name"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.1", grafanaNotificationPolicyGroupBy), "new_new"),
					resource.TestCheckResourceAttr(resourceFullName, grafanaNotificationPolicyGroupWait, "50s"),
					resource.TestCheckResourceAttr(resourceFullName, grafanaNotificationPolicyGroupInterval, "7m"),
					resource.TestCheckResourceAttr(resourceFullName, grafanaNotificationPolicyRepeatInterval, "4h"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.#", grafanaNotificationPolicyPolicy), "1"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s.0.%s",
						grafanaNotificationPolicyPolicy,
						grafanaNotificationPolicyMatcher,
						grafanaNotificationPolicyMatcherLabel), "some_label"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s.0.%s",
						grafanaNotificationPolicyPolicy,
						grafanaNotificationPolicyMatcher,
						grafanaNotificationPolicyMatcherMatch), "="),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s.0.%s",
						grafanaNotificationPolicyPolicy,
						grafanaNotificationPolicyMatcher,
						grafanaNotificationPolicyMatcherValue), "some_value"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s",
						grafanaNotificationPolicyPolicy,
						grafanaNotificationPolicyContinue), "true"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s",
						grafanaNotificationPolicyPolicy,
						grafanaNotificationPolicyContactPoint), grafanaDefaultReceiver),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s",
						grafanaNotificationPolicyPolicy,
						grafanaNotificationPolicyGroupWait), "50s"),
					resource.TestCheckResourceAttr(resourceFullName, fmt.Sprintf("%s.0.%s.0.%s.0.%s",
						grafanaNotificationPolicyPolicy,
						grafanaNotificationPolicyPolicy,
						grafanaNotificationPolicyMatcher,
						grafanaNotificationPolicyMatcherLabel), "another_label"),
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

func TestAccLogzioGrafanaNotificationPolicy_InvalidMatchType(t *testing.T) {
	defer utils.SleepAfterTest()
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      getGrafanaNotificationPolicyConfigInvalidMatchType(),
				ExpectError: regexp.MustCompile("Match type .* is not in the allowed match types list"),
			},
		},
	})
}

func getGrafanaNotificationPolicyConfig() string {
	return fmt.Sprintf(`
resource logzio_grafana_notification_policy test_np {
	contact_point = "%s"
	group_by = ["p8s_logzio_name"]
	group_wait      = "50s"
  	group_interval  = "7m"
  	repeat_interval = "4h"

  	policy {
		matcher {
		  label = "some_label"
		  match = "="
		  value = "some_value"
		}
		contact_point = "%s"
		continue      = true
	
		group_wait      = "50s"
		group_interval  = "7m"
		repeat_interval = "4h"
	
		policy {
		  matcher {
			label = "another_label"
			match = "="
			value = "another_value"
		  }
		  contact_point = "%s"
		}
  }
}
`, grafanaDefaultReceiver, grafanaDefaultReceiver, grafanaDefaultReceiver)
}

func getGrafanaNotificationPolicyConfigUpdate() string {
	return fmt.Sprintf(`
resource logzio_grafana_notification_policy test_np {
	contact_point = "%s"
	group_by = ["p8s_logzio_name", "new_new"]
	group_wait      = "50s"
  	group_interval  = "7m"
  	repeat_interval = "4h"

  	policy {
		matcher {
		  label = "some_label"
		  match = "="
		  value = "some_value"
		}
		contact_point = "%s"
		continue      = true
	
		group_wait      = "50s"
		group_interval  = "7m"
		repeat_interval = "4h"
	
		policy {
		  matcher {
			label = "another_label"
			match = "="
			value = "another_value"
		  }
		  contact_point = "%s"
		}
  }
}
`, grafanaDefaultReceiver, grafanaDefaultReceiver, grafanaDefaultReceiver)
}

func getGrafanaNotificationPolicyConfigInvalidMatchType() string {
	return fmt.Sprintf(`
resource logzio_grafana_notification_policy test_np {
	contact_point = "%s"
	group_by = ["p8s_logzio_name", "new_new"]
	group_wait      = "50s"
  	group_interval  = "7m"
  	repeat_interval = "4h"

  	policy {
		matcher {
		  label = "some_label"
		  match = "@@"
		  value = "some_value"
		}
		contact_point = "%s"
		continue      = true
	
		group_wait      = "50s"
		group_interval  = "7m"
		repeat_interval = "4h"
	
		policy {
		  matcher {
			label = "another_label"
			match = "="
			value = "another_value"
		  }
		  contact_point = "%s"
		}
  }
}
`, grafanaDefaultReceiver, grafanaDefaultReceiver, grafanaDefaultReceiver)
}
