package logzio

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jonboydell/logzio_client/alerts"
)

const (
	alertId                              string = "id"
	alertNotificationEndpoints           string = "alert_notification_endpoints"
	alertCreatedAt                       string = "created_at"
	alertCreatedBy                       string = "created_by"
	alertDescription                     string = "description"
	alertFilter                          string = "filter"
	alert_group_by_aggregation_fields    string = "group_by_aggregation_fields"
	alert_is_enabled                     string = "is_enabled"
	alert_query_string                   string = "query_string"
	alert_last_triggered_at              string = "last_triggered_at"
	alert_last_updated                   string = "last_updated"
	alert_notification_emails            string = "notification_emails"
	alert_operation                      string = "operation"
	alert_search_timeframe_minutes       string = "search_timeframe_minutes"
	alert_severity                       string = "severity"
	alert_severity_threshold_tiers       string = "severity_threshold_tiers"
	alert_suppress_notifications_minutes string = "suppress_notifications_minutes"
	alert_threshold                      string = "threshold"
	alert_title                          string = "title"
	alert_value_aggregation_field        string = "value_aggregation_field"
	alert_value_aggregation_type         string = "value_aggregation_type"
)

/**
 * returns the alert client with the api token from the provider
 */
func alertClient(m interface{}) *alerts.AlertsClient {
	var client *alerts.AlertsClient
	client, _ = alerts.New(m.(Config).apiToken, m.(Config).baseUrl)
	return client
}

