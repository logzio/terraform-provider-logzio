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
	unifiedAlertDescription                         = "description"
	unifiedAlertTags                                = "tags"
	unifiedAlertRunbook                             = "runbook"
	unifiedAlertEnabled                             = "enabled"
	unifiedAlertRca                                 = "rca"
	unifiedAlertRcaNotificationEndpointIds          = "rca_notification_endpoint_ids"
	unifiedAlertUseAlertNotificationEndpointsForRca = "use_alert_notification_endpoints_for_rca"
	unifiedAlertCreatedAt                           = "created_at"
	unifiedAlertUpdatedAt                           = "updated_at"
	unifiedAlertCreatedBy                           = "created_by"
	unifiedAlertUpdatedBy                           = "updated_by"

	// LinkedPanel fields
	unifiedAlertLinkedPanel = "linked_panel"
	linkedPanelFolderId     = "folder_id"
	linkedPanelDashboardId  = "dashboard_id"
	linkedPanelPanelId      = "panel_id"

	// Recipients (top-level)
	unifiedAlertRecipients            = "recipients"
	recipientsEmails                  = "emails"
	recipientsNotificationEndpointIds = "notification_endpoint_ids"

	// AlertConfiguration fields
	unifiedAlertAlertConfiguration          = "alert_configuration"
	alertConfigType                         = "type"
	alertConfigSuppressNotificationsMinutes = "suppress_notifications_minutes"
	alertConfigAlertOutputTemplateType      = "alert_output_template_type"
	alertConfigSearchTimeFrameMinutes       = "search_timeframe_minutes"
	alertConfigSeverity                     = "severity"

	// Sub-component fields (log alerts, inside alert_configuration)
	alertConfigSubComponents = "sub_components"

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

	// Correlations fields (inside alert_configuration)
	alertConfigCorrelations          = "correlations"
	correlationsCorrelationOperators = "correlation_operators"
	correlationsJoins                = "joins"

	// Schedule fields (inside alert_configuration)
	alertConfigSchedule    = "schedule"
	scheduleCronExpression = "cron_expression"
	scheduleTimezone       = "timezone"

	// Metric trigger fields (inside alert_configuration)
	alertConfigTrigger           = "trigger"
	metricTriggerType            = "type"
	metricTriggerCondition       = "condition"
	metricTriggerExpression      = "expression"
	triggerConditionOperatorType = "operator_type"
	triggerConditionThreshold    = "threshold"
	triggerConditionFrom         = "from"
	triggerConditionTo           = "to"

	// Metric queries fields (inside alert_configuration)
	alertConfigQueries               = "queries"
	metricQueryRefId                 = "ref_id"
	metricQueryQueryDefinition       = "query_definition"
	metricQueryDefinitionAccountId   = "account_id"
	metricQueryDefinitionPromqlQuery = "promql_query"

	unifiedAlertRetryAttempts = 8

)

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
			StateContext: resourceUnifiedAlertImport,
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
			unifiedAlertDescription: {
				Type:     schema.TypeString,
				Optional: true,
			},
			unifiedAlertTags: {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			unifiedAlertLinkedPanel: {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						linkedPanelFolderId: {
							Type:     schema.TypeString,
							Optional: true,
						},
						linkedPanelDashboardId: {
							Type:     schema.TypeString,
							Optional: true,
						},
						linkedPanelPanelId: {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
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
			unifiedAlertRecipients: {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     resourceRecipients(),
			},
			unifiedAlertAlertConfiguration: {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     resourceAlertConfiguration(),
			},
			unifiedAlertCreatedAt: {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			unifiedAlertUpdatedAt: {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			unifiedAlertCreatedBy: {
				Type:     schema.TypeString,
				Computed: true,
			},
			unifiedAlertUpdatedBy: {
				Type:     schema.TypeString,
				Computed: true,
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
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
	}
}

func resourceAlertConfiguration() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			alertConfigType: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{unified_alerts.TypeLogAlert, unified_alerts.TypeMetricAlert}, false),
			},
			alertConfigSuppressNotificationsMinutes: {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			alertConfigAlertOutputTemplateType: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{unified_alerts.OutputTypeJson, unified_alerts.OutputTypeTable}, false),
			},
			alertConfigSearchTimeFrameMinutes: {
				Type:     schema.TypeInt,
				Optional: true,
			},
			alertConfigSeverity: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{unified_alerts.SeverityInfo, unified_alerts.SeverityLow, unified_alerts.SeverityMedium, unified_alerts.SeverityHigh, unified_alerts.SeveritySevere}, false),
			},
			alertConfigSubComponents: {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 1,
				Elem:     resourceSubComponent(),
			},
			alertConfigCorrelations: {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						correlationsCorrelationOperators: {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						correlationsJoins: {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
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
			alertConfigSchedule: {
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
			alertConfigTrigger: {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     resourceMetricAlertTrigger(),
			},
			alertConfigQueries: {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 1,
				Elem:     resourceMetricQuery(),
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
							Computed: true,
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
										ValidateFunc: validation.StringInSlice([]string{unified_alerts.AggregationTypeSum, unified_alerts.AggregationTypeMin, unified_alerts.AggregationTypeMax, unified_alerts.AggregationTypeAvg, unified_alerts.AggregationTypeCount, unified_alerts.AggregationTypeUniqueCount, unified_alerts.AggregationTypeNone, unified_alerts.AggregationTypePercentage, unified_alerts.AggregationTypePercentile}, false),
										DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
											// API converts NONE to COUNT silently
											return (old == unified_alerts.AggregationTypeCount && new == unified_alerts.AggregationTypeNone) ||
												(old == unified_alerts.AggregationTypeNone && new == unified_alerts.AggregationTypeCount)
										},
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
							Computed: true,
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
							ValidateFunc: validation.StringInSlice([]string{unified_alerts.OperatorLessThan, unified_alerts.OperatorGreaterThan, unified_alerts.OperatorLessThanOrEquals, unified_alerts.OperatorGreaterThanOrEquals, unified_alerts.OperatorEquals, unified_alerts.OperatorNotEquals}, false),
						},
						triggerSeverityThresholdTiers: {
							Type:     schema.TypeSet,
							Required: true,
							MinItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									severityThresholdTierSeverity: {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{unified_alerts.SeverityInfo, unified_alerts.SeverityLow, unified_alerts.SeverityMedium, unified_alerts.SeverityHigh, unified_alerts.SeveritySevere}, false),
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
							Default:  true,
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
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{unified_alerts.SortAsc, unified_alerts.SortDesc}, false),
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

func resourceMetricAlertTrigger() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			metricTriggerType: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{unified_alerts.TriggerTypeThreshold, unified_alerts.TriggerTypeMath}, false),
			},
			metricTriggerCondition: {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						triggerConditionOperatorType: {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{unified_alerts.OperatorTypeAbove, unified_alerts.OperatorTypeBelow, unified_alerts.OperatorTypeWithinRange, unified_alerts.OperatorTypeOutsideRange}, false),
						},
						triggerConditionThreshold: {
							Type:     schema.TypeFloat,
							Optional: true,
						},
						triggerConditionFrom: {
							Type:     schema.TypeFloat,
							Optional: true,
						},
						triggerConditionTo: {
							Type:     schema.TypeFloat,
							Optional: true,
						},
					},
				},
			},
			metricTriggerExpression: {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceMetricQuery() *schema.Resource {
	return &schema.Resource{
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
						metricQueryDefinitionAccountId: {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						metricQueryDefinitionPromqlQuery: {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceUnifiedAlertCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	createAlert, err := buildCreateUnifiedAlert(d)
	if err != nil {
		return diag.FromErr(err)
	}

	urlType := configTypeToUrlType(createAlert.AlertConfiguration.Type)

	jsonBytes, _ := json.Marshal(createAlert)
	tflog.Debug(ctx, fmt.Sprintf("Creating unified alert: %s", string(jsonBytes)))

	client := unifiedAlertClient(m)
	alert, err := client.CreateUnifiedAlert(urlType, createAlert)
	if err != nil {
		return diag.Errorf("failed to create unified alert: %v", err)
	}

	d.SetId(fmt.Sprintf("%s:%s", urlType, alert.Id))
	tflog.Info(ctx, fmt.Sprintf("Created unified alert with ID: %s, Type: %s", alert.Id, urlType))

	return resourceUnifiedAlertRead(ctx, d, m)
}

func resourceUnifiedAlertRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	urlType, alertId, err := parseUnifiedAlertCompositeId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	client := unifiedAlertClient(m)
	alert, err := client.GetUnifiedAlert(urlType, alertId)
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("Failed to get unified alert: %v", err))
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
			d.SetId("")
			return diag.Diagnostics{}
		}
		return diag.FromErr(err)
	}

	return setUnifiedAlert(d, alert)
}

func resourceUnifiedAlertUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	urlType, alertId, err := parseUnifiedAlertCompositeId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	createAlert, err := buildCreateUnifiedAlert(d)
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

	tflog.Info(ctx, fmt.Sprintf("Updated unified alert with ID: %s, Type: %s", alertId, urlType))
	return diagRet
}

func resourceUnifiedAlertDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	urlType, alertId, err := parseUnifiedAlertCompositeId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	client := unifiedAlertClient(m)
	_, err = client.DeleteUnifiedAlert(urlType, alertId)
	if err != nil {
		return diag.Errorf("failed to delete unified alert: %v", err)
	}

	tflog.Info(ctx, fmt.Sprintf("Deleted unified alert with ID: %s, Type: %s", alertId, urlType))
	return nil
}

func resourceUnifiedAlertImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	compositeId := d.Id()
	urlType, alertId, err := parseUnifiedAlertCompositeId(compositeId)
	if err != nil {
		return nil, fmt.Errorf("invalid import ID format. Expected 'alertType:alertId' (e.g., 'logs:abc-123'), got: %s", compositeId)
	}

	if urlType != unified_alerts.UrlTypeLogs && urlType != unified_alerts.UrlTypeMetrics {
		return nil, fmt.Errorf("invalid alert type '%s'. Must be '%s' or '%s'", urlType, unified_alerts.UrlTypeLogs, unified_alerts.UrlTypeMetrics)
	}

	d.SetId(compositeId)
	d.Set(unifiedAlertId, alertId)

	return []*schema.ResourceData{d}, nil
}

