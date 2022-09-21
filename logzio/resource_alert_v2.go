package logzio

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/logzio/logzio_terraform_client/alerts_v2"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	alertV2Id                          string = "id"
	alertV2Title                       string = "title"
	alertV2Description                 string = "description"
	alertV2Tags                        string = "tags"
	alertV2SearchTimeFrameMinutes      string = "search_timeframe_minutes"
	alertV2IsEnabled                   string = "is_enabled"
	alertV2NotificationEmails          string = "notification_emails"
	alertV2NotificationEndpoints       string = "alert_notification_endpoints"
	alertV2SuppressNotificationMinutes string = "suppress_notifications_minutes"
	alertV2OutputType                  string = "output_type"
	alertV2QueryString                 string = "query_string"
	alertV2FilterMust                  string = "filter_must"
	alertV2FilterMustNot               string = "filter_must_not"
	alertV2GroupBy                     string = "group_by_aggregation_fields"
	alertV2AggregationType             string = "value_aggregation_type"
	alertV2AggregationField            string = "value_aggregation_field"
	alertV2ShouldQueryOnAllAccounts    string = "should_query_on_all_accounts"
	alertV2AccountIdsToQuery           string = "account_ids_to_query_on"
	alertV2Operation                   string = "operation"
	alertV2SeverityThresholdTiers      string = "severity_threshold_tiers"
	alertV2Severity                    string = "severity"
	alertV2Threshold                   string = "threshold"
	alertV2Columns                     string = "columns"
	alertV2ColumnsFieldName            string = "field_name"
	alertV2ColumnsRegex                string = "regex"
	alertV2ColumnSort                  string = "sort"
	alertV2SubComponents               string = "sub_components"
	alertV2CorrelationOperator         string = "correlation_operator"
	alertV2Joins                       string = "joins"

	alertV2CreatedAt string = "created_at"
	alertV2CreatedBy string = "created_by"
	alertV2UpdatedAt string = "updated_at"
	alertV2UpdatedBy string = "updated_by"

	groupByMaxItems int = 3

	delayGetAlertV2 = 1 * time.Second
)

// alertV2Client returns the alert v2 client with the api token from the provider
func alertV2Client(m interface{}) *alerts_v2.AlertsV2Client {
	var client *alerts_v2.AlertsV2Client
	client, _ = alerts_v2.New(m.(Config).apiToken, m.(Config).baseUrl)
	return client
}