func resourceAlert() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlertCreate,
		Read:   resourceAlertRead,
		Update: resourceAlertUpdate,
		Delete: resourceAlertDelete,

		Schema: map[string]*schema.Schema{
			alertNotificationEndpoints: {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			alertDescription: {
				Type:     schema.TypeString,
				Optional: true,
			},
			alertFilter: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "{\"bool\":{\"must\":[], \"must_not\":[]}}",
			},
			alert_group_by_aggregation_fields: {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Default: nil,
			},
			alert_is_enabled: {
				Type:     schema.TypeBool,
				Optional: true,
			},
			alert_query_string: {
				Type:     schema.TypeString,
				Required: true,
			},
			alert_notification_emails: {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			alert_operation: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateOperation,
			},
			alert_search_timeframe_minutes: {
				Type:     schema.TypeInt,
				Required: true,
			},
			alert_severity_threshold_tiers: {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						alert_severity: {
							Type:     schema.TypeString,
							Required: true,
						},
						alert_threshold: {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			alert_suppress_notifications_minutes: {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5,
			},
			alert_title: {
				Type:     schema.TypeString,
				Required: true,
			},
			alert_value_aggregation_field: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  nil,
			},
			alert_value_aggregation_type: {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func prettyprint(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	return out.Bytes(), err
}

/**
 * creates a new alert in logzio
 */
func resourceAlertCreate(d *schema.ResourceData, m interface{}) error {
	alertNotificationEndpoints := d.Get(alertNotificationEndpoints).([]interface{})
	description := d.Get(alertDescription).(string)
	filter := d.Get(alertFilter).(string)

	var isEnabled bool = true
	_, e := d.GetOk(alert_is_enabled)
	if e {
		isEnabled = d.Get(alert_is_enabled).(bool)
	}

	notificationEmails := d.Get(alert_notification_emails).([]interface{})
	operation := d.Get(alert_operation).(string)
	queryString := d.Get(alert_query_string).(string)
	searchTimeFrameMinutes := d.Get(alert_search_timeframe_minutes).(int)

	tiers := d.Get(alert_severity_threshold_tiers).([]interface{})
	severityThresholdTiers := []alerts.SeverityThresholdType{}
	for t := 0; t < len(tiers); t++ {
		tier := tiers[t].(map[string]interface{})
		thresholdTier := alerts.SeverityThresholdType{
			Severity:  tier[alert_severity].(string),
			Threshold: tier[alert_threshold].(int),
		}
		severityThresholdTiers = append(severityThresholdTiers, thresholdTier)
	}

	suppressNotificationsMinutes := d.Get(alert_suppress_notifications_minutes).(int)

	title := d.Get(alert_title).(string)

	valueAggregationField, f := d.GetOk(alert_value_aggregation_field)
	if f {
		valueAggregationField = d.Get(alert_value_aggregation_field).(string)
	} else {
		valueAggregationField = nil
	}

	valueAggregationType := d.Get(alert_value_aggregation_type).(string)

	createAlert := alerts.CreateAlertType{
		AlertNotificationEndpoints:   alertNotificationEndpoints,
		Description:                  description,
		Filter:                       filter,
		IsEnabled:                    isEnabled,
		NotificationEmails:           notificationEmails,
		Operation:                    operation,
		QueryString:                  queryString,
		SearchTimeFrameMinutes:       searchTimeFrameMinutes,
		SeverityThresholdTiers:       severityThresholdTiers,
		Title:                        title,
		SuppressNotificationsMinutes: suppressNotificationsMinutes,
		ValueAggregationField:        valueAggregationField,
		ValueAggregationType:         valueAggregationType,
	}

	_, g := d.GetOk(alert_group_by_aggregation_fields)
	if g {
		createAlert.GroupByAggregationFields = d.Get(alert_group_by_aggregation_fields).([]interface{})
	} else {
		createAlert.GroupByAggregationFields = make([]interface{}, 0)
	}

	jsonBytes, err := json.Marshal(createAlert)
	jsonStr, _ := prettyprint(jsonBytes)
	log.Printf("%s::%s", "resourceAlertCreate", jsonStr)

	client := alertClient(m)
	a, err := client.CreateAlert(createAlert)

	if err != nil {
		switch typedError := err.(type) {
		case alerts.FieldError:
			if typedError.Field == "valueAggregationTypeComposite" {
				return fmt.Errorf("if valueAggregationType is set to None, valueAggregationField and groupByAggregationFields must not be set")
			}
		default:
			return fmt.Errorf("resourceAlertCreate failed: %v", typedError)
		}
		return err
	}

	alertId := strconv.FormatInt(a.AlertId, BASE_10)
	d.SetId(alertId)

	return resourceAlertRead(d, m)
}

/**
 * reads an endpoint from logzio
 */
func resourceAlertRead(d *schema.ResourceData, m interface{}) error {
	alertId, _ := idFromResourceData(d)
	client := alertClient(m)

	var alert *alerts.AlertType
	var err error
	alert, err = client.GetAlert(alertId)
	if err != nil {
		return err
	}
	d.Set(alertNotificationEndpoints, alert.AlertNotificationEndpoints)
	d.Set(alertCreatedAt, alert.CreatedAt)
	d.Set(alertCreatedBy, alert.CreatedBy)
	d.Set(alertDescription, alert.Description)
	d.Set(alertFilter, alert.Filter)
	d.Set(alert_group_by_aggregation_fields, alert.GroupByAggregationFields)
	d.Set(alert_last_triggered_at, alert.LastTriggeredAt)
	d.Set(alert_last_updated, alert.LastUpdated)
	d.Set(alert_notification_emails, alert.NotificationEmails)
	d.Set(alert_operation, alert.Operation)
	d.Set(alert_query_string, alert.QueryString)
	d.Set(alert_title, alert.Title)
	d.Set(alert_search_timeframe_minutes, alert.SearchTimeFrameMinutes)
	d.Set(alert_suppress_notifications_minutes, alert.SuppressNotificationsMinutes)
	d.Set(alert_value_aggregation_field, alert.ValueAggregationField)
	d.Set(alert_value_aggregation_type, alert.ValueAggregationType)

	if len(alert.SeverityThresholdTiers) > 0 {
		var sttList []map[string]interface{}
		for _, stt := range alert.SeverityThresholdTiers {
			mapping := map[string]interface{}{
				alert_severity:  stt.Severity,
				alert_threshold: stt.Threshold,
			}
			sttList = append(sttList, mapping)
		}

		d.Set(alert_severity_threshold_tiers, sttList)
	}

	return nil
}

/**
 * updates an existing alert in logzio, returns an error if it doesn't exist
 */
func resourceAlertUpdate(d *schema.ResourceData, m interface{}) error {
	alertId, _ := idFromResourceData(d)

	alertNotificationEndpoints := d.Get(alertNotificationEndpoints).([]interface{})
	description := d.Get(alertDescription).(string)
	filter := d.Get(alertFilter).(string)

	var isEnabled = true
	_, e := d.GetOk(alert_is_enabled)
	if e {
		isEnabled = d.Get(alert_is_enabled).(bool)
	}

	notificationEmails := d.Get(alert_notification_emails).([]interface{})
	operation := d.Get(alert_operation).(string)
	queryString := d.Get(alert_query_string).(string)
	title := d.Get(alert_title).(string)
	searchTimeFrameMinutes := d.Get(alert_search_timeframe_minutes).(int)

	tiers := d.Get(alert_severity_threshold_tiers).([]interface{})
	severityThresholdTiers := []alerts.SeverityThresholdType{}
	for t := 0; t < len(tiers); t++ {
		tier := tiers[t].(map[string]interface{})
		thresholdTier := alerts.SeverityThresholdType{
			Severity:  tier[alert_severity].(string),
			Threshold: tier[alert_threshold].(int),
		}
		severityThresholdTiers = append(severityThresholdTiers, thresholdTier)
	}

	suppressNotificationsMinutes := d.Get(alert_suppress_notifications_minutes).(int)

	valueAggregationField, v := d.GetOk(alert_value_aggregation_field)
	if v {
		valueAggregationField = d.Get(alert_value_aggregation_field).(string)
	} else {
		valueAggregationField = nil
	}

	valueAggregationType := d.Get(alert_value_aggregation_type).(string)

	updateAlert := alerts.CreateAlertType{
		AlertNotificationEndpoints:   alertNotificationEndpoints,
		Description:                  description,
		Filter:                       filter,
		IsEnabled:                    isEnabled,
		NotificationEmails:           notificationEmails,
		Operation:                    operation,
		QueryString:                  queryString,
		Title:                        title,
		SearchTimeFrameMinutes:       searchTimeFrameMinutes,
		SeverityThresholdTiers:       severityThresholdTiers,
		SuppressNotificationsMinutes: suppressNotificationsMinutes,
		ValueAggregationField:        valueAggregationField,
		ValueAggregationType:         valueAggregationType,
	}

	_, g := d.GetOk(alert_group_by_aggregation_fields)
	if g {
		updateAlert.GroupByAggregationFields = d.Get(alert_group_by_aggregation_fields).([]interface{})
	} else {
		updateAlert.GroupByAggregationFields = nil
	}

	jsonBytes, err := json.Marshal(updateAlert)
	jsonStr, _ := prettyprint(jsonBytes)
	log.Printf("%s::%s", "resourceAlertCreate", jsonStr)

	client := alertClient(m)
	_, err = client.UpdateAlert(alertId, updateAlert)

	if err != nil {
		ferr := err.(alerts.FieldError)
		if ferr.Field == "valueAggregationTypeComposite" {
			return fmt.Errorf("if valueAggregationType is set to None, valueAggregationField and groupByAggregationFields must not be set")
		}
		return err
	}

	return resourceAlertRead(d, m)
}

/**
deletes an existing alert in logzio, returns an error if it doesn't exist
*/
func resourceAlertDelete(d *schema.ResourceData, m interface{}) error {
	client := alertClient(m)
	alertId, _ := idFromResourceData(d)
	err := client.DeleteAlert(alertId)
	return err
}
