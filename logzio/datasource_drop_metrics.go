package logzio

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/go-cty/cty/gocty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/logzio/logzio_terraform_client/drop_metrics"
)

const (
	metricNameLabel           = "__name__"
	dropMetricsSearchPageSize = 200
	nullByte                  = "\x00"
)

func dataSourceDropMetrics() *schema.Resource {
	var filterExprSchema = map[string]*schema.Schema{
		dropMetricsExpressionLabelName: {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		dropMetricsExpressionLabelValue: {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		dropMetricsExpressionCondition: {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
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
		ReadContext: dataSourceDropMetricsRead,
		Schema: map[string]*schema.Schema{
			dropMetricsIdField: {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			dropMetricsAccountId: {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			dropMetricsName: {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			dropMetricsActive: {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			dropMetricsFilters: {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
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
				Computed: true,
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

// dataSourceDropMetricsRead gets a metrics drop filter from Logz.io.
// It first tries to find the filter by ID, and if not found, it searches for a filter that matches the specified attributes and expressions.
func dataSourceDropMetricsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var dropMetricsFilter *drop_metrics.DropMetric
	client := dropMetricsClient(m)
	dropMetricsId := int64Attr(d, dropMetricsIdField)

	if dropMetricsId != 0 {
		dropMetricsFilter, err := client.GetDropMetric(dropMetricsId)
		if err != nil {
			if strings.Contains(strings.ToUpper(err.Error()), "METRICS_DROP_FILTERS/DROP_FILTER_NOT_FOUND") {
				// fallback to search for drop metrics filter with the specified attributes
			} else {
				return diag.FromErr(err)
			}
		}

		if dropMetricsFilter != nil {
			setFoundDatasourceDropMetrics(d, dropMetricsFilter)
			return nil
		}
	}

	// If could not find id, search for drop filter with the specified attributes
	dropMetricsFilter, err := searchMatchingDropMetrics(client, d)
	if err != nil {
		if dropMetricsId != 0 {
			return diag.FromErr(fmt.Errorf("could not find drop metrics filter with id %d", dropMetricsId))
		}
		return diag.FromErr(err)
	}
	setFoundDatasourceDropMetrics(d, dropMetricsFilter)
	return nil
}

// searchMatchingDropMetrics searches for a metrics drop filter that matches the specified attributes in the schema.
// It preforms a pagination search looking for a drop filter that matches the exact expressions specified in the schema.
func searchMatchingDropMetrics(client *drop_metrics.DropMetricsClient, d *schema.ResourceData) (*drop_metrics.DropMetric, error) {
	pageNum := 1
	req, err := buildSearchReqFromDS(d, pageNum)
	if err != nil {
		return nil, err
	}

	results, err := client.SearchDropMetrics(req)
	if err != nil {
		return nil, err
	}

	wanted := wantedExpressionsFromSchema(d)
	if len(wanted) == 0 {
		switch len(results) {
		case 0:
			filterJSON, _ := json.Marshal(req.Filter)
			return nil, fmt.Errorf("no metrics drop filters matched the criteria: %s", string(filterJSON))
		case 1:
			return &results[0], nil
		default:
			return nil, fmt.Errorf("found multiple (%d) metrics drop filters matching the critiria, add more expressions or specify an id", len(results))
		}
	}

	// search for the wanted drop filter in the results, preforms page by page search
	for {
		if len(results) == 0 {
			break
		}

		for i := range results {
			result := &results[i]
			resultFilterExps := result.Filter.Expression
			if len(resultFilterExps) != len(wanted) {
				continue
			}

			resultSet := make(map[string]struct{}, len(resultFilterExps))
			for _, e := range resultFilterExps {
				resultSet[exprKey(e)] = struct{}{}
			}

			allFound := true
			for _, w := range wanted {
				if _, ok := resultSet[exprKey(w)]; !ok {
					allFound = false
					break
				}
			}
			if allFound {
				return result, nil
			}
		}

		pageNum++
		req.Pagination.PageNumber = pageNum
		results, err = client.SearchDropMetrics(req)
		if err != nil {
			return nil, err
		}
	}

	return nil, fmt.Errorf("could not find metrics drop filter with the specified attributes: %v", wanted)
}

// buildSearchReqFromDS builds a search request from the schema.ResourceData.
// It uses the attributes specified in the schema to create a search request for metrics drop filters.
func buildSearchReqFromDS(d *schema.ResourceData, pageNumber int) (drop_metrics.SearchDropMetricsRequest, error) {
	var req drop_metrics.SearchDropMetricsRequest
	var filter drop_metrics.SearchFilter
	setSearchFilters := 0

	accountId := int64Attr(d, dropMetricsAccountId)
	if accountId != 0 {
		filter.AccountIds = []int64{accountId}
		setSearchFilters++
	}

	metricNamesFromFilter := getMetricNamesFromSchema(d)
	if len(metricNamesFromFilter) > 0 {
		filter.MetricNames = metricNamesFromFilter
		setSearchFilters++
	}

	if p, err := optionalBoolPtr(d, dropMetricsActive); err != nil {
		return req, err
	} else if p != nil {
		filter.Active = p
		setSearchFilters++
	}

	if setSearchFilters == 0 {
		return req, fmt.Errorf("could not find drop metrics filter with id %d, and not enough criteria to search", int64Attr(d, dropMetricsIdField))
	}

	req.Filter = &filter

	req.Pagination = &drop_metrics.Pagination{PageNumber: pageNumber, PageSize: dropMetricsSearchPageSize}
	return req, nil
}

// optionalBoolPtr returns a pointer to a boolean value if the key exists and is not null or unknown.
func optionalBoolPtr(d *schema.ResourceData, key string) (*bool, error) {
	v, diags := d.GetRawConfigAt(cty.GetAttrPath(key))
	if diags.HasError() || v.IsNull() || !v.IsKnown() {
		return nil, nil
	}
	var b bool
	if err := gocty.FromCtyValue(v, &b); err != nil {
		return nil, fmt.Errorf("invalid %s: %w", key, err)
	}
	return &b, nil
}

// getMetricNamesFromSchema extracts metric names used in the filter of the given metric drop filter schema.
func getMetricNamesFromSchema(d *schema.ResourceData) []string {
	var result []string
	filterExps, ok := d.GetOk(dropMetricsFilters)
	if !ok {
		return result
	}

	set := filterExps.(*schema.Set)
	for _, raw := range set.List() {
		m := raw.(map[string]interface{})
		name := m[dropMetricsExpressionLabelName].(string)
		val := m[dropMetricsExpressionLabelValue].(string)

		if name == metricNameLabel {
			result = append(result, val)
		}
	}
	return result
}

// setFoundDatasourceDropMetrics sets the schema.ResourceData with the found drop metrics filter.
// It sets the ID and all other attributes of the drop metrics filter.
func setFoundDatasourceDropMetrics(d *schema.ResourceData, dropMetrics *drop_metrics.DropMetric) {
	d.SetId(int64ToStr(dropMetrics.Id))
	setDropMetrics(d, dropMetrics)
}

// wantedExpressionsFromSchema extracts the wanted filter expressions from the schema in drop_metrics.FilterExpression format.
func wantedExpressionsFromSchema(d *schema.ResourceData) []drop_metrics.FilterExpression {
	rawList := d.Get(dropMetricsFilters).(*schema.Set).List()
	exprs := make([]drop_metrics.FilterExpression, 0, len(rawList))

	for _, raw := range rawList {
		m := raw.(map[string]interface{})
		exprs = append(exprs, drop_metrics.FilterExpression{
			Name:             strings.TrimSpace(m[dropMetricsExpressionLabelName].(string)),
			Value:            strings.TrimSpace(m[dropMetricsExpressionLabelValue].(string)),
			ComparisonFilter: strings.TrimSpace(m[dropMetricsExpressionCondition].(string)),
		})
	}
	return exprs
}

// exprKey generates a unique key for a filter expression by combining its name, comparison filter, and value.
// A nullByte ("\x00") separator is used between parts to prevent accidental collisions from
// concatenation and to ensure the key can be used safely in hash-based
// lookups (O(1) search in maps).
func exprKey(e drop_metrics.FilterExpression) string {
	return strings.ToLower(strings.TrimSpace(e.Name)) + nullByte +
		strings.ToUpper(strings.TrimSpace(e.ComparisonFilter)) + nullByte +
		strings.TrimSpace(e.Value)
}
