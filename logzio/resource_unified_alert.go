package logzio

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/avast/retry-go"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/logzio/logzio_terraform_client/unified_alerts"
)

const (
	// Top-level fields
	unifiedAlertId                                  = "alert_id"
	unifiedAlertTitle                               = "title"
	unifiedAlertType                                = "type"
	unifiedAlertDescription                         = "description"
	unifiedAlertTags                                = "tags"
	unifiedAlertFolderId                            = "folder_id"
	unifiedAlertDashboardId                         = "dashboard_id"
	unifiedAlertPanelId                             = "panel_id"
	unifiedAlertRunbook                             = "runbook"
	unifiedAlertEnabled                             = "enabled"
	unifiedAlertRca                                 = "rca"
	unifiedAlertRcaNotificationEndpointIds          = "rca_notification_endpoint_ids"
	unifiedAlertUseAlertNotificationEndpointsForRca = "use_alert_notification_endpoints_for_rca"
	unifiedAlertCreatedAt                           = "created_at"
	unifiedAlertUpdatedAt                           = "updated_at"

	// Log alert fields
	unifiedAlertLogAlert                       = "log_alert"
	logAlertOutput                             = "output"
	logAlertSearchTimeFrameMinutes             = "search_timeframe_minutes"
	logAlertSubComponents                      = "sub_components"
	logAlertCorrelations                       = "correlations"
	logAlertSchedule                           = "schedule"
	logAlertOutputRecipients                   = "recipients"
	logAlertOutputSuppressNotificationsMinutes = "suppress_notifications_minutes"
	logAlertOutputType                         = "type"

	// Recipients fields
	recipientsEmails                  = "emails"
	recipientsNotificationEndpointIds = "notification_endpoint_ids"

	// SubComponent fields
	subComponentQueryDefinition = "query_definition"
	subComponentTrigger         = "trigger"
	subComponentOutput          = "output"

	// QueryDefinition fields
	queryDefinitionQuery                    = "query"
	queryDefinitionFilters                  = "filters"
	queryDefinitionGroupBy                  = "group_by"
	queryDefinitionAggregation              = "aggregation"
	queryDefinitionShouldQueryOnAllAccounts = "should_query_on_all_accounts"
	queryDefinitionAccountIdsToQueryOn      = "account_ids_to_query_on"

	// BoolFilter fields
	boolFilterMust    = "must"
	boolFilterShould  = "should"
	boolFilterFilter  = "filter"
	boolFilterMustNot = "must_not"

	// Aggregation fields
	aggregationAggregationType    = "aggregation_type"
	aggregationFieldToAggregateOn = "field_to_aggregate_on"
	aggregationValueToAggregateOn = "value_to_aggregate_on"

	// SubComponentTrigger fields
	triggerOperator               = "operator"
	triggerSeverityThresholdTiers = "severity_threshold_tiers"

	// SeverityThresholdTier fields
	severityThresholdTierSeverity  = "severity"
	severityThresholdTierThreshold = "threshold"

	// SubComponentOutput fields
	subComponentOutputColumns            = "columns"
	subComponentOutputShouldUseAllFields = "should_use_all_fields"

	// ColumnConfig fields
	columnConfigFieldName = "field_name"
	columnConfigRegex     = "regex"
	columnConfigSort      = "sort"

	// Schedule fields
	scheduleCronExpression = "cron_expression"
	scheduleTimezone       = "timezone"

	// Correlations fields
	correlationsCorrelationOperators = "correlation_operators"
	correlationsJoins                = "joins"

	// Metric alert fields
	unifiedAlertMetricAlert = "metric_alert"
	metricAlertSeverity     = "severity"
	metricAlertTrigger      = "trigger"
	metricAlertQueries      = "queries"
	metricAlertRecipients   = "recipients"

	// MetricTrigger fields
	metricTriggerTriggerType            = "trigger_type"
	metricTriggerMetricOperator         = "metric_operator"
	metricTriggerMinThreshold           = "min_threshold"
	metricTriggerMaxThreshold           = "max_threshold"
	metricTriggerMathExpression         = "math_expression"
	metricTriggerSearchTimeFrameMinutes = "search_timeframe_minutes"

	// MetricQuery fields
	metricQueryRefId           = "ref_id"
	metricQueryQueryDefinition = "query_definition"

	// MetricQueryDefinition fields
	metricQueryDefinitionDatasourceUid = "datasource_uid"
	metricQueryDefinitionPromqlQuery   = "promql_query"

	unifiedAlertRetryAttempts = 8
)

