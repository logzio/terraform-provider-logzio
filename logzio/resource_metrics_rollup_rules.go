package logzio

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/avast/retry-go"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/logzio/logzio_terraform_client/metrics_rollup_rules"
)

const (
	metricsRollupRulesId                      string = "id"
	metricsRollupRulesAccountId               string = "account_id"
	metricsRollupRulesMetricName              string = "metric_name"
	metricsRollupRulesMetricType              string = "metric_type"
	metricsRollupRulesRollupFunction          string = "rollup_function"
	metricsRollupRulesLabelsEliminationMethod string = "labels_elimination_method"
	metricsRollupRulesLabels                  string = "labels"

	metricsRollupRulesRetryAttempts = 8
)

// Returns the metrics rollup rules client with the api token from the provider
func metricsRollupRulesClient(m interface{}) *metrics_rollup_rules.MetricsRollupRulesClient {
	var client *metrics_rollup_rules.MetricsRollupRulesClient
	client, _ = metrics_rollup_rules.New(m.(Config).apiToken, m.(Config).baseUrl)
	return client
}

func resourceMetricsRollupRules() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMetricsRollupRulesCreate,
		ReadContext:   resourceMetricsRollupRulesRead,
		UpdateContext: resourceMetricsRollupRulesUpdate,
		DeleteContext: resourceMetricsRollupRulesDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			metricsRollupRulesId: {
				Type:     schema.TypeString,
				Computed: true,
			},
			metricsRollupRulesAccountId: {
				Type:     schema.TypeInt,
				Required: true,
			},
			metricsRollupRulesMetricName: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 255),
			},
			metricsRollupRulesMetricType: {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice(
					[]string{
						string(metrics_rollup_rules.MetricTypeGauge),
						string(metrics_rollup_rules.MetricTypeCounter),
					}, false),
			},
			metricsRollupRulesRollupFunction: {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice(
					[]string{
						string(metrics_rollup_rules.AggSum),
						string(metrics_rollup_rules.AggMin),
						string(metrics_rollup_rules.AggMax),
						string(metrics_rollup_rules.AggCount),
						string(metrics_rollup_rules.AggLast),
					}, false),
			},
			metricsRollupRulesLabelsEliminationMethod: {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice(
					[]string{
						string(metrics_rollup_rules.LabelsExcludeBy),
					}, false),
			},
			metricsRollupRulesLabels: {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

// resourceMetricsRollupRulesCreate creates a new metrics rollup rule in logzio
func resourceMetricsRollupRulesCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	createRollupRule := createCreateUpdateMetricsRollupRuleFromSchema(d)

	rollupRule, err := metricsRollupRulesClient(m).CreateRollupRule(createRollupRule)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(rollupRule.Id)
	return resourceMetricsRollupRulesRead(ctx, d, m)
}

// resourceMetricsRollupRulesRead gets metrics rollup rule by id
func resourceMetricsRollupRulesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	rollupRuleId := d.Id()
	rollupRule, err := metricsRollupRulesClient(m).GetRollupRule(rollupRuleId)
	if err != nil {
		tflog.Error(ctx, err.Error())
		if strings.Contains(err.Error(), "missing") || strings.Contains(err.Error(), "not found") {
			// If we were not able to find the resource - delete from state
			d.SetId("")
			return diag.Diagnostics{}
		} else {
			return diag.FromErr(err)
		}
	}

	setMetricsRollupRule(d, rollupRule)
	return nil
}

// resourceMetricsRollupRulesUpdate updates an existing metrics rollup rule in logzio
func resourceMetricsRollupRulesUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	rollupRuleId := d.Id()
	updateRule := createCreateUpdateMetricsRollupRuleFromSchema(d)

	client := metricsRollupRulesClient(m)
	_, err := client.UpdateRollupRule(rollupRuleId, updateRule)
	if err != nil {
		return diag.FromErr(err)
	}

	diags := readMetricsRollupRuleUntilConsistent(ctx, d, m, metricsRollupRulesRetryAttempts, "update rollup rule", func() bool {
		createRule := createCreateUpdateMetricsRollupRuleFromSchema(d)
		return reflect.DeepEqual(createRule, updateRule)
	})
	return diags
}

// resourceMetricsRollupRulesDelete deletes a metrics rollup rule in logzio
func resourceMetricsRollupRulesDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	rollupRuleId := d.Id()

	err := metricsRollupRulesClient(m).DeleteRollupRule(rollupRuleId)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

// createCreateUpdateMetricsRollupRuleFromSchema creates a CreateUpdateRollupRule object from the schema
func createCreateUpdateMetricsRollupRuleFromSchema(d *schema.ResourceData) metrics_rollup_rules.CreateUpdateRollupRule {
	accountId := int64(d.Get(metricsRollupRulesAccountId).(int))
	metricName := d.Get(metricsRollupRulesMetricName).(string)
	metricType := metrics_rollup_rules.MetricType(d.Get(metricsRollupRulesMetricType).(string))
	rollupFunction := d.Get(metricsRollupRulesRollupFunction).(string)
	labelsEliminationMethod := d.Get(metricsRollupRulesLabelsEliminationMethod).(string)

	labels := []string{}
	if labelsInterface, ok := d.GetOk(metricsRollupRulesLabels); ok {
		labelsList := labelsInterface.([]interface{})
		for _, label := range labelsList {
			labels = append(labels, label.(string))
		}
	}

	return metrics_rollup_rules.CreateUpdateRollupRule{
		AccountId:               accountId,
		MetricName:              metricName,
		MetricType:              metricType,
		RollupFunction:          metrics_rollup_rules.AggregationFunction(rollupFunction),
		LabelsEliminationMethod: metrics_rollup_rules.LabelsRemovalMethod(labelsEliminationMethod),
		Labels:                  labels,
	}
}

// setMetricsRollupRule sets the resource data from a RollupRule object
func setMetricsRollupRule(d *schema.ResourceData, rollupRule *metrics_rollup_rules.RollupRule) {
	d.Set(metricsRollupRulesAccountId, rollupRule.AccountId)
	d.Set(metricsRollupRulesMetricName, rollupRule.MetricName)
	d.Set(metricsRollupRulesMetricType, string(rollupRule.MetricType))
	d.Set(metricsRollupRulesRollupFunction, rollupRule.RollupFunction)
	d.Set(metricsRollupRulesLabelsEliminationMethod, rollupRule.LabelsEliminationMethod)
	d.Set(metricsRollupRulesLabels, rollupRule.Labels)
}

// readMetricsRollupRuleUntilConsistent retries reading until the resource state is consistent
func readMetricsRollupRuleUntilConsistent(ctx context.Context, d *schema.ResourceData, m interface{}, retryAttempts int, operation string, isConsistent func() bool) diag.Diagnostics {
	err := retry.Do(
		func() error {
			err := resourceMetricsRollupRulesRead(ctx, d, m)
			if err != nil && len(err) > 0 {
				return fmt.Errorf("failed to read after %s: %v", operation, err)
			}
			if !isConsistent() {
				return fmt.Errorf("resource state not consistent after %s", operation)
			}
			return nil
		},
		retry.RetryIf(func(err error) bool {
			return err != nil
		}),
		retry.Attempts(uint(retryAttempts)),
		retry.DelayType(retry.BackOffDelay),
	)

	if err != nil {
		tflog.Warn(ctx, fmt.Sprintf("Failed to achieve consistency after %s: %v", operation, err))
	}

	return nil
}