func parseUnifiedAlertCompositeId(id string) (urlType string, alertId string, err error) {
	parts := strings.SplitN(id, ":", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid composite ID format: %s. Expected 'alertType:alertId'", id)
	}
	return parts[0], parts[1], nil
}

func configTypeToUrlType(configType string) string {
	if configType == unified_alerts.TypeLogAlert {
		return unified_alerts.UrlTypeLogs
	}
	return unified_alerts.UrlTypeMetrics
}

// urlTypeToConfigType converts URL type to AlertConfiguration.Type.
// "logs" -> "LOG_ALERT", "metrics" -> "METRIC_ALERT"
func urlTypeToConfigType(urlType string) string {
	if urlType == unified_alerts.UrlTypeLogs {
		return unified_alerts.TypeLogAlert
	}
	return unified_alerts.TypeMetricAlert
}

func float64Ptr(v float64) *float64 {
	return &v
}

func buildCreateUnifiedAlert(d *schema.ResourceData) (unified_alerts.CreateUnifiedAlert, error) {
	enabled := d.Get(unifiedAlertEnabled).(bool)

	alert := unified_alerts.CreateUnifiedAlert{
		Title:                               d.Get(unifiedAlertTitle).(string),
		Description:                         d.Get(unifiedAlertDescription).(string),
		Tags:                                setToStringSlice(d.Get(unifiedAlertTags).(*schema.Set)),
		Runbook:                             d.Get(unifiedAlertRunbook).(string),
		Enabled:                             &enabled,
		Rca:                                 d.Get(unifiedAlertRca).(bool),
		RcaNotificationEndpointIds:          interfaceSliceToIntSlice(d.Get(unifiedAlertRcaNotificationEndpointIds).([]interface{})),
		UseAlertNotificationEndpointsForRca: d.Get(unifiedAlertUseAlertNotificationEndpointsForRca).(bool),
	}

	linkedPanelList := d.Get(unifiedAlertLinkedPanel).([]interface{})
	if len(linkedPanelList) > 0 && linkedPanelList[0] != nil {
		lpMap := linkedPanelList[0].(map[string]interface{})
		alert.LinkedPanel = &unified_alerts.LinkedPanel{
			FolderId:    lpMap[linkedPanelFolderId].(string),
			DashboardId: lpMap[linkedPanelDashboardId].(string),
			PanelId:     lpMap[linkedPanelPanelId].(string),
		}
	}

	recipientsList := d.Get(unifiedAlertRecipients).([]interface{})
	if len(recipientsList) > 0 && recipientsList[0] != nil {
		recipientsMap := recipientsList[0].(map[string]interface{})
		alert.Recipients = &unified_alerts.Recipients{
			Emails:                  interfaceSliceToStringSlice(recipientsMap[recipientsEmails].([]interface{})),
			NotificationEndpointIds: interfaceSliceToIntSlice(recipientsMap[recipientsNotificationEndpointIds].([]interface{})),
		}
	}

	alertConfig, err := buildAlertConfiguration(d)
	if err != nil {
		return alert, err
	}
	alert.AlertConfiguration = alertConfig

	return alert, nil
}

