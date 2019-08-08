package logzio

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jonboydell/logzio_client/alerts"
)

func dataSourceAlert() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlertRead,
		Schema: map[string]*schema.Schema{
			alertId: {
				Type:     schema.TypeInt,
				Optional: true,
			},
			alert_title: {
				Type:     schema.TypeString,
				Optional: true,
			},
			alertNotificationEndpoints: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			alertDescription: {
				Type:     schema.TypeString,
				Computed: true,
			},
			alertFilter: {
				Type:     schema.TypeString,
				Computed: true,
			},
			alert_group_by_aggregation_fields: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			alert_is_enabled: {
				Type:     schema.TypeBool,
				Computed: true,
			},
			alert_query_string: {
				Type:     schema.TypeString,
				Computed: true,
			},
			alert_notification_emails: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			alert_operation: {
				Type:     schema.TypeString,
				Computed: true,
			},
			alert_search_timeframe_minutes: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			alert_severity_threshold_tiers: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						alert_severity: {
							Type:     schema.TypeString,
							Computed: true,
						},
						alert_threshold: {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			alert_suppress_notifications_minutes: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			alert_value_aggregation_field: {
				Type:     schema.TypeString,
				Computed: true,
			},
			alert_value_aggregation_type: {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAlertRead(d *schema.ResourceData, m interface{}) error {
	var client *alerts.AlertsClient
	client, _ = alerts.New(m.(Config).apiToken, m.(Config).baseUrl)
	alertIdString, ok := d.GetOk(alertId);
	if ok {
		id := int64(alertIdString.(int))
		alert, err := client.GetAlert(id)
		if err != nil {
			return err
		}
		d.SetId(fmt.Sprintf("%d", id))
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
		return nil
	}

	alertTitle, ok := d.GetOk(alert_title)
	if ok {
		list, err := client.ListAlerts()
		if err != nil {
			return err
		}
		for i := 0; i < len(list); i++ {
			alert := list[i]
			if alert.Title == alertTitle {
				d.SetId(fmt.Sprintf("%d", alert.AlertId))
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
				return nil
			}
		}
	}
	return fmt.Errorf("couldn't find alert with specified attributes")
}
