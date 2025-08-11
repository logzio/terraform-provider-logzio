package logzio

import (
	"context"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/logzio/logzio_terraform_client/drop_metrics"
	"reflect"
	"strconv"
)

const (
	dropMetricsIdField              = "drop_metric_id"
	dropMetricsAccountId            = "account_id"
	dropMetricsActive               = "active"
	dropMetricsFilters              = "filters"
	dropMetricsFilterOperator       = "operator"
	dropMetricsExpressionLabelName  = "name"
	dropMetricsExpressionLabelValue = "value"
	dropMetricsExpressionCondition  = "condition"
	dropMetricsCreatedAt            = "created_at"
	dropMetricsCreatedBy            = "created_by"
	dropMetricsModifiedAt           = "modified_at"
	dropMetricsModifiedBy           = "modified_by"

	dropMetricsRetryAttempts = 8
)

// Returns the drop metrics client with the api token from the provider
func dropMetricsClient(m interface{}) *drop_metrics.DropMetricsClient {
	var client *drop_metrics.DropMetricsClient
	client, _ = drop_metrics.New(m.(Config).apiToken, m.(Config).baseUrl)
	return client
}

func resourceDropMetrics() *schema.Resource {
	var filterExprSchema = map[string]*schema.Schema{
		dropMetricsExpressionLabelName: {
			Type:     schema.TypeString,
			Required: true,
		},
		dropMetricsExpressionLabelValue: {
			Type:     schema.TypeString,
			Required: true,
		},
		dropMetricsExpressionCondition: {
			Type:     schema.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice(
				[]string{
					drop_metrics.ComparisonEq,
					drop_metrics.ComparisonNotEq,
					drop_metrics.ComparisonRegexMatch,
					drop_metrics.ComparisonRegexNoMatch,
				}, false),
		},
	}

	return &schema.Resource{
		CreateContext: resourceDropMetricsCreate,
		ReadContext:   resourceDropMetricsRead,
		UpdateContext: resourceDropMetricsUpdate,
		DeleteContext: resourceDropMetricsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			dropMetricsIdField: {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			dropMetricsAccountId: {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			dropMetricsActive: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			dropMetricsFilters: {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: filterExprSchema,
				},
				Set: schema.HashResource(&schema.Resource{
					Schema: filterExprSchema,
				}),
			},
			dropMetricsFilterOperator: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  drop_metrics.OperatorAnd,
				ValidateFunc: validation.StringInSlice(
					[]string{drop_metrics.OperatorAnd},
					false),
			},
			dropMetricsCreatedAt: {
				Type:     schema.TypeString,
				Computed: true,
			},
			dropMetricsCreatedBy: {
				Type:     schema.TypeString,
				Computed: true,
			},
			dropMetricsModifiedAt: {
				Type:     schema.TypeString,
				Computed: true,
			},
			dropMetricsModifiedBy: {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// resourceDropMetricsCreate creates a new metrics drop filter in logzio
func resourceDropMetricsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	createDropMetrics := createCreateUpdateDropMetricsFromSchema(d)

	// dropFilter, err := dropFilterClient(m).CreateDropFilter(createDropFilter)
	dropMetrics, err := dropMetricsClient(m).CreateDropMetric(createDropMetrics)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(dropMetrics.Id, 10))

	return resourceDropMetricsRead(ctx, d, m)
}

// resourceDropMetricsRead gets metrics drop filter by id
func resourceDropMetricsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	dropId, err := stateIDInt64(d)
	if err != nil {
		return diag.FromErr(err)
	}
	dropMetrics, err := dropMetricsClient(m).GetDropMetric(dropId)
	if err != nil {
		return diag.FromErr(err)
	}
	if dropMetrics == nil {
		// If we were not able to find the resource - delete from state
		tflog.Error(ctx, "could not find metrics drop filter with id: "+d.Id())
		d.SetId("")
		return diag.Diagnostics{}
	}
	setDropMetrics(d, dropMetrics)
	return nil
}

// resourceDropMetricsUpdate updates an existing metrics drop filter in logzio
func resourceDropMetricsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	dropId, err := stateIDInt64(d)
	if err != nil {
		return diag.FromErr(err)
	}

	onlyActiveChanged := d.HasChange(dropMetricsActive) && !d.HasChanges(dropMetricsAccountId, dropMetricsFilters, dropMetricsFilterOperator)

	if onlyActiveChanged {
		return updateActiveState(ctx, d, m, dropId)
	}
	return updateDropMetrics(ctx, d, m, dropId)
}