// unifiedAlertClient returns the unified alert client with the api token from the provider
func unifiedAlertClient(m interface{}) *unified_alerts.UnifiedAlertsClient {
	var client *unified_alerts.UnifiedAlertsClient
	client, _ = unified_alerts.New(m.(Config).apiToken, m.(Config).baseUrl)
	return client
}

func resourceUnifiedAlert() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUnifiedAlertCreate,
		ReadContext:   resourceUnifiedAlertRead,
		UpdateContext: resourceUnifiedAlertUpdate,
		DeleteContext: resourceUnifiedAlertDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			unifiedAlertId: {
				Type:     schema.TypeString,
				Computed: true,
			},
			unifiedAlertTitle: {
				Type:     schema.TypeString,
				Required: true,
			},
			unifiedAlertType: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"LOG_ALERT", "METRIC_ALERT"}, false),
			},
			unifiedAlertDescription: {
				Type:     schema.TypeString,
				Optional: true,
			},
			unifiedAlertTags: {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			unifiedAlertFolderId: {
				Type:     schema.TypeString,
				Optional: true,
			},
			unifiedAlertDashboardId: {
				Type:     schema.TypeString,
				Optional: true,
			},
			unifiedAlertPanelId: {
				Type:     schema.TypeString,
				Optional: true,
			},
			unifiedAlertRunbook: {
				Type:     schema.TypeString,
				Optional: true,
			},
			unifiedAlertEnabled: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			unifiedAlertRca: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			unifiedAlertRcaNotificationEndpointIds: {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			unifiedAlertUseAlertNotificationEndpointsForRca: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			unifiedAlertCreatedAt: {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			unifiedAlertUpdatedAt: {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			unifiedAlertLogAlert: {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{unifiedAlertMetricAlert},
				Elem:          resourceLogAlertConfig(),
			},
			unifiedAlertMetricAlert: {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{unifiedAlertLogAlert},
				Elem:          resourceMetricAlertConfig(),
			},
		},
	}
}

func resourceLogAlertConfig() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			logAlertSearchTimeFrameMinutes: {
				Type:     schema.TypeInt,
				Required: true,
			},
			logAlertOutput: {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						logAlertOutputType: {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"JSON", "TABLE"}, false),
						},
						logAlertOutputSuppressNotificationsMinutes: {
							Type:     schema.TypeInt,
							Optional: true,
						},
						logAlertOutputRecipients: {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem:     resourceRecipients(),
						},
					},
				},
			},
			logAlertSubComponents: {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem:     resourceSubComponent(),
			},
			logAlertCorrelations: {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						correlationsCorrelationOperators: {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						correlationsJoins: {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeMap,
								Elem: &schema.Schema{
									Type: schema.TypeString,
								},
							},
						},
					},
				},
			},
			logAlertSchedule: {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						scheduleCronExpression: {
							Type:     schema.TypeString,
							Required: true,
						},
						scheduleTimezone: {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "UTC",
						},
					},
				},
			},
		},
	}
}

func resourceRecipients() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			recipientsEmails: {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			recipientsNotificationEndpointIds: {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
	}
}