func buildAlertConfiguration(d *schema.ResourceData) (*unified_alerts.AlertConfiguration, error) {
	configList := d.Get(unifiedAlertAlertConfiguration).([]interface{})
	if len(configList) == 0 || configList[0] == nil {
		return nil, fmt.Errorf("alert_configuration is required")
	}
	configMap := configList[0].(map[string]interface{})

	config := &unified_alerts.AlertConfiguration{
		Type:                         configMap[alertConfigType].(string),
		SuppressNotificationsMinutes: configMap[alertConfigSuppressNotificationsMinutes].(int),
		AlertOutputTemplateType:      configMap[alertConfigAlertOutputTemplateType].(string),
		SearchTimeFrameMinutes:       configMap[alertConfigSearchTimeFrameMinutes].(int),
		Severity:                     configMap[alertConfigSeverity].(string),
	}

	// Build sub_components (log alerts)
	subComponentsList := configMap[alertConfigSubComponents].([]interface{})
	if len(subComponentsList) > 0 {
		subComponents, err := buildSubComponents(subComponentsList)
		if err != nil {
			return nil, err
		}
		config.SubComponents = subComponents
	}

	correlationsList := configMap[alertConfigCorrelations].([]interface{})
	if len(correlationsList) > 0 && correlationsList[0] != nil {
		config.Correlations = buildCorrelations(correlationsList)
	}

	scheduleList := configMap[alertConfigSchedule].([]interface{})
	if len(scheduleList) > 0 && scheduleList[0] != nil {
		config.Schedule = buildSchedule(scheduleList)
	}

	// Build metric trigger
	triggerList := configMap[alertConfigTrigger].([]interface{})
	if len(triggerList) > 0 && triggerList[0] != nil {
		trigger, err := buildMetricAlertTrigger(triggerList)
		if err != nil {
			return nil, err
		}
		config.Trigger = trigger
	}

	// Build metric queries
	queriesList := configMap[alertConfigQueries].([]interface{})
	if len(queriesList) > 0 {
		config.Queries = buildMetricQueries(queriesList)
	}

	return config, nil
}

