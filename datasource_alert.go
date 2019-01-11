package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jonboydell/logzio_client/alerts"
)

func dataSourceAlert() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlertRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"title": {
				Type: schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceAlertRead(d *schema.ResourceData, m interface{}) error {
	api_token := m.(Config).api_token

	var client *alerts.Alerts
	client, _ = alerts.New(api_token)

	alertId, ok := d.GetOk("id")
	if ok {
		alert, err := client.GetAlert(alertId.(int64))
		if err != nil {
			return err
		}
		d.SetId(fmt.Sprintf("%d", alertId.(int64)))
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

	alertTitle, ok := d.GetOk("title")
	if ok {
		list, err := client.ListAlerts()
		if err != nil {
			return err
		}
		for i := 0; i < len(list); i++ {
			alert := list[i]
			if alert.Title == alertTitle {
				d.SetId(fmt.Sprintf("%d", alert.AlertId))
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
		}
	}
	return fmt.Errorf("Couldn't find alert with specified attributes")
}