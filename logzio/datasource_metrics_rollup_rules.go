package logzio

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/logzio/logzio_terraform_client/metrics_rollup_rules"
)

func dataSourceMetricsRollupRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMetricsRollupRulesRead,
		Schema: map[string]*schema.Schema{
			metricsRollupRulesId: {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			metricsRollupRulesAccountId: {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			metricsRollupRulesName: {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			metricsRollupRulesMetricName: {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			metricsRollupRulesMetricType: {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			metricsRollupRulesRollupFunction: {
				Type:     schema.TypeString,
				Computed: true,
			},
			metricsRollupRulesLabelsEliminationMethod: {
				Type:     schema.TypeString,
				Computed: true,
			},
			metricsRollupRulesLabels: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			metricsRollupRulesNewMetricNameTemplate: {
				Type:     schema.TypeString,
				Computed: true,
			},
			metricsRollupRulesDropOriginalMetric: {
				Type:     schema.TypeBool,
				Computed: true,
			},
			metricsRollupRulesFilter: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						metricsRollupRulesFilterExpression: {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									metricsRollupRulesFilterComparison: {
										Type:     schema.TypeString,
										Computed: true,
									},
									metricsRollupRulesFilterName: {
										Type:     schema.TypeString,
										Computed: true,
									},
									metricsRollupRulesFilterValue: {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			metricsRollupRulesNamespaces: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			metricsRollupRulesClusterId: {
				Type:     schema.TypeString,
				Computed: true,
			},
			metricsRollupRulesIsDeleted: {
				Type:     schema.TypeBool,
				Computed: true,
			},
			metricsRollupRulesDropPolicyRuleId: {
				Type:     schema.TypeString,
				Computed: true,
			},
			metricsRollupRulesVersion: {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceMetricsRollupRulesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := metricsRollupRulesClient(m)

	rollupRuleId, hasId := d.GetOk(metricsRollupRulesId)

	if hasId {
		rollupRule, err := client.GetRollupRule(rollupRuleId.(string))
		if err != nil {
			return diag.Errorf(errorRollupRuleNotFound, rollupRuleId.(string))
		}

		d.SetId(rollupRule.Id)
		setMetricsRollupRule(d, rollupRule)
		return nil
	}

	rollupRule, err := searchMatchingRollupRule(client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(rollupRule.Id)
	setMetricsRollupRule(d, rollupRule)
	return nil
}

// searchMatchingRollupRule searches for a metrics rollup rule that matches the specified attributes in the schema
func searchMatchingRollupRule(client *metrics_rollup_rules.MetricsRollupRulesClient, d *schema.ResourceData) (*metrics_rollup_rules.RollupRule, error) {
	// Build search request from schema attributes
	searchReq := buildSearchRequestFromSchema(d)

	results, err := client.SearchRollupRules(searchReq)
	if err != nil {
		return nil, err
	}

	// Handle search results
	switch len(results) {
	case 0:
		return nil, fmt.Errorf(errorNoMatchingRollupRules)
	case 1:
		return &results[0], nil
	default:
		exactMatch := findExactMatch(results, d)
		if exactMatch != nil {
			return exactMatch, nil
		}
		return nil, fmt.Errorf(errorMultipleMatchingRules, len(results))
	}
}

// buildSearchRequestFromSchema creates a search request from the datasource schema attributes
func buildSearchRequestFromSchema(d *schema.ResourceData) metrics_rollup_rules.SearchRollupRulesRequest {
	req := metrics_rollup_rules.SearchRollupRulesRequest{
		Filter: &metrics_rollup_rules.SearchFilter{},
	}

	if accountId, ok := d.GetOk(metricsRollupRulesAccountId); ok {
		req.Filter.AccountIds = []int64{int64(accountId.(int))}
	}

	if metricName, ok := d.GetOk(metricsRollupRulesMetricName); ok {
		req.Filter.MetricNames = []string{metricName.(string)}
	}

	if name, ok := d.GetOk(metricsRollupRulesName); ok {
		req.Filter.SearchTerm = name.(string)
	}

	return req
}

// findExactMatch finds an exact match among multiple results based on additional criteria
func findExactMatch(results []metrics_rollup_rules.RollupRule, d *schema.ResourceData) *metrics_rollup_rules.RollupRule {
	for i := range results {
		rule := &results[i]

		if metricType, ok := d.GetOk(metricsRollupRulesMetricType); ok {
			if string(rule.MetricType) != metricType.(string) {
				continue
			}
		}

		if name, ok := d.GetOk(metricsRollupRulesName); ok {
			if rule.Name != name.(string) {
				continue
			}
		}

		if metricName, ok := d.GetOk(metricsRollupRulesMetricName); ok {
			if rule.MetricName != metricName.(string) {
				continue
			}
		}

		return rule
	}

	return nil
}
