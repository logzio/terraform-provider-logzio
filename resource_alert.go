package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jonboydell/logzio_client"
	"log"
	"strconv"
)

const BASE_10 int = 10
const BITSIZE_64 int = 64

const title string = "title"
const query_string string = "query_string"

func resourceAlert() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlertCreate,
		Read:   resourceAlertRead,
		Update: resourceAlertUpdate,
		Delete: resourceAlertDelete,

		Schema: map[string]*schema.Schema{
			title: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			query_string: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"is_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"filter": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"operation": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateOperation,
			},
			"search_timeframe_minutes": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"notification_emails": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"value_aggregation_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"value_aggregation_field": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  nil,
			},
			"suppress_notification_minutes": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"group_by_aggregation_fields": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Default: nil,
			},
			"alert_notification_endpoints": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"severity_threshold_tiers": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"severity": {
							Type:     schema.TypeString,
							Required: true,
						},
						"threshold": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
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

	title := d.Get(title).(string)
	description := d.Get("description").(string)
	queryString := d.Get(query_string).(string)
	filter := d.Get("filter").(string)
	operation := d.Get("operation").(string)
	searchTimeFrameMinutes := d.Get("search_timeframe_minutes").(int)
	suppressNotificationMinutes := d.Get("suppress_notification_minutes").(int)
	notificationEmails := d.Get("notification_emails").([]interface{})
	valueAggregationType := d.Get("value_aggregation_type").(string)

	valueAggregationField, e := d.GetOk("value_aggregation_field")
	if e {
		valueAggregationField = d.Get("value_aggregation_field").(string)
	} else {
		valueAggregationField = nil
	}

	alertNotificationEndpoints := d.Get("alert_notification_endpoints").([]interface{})

	var isEnabled bool = true
	_, e = d.GetOk("is_enabled")
	if e {
		isEnabled = d.Get("is_enabled").(bool)
	}

	tiers := d.Get("severity_threshold_tiers").([]interface{})
	severityThresholdTiers := []logzio_client.SeverityThresholdType{}

	for x := 0; x < len(tiers); x++ {
		tier := tiers[x].(map[string]interface{})
		thresholdTier := logzio_client.SeverityThresholdType{
			Severity:  tier["severity"].(string),
			Threshold: tier["threshold"].(int),
		}
		severityThresholdTiers = append(severityThresholdTiers, thresholdTier)
	}

	createAlert := logzio_client.CreateAlertType{
		Title:                       title,
		Description:                 description,
		QueryString:                 queryString,
		Filter:                      filter,
		Operation:                   operation,
		SeverityThresholdTiers:      severityThresholdTiers,
		SearchTimeFrameMinutes:      searchTimeFrameMinutes,
		NotificationEmails:          notificationEmails,
		IsEnabled:                   isEnabled,
		SuppressNotificationMinutes: suppressNotificationMinutes,
		ValueAggregationType:        valueAggregationType,
		ValueAggregationField:       valueAggregationField,
		AlertNotificationEndpoints:  alertNotificationEndpoints,
	}

	_, e = d.GetOk("group_by_aggregation_fields")
	if e {
		createAlert.GroupByAggregationFields = d.Get("group_by_aggregation_fields").([]interface{})
	} else {
		createAlert.GroupByAggregationFields = nil
	}

	jsonBytes, err := json.Marshal(createAlert)
	jsonStr, _ := prettyprint(jsonBytes)
	log.Printf("%s::%s", "resourceAlertCreate", jsonStr)

	api_token := m.(Config).api_token

	var client *logzio_client.Client
	client = logzio_client.New(api_token)

	a, err := client.CreateAlert(createAlert)

	if err != nil {
		ferr := err.(logzio_client.FieldError)
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

	var client *logzio_client.Client
	client = logzio_client.New(api_token)

	alertId, _ := strconv.ParseInt(d.Id(), BASE_10, BITSIZE_64)

	var alert *logzio_client.AlertType
	var err error
	alert, err = client.GetAlert(alertId)
	if err != nil {
		return err
	}
	d.Set(query_string, alert.QueryString)
	d.Set(title, alert.Title)
	d.Set("description", alert.Description)
	d.Set("filter", alert.Filter)
	d.Set("operation", alert.Operation)
	d.Set("search_timeframe_minutes", alert.SearchTimeFrameMinutes)
	d.Set("notification_emails", alert.NotificationEmails)
	d.Set("value_aggregation_type", alert.ValueAggregationType)
	d.Set("value_aggregation_field", alert.ValueAggregationField)
	d.Set("group_by_aggregation_fields", alert.GroupByAggregationFields)
	d.Set("alert_notification_endpoints", alert.AlertNotificationEndpoints)
	d.Set("suppress_notification_minutes", alert.SuppressNotificationMinutes)
	d.Set("created_at", alert.CreatedAt)
	d.Set("created_by", alert.CreatedBy)
	d.Set("last_triggered_at", alert.LastTriggeredAt)

	return nil
}

func resourceAlertUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("%s", "resourceAlertUpdate")

	alertId, _ := strconv.ParseInt(d.Id(), BASE_10, BITSIZE_64)

	title := d.Get(title).(string)
	description := d.Get("description").(string)
	queryString := d.Get(query_string).(string)
	filter := d.Get("filter").(string)
	operation := d.Get("operation").(string)
	searchTimeFrameMinutes := d.Get("search_timeframe_minutes").(int)
	suppressNotificationMinutes := d.Get("suppress_notification_minutes").(int)
	notificationEmails := d.Get("notification_emails").([]interface{})
	valueAggregationType := d.Get("value_aggregation_type").(string)

	valueAggregationField, e := d.GetOk("value_aggregation_field")
	if e {
		valueAggregationField = d.Get("value_aggregation_field").(string)
	} else {
		valueAggregationField = nil
	}

	alertNotificationEndpoints := d.Get("alert_notification_endpoints").([]interface{})

	var isEnabled bool = true
	_, e = d.GetOk("is_enabled")
	if e {
		isEnabled = d.Get("is_enabled").(bool)
	}

	tiers := d.Get("severity_threshold_tiers").([]interface{})
	severityThresholdTiers := []logzio_client.SeverityThresholdType{}

	for x := 0; x < len(tiers); x++ {
		tier := tiers[x].(map[string]interface{})
		thresholdTier := logzio_client.SeverityThresholdType{
			Severity:  tier["severity"].(string),
			Threshold: tier["threshold"].(int),
		}
		severityThresholdTiers = append(severityThresholdTiers, thresholdTier)
	}

	updateAlert := logzio_client.CreateAlertType{
		Title:                       title,
		Description:                 description,
		QueryString:                 queryString,
		Filter:                      filter,
		Operation:                   operation,
		SeverityThresholdTiers:      severityThresholdTiers,
		SearchTimeFrameMinutes:      searchTimeFrameMinutes,
		NotificationEmails:          notificationEmails,
		IsEnabled:                   isEnabled,
		SuppressNotificationMinutes: suppressNotificationMinutes,
		ValueAggregationType:        valueAggregationType,
		ValueAggregationField:       valueAggregationField,
		AlertNotificationEndpoints:  alertNotificationEndpoints,
	}

	_, e = d.GetOk("group_by_aggregation_fields")
	if e {
		updateAlert.GroupByAggregationFields = d.Get("group_by_aggregation_fields").([]interface{})
	} else {
		updateAlert.GroupByAggregationFields = nil
	}

	jsonBytes, err := json.Marshal(updateAlert)
	jsonStr, _ := prettyprint(jsonBytes)
	log.Printf("%s::%s", "resourceAlertCreate", jsonStr)

	api_token := m.(Config).api_token

	var client *logzio_client.Client
	client = logzio_client.New(api_token)

	_, err = client.UpdateAlert(alertId, updateAlert)

	if err != nil {
		ferr := err.(logzio_client.FieldError)
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

	var client *logzio_client.Client
	client = logzio_client.New(api_token)

	alertId, _ := strconv.ParseInt(d.Id(), BASE_10, BITSIZE_64)
	err := client.DeleteAlert(alertId)
	return err
}
