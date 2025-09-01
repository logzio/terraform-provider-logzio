package logzio

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMetricsRollupRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMetricsRollupRulesRead,
		Schema: map[string]*schema.Schema{
			metricsRollupRulesId: {
				Type:     schema.TypeString,
				Required: true,
			},
			metricsRollupRulesAccountId: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			metricsRollupRulesName: {
				Type:     schema.TypeString,
				Computed: true,
			},
			metricsRollupRulesMetricName: {
				Type:     schema.TypeString,
				Computed: true,
			},
			metricsRollupRulesMetricType: {
				Type:     schema.TypeString,
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
						"expression": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"comparison": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
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

	if !hasId {
		return diag.Errorf("id must be specified for data source")
	}

	// Get by ID
	rollupRule, err := client.GetRollupRule(rollupRuleId.(string))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(rollupRule.Id)
	setMetricsRollupRule(d, rollupRule)
	return nil
}