func resourceSubComponent() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			subComponentQueryDefinition: {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						queryDefinitionQuery: {
							Type:     schema.TypeString,
							Required: true,
						},
						queryDefinitionFilters: {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsJSON,
						},
						queryDefinitionGroupBy: {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						queryDefinitionAggregation: {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									aggregationAggregationType: {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"SUM", "MIN", "MAX", "AVG", "COUNT", "UNIQUE_COUNT", "NONE"}, false),
									},
									aggregationFieldToAggregateOn: {
										Type:     schema.TypeString,
										Optional: true,
									},
									aggregationValueToAggregateOn: {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						queryDefinitionShouldQueryOnAllAccounts: {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						queryDefinitionAccountIdsToQueryOn: {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
					},
				},
			},
			subComponentTrigger: {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						triggerOperator: {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"LESS_THAN", "GREATER_THAN", "LESS_THAN_OR_EQUALS", "GREATER_THAN_OR_EQUALS", "EQUALS", "NOT_EQUALS"}, false),
						},
						triggerSeverityThresholdTiers: {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									severityThresholdTierSeverity: {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"INFO", "LOW", "MEDIUM", "HIGH", "SEVERE"}, false),
									},
									severityThresholdTierThreshold: {
										Type:     schema.TypeFloat,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			subComponentOutput: {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						subComponentOutputShouldUseAllFields: {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						subComponentOutputColumns: {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									columnConfigFieldName: {
										Type:     schema.TypeString,
										Required: true,
									},
									columnConfigRegex: {
										Type:     schema.TypeString,
										Optional: true,
									},
									columnConfigSort: {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"ASC", "DESC"}, false),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceMetricAlertConfig() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			metricAlertSeverity: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"INFO", "LOW", "MEDIUM", "HIGH", "SEVERE"}, false),
			},
			metricAlertTrigger: {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						metricTriggerTriggerType: {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"THRESHOLD", "MATH"}, false),
						},
						metricTriggerMetricOperator: {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"ABOVE", "BELOW", "WITHIN_RANGE", "OUTSIDE_RANGE"}, false),
						},
						metricTriggerMinThreshold: {
							Type:     schema.TypeFloat,
							Optional: true,
						},
						metricTriggerMaxThreshold: {
							Type:     schema.TypeFloat,
							Optional: true,
						},
						metricTriggerMathExpression: {
							Type:     schema.TypeString,
							Optional: true,
						},
						metricTriggerSearchTimeFrameMinutes: {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			metricAlertQueries: {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						metricQueryRefId: {
							Type:     schema.TypeString,
							Required: true,
						},
						metricQueryQueryDefinition: {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									metricQueryDefinitionDatasourceUid: {
										Type:     schema.TypeString,
										Required: true,
									},
									metricQueryDefinitionPromqlQuery: {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			metricAlertRecipients: {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     resourceRecipients(),
			},
		},
	}
}

// resourceUnifiedAlertCreate creates a new unified alert in logzio
func resourceUnifiedAlertCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	alertType := d.Get(unifiedAlertType).(string)
	createAlert, urlType, err := buildCreateUnifiedAlert(d)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonBytes, _ := json.Marshal(createAlert)
	tflog.Debug(ctx, fmt.Sprintf("Creating unified alert: %s", string(jsonBytes)))

	client := unifiedAlertClient(m)
	alert, err := client.CreateUnifiedAlert(urlType, createAlert)
	if err != nil {
		return diag.Errorf("failed to create unified alert: %v", err)
	}

	d.SetId(alert.Id)
	tflog.Info(ctx, fmt.Sprintf("Created unified alert with ID: %s, Type: %s", alert.Id, alertType))

	return resourceUnifiedAlertRead(ctx, d, m)
}

// resourceUnifiedAlertRead reads a unified alert from logzio
func resourceUnifiedAlertRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	alertId := d.Id()
	alertType := d.Get(unifiedAlertType).(string)
	urlType := getUrlTypeFromAlertType(alertType)

	client := unifiedAlertClient(m)
	alert, err := client.GetUnifiedAlert(urlType, alertId)
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("Failed to get unified alert: %v", err))
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
			// If we were not able to find the resource - delete from state
			d.SetId("")
			return diag.Diagnostics{}
		}
		return diag.FromErr(err)
	}

	return setUnifiedAlert(d, alert)
}

// resourceUnifiedAlertUpdate updates an existing unified alert in logzio
func resourceUnifiedAlertUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	alertId := d.Id()
	alertType := d.Get(unifiedAlertType).(string)
	createAlert, urlType, err := buildCreateUnifiedAlert(d)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonBytes, _ := json.Marshal(createAlert)
	tflog.Debug(ctx, fmt.Sprintf("Updating unified alert %s: %s", alertId, string(jsonBytes)))

	client := unifiedAlertClient(m)
	_, err = client.UpdateUnifiedAlert(urlType, alertId, createAlert)
	if err != nil {
		return diag.Errorf("failed to update unified alert: %v", err)
	}

	var diagRet diag.Diagnostics
	readErr := retry.Do(
		func() error {
			diagRet = resourceUnifiedAlertRead(ctx, d, m)
			if diagRet.HasError() {
				return fmt.Errorf("received error from read unified alert")
			}
			return nil
		},
		retry.RetryIf(func(err error) bool {
			return err != nil
		}),
		retry.Attempts(unifiedAlertRetryAttempts),
		retry.DelayType(retry.BackOffDelay),
	)

	if readErr != nil {
		tflog.Warn(ctx, fmt.Sprintf("Failed to read unified alert after update: %v", readErr))
	}

	tflog.Info(ctx, fmt.Sprintf("Updated unified alert with ID: %s, Type: %s", alertId, alertType))
	return diagRet
}

