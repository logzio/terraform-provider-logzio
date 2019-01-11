package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jonboydell/logzio_client/alerts"
	"log"
	"strconv"
)

const BASE_10 int = 10
const BITSIZE_64 int = 64

const (
	alert_notification_endpoints   string = "alert_notification_endpoints"
	created_at                     string = "created_at"
	created_by                     string = "created_by"
	description                    string = "description"
	filter                         string = "filter"
	group_by_aggregation_fields    string = "group_by_aggregation_fields"
	is_enabled                     string = "is_enabled"
	query_string                   string = "query_string"
	last_triggered_at              string = "last_triggered_at"
	last_updated                   string = "last_updated"
	notification_emails            string = "notification_emails"
	operation                      string = "operation"
	search_timeframe_minutes       string = "search_timeframe_minutes"
	severity                       string = "severity"
	severity_threshold_tiers       string = "severity_threshold_tiers"
	suppress_notifications_minutes string = "suppress_notifications_minutes"
	threshold                      string = "threshold"
	title                          string = "title"
	value_aggregation_field        string = "value_aggregation_field"
	value_aggregation_type         string = "value_aggregation_type"
)

func resourceAlert() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlertCreate,
		Read:   resourceAlertRead,
		Update: resourceAlertUpdate,
		Delete: resourceAlertDelete,

		Schema: map[string]*schema.Schema{
			alert_notification_endpoints: &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			description: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			filter: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "{\"bool\":{\"must\":[], \"must_not\":[]}}",
			},
			group_by_aggregation_fields: &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Default: nil,
			},
			is_enabled: &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			query_string: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			notification_emails: &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			operation: &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateOperation,
			},
			search_timeframe_minutes: &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			severity_threshold_tiers: &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						severity: {
							Type:     schema.TypeString,
							Required: true,
						},
						threshold: {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			suppress_notifications_minutes: &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5,
			},
			title: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			value_aggregation_field: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  nil,
			},
			value_aggregation_type: &schema.Schema{
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

func resourceAlertCreate(d *schema.ResourceData, m interface{}) error {
	alertNotificationEndpoints := d.Get(alert_notification_endpoints).([]interface{})
	description := d.Get(description).(string)
	filter := d.Get(filter).(string)

	var isEnabled bool = true
	_, e := d.GetOk(is_enabled)
	if e {
		isEnabled = d.Get(is_enabled).(bool)
	}

	notificationEmails := d.Get(notification_emails).([]interface{})
	operation := d.Get(operation).(string)
	queryString := d.Get(query_string).(string)
	searchTimeFrameMinutes := d.Get(search_timeframe_minutes).(int)

	tiers := d.Get(severity_threshold_tiers).([]interface{})
	severityThresholdTiers := []alerts.SeverityThresholdType{}
	for t := 0; t < len(tiers); t++ {
		tier := tiers[t].(map[string]interface{})
		thresholdTier := alerts.SeverityThresholdType{
			Severity:  tier[severity].(string),
			Threshold: tier[threshold].(int),
		}
		severityThresholdTiers = append(severityThresholdTiers, thresholdTier)
	}

	suppressNotificationsMinutes := d.Get(suppress_notifications_minutes).(int)

	title := d.Get(title).(string)

	valueAggregationField, f := d.GetOk(value_aggregation_field)
	if f {
		valueAggregationField = d.Get(value_aggregation_field).(string)
	} else {
		valueAggregationField = nil
	}

	valueAggregationType := d.Get(value_aggregation_type).(string)

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

	_, g := d.GetOk(group_by_aggregation_fields)
	if g {
		createAlert.GroupByAggregationFields = d.Get(group_by_aggregation_fields).([]interface{})
	} else {
		createAlert.GroupByAggregationFields = nil
	}

	jsonBytes, err := json.Marshal(createAlert)
	jsonStr, _ := prettyprint(jsonBytes)
	log.Printf("%s::%s", "resourceAlertCreate", jsonStr)

	api_token := m.(Config).api_token

	var client *alerts.Alerts
	client, err = alerts.New(api_token)
	a, err := client.CreateAlert(createAlert)

	if err != nil {
		ferr := err.(alerts.FieldError)
		if ferr.Field == "valueAggregationTypeComposite" {
			return fmt.Errorf("if valueAggregationType is set to None, valueAggregationField and groupByAggregationFields must not be set")
		}
		return err
	}

	alertId := strconv.FormatInt(a.AlertId, BASE_10)
	d.SetId(alertId)

	return resourceAlertRead(d, m)
}

func resourceAlertRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("%s::%s", "resourceAlertRead", d.Id())
	api_token := m.(Config).api_token

	var client *alerts.Alerts
	client, _ = alerts.New(api_token)

	alertId, _ := strconv.ParseInt(d.Id(), BASE_10, BITSIZE_64)

	var alert *alerts.AlertType
	var err error
	alert, err = client.GetAlert(alertId)
	if err != nil {
		return err
	}
	d.Set(alert_notification_endpoints, alert.AlertNotificationEndpoints)
	d.Set(created_at, alert.CreatedAt)
	d.Set(created_by, alert.CreatedBy)
	d.Set(description, alert.Description)
	d.Set(filter, alert.Filter)
	d.Set(group_by_aggregation_fields, alert.GroupByAggregationFields)
	d.Set(last_triggered_at, alert.LastTriggeredAt)
	d.Set(last_updated, alert.LastUpdated)
	d.Set(notification_emails, alert.NotificationEmails)
	d.Set(operation, alert.Operation)
	d.Set(query_string, alert.QueryString)
	d.Set(title, alert.Title)
	d.Set(search_timeframe_minutes, alert.SearchTimeFrameMinutes)
	d.Set(suppress_notifications_minutes, alert.SuppressNotificationsMinutes)
	d.Set(value_aggregation_field, alert.ValueAggregationField)
	d.Set(value_aggregation_type, alert.ValueAggregationType)

	return nil
}

func resourceAlertUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("%s", "resourceAlertUpdate")

	alertId, _ := strconv.ParseInt(d.Id(), BASE_10, BITSIZE_64)

	alertNotificationEndpoints := d.Get(alert_notification_endpoints).([]interface{})
	description := d.Get(description).(string)
	filter := d.Get(filter).(string)

	var isEnabled bool = true
	_, e := d.GetOk(is_enabled)
	if e {
		isEnabled = d.Get(is_enabled).(bool)
	}

	notificationEmails := d.Get(notification_emails).([]interface{})
	operation := d.Get(operation).(string)
	queryString := d.Get(query_string).(string)
	title := d.Get(title).(string)
	searchTimeFrameMinutes := d.Get(search_timeframe_minutes).(int)

	tiers := d.Get(severity_threshold_tiers).([]interface{})
	severityThresholdTiers := []alerts.SeverityThresholdType{}
	for t := 0; t < len(tiers); t++ {
		tier := tiers[t].(map[string]interface{})
		thresholdTier := alerts.SeverityThresholdType{
			Severity:  tier[severity].(string),
			Threshold: tier[threshold].(int),
		}
		severityThresholdTiers = append(severityThresholdTiers, thresholdTier)
	}

	suppressNotificationsMinutes := d.Get(suppress_notifications_minutes).(int)

	valueAggregationField, v := d.GetOk(value_aggregation_field)
	if v {
		valueAggregationField = d.Get(value_aggregation_field).(string)
	} else {
		valueAggregationField = nil
	}

	valueAggregationType := d.Get(value_aggregation_type).(string)

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

	_, g := d.GetOk(group_by_aggregation_fields)
	if g {
		updateAlert.GroupByAggregationFields = d.Get(group_by_aggregation_fields).([]interface{})
	} else {
		updateAlert.GroupByAggregationFields = nil
	}

	jsonBytes, err := json.Marshal(updateAlert)
	jsonStr, _ := prettyprint(jsonBytes)
	log.Printf("%s::%s", "resourceAlertCreate", jsonStr)

	api_token := m.(Config).api_token

	var client *alerts.Alerts
	client, _ = alerts.New(api_token)

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

func resourceAlertDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("%s::%s", "resourceAlertDelete", d.Id())
	api_token := m.(Config).api_token

	var client *alerts.Alerts
	client, _ = alerts.New(api_token)

	alertId, _ := strconv.ParseInt(d.Id(), BASE_10, BITSIZE_64)
	err := client.DeleteAlert(alertId)
	return err
}
