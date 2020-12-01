Provides a Logz.io alert resource. This can be used to create and manage Logz.io log monitoring alerts. [Learn more](https://docs.logz.io/user-guide/alerts/)

| Argument | Description| 
|	id | Logz.io alert ID. |
|	alert_notification_endpoints | | 
|	created_at | Date and time in UTC when the alert was first created. | 
|	created_by | Email of the user who first created the alert.| 
|	description | A description of the event, its significance, and suggested next steps or instructions for the team. | 
|	filter | | 
|	tags | Tags for filtering alerts and triggered alerts. Can be used in Kibana Discover, dashboards, and more. | 
|	group_by_aggregation_fields | | 
|	is_enabled | boolean field. If `true`, the alert is currently active.| 
|	query_string | | 
|	last_triggered_at | | 
|	last_updated | Date and time in UTC when the alert was last updated. | 
|	notification_emails | | 
|	operation | | 
|	search_timeframe_minutes | The time frame for evaluating the log data is a sliding window, with 1 minute granularity.
          
          
          The recommended minimum and maximum values are not validated, but needed to guarantee the alert's accuracy.


          The minimum recommended time frame is 5 minutes, as anything shorter will be less reliable and unnecessarily resource-heavy.


          The maximum recommended time frame is 1440 minutes (24 hours). The alert runs on the index from today and yesterday (in UTC) and the maximum time frame increases throughout the day, reaching 48 hours exactly before midnight UTC. | 
|	severity | | 
|	severity_threshold_tiers | | 
|	suppress_notifications_minutes | | 
|	threshold | | 
|	title | Alert title. | 
|	value_aggregation_field | | 
	value_aggregation_type | | 