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
	metricsRollupRulesName                    string = "name"
	metricsRollupRulesMetricName              string = "metric_name"
	metricsRollupRulesMetricType              string = "metric_type"
	metricsRollupRulesRollupFunction          string = "rollup_function"
	metricsRollupRulesLabelsEliminationMethod string = "labels_elimination_method"
	metricsRollupRulesLabels                  string = "labels"
	metricsRollupRulesNewMetricNameTemplate   string = "new_metric_name_template"
	metricsRollupRulesDropOriginalMetric      string = "drop_original_metric"
	metricsRollupRulesFilter                  string = "filter"
	metricsRollupRulesNamespaces              string = "namespaces"
	metricsRollupRulesClusterId               string = "cluster_id"
	metricsRollupRulesIsDeleted               string = "is_deleted"
	metricsRollupRulesDropPolicyRuleId        string = "drop_policy_rule_id"
	metricsRollupRulesVersion                 string = "version"

	metricsRollupRulesFilterExpression string = "expression"
	metricsRollupRulesFilterComparison string = "comparison"
	metricsRollupRulesFilterName       string = "name"
	metricsRollupRulesFilterValue      string = "value"

	comparisonEq           = "EQ"
	comparisonNotEq        = "NOT_EQ"
	comparisonRegexMatch   = "REGEX_MATCH"
	comparisonRegexNoMatch = "REGEX_NO_MATCH"

	errorIdMustBeSpecified     = "id must be specified for data source"
	errorNoMatchingRollupRules = "couldn't find metrics rollup rule with specified attributes"
	errorMultipleMatchingRules = "found multiple (%d) metrics rollup rules matching the criteria, please specify an id or add more search criteria"
	errorRollupRuleNotFound    = "could not find metrics rollup rule with id %s"

	metricsRollupRulesRetryAttempts = 8
)

// Helper functions to convert client enums to string slices for validation
func convertMetricTypesToStrings(types []metrics_rollup_rules.MetricType) []string {
	result := make([]string, len(types))
	for i, t := range types {
		result[i] = string(t)
	}
	return result
}

func convertAggregationFunctionsToStrings(funcs []metrics_rollup_rules.AggregationFunction) []string {
	result := make([]string, len(funcs))
	for i, f := range funcs {
		result[i] = string(f)
	}
	return result
}