// resourceUnifiedAlertDelete deletes a unified alert from logzio
func resourceUnifiedAlertDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	alertId := d.Id()
	alertType := d.Get(unifiedAlertType).(string)
	urlType := getUrlTypeFromAlertType(alertType)

	client := unifiedAlertClient(m)
	_, err := client.DeleteUnifiedAlert(urlType, alertId)
	if err != nil {
		return diag.Errorf("failed to delete unified alert: %v", err)
	}

	tflog.Info(ctx, fmt.Sprintf("Deleted unified alert with ID: %s, Type: %s", alertId, alertType))
	return nil
}

// Helper functions

func getUrlTypeFromAlertType(alertType string) string {
	if alertType == unified_alerts.TypeLogAlert {
		return unified_alerts.UrlTypeLogs
	}
	return unified_alerts.UrlTypeMetrics
}

func buildCreateUnifiedAlert(d *schema.ResourceData) (unified_alerts.CreateUnifiedAlert, string, error) {
	alertType := d.Get(unifiedAlertType).(string)
	urlType := getUrlTypeFromAlertType(alertType)

	enabled := d.Get(unifiedAlertEnabled).(bool)

	alert := unified_alerts.CreateUnifiedAlert{
		Title:                               d.Get(unifiedAlertTitle).(string),
		Type:                                alertType,
		Description:                         d.Get(unifiedAlertDescription).(string),
		Tags:                                interfaceSliceToStringSlice(d.Get(unifiedAlertTags).([]interface{})),
		FolderId:                            d.Get(unifiedAlertFolderId).(string),
		DashboardId:                         d.Get(unifiedAlertDashboardId).(string),
		PanelId:                             d.Get(unifiedAlertPanelId).(string),
		Runbook:                             d.Get(unifiedAlertRunbook).(string),
		Enabled:                             &enabled,
		Rca:                                 d.Get(unifiedAlertRca).(bool),
		RcaNotificationEndpointIds:          interfaceSliceToIntSlice(d.Get(unifiedAlertRcaNotificationEndpointIds).([]interface{})),
		UseAlertNotificationEndpointsForRca: d.Get(unifiedAlertUseAlertNotificationEndpointsForRca).(bool),
	}

	if alertType == unified_alerts.TypeLogAlert {
		logAlert, err := buildLogAlertConfig(d)
		if err != nil {
			return alert, urlType, err
		}
		alert.LogAlert = logAlert
	} else if alertType == unified_alerts.TypeMetricAlert {
		metricAlert, err := buildMetricAlertConfig(d)
		if err != nil {
			return alert, urlType, err
		}
		alert.MetricAlert = metricAlert
	}

	return alert, urlType, nil
}

