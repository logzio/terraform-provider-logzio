package logzio

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/logzio/logzio_terraform_client/grafana_notification_policies"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
)

const (
	grafanaNotificationPolicyContactPoint   = "contact_point"
	grafanaNotificationPolicyGroupBy        = "group_by"
	grafanaNotificationPolicyGroupInterval  = "group_interval"
	grafanaNotificationPolicyGroupWait      = "group_wait"
	grafanaNotificationPolicyRepeatInterval = "repeat_interval"
	grafanaNotificationPolicyPolicy         = "policy"

	grafanaNotificationPolicyMatcher      = "matcher"
	grafanaNotificationPolicyMatcherLabel = "label"
	grafanaNotificationPolicyMatcherMatch = "match"
	grafanaNotificationPolicyMatcherValue = "value"
	grafanaNotificationPolicyMuteTiming   = "mute_timing"
	grafanaNotificationPolicyContinue     = "continue"

	grafanaNotificationPolicyTreeDepth = 4
)

// resourceGrafanaNotificationPolicy represents a Grafana notification policy tree. Note that one resource represents the entire policy tree
func resourceGrafanaNotificationPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGrafanaNotificationPolicyCreate,
		//ReadContext:   resourceGrafanaNotificationPolicyRead,
		//UpdateContext: resourceGrafanaNotificationPolicyUpdate,
		//DeleteContext: resourceGrafanaNotificationPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			grafanaNotificationPolicyContactPoint: {
				Type:     schema.TypeString,
				Required: true,
			},
			grafanaNotificationPolicyGroupBy: {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
			grafanaNotificationPolicyGroupInterval: {
				Type:     schema.TypeString,
				Optional: true,
			},
			grafanaNotificationPolicyGroupWait: {
				Type:     schema.TypeString,
				Optional: true,
			},
			grafanaNotificationPolicyRepeatInterval: {
				Type:     schema.TypeString,
				Optional: true,
			},
			grafanaNotificationPolicyPolicy: {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     buildPolicySchema(grafanaNotificationPolicyTreeDepth),
			},
		},
	}
}