func convertLabelsRemovalMethodsToStrings(methods []metrics_rollup_rules.LabelsRemovalMethod) []string {
	result := make([]string, len(methods))
	for i, m := range methods {
		result[i] = string(m)
	}
	return result
}

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
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
			metricType := d.Get(metricsRollupRulesMetricType).(string)
			rollupFunction := d.Get(metricsRollupRulesRollupFunction).(string)

			switch metricType {
			case string(metrics_rollup_rules.MetricTypeGauge):
				if rollupFunction == "" {
					return fmt.Errorf("rollup_function must be set for GAUGE metrics")
				}
			case string(metrics_rollup_rules.MetricTypeMeasurement):
				if rollupFunction == "" {
					return fmt.Errorf("rollup_function must be set for MEASUREMENT metrics")
				}
				// Validate MEASUREMENT allows only specific aggregation functions
				allowedForMeasurement := map[string]bool{
					string(metrics_rollup_rules.AggSum):   true,
					string(metrics_rollup_rules.AggMin):   true,
					string(metrics_rollup_rules.AggMax):   true,
					string(metrics_rollup_rules.AggCount): true,
					string(metrics_rollup_rules.AggSumSq): true,
					string(metrics_rollup_rules.AggMean):  true,
					string(metrics_rollup_rules.AggLast):  true,
				}
				if !allowedForMeasurement[rollupFunction] {
					return fmt.Errorf("invalid aggregation function %q for MEASUREMENT metric type. Allowed functions: SUM, MIN, MAX, COUNT, SUMSQ, MEAN, LAST", rollupFunction)
				}
			case string(metrics_rollup_rules.MetricTypeCounter),
				string(metrics_rollup_rules.MetricTypeDeltaCounter),
				string(metrics_rollup_rules.MetricTypeCumulativeCounter):
				if rollupFunction == "" {
					return fmt.Errorf("rollup_function must be set for %s metrics", metricType)
				}
				if rollupFunction != string(metrics_rollup_rules.AggSum) {
					return fmt.Errorf("for %s metrics, rollup_function must be SUM", metricType)
				}
			}
			return nil
		},
		Schema: map[string]*schema.Schema{
			metricsRollupRulesId: {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			metricsRollupRulesAccountId: {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			metricsRollupRulesName: {
				Type:     schema.TypeString,
				Optional: true,
			},
			metricsRollupRulesMetricName: {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{metricsRollupRulesMetricName, metricsRollupRulesFilter},
			},
			metricsRollupRulesMetricType: {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice(
					convertMetricTypesToStrings(metrics_rollup_rules.GetValidMetricType()), false),
			},
			metricsRollupRulesRollupFunction: {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice(
					convertAggregationFunctionsToStrings(metrics_rollup_rules.GetValidAggregationFunctions()), false),
			},
			metricsRollupRulesLabelsEliminationMethod: {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice(
					convertLabelsRemovalMethodsToStrings(metrics_rollup_rules.GetValidLabelsRemovalMethods()), false),
			},
			metricsRollupRulesLabels: {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			metricsRollupRulesNewMetricNameTemplate: {
				Type:     schema.TypeString,
				Optional: true,
			},
			metricsRollupRulesDropOriginalMetric: {
				Type:     schema.TypeBool,
				Optional: true,
			},
			metricsRollupRulesFilter: {
				Type:         schema.TypeList,
				Optional:     true,
				MaxItems:     1,
				ExactlyOneOf: []string{metricsRollupRulesMetricName, metricsRollupRulesFilter},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						metricsRollupRulesFilterExpression: {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									metricsRollupRulesFilterComparison: {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice(
											[]string{comparisonEq, comparisonNotEq, comparisonRegexMatch, comparisonRegexNoMatch}, false),
									},
									metricsRollupRulesFilterName: {
										Type:     schema.TypeString,
										Required: true,
									},
									metricsRollupRulesFilterValue: {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			// Computed-only fields
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
	metricType := metrics_rollup_rules.MetricType(d.Get(metricsRollupRulesMetricType).(string))
	labelsEliminationMethod := d.Get(metricsRollupRulesLabelsEliminationMethod).(string)

	var rollupFunction metrics_rollup_rules.AggregationFunction
	if v, ok := d.GetOk(metricsRollupRulesRollupFunction); ok {
		rollupFunction = metrics_rollup_rules.AggregationFunction(v.(string))
	}

	var name string
	if v, ok := d.GetOk(metricsRollupRulesName); ok {
		name = v.(string)
	}

	var metricName string
	if v, ok := d.GetOk(metricsRollupRulesMetricName); ok {
		metricName = v.(string)
	}

	var newMetricNameTemplate *string
	if v, ok := d.GetOk(metricsRollupRulesNewMetricNameTemplate); ok {
		s := v.(string)
		newMetricNameTemplate = &s
	}

	var dropOriginalMetric *bool
	if v, ok := d.GetOkExists(metricsRollupRulesDropOriginalMetric); ok {
		b := v.(bool)
		dropOriginalMetric = &b
	}

	labels := []string{}
	if labelsInterface, ok := d.GetOk(metricsRollupRulesLabels); ok {
		labelsList := labelsInterface.([]interface{})
		for _, label := range labelsList {
			labels = append(labels, label.(string))
		}
	}

	var filter *metrics_rollup_rules.ComplexFilter
	if f, ok := d.GetOk(metricsRollupRulesFilter); ok {
		filterList := f.([]interface{})
		if len(filterList) > 0 && filterList[0] != nil {
			filterMap := filterList[0].(map[string]interface{})
			expressionInterface := filterMap[metricsRollupRulesFilterExpression].([]interface{})
			var expressions []metrics_rollup_rules.SingleFilter
			for _, e := range expressionInterface {
				expressionMap := e.(map[string]interface{})
				expressions = append(expressions, metrics_rollup_rules.SingleFilter{
					Comparison: metrics_rollup_rules.Comparison(expressionMap[metricsRollupRulesFilterComparison].(string)),
					Name:       expressionMap[metricsRollupRulesFilterName].(string),
					Value:      expressionMap[metricsRollupRulesFilterValue].(string),
				})
			}
			filter = &metrics_rollup_rules.ComplexFilter{Expression: expressions}
		}
	}

	return metrics_rollup_rules.CreateUpdateRollupRule{
		AccountId:               accountId,
		Name:                    name,
		MetricName:              metricName,
		MetricType:              metricType,
		RollupFunction:          rollupFunction,
		LabelsEliminationMethod: metrics_rollup_rules.LabelsRemovalMethod(labelsEliminationMethod),
		Labels:                  labels,
		Filter:                  filter,
		NewMetricNameTemplate:   newMetricNameTemplate,
		DropOriginalMetric:      dropOriginalMetric,
	}
}

// setMetricsRollupRule sets the resource data from a RollupRule object
func setMetricsRollupRule(d *schema.ResourceData, rollupRule *metrics_rollup_rules.RollupRule) {
	d.Set(metricsRollupRulesAccountId, rollupRule.AccountId)
	d.Set(metricsRollupRulesName, rollupRule.Name)
	d.Set(metricsRollupRulesMetricName, rollupRule.MetricName)
	d.Set(metricsRollupRulesMetricType, string(rollupRule.MetricType))
	d.Set(metricsRollupRulesRollupFunction, rollupRule.RollupFunction)
	d.Set(metricsRollupRulesLabelsEliminationMethod, rollupRule.LabelsEliminationMethod)
	d.Set(metricsRollupRulesLabels, rollupRule.Labels)
	d.Set(metricsRollupRulesNamespaces, rollupRule.Namespaces)
	d.Set(metricsRollupRulesClusterId, rollupRule.ClusterId)
	d.Set(metricsRollupRulesIsDeleted, rollupRule.IsDeleted)
	d.Set(metricsRollupRulesVersion, rollupRule.Version)
	d.Set(metricsRollupRulesDropOriginalMetric, rollupRule.DropOriginalMetric)

	if rollupRule.DropPolicyRuleId != nil {
		d.Set(metricsRollupRulesDropPolicyRuleId, *rollupRule.DropPolicyRuleId)
	} else {
		d.Set(metricsRollupRulesDropPolicyRuleId, "")
	}

	if rollupRule.NewMetricNameTemplate != nil {
		d.Set(metricsRollupRulesNewMetricNameTemplate, *rollupRule.NewMetricNameTemplate)
	} else {
		d.Set(metricsRollupRulesNewMetricNameTemplate, "")
	}

	// Set filter if present
	if rollupRule.Filter != nil {
		filterList := []interface{}{
			map[string]interface{}{
				metricsRollupRulesFilterExpression: func() []interface{} {
					expressions := make([]interface{}, len(rollupRule.Filter.Expression))
					for i, expr := range rollupRule.Filter.Expression {
						expressions[i] = map[string]interface{}{
							metricsRollupRulesFilterComparison: string(expr.Comparison),
							metricsRollupRulesFilterName:       expr.Name,
							metricsRollupRulesFilterValue:      expr.Value,
						}
					}
					return expressions
				}(),
			},
		}
		d.Set(metricsRollupRulesFilter, filterList)
	} else {
		d.Set(metricsRollupRulesFilter, nil)
	}
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
