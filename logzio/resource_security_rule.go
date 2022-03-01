package logzio

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/logzio/logzio_terraform_client/alerts_v2"
	"github.com/logzio/logzio_terraform_client/security_rules"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
	"strconv"
	"strings"
	"time"
)

const (
	securityRuleTitle                       = "title"
	securityRuleDescription                 = "description"
	securityRuleTags                        = "tags"
	SecurityRuleEmails                      = "notification_emails"
	securityRuleNotificationEndpoints       = "rule_notification_endpoints"
	securityRuleSuppressNotificationMinutes = "suppress_notifications_minutes"
	securityRuleOutputType                  = "output_type"
	securityRuleSearchTimeFrameMinutes      = "search_timeframe_minutes"
	securityRuleSubComponents               = "sub_components"
	securityRuleQueryString                 = "query_string"
	securityRuleFilterMust                  = "filter_must"
	securityRuleFilterMustNot               = "filter_must_not"
	securityRuleGroupBy                     = "group_by_aggregation_fields"
	securityRuleAggregationType             = "value_aggregation_type"
	securityRuleAggregationField            = "value_aggregation_field"
	securityRuleShouldQueryOnAllAccounts    = "should_query_on_all_accounts"
	securityRuleAccountIdsToQuery           = "account_ids_to_query_on"
	securityRuleOperation                   = "operation"
	securityRuleSeverityThresholdTiers      = "severity_threshold_tiers"
	securityRuleSeverity                    = "severity"
	securityRuleThreshold                   = "threshold"
	securityRuleColumns                     = "columns"
	securityRuleColumnsFieldName            = "field_name"
	securityRuleColumnsRegex                = "regex"
	securityRuleColumnSort                  = "sort"
	securityRuleCorrelationOperator         = "correlation_operator"
	securityRuleJoins                       = "joins"
	securityRuleIsEnabled                   = "is_enabled"
	securityRuleCreatedAt                   = "created_at"
	securityRuleCreatedBy                   = "created_by"
	securityRuleUpdatedAt                   = "updated_at"
	securityRuleUpdatedBy                   = "updated_by"

	securityRuleGroupByMaxItems int = 3
)

func resourceSecurityRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceEndpointCreate,
		//Read:   resourceEndpointRead,
		//Update: resourceEndpointUpdate,
		Delete: resourceEndpointDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			securityRuleTitle: {
				Type:     schema.TypeString,
				Required: true,
			},
			securityRuleDescription: {
				Type:     schema.TypeString,
				Optional: true,
			},
			securityRuleTags: {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			SecurityRuleEmails: {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			securityRuleNotificationEndpoints: {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			securityRuleSuppressNotificationMinutes: {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5,
			},
			securityRuleOutputType: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: utils.ValidateOutputType,
			},
			securityRuleSearchTimeFrameMinutes: {
				Type:     schema.TypeInt,
				Required: true,
			},
			securityRuleSubComponents: {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						securityRuleQueryString: {
							Type:     schema.TypeString,
							Required: true,
						},
						securityRuleFilterMust: {
							Type:     schema.TypeString,
							Optional: true,
						},
						securityRuleFilterMustNot: {
							Type:     schema.TypeString,
							Optional: true,
						},
						securityRuleGroupBy: {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: securityRuleGroupByMaxItems,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						securityRuleAggregationType: {
							Type:     schema.TypeString,
							Required: true,
						},
						securityRuleAggregationField: {
							Type:     schema.TypeString,
							Optional: true,
						},
						securityRuleShouldQueryOnAllAccounts: {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						securityRuleAccountIdsToQuery: {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						securityRuleOperation: {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: utils.ValidateOperationV2,
						},
						securityRuleSeverityThresholdTiers: {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									securityRuleSeverity: {
										Type:     schema.TypeString,
										Required: true,
									},
									securityRuleThreshold: {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},
						securityRuleColumns: {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									securityRuleColumnsFieldName: {
										Type:     schema.TypeString,
										Optional: true,
									},
									securityRuleColumnsRegex: {
										Type:     schema.TypeString,
										Optional: true,
									},
									securityRuleColumnSort: {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: utils.ValidateSortTypes,
									},
								},
							},
						},
					},
				},
			},
			securityRuleCorrelationOperator: {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: setSecurityRulesCorrelationDefault,
			},
			securityRuleJoins: {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
			securityRuleIsEnabled: {
				Type:     schema.TypeBool,
				Optional: true,
			},
			securityRuleCreatedAt: {
				Type:     schema.TypeString,
				Computed: true,
			},
			securityRuleCreatedBy: {
				Type:     schema.TypeString,
				Computed: true,
			},
			securityRuleUpdatedAt: {
				Type:     schema.TypeString,
				Computed: true,
			},
			securityRuleUpdatedBy: {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Second),
			Update: schema.DefaultTimeout(5 * time.Second),
			Delete: schema.DefaultTimeout(5 * time.Second),
		},
	}
}

func resourceSecurityRuleCreate(d *schema.ResourceData, m interface{}) error {
	createRule := createCreateUpdateSecurityRule(d)
}

