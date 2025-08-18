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