func buildLogAlertConfig(d *schema.ResourceData) (*unified_alerts.LogAlertConfig, error) {
	logAlertList := d.Get(unifiedAlertLogAlert).([]interface{})
	if len(logAlertList) == 0 {
		return nil, fmt.Errorf("log_alert configuration is required for LOG_ALERT type")
	}

	logAlertMap := logAlertList[0].(map[string]interface{})

	config := &unified_alerts.LogAlertConfig{
		SearchTimeFrameMinutes: logAlertMap[logAlertSearchTimeFrameMinutes].(int),
	}

	// Parse output
	outputList := logAlertMap[logAlertOutput].([]interface{})
	if len(outputList) > 0 {
		outputMap := outputList[0].(map[string]interface{})
		config.Output = unified_alerts.LogAlertOutput{
			Type:                         outputMap[logAlertOutputType].(string),
			SuppressNotificationsMinutes: outputMap[logAlertOutputSuppressNotificationsMinutes].(int),
		}

		// Parse recipients
		recipientsList := outputMap[logAlertOutputRecipients].([]interface{})
		if len(recipientsList) > 0 {
			recipientsMap := recipientsList[0].(map[string]interface{})
			config.Output.Recipients = unified_alerts.Recipients{
				Emails:                  interfaceSliceToStringSlice(recipientsMap[recipientsEmails].([]interface{})),
				NotificationEndpointIds: interfaceSliceToIntSlice(recipientsMap[recipientsNotificationEndpointIds].([]interface{})),
			}
		}
	}

	// Parse sub components
	subComponentsList := logAlertMap[logAlertSubComponents].([]interface{})
	config.SubComponents = make([]unified_alerts.SubComponent, len(subComponentsList))
	for i, scItem := range subComponentsList {
		scMap := scItem.(map[string]interface{})

		subComp := unified_alerts.SubComponent{}

		// Parse query definition
		queryDefList := scMap[subComponentQueryDefinition].([]interface{})
		if len(queryDefList) > 0 {
			queryDefMap := queryDefList[0].(map[string]interface{})
			subComp.QueryDefinition = unified_alerts.QueryDefinition{
				Query:                    queryDefMap[queryDefinitionQuery].(string),
				GroupBy:                  interfaceSliceToStringSlice(queryDefMap[queryDefinitionGroupBy].([]interface{})),
				ShouldQueryOnAllAccounts: queryDefMap[queryDefinitionShouldQueryOnAllAccounts].(bool),
				AccountIdsToQueryOn:      interfaceSliceToIntSlice(queryDefMap[queryDefinitionAccountIdsToQueryOn].([]interface{})),
			}

			// Parse filters if present
			if filtersStr, ok := queryDefMap[queryDefinitionFilters].(string); ok && filtersStr != "" {
				var filters unified_alerts.BoolFilter
				if err := json.Unmarshal([]byte(filtersStr), &filters); err != nil {
					return nil, fmt.Errorf("failed to parse filters JSON: %v", err)
				}
				subComp.QueryDefinition.Filters = filters
			}

			// Parse aggregation if present
			aggregationList := queryDefMap[queryDefinitionAggregation].([]interface{})
			if len(aggregationList) > 0 {
				aggMap := aggregationList[0].(map[string]interface{})
				subComp.QueryDefinition.Aggregation = unified_alerts.Aggregation{
					AggregationType:    aggMap[aggregationAggregationType].(string),
					FieldToAggregateOn: aggMap[aggregationFieldToAggregateOn].(string),
					ValueToAggregateOn: aggMap[aggregationValueToAggregateOn].(string),
				}
			}
		}

		// Parse trigger
		triggerList := scMap[subComponentTrigger].([]interface{})
		if len(triggerList) > 0 {
			triggerMap := triggerList[0].(map[string]interface{})
			subComp.Trigger = unified_alerts.SubComponentTrigger{
				Operator:               triggerMap[triggerOperator].(string),
				SeverityThresholdTiers: make(map[string]float32),
			}

			// Parse severity threshold tiers
			tiersList := triggerMap[triggerSeverityThresholdTiers].([]interface{})
			for _, tierItem := range tiersList {
				tierMap := tierItem.(map[string]interface{})
				severity := tierMap[severityThresholdTierSeverity].(string)
				threshold := float32(tierMap[severityThresholdTierThreshold].(float64))
				subComp.Trigger.SeverityThresholdTiers[severity] = threshold
			}
		}

		// Parse output if present
		outputList := scMap[subComponentOutput].([]interface{})
		if len(outputList) > 0 {
			outputMap := outputList[0].(map[string]interface{})
			subComp.Output = unified_alerts.SubComponentOutput{
				ShouldUseAllFields: outputMap[subComponentOutputShouldUseAllFields].(bool),
			}

			// Parse columns if present
			columnsList := outputMap[subComponentOutputColumns].([]interface{})
			if len(columnsList) > 0 {
				subComp.Output.Columns = make([]unified_alerts.ColumnConfig, len(columnsList))
				for j, colItem := range columnsList {
					colMap := colItem.(map[string]interface{})
					subComp.Output.Columns[j] = unified_alerts.ColumnConfig{
						FieldName: colMap[columnConfigFieldName].(string),
						Regex:     colMap[columnConfigRegex].(string),
						Sort:      colMap[columnConfigSort].(string),
					}
				}
			}
		}

		config.SubComponents[i] = subComp
	}

	// Parse correlations if present
	correlationsList := logAlertMap[logAlertCorrelations].([]interface{})
	if len(correlationsList) > 0 {
		correlationsMap := correlationsList[0].(map[string]interface{})
		config.Correlations = unified_alerts.Correlations{
			CorrelationOperators: interfaceSliceToStringSlice(correlationsMap[correlationsCorrelationOperators].([]interface{})),
		}

		// Parse joins if present
		joinsList := correlationsMap[correlationsJoins].([]interface{})
		if len(joinsList) > 0 {
			config.Correlations.Joins = make([]map[string]string, len(joinsList))
			for i, joinItem := range joinsList {
				joinMapInterface := joinItem.(map[string]interface{})
				joinMap := make(map[string]string)
				for k, v := range joinMapInterface {
					joinMap[k] = v.(string)
				}
				config.Correlations.Joins[i] = joinMap
			}
		}
	}

	// Parse schedule if present
	scheduleList := logAlertMap[logAlertSchedule].([]interface{})
	if len(scheduleList) > 0 {
		scheduleMap := scheduleList[0].(map[string]interface{})
		config.Schedule = unified_alerts.Schedule{
			CronExpression: scheduleMap[scheduleCronExpression].(string),
			Timezone:       scheduleMap[scheduleTimezone].(string),
		}
	}

	return config, nil
}