// buildPolicySchema builds the policy tree's nodes schema.
// Since TF does not support infinite recursive schema, we restrict the tree's depth with const grafanaNotificationPolicyTreeDepth
func buildPolicySchema(treeDepth uint) *schema.Resource {
	if treeDepth == 0 {
		// we should never get here
		panic("cannot create policy schema with depth 0")
	}

	policy := &schema.Resource{
		Schema: map[string]*schema.Schema{
			grafanaNotificationPolicyContactPoint: {
				Type:     schema.TypeString,
				Required: true,
			},
			grafanaNotificationPolicyGroupBy: {
				Type:     schema.TypeList,
				Required: treeDepth == 1,
				Optional: treeDepth > 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			grafanaNotificationPolicyMatcher: {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						grafanaNotificationPolicyMatcherLabel: {
							Type:     schema.TypeString,
							Required: true,
						},
						grafanaNotificationPolicyMatcherMatch: {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: utils.ValidateGrafanaNotificationPolicyMatcherMatch,
						},
						grafanaNotificationPolicyMatcherValue: {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			grafanaNotificationPolicyMuteTiming: {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			grafanaNotificationPolicyContinue: {
				Type:     schema.TypeBool,
				Optional: true,
			},
			grafanaNotificationPolicyGroupWait: {
				Type:     schema.TypeString,
				Optional: true,
			},
			grafanaNotificationPolicyGroupInterval: {
				Type:     schema.TypeString,
				Optional: true,
			},
			grafanaNotificationPolicyRepeatInterval: {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}

	if treeDepth > 1 {
		policy.Schema[grafanaNotificationPolicyPolicy] = &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem:     buildPolicySchema(treeDepth - 1),
		}
	}

	return policy
}

func grafanaNotificationPolicyClient(m interface{}) *grafana_notification_policies.GrafanaNotificationPolicyClient {
	var client *grafana_notification_policies.GrafanaNotificationPolicyClient
	client, _ = grafana_notification_policies.New(m.(Config).apiToken, m.(Config).baseUrl)
	return client
}

func resourceGrafanaNotificationPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	grafanaNotificationPolicy, err := createGrafanaNotificationPolicyFromSchema(d)
}

func createGrafanaNotificationPolicyFromSchema(d *schema.ResourceData) (grafana_notification_policies.GrafanaNotificationPolicyTree, error) {
	var grafanaNotificationPolicyTree = grafana_notification_policies.GrafanaNotificationPolicyTree{
		GroupInterval:  d.Get(grafanaNotificationPolicyGroupInterval).(string),
		GroupWait:      d.Get(grafanaNotificationPolicyGroupWait).(string),
		Receiver:       d.Get(grafanaNotificationPolicyContactPoint).(string),
		RepeatInterval: d.Get(grafanaNotificationPolicyRepeatInterval).(string),
	}

	groupByInterface := d.Get(grafanaNotificationPolicyGroupBy).([]interface{})
	groupBy := make([]string, 0, len(groupByInterface))
	for _, group := range groupByInterface {
		groupBy = append(groupBy, group.(string))
	}

	var policies []grafana_notification_policies.GrafanaNotificationPolicy
	if policiesFromSchema, ok := d.GetOk(grafanaNotificationPolicyPolicy); ok {
		routes := policiesFromSchema.([]interface{})
		for _, route := range routes {
			policy, err := getPolicyFromSchema(route)
			if err != nil {
				return grafana_notification_policies.GrafanaNotificationPolicyTree{}, err
			}
			policies = append(policies, policy)
		}
	}

	return grafanaNotificationPolicyTree, nil
}

func getPolicyFromSchema(policyFromSchema interface{}) (grafana_notification_policies.GrafanaNotificationPolicy, error) {
	var policy grafana_notification_policies.GrafanaNotificationPolicy
	policyMap := policyFromSchema.(map[string]interface{})

	if v, ok := policyMap[grafanaNotificationPolicyContinue]; ok && v != nil {
		policy.Continue = v.(bool)
	}

	if v, ok := policyMap[grafanaNotificationPolicyGroupBy]; ok {
		policy.GroupBy = utils.ParseInterfaceSliceToStringSlice(v.([]interface{}))
	}

	if v, ok := policyMap[grafanaNotificationPolicyGroupInterval]; ok && v != nil {
		policy.GroupInterval = v.(string)
	}

	if v, ok := policyMap[grafanaNotificationPolicyGroupWait]; ok && v != nil {
		policy.GroupWait = v.(string)
	}

	if v, ok := policyMap[grafanaNotificationPolicyMuteTiming]; ok && v != nil {
		policy.MuteTimeIntervals = utils.ParseInterfaceSliceToStringSlice(v.([]interface{}))
	}

	if v, ok := policyMap[grafanaNotificationPolicyMatcher]; ok && v != nil {
		ms := v.([]interface{})
		matchers := make([]grafana_notification_policies.MatcherObj, 0, len(ms))
		for _, m := range ms {
			matcher, err := getMatcherFromSchema(m)
			if err != nil {
				return grafana_notification_policies.GrafanaNotificationPolicy{}, err
			}
			matchers = append(matchers, matcher)
		}

		policy.ObjectMatchers = matchers
	}

	policy.Receiver = policyMap[grafanaNotificationPolicyContactPoint].(string)

	if v, ok := policyMap[grafanaNotificationPolicyRepeatInterval]; ok && v != nil {
		policy.RepeatInterval = v.(string)
	}

	if v, ok := policyMap[grafanaNotificationPolicyPolicy]; ok && v != nil {
		ps := v.([]interface{})
		policies := make([]grafana_notification_policies.GrafanaNotificationPolicy, 0, len(ps))
		for _, p := range ps {
			unpacked, err := getPolicyFromSchema(p)
			if err != nil {
				return grafana_notification_policies.GrafanaNotificationPolicy{}, err
			}
			policies = append(policies, unpacked)
		}
		policy.Routes = policies
	}

	return policy, nil
}

func getMatcherFromSchema(m interface{}) (grafana_notification_policies.MatcherObj, error) {
	matcherMap := m.(map[string]interface{})

	matcher := []string{matcherMap[grafanaNotificationPolicyMatcherLabel].(string),
		matcherMap[grafanaNotificationPolicyMatcherMatch].(string),
		matcherMap[grafanaNotificationPolicyMatcherValue].(string)}

	return matcher, nil
}
