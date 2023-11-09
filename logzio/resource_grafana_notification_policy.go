package logzio

import (
	"context"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/logzio/logzio_terraform_client/grafana_notification_policies"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
	"strings"
	"time"
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
	grafanaNotificationPolicyMuteTimings  = "mute_timings"
	grafanaNotificationPolicyContinue     = "continue"

	grafanaNotificationPolicyTreeDepth = 4

	// Since one resource manages the entire tree, and does not create an id, we'll use this id for Terraform
	grafanaNotificationPolicyStaticId = "logzio_policy"

	grafanaNotificationPolicyUpdateDelaySeconds = 4
)

// resourceGrafanaNotificationPolicy represents a Grafana notification policy tree. Note that one resource represents the entire policy tree
func resourceGrafanaNotificationPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGrafanaNotificationPolicyCreate,
		ReadContext:   resourceGrafanaNotificationPolicyRead,
		UpdateContext: resourceGrafanaNotificationPolicyUpdate,
		DeleteContext: resourceGrafanaNotificationPolicyDelete,
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
			grafanaNotificationPolicyMuteTimings: {
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
	if err != nil {
		return diag.FromErr(err)
	}

	err = grafanaNotificationPolicyClient(m).SetupGrafanaNotificationPolicyTree(grafanaNotificationPolicy)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(grafanaNotificationPolicyStaticId)
	return resourceGrafanaNotificationPolicyRead(ctx, d, m)
}

func resourceGrafanaNotificationPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	grafanaNotificationPolicy, err := grafanaNotificationPolicyClient(m).GetGrafanaNotificationPolicyTree()
	if err != nil {
		tflog.Error(ctx, err.Error())
		if strings.Contains(err.Error(), "missing grafana notification policy") {
			// If we were not able to find the resource - delete from state
			d.SetId("")
			return diag.Diagnostics{}
		} else {
			return diag.FromErr(err)
		}
	}

	setGrafanaNotificationPolicy(d, grafanaNotificationPolicy)
	return nil
}

func resourceGrafanaNotificationPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	grafanaNotificationPolicy, err := createGrafanaNotificationPolicyFromSchema(d)
	if err != nil {
		return diag.FromErr(err)
	}

	err = grafanaNotificationPolicyClient(m).SetupGrafanaNotificationPolicyTree(grafanaNotificationPolicy)
	if err != nil {
		return diag.FromErr(err)
	}

	time.Sleep(grafanaNotificationPolicyUpdateDelaySeconds * time.Second)

	return resourceGrafanaNotificationPolicyRead(ctx, d, m)
}

// resourceGrafanaNotificationPolicyDelete only RESETS the notification policy tree (because of the way the Grafana Notification Policy API works).
// Using this endpoint will reset the entire notification policy tree.
func resourceGrafanaNotificationPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := grafanaNotificationPolicyClient(m).ResetGrafanaNotificationPolicyTree()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func setGrafanaNotificationPolicy(d *schema.ResourceData, grafanaNotificationPolicy grafana_notification_policies.GrafanaNotificationPolicyTree) {
	d.Set(grafanaNotificationPolicyContactPoint, grafanaNotificationPolicy.Receiver)
	d.Set(grafanaNotificationPolicyGroupBy, grafanaNotificationPolicy.GroupBy)
	d.Set(grafanaNotificationPolicyGroupWait, grafanaNotificationPolicy.GroupWait)
	d.Set(grafanaNotificationPolicyGroupInterval, grafanaNotificationPolicy.GroupInterval)
	d.Set(grafanaNotificationPolicyRepeatInterval, grafanaNotificationPolicy.RepeatInterval)

	if len(grafanaNotificationPolicy.Routes) > 0 {
		policies := make([]interface{}, 0, len(grafanaNotificationPolicy.Routes))
		for _, route := range grafanaNotificationPolicy.Routes {
			policies = append(policies, getPolicyFromObject(route, grafanaNotificationPolicyTreeDepth))
		}

		d.Set(grafanaNotificationPolicyPolicy, policies)
	}
}

func getPolicyFromObject(policy grafana_notification_policies.GrafanaNotificationPolicy, treeDepth uint) interface{} {
	policyMap := map[string]interface{}{}

	policyMap[grafanaNotificationPolicyContinue] = policy.Continue

	if len(policy.GroupBy) > 0 {
		policyMap[grafanaNotificationPolicyGroupBy] = policy.GroupBy
	}

	if len(policy.GroupInterval) > 0 {
		policyMap[grafanaNotificationPolicyGroupInterval] = policy.GroupInterval
	}

	if len(policy.GroupWait) > 0 {
		policyMap[grafanaNotificationPolicyGroupWait] = policy.GroupWait
	}

	if policy.MuteTimeIntervals != nil && len(policy.MuteTimeIntervals) > 0 {
		policyMap[grafanaNotificationPolicyMuteTimings] = policy.MuteTimeIntervals
	}

	if policy.ObjectMatchers != nil && len(policy.ObjectMatchers) > 0 {
		matchers := make([]interface{}, 0, len(policy.ObjectMatchers))
		for _, matcher := range policy.ObjectMatchers {
			matchers = append(matchers, getMatcherFromObject(matcher))
		}
		policyMap[grafanaNotificationPolicyMatcher] = matchers
	}

	policyMap[grafanaNotificationPolicyContactPoint] = policy.Receiver

	if len(policy.RepeatInterval) > 0 {
		policyMap[grafanaNotificationPolicyRepeatInterval] = policy.RepeatInterval
	}

	if treeDepth > 1 && policy.Routes != nil && len(policy.Routes) > 0 {
		policies := make([]interface{}, 0, len(policy.Routes))
		for _, route := range policy.Routes {
			policies = append(policies, getPolicyFromObject(route, treeDepth-1))
		}
		policyMap[grafanaNotificationPolicyPolicy] = policies
	}

	return policyMap
}

func getMatcherFromObject(matcherObject grafana_notification_policies.MatcherObj) interface{} {
	const (
		labelIndex = iota
		matchIndex
		valueIndex
	)

	return map[string]interface{}{
		grafanaNotificationPolicyMatcherLabel: matcherObject[labelIndex],
		grafanaNotificationPolicyMatcherMatch: matcherObject[matchIndex],
		grafanaNotificationPolicyMatcherValue: matcherObject[valueIndex],
	}
}

func createGrafanaNotificationPolicyFromSchema(d *schema.ResourceData) (grafana_notification_policies.GrafanaNotificationPolicyTree, error) {
	var grafanaNotificationPolicyTree = grafana_notification_policies.GrafanaNotificationPolicyTree{
		GroupInterval:  d.Get(grafanaNotificationPolicyGroupInterval).(string),
		GroupWait:      d.Get(grafanaNotificationPolicyGroupWait).(string),
		Receiver:       d.Get(grafanaNotificationPolicyContactPoint).(string),
		RepeatInterval: d.Get(grafanaNotificationPolicyRepeatInterval).(string),
	}

	groupByInterface := d.Get(grafanaNotificationPolicyGroupBy).([]interface{})
	for _, group := range groupByInterface {
		grafanaNotificationPolicyTree.GroupBy = append(grafanaNotificationPolicyTree.GroupBy, group.(string))
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

	grafanaNotificationPolicyTree.Routes = policies

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

	if v, ok := policyMap[grafanaNotificationPolicyMuteTimings]; ok && v != nil {
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