func buildMetricAlertConfig(d *schema.ResourceData) (*unified_alerts.MetricAlertConfig, error) {
	metricAlertList := d.Get(unifiedAlertMetricAlert).([]interface{})
	if len(metricAlertList) == 0 {
		return nil, fmt.Errorf("metric_alert configuration is required for METRIC_ALERT type")
	}

	metricAlertMap := metricAlertList[0].(map[string]interface{})

	config := &unified_alerts.MetricAlertConfig{
		Severity: metricAlertMap[metricAlertSeverity].(string),
	}

	// Parse trigger
	triggerList := metricAlertMap[metricAlertTrigger].([]interface{})
	if len(triggerList) > 0 {
		triggerMap := triggerList[0].(map[string]interface{})
		config.Trigger = unified_alerts.MetricTrigger{
			TriggerType:            triggerMap[metricTriggerTriggerType].(string),
			MetricOperator:         triggerMap[metricTriggerMetricOperator].(string),
			MinThreshold:           triggerMap[metricTriggerMinThreshold].(float64),
			MaxThreshold:           triggerMap[metricTriggerMaxThreshold].(float64),
			MathExpression:         triggerMap[metricTriggerMathExpression].(string),
			SearchTimeFrameMinutes: triggerMap[metricTriggerSearchTimeFrameMinutes].(int),
		}
	}

	// Parse queries
	queriesList := metricAlertMap[metricAlertQueries].([]interface{})
	config.Queries = make([]unified_alerts.MetricQuery, len(queriesList))
	for i, queryItem := range queriesList {
		queryMap := queryItem.(map[string]interface{})

		metricQuery := unified_alerts.MetricQuery{
			RefId: queryMap[metricQueryRefId].(string),
		}

		// Parse query definition
		queryDefList := queryMap[metricQueryQueryDefinition].([]interface{})
		if len(queryDefList) > 0 {
			queryDefMap := queryDefList[0].(map[string]interface{})
			metricQuery.QueryDefinition = unified_alerts.MetricQueryDefinition{
				DatasourceUid: queryDefMap[metricQueryDefinitionDatasourceUid].(string),
				PromqlQuery:   queryDefMap[metricQueryDefinitionPromqlQuery].(string),
			}
		}

		config.Queries[i] = metricQuery
	}

	// Parse recipients
	recipientsList := metricAlertMap[metricAlertRecipients].([]interface{})
	if len(recipientsList) > 0 {
		recipientsMap := recipientsList[0].(map[string]interface{})
		config.Recipients = unified_alerts.Recipients{
			Emails:                  interfaceSliceToStringSlice(recipientsMap[recipientsEmails].([]interface{})),
			NotificationEndpointIds: interfaceSliceToIntSlice(recipientsMap[recipientsNotificationEndpointIds].([]interface{})),
		}
	}

	return config, nil
}