// resourceDropMetricsDelete deletes a metrics drop filter in logzio
func resourceDropMetricsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	dropId, err := stateIDInt64(d)
	if err != nil {
		return diag.FromErr(err)
	}

	err = dropMetricsClient(m).DeleteDropMetric(dropId)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func createCreateUpdateDropMetricsFromSchema(d *schema.ResourceData) drop_metrics.CreateUpdateDropMetric {
	activeVal := d.Get(dropMetricsActive).(bool)
	activePtr := &activeVal

	return drop_metrics.CreateUpdateDropMetric{
		AccountId: int64Attr(d, dropMetricsAccountId),
		Active:    activePtr,
		Filter: drop_metrics.FilterObject{
			Operator:   d.Get(dropMetricsFilterOperator).(string),
			Expression: schemaToDropMetricsFilterExpression(d),
		},
	}
}

func setDropMetrics(d *schema.ResourceData, dropMetric *drop_metrics.DropMetric) {
	d.Set(dropMetricsIdField, int(dropMetric.Id))
	d.Set(dropMetricsAccountId, int(dropMetric.AccountId))
	d.Set(dropMetricsActive, dropMetric.Active)
	d.Set(dropMetricsCreatedAt, dropMetric.CreatedAt)
	d.Set(dropMetricsCreatedBy, dropMetric.CreatedBy)
	d.Set(dropMetricsModifiedAt, dropMetric.ModifiedAt)
	d.Set(dropMetricsModifiedBy, dropMetric.ModifiedBy)
	d.Set(dropMetricsFilters, dropMetricsFilterExpressionToInterface(dropMetric.Filter.Expression))
	d.Set(dropMetricsFilterOperator, dropMetric.Filter.Operator)
}

func schemaToDropMetricsFilterExpression(d *schema.ResourceData) []drop_metrics.FilterExpression {
	raw := d.Get(dropMetricsFilters).(*schema.Set).List()
	result := make([]drop_metrics.FilterExpression, 0, len(raw))
	for _, item := range raw {
		m := item.(map[string]interface{})
		result = append(result, drop_metrics.FilterExpression{
			Name:             m[dropMetricsExpressionLabelName].(string),
			Value:            m[dropMetricsExpressionLabelValue].(string),
			ComparisonFilter: m[dropMetricsExpressionCondition].(string),
		})
	}
	return result
}

func dropMetricsFilterExpressionToInterface(expressions []drop_metrics.FilterExpression) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(expressions))
	for _, e := range expressions {
		result = append(result, map[string]interface{}{
			dropMetricsExpressionLabelName:  e.Name,
			dropMetricsExpressionLabelValue: e.Value,
			dropMetricsExpressionCondition:  e.ComparisonFilter,
		})
	}
	return result
}

func int64Attr(d *schema.ResourceData, key string) int64 {
	return int64(d.Get(key).(int))
}

func stateIDInt64(d *schema.ResourceData) (int64, error) {
	id := d.Id()
	if id == "" {
		return 0, fmt.Errorf("resource ID is empty")
	}
	n, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid resource ID %q: %v", id, err)
	}
	return n, nil
}

func updateActiveState(ctx context.Context, d *schema.ResourceData, m interface{}, dropId int64) diag.Diagnostics {
	activate := d.Get(dropMetricsActive).(bool)
	var err error
	if activate {
		err = dropMetricsClient(m).EnableDropMetric(dropId)
	} else {
		err = dropMetricsClient(m).DisableDropMetric(dropId)
	}
	if err != nil {
		return diag.FromErr(err)
	}

	diags := readUntilConsistent(ctx, d, m, dropMetricsRetryAttempts, "toggle active", func() bool {
		return d.Get(dropMetricsActive).(bool) == activate
	})

	return diags
}

func updateDropMetrics(ctx context.Context, d *schema.ResourceData, m interface{}, dropId int64) diag.Diagnostics {
	updateFilter := createCreateUpdateDropMetricsFromSchema(d)

	client := dropMetricsClient(m)
	_, err := client.UpdateDropMetric(dropId, updateFilter)
	if err != nil {
		return diag.FromErr(err)
	}

	diags := readUntilConsistent(ctx, d, m, dropMetricsRetryAttempts, "update filters", func() bool {
		createFilter := createCreateUpdateDropMetricsFromSchema(d)
		return !reflect.DeepEqual(createFilter, updateFilter)
	})
	return diags
}

func readUntilConsistent(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
	attempts uint,
	tag string,
	consistent func() bool,
) diag.Diagnostics {
	var ret diag.Diagnostics
	readErr := retry.Do(
		func() error {
			ret = resourceDropMetricsRead(ctx, d, m)
			if ret.HasError() {
				return fmt.Errorf("%s: read failed", tag)
			}
			if !consistent() {
				return fmt.Errorf("%s: not yet consistent", tag)
			}
			return nil
		},
		retry.DelayType(retry.BackOffDelay),
		retry.Attempts(attempts),
	)
	if readErr != nil {
		tflog.Warn(ctx, tag+" not reflected yet; returning last read")
	}
	return ret
}