func buildSubComponents(subComponentsList []interface{}) ([]unified_alerts.SubComponent, error) {
	subComponents := make([]unified_alerts.SubComponent, len(subComponentsList))
	for i, scItem := range subComponentsList {
		scMap := scItem.(map[string]interface{})
		subComp := unified_alerts.SubComponent{}

		// Parse query definition
		queryDefList := scMap[subComponentQueryDefinition].([]interface{})
		if len(queryDefList) > 0 && queryDefList[0] != nil {
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
			if len(aggregationList) > 0 && aggregationList[0] != nil {
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
		if len(triggerList) > 0 && triggerList[0] != nil {
			triggerMap := triggerList[0].(map[string]interface{})
			subComp.Trigger = unified_alerts.SubComponentTrigger{
				Operator:               triggerMap[triggerOperator].(string),
				SeverityThresholdTiers: make(map[string]float32),
			}

			tiersSet := triggerMap[triggerSeverityThresholdTiers].(*schema.Set)
			for _, tierItem := range tiersSet.List() {
				tierMap := tierItem.(map[string]interface{})
				severity := tierMap[severityThresholdTierSeverity].(string)
				threshold := float32(tierMap[severityThresholdTierThreshold].(float64))
				subComp.Trigger.SeverityThresholdTiers[severity] = threshold
			}
		}

		// Parse output if present
		outputList := scMap[subComponentOutput].([]interface{})
		if len(outputList) > 0 && outputList[0] != nil {
			outputMap := outputList[0].(map[string]interface{})
			subComp.Output = unified_alerts.SubComponentOutput{
				ShouldUseAllFields: outputMap[subComponentOutputShouldUseAllFields].(bool),
			}

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

		subComponents[i] = subComp
	}
	return subComponents, nil
}

func buildCorrelations(correlationsList []interface{}) *unified_alerts.Correlations {
	correlationsMap := correlationsList[0].(map[string]interface{})
	correlations := &unified_alerts.Correlations{
		CorrelationOperators: interfaceSliceToStringSlice(correlationsMap[correlationsCorrelationOperators].([]interface{})),
	}

	joinsList := correlationsMap[correlationsJoins].([]interface{})
	if len(joinsList) > 0 {
		correlations.Joins = make([]map[string]string, len(joinsList))
		for i, joinItem := range joinsList {
			joinMapInterface := joinItem.(map[string]interface{})
			joinMap := make(map[string]string)
			for k, v := range joinMapInterface {
				joinMap[k] = v.(string)
			}
			correlations.Joins[i] = joinMap
		}
	}

	return correlations
}

func buildSchedule(scheduleList []interface{}) *unified_alerts.Schedule {
	scheduleMap := scheduleList[0].(map[string]interface{})
	return &unified_alerts.Schedule{
		CronExpression: scheduleMap[scheduleCronExpression].(string),
		Timezone:       scheduleMap[scheduleTimezone].(string),
	}
}

func buildMetricAlertTrigger(triggerList []interface{}) (*unified_alerts.MetricAlertTrigger, error) {
	triggerMap := triggerList[0].(map[string]interface{})

	trigger := &unified_alerts.MetricAlertTrigger{
		Type:       triggerMap[metricTriggerType].(string),
		Expression: triggerMap[metricTriggerExpression].(string),
	}

	conditionList := triggerMap[metricTriggerCondition].([]interface{})
	if len(conditionList) > 0 && conditionList[0] != nil {
		condMap := conditionList[0].(map[string]interface{})
		operatorType := condMap[triggerConditionOperatorType].(string)
		condition := &unified_alerts.TriggerCondition{
			OperatorType: operatorType,
		}

		// Set pointer fields based on operator type to avoid sending zero-value pointers
		// for irrelevant fields (e.g., from/to for "above"/"below" operators)
		switch operatorType {
		case unified_alerts.OperatorTypeAbove, unified_alerts.OperatorTypeBelow:
			threshold := condMap[triggerConditionThreshold].(float64)
			condition.Threshold = float64Ptr(threshold)
		case unified_alerts.OperatorTypeWithinRange, unified_alerts.OperatorTypeOutsideRange:
			from := condMap[triggerConditionFrom].(float64)
			to := condMap[triggerConditionTo].(float64)
			condition.From = float64Ptr(from)
			condition.To = float64Ptr(to)
		}

		trigger.Condition = condition
	}

	return trigger, nil
}

func buildMetricQueries(queriesList []interface{}) []unified_alerts.MetricQuery {
	queries := make([]unified_alerts.MetricQuery, len(queriesList))
	for i, queryItem := range queriesList {
		queryMap := queryItem.(map[string]interface{})
		metricQuery := unified_alerts.MetricQuery{
			RefId: queryMap[metricQueryRefId].(string),
		}

		queryDefList := queryMap[metricQueryQueryDefinition].([]interface{})
		if len(queryDefList) > 0 && queryDefList[0] != nil {
			queryDefMap := queryDefList[0].(map[string]interface{})
			metricQuery.QueryDefinition = unified_alerts.MetricQueryDefinition{
				AccountId:   int32(queryDefMap[metricQueryDefinitionAccountId].(int)),
				PromqlQuery: queryDefMap[metricQueryDefinitionPromqlQuery].(string),
			}
		}

		queries[i] = metricQuery
	}
	return queries
}

func setUnifiedAlert(d *schema.ResourceData, alert *unified_alerts.UnifiedAlert) diag.Diagnostics {
	d.Set(unifiedAlertId, alert.Id)
	d.Set(unifiedAlertTitle, alert.Title)
	d.Set(unifiedAlertDescription, alert.Description)
	d.Set(unifiedAlertTags, alert.Tags)
	d.Set(unifiedAlertRunbook, alert.Runbook)
	d.Set(unifiedAlertEnabled, alert.Enabled)
	d.Set(unifiedAlertRca, alert.Rca)
	d.Set(unifiedAlertRcaNotificationEndpointIds, alert.RcaNotificationEndpointIds)
	d.Set(unifiedAlertUseAlertNotificationEndpointsForRca, alert.UseAlertNotificationEndpointsForRca)
	d.Set(unifiedAlertCreatedAt, alert.CreatedAt)
	d.Set(unifiedAlertUpdatedAt, alert.UpdatedAt)
	d.Set(unifiedAlertCreatedBy, alert.CreatedBy)
	d.Set(unifiedAlertUpdatedBy, alert.UpdatedBy)

	// Set LinkedPanel
	if alert.LinkedPanel != nil {
		d.Set(unifiedAlertLinkedPanel, []interface{}{
			map[string]interface{}{
				linkedPanelFolderId:    alert.LinkedPanel.FolderId,
				linkedPanelDashboardId: alert.LinkedPanel.DashboardId,
				linkedPanelPanelId:     alert.LinkedPanel.PanelId,
			},
		})
	}

	// Set Recipients
	if alert.Recipients != nil {
		recipientsMap := map[string]interface{}{}
		if len(alert.Recipients.Emails) > 0 {
			recipientsMap[recipientsEmails] = alert.Recipients.Emails
		}
		if len(alert.Recipients.NotificationEndpointIds) > 0 {
			recipientsMap[recipientsNotificationEndpointIds] = alert.Recipients.NotificationEndpointIds
		}
		d.Set(unifiedAlertRecipients, []interface{}{recipientsMap})
	}

	// Set AlertConfiguration
	if alert.AlertConfiguration != nil {
		if err := setAlertConfiguration(d, alert.AlertConfiguration); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func setAlertConfiguration(d *schema.ResourceData, config *unified_alerts.AlertConfiguration) error {
	configMap := make(map[string]interface{})

	configMap[alertConfigType] = config.Type
	configMap[alertConfigSuppressNotificationsMinutes] = config.SuppressNotificationsMinutes
	configMap[alertConfigAlertOutputTemplateType] = config.AlertOutputTemplateType
	configMap[alertConfigSearchTimeFrameMinutes] = config.SearchTimeFrameMinutes
	configMap[alertConfigSeverity] = config.Severity

	// Set sub_components
	if len(config.SubComponents) > 0 {
		configMap[alertConfigSubComponents] = flattenSubComponents(config.SubComponents)
	}

	// Set correlations
	if config.Correlations != nil {
		configMap[alertConfigCorrelations] = flattenCorrelations(config.Correlations)
	}

	// Set schedule - only if user configured one (API returns default {timezone:"UTC"} even when not set)
	if config.Schedule != nil && config.Schedule.CronExpression != "" {
		configMap[alertConfigSchedule] = flattenSchedule(config.Schedule)
	}

	// Set metric trigger
	if config.Trigger != nil {
		configMap[alertConfigTrigger] = flattenMetricAlertTrigger(config.Trigger)
	}

	// Set metric queries
	if len(config.Queries) > 0 {
		configMap[alertConfigQueries] = flattenMetricQueries(config.Queries)
	}

	return d.Set(unifiedAlertAlertConfiguration, []interface{}{configMap})
}

func flattenSubComponents(subComponents []unified_alerts.SubComponent) []interface{} {
	result := make([]interface{}, len(subComponents))
	for i, sc := range subComponents {
		scMap := make(map[string]interface{})

		// Set query definition
		queryDefMap := map[string]interface{}{
			queryDefinitionQuery:                    sc.QueryDefinition.Query,
			queryDefinitionShouldQueryOnAllAccounts: sc.QueryDefinition.ShouldQueryOnAllAccounts,
		}
		if len(sc.QueryDefinition.GroupBy) > 0 {
			queryDefMap[queryDefinitionGroupBy] = sc.QueryDefinition.GroupBy
		}
		if len(sc.QueryDefinition.AccountIdsToQueryOn) > 0 {
			queryDefMap[queryDefinitionAccountIdsToQueryOn] = sc.QueryDefinition.AccountIdsToQueryOn
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
				severityThresholdTierThreshold: float64(threshold),
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
					sortVal := col.Sort
					if sortVal == "" {
						sortVal = "ASC"
					}

					colMap := map[string]interface{}{
						columnConfigFieldName: col.FieldName,
						columnConfigSort:      sortVal,
					}
					if col.Regex != "" {
						colMap[columnConfigRegex] = col.Regex
					}
					columnsList[j] = colMap
				}
				outputMap[subComponentOutputColumns] = columnsList
			}

			scMap[subComponentOutput] = []interface{}{outputMap}
		}

		result[i] = scMap
	}
	return result
}

func flattenCorrelations(correlations *unified_alerts.Correlations) []interface{} {
	correlationsMap := map[string]interface{}{
		correlationsCorrelationOperators: correlations.CorrelationOperators,
	}

	if len(correlations.Joins) > 0 {
		joinsList := make([]interface{}, len(correlations.Joins))
		for i, join := range correlations.Joins {
			joinMapInterface := make(map[string]interface{})
			for k, v := range join {
				joinMapInterface[k] = v
			}
			joinsList[i] = joinMapInterface
		}
		correlationsMap[correlationsJoins] = joinsList
	}

	return []interface{}{correlationsMap}
}

func flattenSchedule(schedule *unified_alerts.Schedule) []interface{} {
	return []interface{}{
		map[string]interface{}{
			scheduleCronExpression: schedule.CronExpression,
			scheduleTimezone:       schedule.Timezone,
		},
	}
}

func flattenMetricAlertTrigger(trigger *unified_alerts.MetricAlertTrigger) []interface{} {
	triggerMap := map[string]interface{}{
		metricTriggerType:       trigger.Type,
		metricTriggerExpression: trigger.Expression,
	}

	if trigger.Condition != nil {
		condMap := map[string]interface{}{
			triggerConditionOperatorType: trigger.Condition.OperatorType,
		}
		if trigger.Condition.Threshold != nil {
			condMap[triggerConditionThreshold] = *trigger.Condition.Threshold
		}
		if trigger.Condition.From != nil {
			condMap[triggerConditionFrom] = *trigger.Condition.From
		}
		if trigger.Condition.To != nil {
			condMap[triggerConditionTo] = *trigger.Condition.To
		}
		triggerMap[metricTriggerCondition] = []interface{}{condMap}
	}

	return []interface{}{triggerMap}
}

func flattenMetricQueries(queries []unified_alerts.MetricQuery) []interface{} {
	result := make([]interface{}, len(queries))
	for i, query := range queries {
		queryDefMap := map[string]interface{}{
			metricQueryDefinitionAccountId:   int(query.QueryDefinition.AccountId),
			metricQueryDefinitionPromqlQuery: query.QueryDefinition.PromqlQuery,
		}

		result[i] = map[string]interface{}{
			metricQueryRefId:           query.RefId,
			metricQueryQueryDefinition: []interface{}{queryDefMap},
		}
	}
	return result
}

// Utility functions

func setToStringSlice(set *schema.Set) []string {
	list := set.List()
	result := make([]string, len(list))
	for i, v := range list {
		result[i] = v.(string)
	}
	return result
}

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