func setUnifiedAlert(d *schema.ResourceData, alert *unified_alerts.UnifiedAlert) diag.Diagnostics {
	d.Set(unifiedAlertId, alert.Id)
	d.Set(unifiedAlertTitle, alert.Title)
	d.Set(unifiedAlertType, alert.Type)
	d.Set(unifiedAlertDescription, alert.Description)
	d.Set(unifiedAlertTags, alert.Tags)
	d.Set(unifiedAlertFolderId, alert.FolderId)
	d.Set(unifiedAlertDashboardId, alert.DashboardId)
	d.Set(unifiedAlertPanelId, alert.PanelId)
	d.Set(unifiedAlertRunbook, alert.Runbook)
	d.Set(unifiedAlertEnabled, alert.Enabled)
	d.Set(unifiedAlertRca, alert.Rca)
	d.Set(unifiedAlertRcaNotificationEndpointIds, alert.RcaNotificationEndpointIds)
	d.Set(unifiedAlertUseAlertNotificationEndpointsForRca, alert.UseAlertNotificationEndpointsForRca)
	d.Set(unifiedAlertCreatedAt, alert.CreatedAt)
	d.Set(unifiedAlertUpdatedAt, alert.UpdatedAt)

	if alert.Type == unified_alerts.TypeLogAlert && alert.LogAlert != nil {
		if err := setLogAlert(d, alert.LogAlert); err != nil {
			return diag.FromErr(err)
		}
	}

	if alert.Type == unified_alerts.TypeMetricAlert && alert.MetricAlert != nil {
		if err := setMetricAlert(d, alert.MetricAlert); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func setLogAlert(d *schema.ResourceData, logAlert *unified_alerts.LogAlertConfig) error {
	logAlertMap := make(map[string]interface{})

	logAlertMap[logAlertSearchTimeFrameMinutes] = logAlert.SearchTimeFrameMinutes

	// Set output
	outputMap := map[string]interface{}{
		logAlertOutputType:                         logAlert.Output.Type,
		logAlertOutputSuppressNotificationsMinutes: logAlert.Output.SuppressNotificationsMinutes,
	}

	// Set recipients
	recipientsMap := map[string]interface{}{
		recipientsEmails:                  logAlert.Output.Recipients.Emails,
		recipientsNotificationEndpointIds: logAlert.Output.Recipients.NotificationEndpointIds,
	}
	outputMap[logAlertOutputRecipients] = []interface{}{recipientsMap}
	logAlertMap[logAlertOutput] = []interface{}{outputMap}

	// Set sub components
	subComponents := make([]interface{}, len(logAlert.SubComponents))
	for i, sc := range logAlert.SubComponents {
		scMap := make(map[string]interface{})

		// Set query definition
		queryDefMap := map[string]interface{}{
			queryDefinitionQuery:                    sc.QueryDefinition.Query,
			queryDefinitionGroupBy:                  sc.QueryDefinition.GroupBy,
			queryDefinitionShouldQueryOnAllAccounts: sc.QueryDefinition.ShouldQueryOnAllAccounts,
			queryDefinitionAccountIdsToQueryOn:      sc.QueryDefinition.AccountIdsToQueryOn,
		}

		// Set filters if present
		if len(sc.QueryDefinition.Filters.Bool.Must) > 0 || len(sc.QueryDefinition.Filters.Bool.Should) > 0 ||
			len(sc.QueryDefinition.Filters.Bool.Filter) > 0 || len(sc.QueryDefinition.Filters.Bool.MustNot) > 0 {
			filtersBytes, _ := json.Marshal(sc.QueryDefinition.Filters)
			queryDefMap[queryDefinitionFilters] = string(filtersBytes)
		}

		// Set aggregation if present
		if sc.QueryDefinition.Aggregation.AggregationType != "" {
			aggMap := map[string]interface{}{
				aggregationAggregationType:    sc.QueryDefinition.Aggregation.AggregationType,
				aggregationFieldToAggregateOn: sc.QueryDefinition.Aggregation.FieldToAggregateOn,
				aggregationValueToAggregateOn: sc.QueryDefinition.Aggregation.ValueToAggregateOn,
			}
			queryDefMap[queryDefinitionAggregation] = []interface{}{aggMap}
		}

		scMap[subComponentQueryDefinition] = []interface{}{queryDefMap}

		// Set trigger
		tiersList := make([]interface{}, 0, len(sc.Trigger.SeverityThresholdTiers))
		for severity, threshold := range sc.Trigger.SeverityThresholdTiers {
			tiersList = append(tiersList, map[string]interface{}{
				severityThresholdTierSeverity:  severity,
				severityThresholdTierThreshold: threshold,
			})
		}

		triggerMap := map[string]interface{}{
			triggerOperator:               sc.Trigger.Operator,
			triggerSeverityThresholdTiers: tiersList,
		}
		scMap[subComponentTrigger] = []interface{}{triggerMap}

		// Set output if present
		if len(sc.Output.Columns) > 0 || sc.Output.ShouldUseAllFields {
			outputMap := map[string]interface{}{
				subComponentOutputShouldUseAllFields: sc.Output.ShouldUseAllFields,
			}

			if len(sc.Output.Columns) > 0 {
				columnsList := make([]interface{}, len(sc.Output.Columns))
				for j, col := range sc.Output.Columns {
					columnsList[j] = map[string]interface{}{
						columnConfigFieldName: col.FieldName,
						columnConfigRegex:     col.Regex,
						columnConfigSort:      col.Sort,
					}
				}
				outputMap[subComponentOutputColumns] = columnsList
			}

			scMap[subComponentOutput] = []interface{}{outputMap}
		}

		subComponents[i] = scMap
	}
	logAlertMap[logAlertSubComponents] = subComponents

	// Set correlations if present
	if len(logAlert.Correlations.CorrelationOperators) > 0 || len(logAlert.Correlations.Joins) > 0 {
		correlationsMap := map[string]interface{}{
			correlationsCorrelationOperators: logAlert.Correlations.CorrelationOperators,
		}

		if len(logAlert.Correlations.Joins) > 0 {
			joinsList := make([]interface{}, len(logAlert.Correlations.Joins))
			for i, join := range logAlert.Correlations.Joins {
				joinMapInterface := make(map[string]interface{})
				for k, v := range join {
					joinMapInterface[k] = v
				}
				joinsList[i] = joinMapInterface
			}
			correlationsMap[correlationsJoins] = joinsList
		}

		logAlertMap[logAlertCorrelations] = []interface{}{correlationsMap}
	}

	// Set schedule if present
	if logAlert.Schedule.CronExpression != "" {
		scheduleMap := map[string]interface{}{
			scheduleCronExpression: logAlert.Schedule.CronExpression,
			scheduleTimezone:       logAlert.Schedule.Timezone,
		}
		logAlertMap[logAlertSchedule] = []interface{}{scheduleMap}
	}

	return d.Set(unifiedAlertLogAlert, []interface{}{logAlertMap})
}

func setMetricAlert(d *schema.ResourceData, metricAlert *unified_alerts.MetricAlertConfig) error {
	metricAlertMap := make(map[string]interface{})

	metricAlertMap[metricAlertSeverity] = metricAlert.Severity

	// Set trigger
	triggerMap := map[string]interface{}{
		metricTriggerTriggerType:            metricAlert.Trigger.TriggerType,
		metricTriggerMetricOperator:         metricAlert.Trigger.MetricOperator,
		metricTriggerMinThreshold:           metricAlert.Trigger.MinThreshold,
		metricTriggerMaxThreshold:           metricAlert.Trigger.MaxThreshold,
		metricTriggerMathExpression:         metricAlert.Trigger.MathExpression,
		metricTriggerSearchTimeFrameMinutes: metricAlert.Trigger.SearchTimeFrameMinutes,
	}
	metricAlertMap[metricAlertTrigger] = []interface{}{triggerMap}

	// Set queries
	queries := make([]interface{}, len(metricAlert.Queries))
	for i, query := range metricAlert.Queries {
		queryDefMap := map[string]interface{}{
			metricQueryDefinitionDatasourceUid: query.QueryDefinition.DatasourceUid,
			metricQueryDefinitionPromqlQuery:   query.QueryDefinition.PromqlQuery,
		}

		queries[i] = map[string]interface{}{
			metricQueryRefId:           query.RefId,
			metricQueryQueryDefinition: []interface{}{queryDefMap},
		}
	}
	metricAlertMap[metricAlertQueries] = queries

	// Set recipients
	recipientsMap := map[string]interface{}{
		recipientsEmails:                  metricAlert.Recipients.Emails,
		recipientsNotificationEndpointIds: metricAlert.Recipients.NotificationEndpointIds,
	}
	metricAlertMap[metricAlertRecipients] = []interface{}{recipientsMap}

	return d.Set(unifiedAlertMetricAlert, []interface{}{metricAlertMap})
}

// Utility functions

func interfaceSliceToStringSlice(slice []interface{}) []string {
	if slice == nil {
		return nil
	}
	result := make([]string, len(slice))
	for i, v := range slice {
		if v != nil {
			result[i] = v.(string)
		}
	}
	return result
}

func interfaceSliceToIntSlice(slice []interface{}) []int {
	if slice == nil {
		return nil
	}
	result := make([]int, len(slice))
	for i, v := range slice {
		if v != nil {
			result[i] = v.(int)
		}
	}
	return result
}