func createCreateUpdateSecurityRule(d *schema.ResourceData) security_rules.CreateUpdateSecurityRule {
	mappedFlatComponents := getVariousSecurityRuleFields(d)
	tags := utils.GetTags(d, securityRuleTags)
	subComponentsFromConfig := d.Get(securityRuleSubComponents).([]interface{})
	subComponents := getSecurityRuleSubComponents(subComponentsFromConfig)

}

func getSecurityRuleSubComponents(subComponentsFromConfig []interface{}) []security_rules.SubRule {
	var retArray []security_rules.SubRule

	for _, subComponentElement := range subComponentsFromConfig {
		var subAlertElement security_rules.SubRule
		element := subComponentElement.(map[string]interface{})

		subAlertElement.QueryDefinition.Query = element[alertV2QueryString].(string)
		subAlertElement.QueryDefinition.Aggregation.AggregationType = element[alertV2AggregationType].(string)
		subAlertElement.QueryDefinition.Aggregation.FieldToAggregateOn = element[alertV2AggregationField].(string)
		subAlertElement.QueryDefinition.ShouldQueryOnAllAccounts = element[alertV2ShouldQueryOnAllAccounts].(bool)
		subAlertElement.Trigger.Operator = element[alertV2Operation].(string)

		if _, ok := element[alertV2FilterMust]; ok {
			mustToAppend := utils.ParseStringToMapList(element[alertV2FilterMust].(string))
			subAlertElement.QueryDefinition.Filters.Bool.Must = mustToAppend
		}

		if _, ok := element[alertV2FilterMustNot]; ok {
			mustNotToAppend := utils.ParseStringToMapList(element[alertV2FilterMustNot].(string))
			subAlertElement.QueryDefinition.Filters.Bool.MustNot = mustNotToAppend
		}

		if _, ok := element[alertV2GroupBy]; ok {
			groupByInterface := element[alertV2GroupBy].([]interface{})
			for _, gb := range groupByInterface {
				subAlertElement.QueryDefinition.GroupBy = append(subAlertElement.QueryDefinition.GroupBy, gb.(string))
			}
		}

		if _, ok := element[alertV2AccountIdsToQuery]; ok {
			idsInterface := element[alertV2AccountIdsToQuery].([]interface{})
			for _, i := range idsInterface {
				subAlertElement.QueryDefinition.AccountIdsToQueryOn = append(subAlertElement.QueryDefinition.AccountIdsToQueryOn, i.(int))
			}
		}

		var columnsCreateAlert []alerts_v2.ColumnConfig
		if _, ok := element[alertV2Columns]; ok {
			columns := element[alertV2Columns].([]interface{})
			for _, columnElement := range columns {
				column := columnElement.(map[string]interface{})
				var columnCreateAlert alerts_v2.ColumnConfig
				if _, ok := column[alertV2ColumnsFieldName]; ok {
					columnCreateAlert.FieldName = column[alertV2ColumnsFieldName].(string)
				}

				if _, ok := column[alertV2ColumnsRegex]; ok {
					columnCreateAlert.Regex = column[alertV2ColumnsRegex].(string)
				}

				if _, ok := column[alertV2ColumnSort]; ok {
					columnCreateAlert.Sort = column[alertV2ColumnSort].(string)
				}

				columnsCreateAlert = append(columnsCreateAlert, columnCreateAlert)
			}

			subAlertElement.Output.Columns = columnsCreateAlert
		}

		tiers := element[alertV2SeverityThresholdTiers].(*schema.Set).List()

		subAlertElement.Trigger.SeverityThresholdTiers = make(map[string]float32)
		for _, t := range tiers {
			tier := t.(map[string]interface{})
			subAlertElement.Trigger.SeverityThresholdTiers[tier[alertV2Severity].(string)] = float32(tier[alertV2Threshold].(int))
		}

		retArray = append(retArray, subAlertElement)
	}

	return retArray
}

func getVariousSecurityRuleFields(d *schema.ResourceData) map[string]interface{} {
	correlationsString := d.Get(securityRuleCorrelationOperator).(string)
	correlations := strings.Split(correlationsString, ",")

	var joins []map[string]string
	joinsInterface := d.Get(securityRuleJoins).([]interface{})
	for _, j := range joinsInterface {
		joins = append(joins, j.(map[string]string))
	}

	mappedComponents := map[string]interface{}{
		securityRuleTitle:                       d.Get(securityRuleTitle).(string),
		securityRuleDescription:                 d.Get(securityRuleDescription).(string),
		securityRuleSearchTimeFrameMinutes:      d.Get(securityRuleSearchTimeFrameMinutes).(int),
		securityRuleIsEnabled:                   strconv.FormatBool(d.Get(securityRuleIsEnabled).(bool)),
		securityRuleSuppressNotificationMinutes: d.Get(securityRuleSuppressNotificationMinutes).(int),
		securityRuleOutputType:                  d.Get(securityRuleOutputType).(string),
		securityRuleCorrelationOperator:         correlations,
		securityRuleJoins:                       joins,
	}

	return mappedComponents
}

func setSecurityRulesCorrelationDefault() (interface{}, error) {
	operators := []string{security_rules.SecurityRuleCorrelationOperatorAnd}
	correlationsOperators := strings.Join(operators, ",")

	return correlationsOperators, nil
}