func resourceAlertV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAlertV2Create,
		ReadContext:   resourceAlertV2Read,
		UpdateContext: resourceAlertV2Update,
		DeleteContext: resourceAlertV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			alertV2Title: {
				Type:     schema.TypeString,
				Required: true,
			},
			alertV2Description: {
				Type:     schema.TypeString,
				Optional: true,
			},
			alertV2Tags: {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			alertV2SearchTimeFrameMinutes: {
				Type:     schema.TypeInt,
				Required: true,
			},
			alertV2IsEnabled: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			alertV2NotificationEmails: {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			alertV2NotificationEndpoints: {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			alertV2SuppressNotificationMinutes: {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5,
			},
			alertV2OutputType: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: utils.ValidateOutputType,
			},
			alertV2SubComponents: {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						alertV2QueryString: {
							Type:     schema.TypeString,
							Required: true,
						},
						alertV2FilterMust: {
							Type:     schema.TypeString,
							Optional: true,
						},
						alertV2FilterMustNot: {
							Type:     schema.TypeString,
							Optional: true,
						},
						alertV2GroupBy: {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: groupByMaxItems,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						alertV2AggregationType: {
							Type:     schema.TypeString,
							Required: true,
						},
						alertV2AggregationField: {
							Type:     schema.TypeString,
							Optional: true,
						},
						alertV2ShouldQueryOnAllAccounts: {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						alertV2AccountIdsToQuery: {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						alertV2Operation: {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: utils.ValidateOperationV2,
						},
						alertV2SeverityThresholdTiers: {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									alertV2Severity: {
										Type:     schema.TypeString,
										Required: true,
									},
									alertV2Threshold: {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},
						alertV2Columns: {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									alertV2ColumnsFieldName: {
										Type:     schema.TypeString,
										Optional: true,
									},
									alertV2ColumnsRegex: {
										Type:     schema.TypeString,
										Optional: true,
									},
									alertV2ColumnSort: {
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
			alertV2CorrelationOperator: {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: setCorrelationDefault,
			},
			alertV2Joins: {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
			alertV2CreatedAt: {
				Type:     schema.TypeString,
				Computed: true,
			},
			alertV2CreatedBy: {
				Type:     schema.TypeString,
				Computed: true,
			},
			alertV2UpdatedAt: {
				Type:     schema.TypeString,
				Computed: true,
			},
			alertV2UpdatedBy: {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// resourceAlertV2Create creates a new alert (v2) in logzio
func resourceAlertV2Create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	createAlert := createCreateAlertType(d)

	jsonBytes, err := json.Marshal(createAlert)
	tflog.Debug(ctx, fmt.Sprintf("%s::%s", "resourceAlertCreate", string(jsonBytes)))
	client := alertV2Client(m)
	a, err := client.CreateAlert(createAlert)

	if err != nil {
		switch typedError := err.(type) {
		case alerts_v2.FieldError:
			if typedError.Field == "valueAggregationTypeComposite" {
				return diag.Errorf("if valueAggregationType is set to None, valueAggregationField and groupByAggregationFields must not be set")
			}
		default:
			return diag.Errorf("resourceAlertCreate failed: %v", typedError)
		}

		return diag.FromErr(err)
	}

	alertId := strconv.FormatInt(a.AlertId, utils.BASE_10)
	d.SetId(alertId)

	return resourceAlertV2Read(ctx, d, m)
}

// resourceAlertV2Read reads an alert (v2) from logzio
func resourceAlertV2Read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var alert *alerts_v2.AlertType
	var err error
	alertId, _ := utils.IdFromResourceData(d)
	client := alertV2Client(m)

	readErr := retry.Do(
		func() error {
			alert, err = client.GetAlert(alertId)
			if err != nil {
				return err
			}

			return nil
		},
		retry.RetryIf(
			func(err error) bool {
				if err != nil {
					if strings.Contains(err.Error(), "missing alert") {
						return true
					}
				}
				return false
			}),
		retry.DelayType(retry.BackOffDelay),
		retry.Attempts(15),
	)

	if readErr != nil {
		// If we were not able to find the resource - delete from state
		d.SetId("")
		return diag.FromErr(readErr)
	}

	setValuesAlertV2(d, alert)
	setCreatedUpdatedFields(d, alert)

	return nil
}

// resourceAlertV2Update updates an existing alert in logzio, returns an error if it doesn't exist
func resourceAlertV2Update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	alertId, _ := utils.IdFromResourceData(d)
	updateAlert := createCreateAlertType(d)

	jsonBytes, err := json.Marshal(updateAlert)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonStr, _ := printFormatted(jsonBytes)
	tflog.Debug(ctx, fmt.Sprintf("%s::%s", "resourceAlertCreate", jsonStr))

	client := alertV2Client(m)
	_, err = client.UpdateAlert(alertId, updateAlert)

	if err != nil {
		if strings.Contains(err.Error(), "valueAggregationTypeComposite") {
			return diag.Errorf("if valueAggregationType is set to None, valueAggregationField and groupByAggregationFields must not be set")
		}

		return diag.FromErr(err)
	}

	var diagRet diag.Diagnostics
	readErr := retry.Do(
		func() error {
			diagRet = resourceAlertV2Read(ctx, d, m)
			if diagRet.HasError() {
				return fmt.Errorf("received error from read alert v2")
			}

			return nil
		},
		retry.RetryIf(
			func(err error) bool {
				if err != nil {
					return true
				} else {
					// Check if the update shows on read
					// if not updated yet - retry
					createAlert := createCreateAlertType(d)
					return !reflect.DeepEqual(createAlert, updateAlert)
				}
			}),
		retry.DelayType(retry.BackOffDelay),
		retry.Attempts(15),
	)

	if readErr != nil {
		tflog.Error(ctx, "could not update schema")
		return diagRet
	}

	return nil
}

// resourceAlertV2Delete deletes an existing alert in logzio, returns an error if it doesn't exist
func resourceAlertV2Delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	alertId, err := utils.IdFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	deleteErr := retry.Do(
		func() error {
			return alertV2Client(m).DeleteAlert(alertId)
		},
		retry.RetryIf(
			func(err error) bool {
				return err != nil
			}),
		retry.DelayType(retry.BackOffDelay),
		retry.Attempts(15),
	)

	if deleteErr != nil {
		return diag.FromErr(deleteErr)
	}

	d.SetId("")
	return nil
}

func getSubComponentMapping(sc []alerts_v2.SubAlert) []map[string]interface{} {
	var subComponentsMapping []map[string]interface{}
	for _, subComponent := range sc {
		var columns []map[string]string
		var severityThreshold []interface{}
		for _, column := range subComponent.Output.Columns {
			columnMapping := map[string]string{
				alertV2ColumnsFieldName: column.FieldName,
				alertV2ColumnsRegex:     column.Regex,
				alertV2ColumnSort:       column.Sort,
			}

			columns = append(columns, columnMapping)
		}

		for key, val := range subComponent.Trigger.SeverityThresholdTiers {
			severityElement := map[string]interface{}{alertV2Severity: key, alertV2Threshold: val}
			severityThreshold = append(severityThreshold, severityElement)
		}

		mapping := map[string]interface{}{
			alertV2QueryString:              subComponent.QueryDefinition.Query,
			alertV2FilterMust:               utils.ParseObjectToString(subComponent.QueryDefinition.Filters.Bool.Must),
			alertV2FilterMustNot:            utils.ParseObjectToString(subComponent.QueryDefinition.Filters.Bool.MustNot),
			alertV2GroupBy:                  subComponent.QueryDefinition.GroupBy,
			alertV2AggregationField:         subComponent.QueryDefinition.Aggregation.FieldToAggregateOn,
			alertV2AggregationType:          subComponent.QueryDefinition.Aggregation.AggregationType,
			alertV2ShouldQueryOnAllAccounts: subComponent.QueryDefinition.ShouldQueryOnAllAccounts,
			alertV2AccountIdsToQuery:        subComponent.QueryDefinition.AccountIdsToQueryOn,
			alertV2Operation:                subComponent.Trigger.Operator,
			alertV2SeverityThresholdTiers:   severityThreshold,
			alertV2Columns:                  columns,
		}

		subComponentsMapping = append(subComponentsMapping, mapping)
	}

	return subComponentsMapping
}

func printFormatted(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	return out.Bytes(), err
}

func getRecipients(emails []interface{}, endpointIds []interface{}) *alerts_v2.AlertRecipients {
	if emails == nil && endpointIds == nil {
		return nil
	}

	var emailsArrayString []string
	for _, email := range emails {
		emailsArrayString = append(emailsArrayString, email.(string))
	}

	var endpointsArrayString []int
	for _, endpoint := range endpointIds {
		endpointsArrayString = append(endpointsArrayString, endpoint.(int))
	}

	recipients := alerts_v2.AlertRecipients{
		Emails:                  emailsArrayString,
		NotificationEndpointIds: endpointsArrayString,
	}

	return &recipients
}

func getSubComponents(subComponentsFromConfig []interface{}) []alerts_v2.SubAlert {
	var retArray []alerts_v2.SubAlert

	for _, subComponentElement := range subComponentsFromConfig {
		var subAlertElement alerts_v2.SubAlert
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

func getVariousFields(d *schema.ResourceData) map[string]interface{} {
	correlationsString := d.Get(alertV2CorrelationOperator).(string)
	correlations := strings.Split(correlationsString, ",")

	var joins []map[string]string
	joinsInterface := d.Get(alertV2Joins).([]interface{})
	for _, j := range joinsInterface {
		joins = append(joins, j.(map[string]string))
	}

	mappedComponents := map[string]interface{}{
		alertV2Title:                       d.Get(alertV2Title).(string),
		alertV2Description:                 d.Get(alertV2Description).(string),
		alertV2SearchTimeFrameMinutes:      d.Get(alertV2SearchTimeFrameMinutes).(int),
		alertV2IsEnabled:                   strconv.FormatBool(d.Get(alertV2IsEnabled).(bool)),
		alertV2SuppressNotificationMinutes: d.Get(alertV2SuppressNotificationMinutes).(int),
		alertV2OutputType:                  d.Get(alertV2OutputType).(string),
		alertV2CorrelationOperator:         correlations,
		alertV2Joins:                       joins,
	}

	return mappedComponents
}

func getTags(d *schema.ResourceData) []string {
	var tags []string
	if alertTags, ok := d.GetOk(alertV2Tags); ok {
		for _, tag := range alertTags.(*schema.Set).List() {
			tags = append(tags, tag.(string))
		}
	}

	return tags
}

func createCreateAlertType(d *schema.ResourceData) alerts_v2.CreateAlertType {
	mappedFlatComponents := getVariousFields(d)
	tags := getTags(d)
	subComponentsFromConfig := d.Get(alertV2SubComponents).([]interface{})
	subComponents := getSubComponents(subComponentsFromConfig)
	emails := d.Get(alertV2NotificationEmails).(*schema.Set).List()
	endpoints := d.Get(alertV2NotificationEndpoints).(*schema.Set).List()
	recipients := getRecipients(emails, endpoints)

	alertOutput := alerts_v2.AlertOutput{
		Recipients:                   *recipients,
		SuppressNotificationsMinutes: mappedFlatComponents[alertV2SuppressNotificationMinutes].(int),
		Type:                         mappedFlatComponents[alertV2OutputType].(string),
	}

	correlations := alerts_v2.SubAlertCorrelation{
		CorrelationOperators: mappedFlatComponents[alertV2CorrelationOperator].([]string),
		Joins:                mappedFlatComponents[alertV2Joins].([]map[string]string),
	}

	createAlert := alerts_v2.CreateAlertType{
		Title:                  mappedFlatComponents[alertV2Title].(string),
		Description:            mappedFlatComponents[alertV2Description].(string),
		Tags:                   tags,
		SearchTimeFrameMinutes: mappedFlatComponents[alertV2SearchTimeFrameMinutes].(int),
		Enabled:                mappedFlatComponents[alertV2IsEnabled].(string),
		Output:                 alertOutput,
		SubComponents:          subComponents,
		Correlations:           correlations,
	}

	return createAlert
}

func setValuesAlertV2(d *schema.ResourceData, alert *alerts_v2.AlertType) {
	d.Set(alertV2Title, alert.Title)
	d.Set(alertV2Description, alert.Description)
	d.Set(alertV2Tags, alert.Tags)
	d.Set(alertV2SearchTimeFrameMinutes, alert.SearchTimeFrameMinutes)
	d.Set(alertV2IsEnabled, alert.Enabled)
	d.Set(alertV2NotificationEmails, alert.Output.Recipients.Emails)
	d.Set(alertV2NotificationEndpoints, alert.Output.Recipients.NotificationEndpointIds)
	d.Set(alertV2SuppressNotificationMinutes, alert.Output.SuppressNotificationsMinutes)
	d.Set(alertV2OutputType, alert.Output.Type)
	d.Set(alertV2Joins, alert.Correlations.Joins)

	correlationString := strings.Join(alert.Correlations.CorrelationOperators, ",")
	d.Set(alertV2CorrelationOperator, correlationString)

	subComponentsMapping := getSubComponentMapping(alert.SubComponents)

	d.Set(alertV2SubComponents, subComponentsMapping)
}

func setCreatedUpdatedFields(d *schema.ResourceData, alert *alerts_v2.AlertType) {
	d.Set(alertV2CreatedAt, alert.CreatedAt)
	d.Set(alertV2CreatedBy, alert.CreatedBy)
	d.Set(alertV2UpdatedAt, alert.UpdatedAt)
	d.Set(alertV2UpdatedBy, alert.UpdatedBy)
}

func setCorrelationDefault() (interface{}, error) {
	operators := []string{alerts_v2.CorrelationOperatorAnd}
	correlationsOperators := strings.Join(operators, ",")

	return correlationsOperators, nil
}
