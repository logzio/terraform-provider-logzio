# Unified Alert Resource

Provides a Logz.io unified alert resource. This resource allows you to create and manage both log-based and metric-based alerts through a single unified API.

**Note:** This is a POC (Proof of Concept) endpoint for the unified alerts API.

The unified alert resource models either a log alert or a metric alert:
- When `type = "LOG_ALERT"`, configure the `log_alert` block and do not set `metric_alert`.
- When `type = "METRIC_ALERT"`, configure the `metric_alert` block and do not set `log_alert`.

## Example Usage - Log Alert (full)

```hcl
resource "logzio_unified_alert" "log_alert_example" {
  title       = "High error rate in checkout service"
  type        = "LOG_ALERT"
  description = "Triggers when the error rate of the checkout service exceeds the defined threshold."
  tags        = ["environment:production", "service:checkout"]
  enabled     = true

  # Optional: Link to dashboards/runbooks
  folder_id    = "unified-folder-uid"
  dashboard_id = "unified-dashboard-uid"
  panel_id     = "A"
  runbook      = "If this alert fires, check checkout pods and logs, verify recent deployments, and roll back if necessary."

  # Optional: RCA configuration
  rca                                      = true
  rca_notification_endpoint_ids            = [101, 102]
  use_alert_notification_endpoints_for_rca = true

  log_alert {
    search_timeframe_minutes = 15

    output {
      type                           = "JSON"
      suppress_notifications_minutes = 30

      recipients {
        emails                    = ["devops@company.com", "oncall@company.com"]
        notification_endpoint_ids = [11, 12]
      }
    }

    sub_components {
      query_definition {
        query = "kubernetes.container_name:checkout AND level:error"

        # Optional boolean filters in JSON format
        # See the Filters example under Argument Reference below
        filters = jsonencode({
          bool = {
            must     = []
            should   = []
            filter   = []
            must_not = []
          }
        })

        group_by = ["kubernetes.pod_name"]

        aggregation {
          aggregation_type        = "SUM"
          field_to_aggregate_on   = "error_count"
          # value_to_aggregate_on is optional
        }

        should_query_on_all_accounts = false
        account_ids_to_query_on      = [12345]
      }

      trigger {
        operator = "GREATER_THAN"

        severity_threshold_tiers {
          severity  = "HIGH"
          threshold = 100
        }

        severity_threshold_tiers {
          severity  = "MEDIUM"
          threshold = 50
        }
      }

      output {
        should_use_all_fields = false

        columns {
          field_name = "kubernetes.pod_name"
          sort       = "DESC"
        }
      }
    }

    correlations {
      correlation_operators = ["AND"]
      # joins = [] # Optional
    }

    schedule {
      cron_expression = "*/1 * * * *"
      timezone        = "UTC"
    }
  }
}
```

## Example Usage - Metric Alert (threshold)

```hcl
resource "logzio_unified_alert" "metric_alert_example" {
  title       = "High 5xx rate (absolute)"
  type        = "METRIC_ALERT"
  description = "Fire when 5xx requests exceed 5 req/min over 5 minutes."
  tags        = ["environment:production", "service:checkout"]
  enabled     = true

  # Optional: Link to dashboards/runbooks
  folder_id    = "unified-folder-uid"
  dashboard_id = "unified-dashboard-uid"
  panel_id     = "A"
  runbook      = "RCA: inspect ingress errors by pod and compare to last deploy. Check DB health. Propose rollback if sustained."

  metric_alert {
    severity = "INFO"

    trigger {
      trigger_type            = "THRESHOLD"
      metric_operator         = "ABOVE"
      min_threshold           = 5
      search_timeframe_minutes = 5
    }

    queries {
      ref_id = "A"

      query_definition {
        datasource_uid = "prometheus"
        promql_query   = "sum(rate(http_requests_total{status=~\"5..\"}[5m]))"
      }
    }

    recipients {
      emails                    = ["devops@company.com"]
      notification_endpoint_ids = [11, 12]
    }
  }
}
```

## Example Usage - Metric Alert with Math Expression

```hcl
resource "logzio_unified_alert" "metric_math_alert" {
  title       = "5xx error rate percentage is high"
  type        = "METRIC_ALERT"
  description = "Fire when 5xx responses exceed 2% of total requests over 5 minutes."
  tags        = ["environment:production", "service:checkout"]
  enabled     = true

  metric_alert {
    severity = "INFO"

    trigger {
      trigger_type             = "MATH"
      math_expression          = "(A / B) * 100"
      metric_operator          = "ABOVE"
      min_threshold            = 2
      search_timeframe_minutes = 5
    }

    queries {
      ref_id = "A"
      query_definition {
        datasource_uid = "prometheus"
        promql_query   = "sum(rate(http_requests_total{status=~\"5..\"}[5m]))"
      }
    }

    queries {
      ref_id = "B"
      query_definition {
        datasource_uid = "prometheus"
        promql_query   = "sum(rate(http_requests_total[5m]))"
      }
    }

    recipients {
      emails                    = ["devops@company.com"]
      notification_endpoint_ids = [11, 12]
    }
  }
}
```

## Argument Reference

### Required Arguments

* `title` - (String) Alert name displayed in UI and notifications.
* `type` - (String) Alert type. Must be either `LOG_ALERT` or `METRIC_ALERT`.

### Optional Arguments

* `description` - (String) Explanation of alert purpose and firing conditions.
* `tags` - (List of String) Labels for grouping and searching alerts.
* `folder_id` - (String) UID of the unified folder in Logz.io.
* `dashboard_id` - (String) UID of the unified dashboard for context linking.
* `panel_id` - (String) Specific panel ID on the dashboard.
* `runbook` - (String) Operational instructions for responders; also used as RCA instruction text when RCA is enabled.
* `enabled` - (Boolean) Alert activation status. Default: `true`.
* `rca` - (Boolean) Enable Root Cause Analysis. Default: `false`.
* `rca_notification_endpoint_ids` - (List of Integer) Notification endpoint IDs for RCA results.
* `use_alert_notification_endpoints_for_rca` - (Boolean) When true, RCA uses same endpoints as alert. Default: `false`.
* `log_alert` - (Block, Max: 1) Log alert configuration. Required when `type = "LOG_ALERT"`. See [Log Alert](#log-alert) below.
* `metric_alert` - (Block, Max: 1) Metric alert configuration. Required when `type = "METRIC_ALERT"`. See [Metric Alert](#metric-alert) below.

### Log Alert

The `log_alert` block supports:

* `search_timeframe_minutes` - (Required, Integer) Time window in minutes for log evaluation.
* `output` - (Required, Block) Notification configuration. See [Log Alert Output](#log-alert-output) below.
* `sub_components` - (Required, List of Block) Detection rules. At least one required. See [Sub Component](#sub-component) below.
* `correlations` - (Optional, Block) Correlation logic between sub-components. See [Correlations](#correlations) below.
* `schedule` - (Optional, Block) Cron-based evaluation schedule. See [Schedule](#schedule) below.

#### Log Alert Output

The `output` block supports:

* `type` - (Required, String) Output format. Must be `JSON` or `TABLE`.
* `suppress_notifications_minutes` - (Optional, Integer) Mute period after alert fires.
* `recipients` - (Required, Block) Who receives notifications. See [Recipients](#recipients) below.

#### Recipients

The `recipients` block supports:

* `emails` - (Optional, List of String) Email addresses for notifications.
* `notification_endpoint_ids` - (Optional, List of Integer) IDs of configured notification endpoints.

**Note:** At least one of `emails` or `notification_endpoint_ids` should be provided.

#### Sub Component

The `sub_components` block supports:

* `query_definition` - (Required, Block) The query configuration. See [Query Definition](#query-definition) below.
* `trigger` - (Required, Block) Trigger conditions. See [Sub Component Trigger](#sub-component-trigger) below.
* `output` - (Optional, Block) Output configuration. See [Sub Component Output](#sub-component-output) below.

#### Query Definition

The `query_definition` block supports:

* `query` - (Required, String) Lucene/Elasticsearch query string (e.g., `"level:ERROR AND service:checkout"`).
* `filters` - (Optional, String) Boolean filters as JSON string. Example shape:

```json
{
  "bool": {
    "must": [],
    "should": [],
    "filter": [],
    "must_not": []
  }
}
```
* `group_by` - (Optional, List of String) Fields to group results by.
* `aggregation` - (Optional, Block) How to aggregate matching logs. See [Aggregation](#aggregation) below.
* `should_query_on_all_accounts` - (Optional, Boolean) Whether to query all accessible accounts. Default: `true`.
* `account_ids_to_query_on` - (Optional, List of Integer) Required if `should_query_on_all_accounts = false`.

#### Aggregation

The `aggregation` block supports:

* `aggregation_type` - (Required, String) Type of aggregation. Valid values: `SUM`, `MIN`, `MAX`, `AVG`, `COUNT`, `UNIQUE_COUNT`, `NONE`.
* `field_to_aggregate_on` - (Optional, String) Field to aggregate on.
* `value_to_aggregate_on` - (Optional, String) Value to aggregate on.

#### Sub Component Trigger

The `trigger` block supports:

* `operator` - (Required, String) Comparison operator. Valid values: `LESS_THAN`, `GREATER_THAN`, `LESS_THAN_OR_EQUALS`, `GREATER_THAN_OR_EQUALS`, `EQUALS`, `NOT_EQUALS`.
* `severity_threshold_tiers` - (Required, List of Block) Severity thresholds. At least one required. See [Severity Threshold Tier](#severity-threshold-tier) below.

#### Severity Threshold Tier

The `severity_threshold_tiers` block supports:

* `severity` - (Required, String) Severity level. Valid values: `INFO`, `LOW`, `MEDIUM`, `HIGH`, `SEVERE`.
* `threshold` - (Required, Float) Threshold value.

**Important:** Threshold ordering depends on the operator:
- **For `GREATER_THAN`/`GREATER_THAN_OR_EQUALS`:** Higher severity must have higher thresholds (e.g., HIGH: 100, MEDIUM: 50, LOW: 10)
- **For `LESS_THAN`/`LESS_THAN_OR_EQUALS`:** Higher severity must have lower thresholds (e.g., HIGH: 10, MEDIUM: 50, LOW: 100)
- **For `EQUALS`/`NOT_EQUALS`:** Standard ordering applies

Example for `LESS_THAN` (detecting low values):
```hcl
trigger {
  operator = "LESS_THAN"
  
  severity_threshold_tiers {
    severity  = "HIGH"
    threshold = 10    # More critical: value < 10
  }
  
  severity_threshold_tiers {
    severity  = "MEDIUM"
    threshold = 50    # Less critical: value < 50
  }
}
```

#### Sub Component Output

The `output` block supports:

* `should_use_all_fields` - (Optional, Boolean) Whether to use all fields in output. Default: `false`.
* `columns` - (Optional, List of Block) Column configurations. See [Column Config](#column-config) below.

**Important:** Custom `columns` are **only valid when `aggregation_type = "NONE"`**. 

- If using any aggregation (`COUNT`, `SUM`, `AVG`, `MIN`, `MAX`, `UNIQUE_COUNT`): **Must set** `should_use_all_fields = true` and **cannot specify** `columns`.
- If using `aggregation_type = "NONE"`: Can set `should_use_all_fields = false` and specify custom `columns`.

Example with aggregation (no custom columns):
```hcl
query_definition {
  aggregation {
    aggregation_type = "COUNT"
  }
}

output {
  should_use_all_fields = true  # Required with aggregation
  # Cannot specify columns here
}
```

Example with NONE aggregation (custom columns allowed):
```hcl
query_definition {
  aggregation {
    aggregation_type = "NONE"
  }
}

output {
  should_use_all_fields = false
  
  columns {
    field_name = "timestamp"
    sort       = "DESC"
  }
  
  columns {
    field_name = "message"
  }
}
```

#### Column Config

The `columns` block supports:

* `field_name` - (Required, String) Field name.
* `regex` - (Optional, String) Regular expression for field extraction.
* `sort` - (Optional, String) Sort direction. Valid values: `ASC`, `DESC`.

#### Schedule

The `schedule` block supports:

* `cron_expression` - (Required, String) Standard cron expression (e.g., `"*/5 * * * *"` = every 5 minutes).
* `timezone` - (Optional, String) Timezone for the cron expression. Default: `UTC`.

#### Correlations

The `correlations` block supports:

* `correlation_operators` - (Optional, List of String) Correlation operators (e.g., `["AND"]`).
* `joins` - (Optional, List of Map) Join configurations.

### Metric Alert

The `metric_alert` block supports:

* `severity` - (Required, String) Alert severity level. Valid values: `INFO`, `LOW`, `MEDIUM`, `HIGH`, `SEVERE`.
* `trigger` - (Required, Block) Trigger configuration. See [Metric Trigger](#metric-trigger) below.
* `queries` - (Required, List of Block) Metric queries. At least one required. See [Metric Query](#metric-query) below.
* `recipients` - (Required, Block) Who receives notifications. See [Recipients](#recipients) above.

#### Metric Trigger

The `trigger` block supports:

* `trigger_type` - (Required, String) Trigger type. Valid values: `THRESHOLD`, `MATH`.
* `metric_operator` - (Optional, String) Required for threshold triggers. Valid values: `ABOVE`, `BELOW`, `WITHIN_RANGE`, `OUTSIDE_RANGE`.
* `min_threshold` - (Optional, Float) Minimum threshold value.
* `max_threshold` - (Optional, Float) Maximum threshold value (required for `WITHIN_RANGE` and `OUTSIDE_RANGE`).
* `math_expression` - (Optional, String) Required when `trigger_type = "MATH"`. Expression using query ref_ids (e.g., `"$A / $B * 100"`).
* `search_timeframe_minutes` - (Required, Integer) Evaluation time window in minutes.

#### Metric Query

The `queries` block supports:

* `ref_id` - (Required, String) Query identifier (e.g., "A", "B") for use in math expressions.
* `query_definition` - (Required, Block) The query configuration. See [Metric Query Definition](#metric-query-definition) below.

#### Metric Query Definition

The `query_definition` block supports:

* `datasource_uid` - (Required, String) UID of the Prometheus/metrics datasource in Logz.io.
* `promql_query` - (Required, String) PromQL query string (e.g., `"rate(http_requests_total[5m])"`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `alert_id` - The unique alert identifier assigned by Logz.io.
* `created_at` - Unix timestamp (float) of alert creation.
* `updated_at` - Unix timestamp (float) of last update.

## Import

Unified alerts can be imported using the alert type and ID, separated by a colon:

```bash
$ terraform import logzio_unified_alert.my_log_alert LOG_ALERT:alert-id-here
$ terraform import logzio_unified_alert.my_metric_alert METRIC_ALERT:alert-id-here
```

**Note:** When importing, you must specify both the alert type (`LOG_ALERT` or `METRIC_ALERT`) and the alert ID.

`Africa/Abidjan`,
`Africa/Accra`,
`Africa/Addis_Ababa`,
`Africa/Algiers`,
`Africa/Asmara`,
`Africa/Asmera`,
`Africa/Bamako`,
`Africa/Bangui`,
`Africa/Banjul`,
`Africa/Bissau`,
`Africa/Blantyre`,
`Africa/Brazzaville`,
`Africa/Bujumbura`,
`Africa/Cairo`,
`Africa/Casablanca`,
`Africa/Ceuta`,
`Africa/Conakry`,
`Africa/Dakar`,
`Africa/Dar_es_Salaam`,
`Africa/Djibouti`,
`Africa/Douala`,
`Africa/El_Aaiun`,
`Africa/Freetown`,
`Africa/Gaborone`,
`Africa/Harare`,
`Africa/Johannesburg`,
`Africa/Juba`,
`Africa/Kampala`,
`Africa/Khartoum`,
`Africa/Kigali`,
`Africa/Kinshasa`,
`Africa/Lagos`,
`Africa/Libreville`,
`Africa/Lome`,
`Africa/Luanda`,
`Africa/Lubumbashi`,
`Africa/Lusaka`,
`Africa/Malabo`,
`Africa/Maputo`,
`Africa/Maseru`,
`Africa/Mbabane`,
`Africa/Mogadishu`,
`Africa/Monrovia`,
`Africa/Nairobi`,
`Africa/Ndjamena`,
`Africa/Niamey`,
`Africa/Nouakchott`,
`Africa/Ouagadougou`,
`Africa/Porto-Novo`,
`Africa/Sao_Tome`,
`Africa/Timbuktu`,
`Africa/Tripoli`,
`Africa/Tunis`,
`Africa/Windhoek`,
`America/Adak`,
`America/Anchorage`,
`America/Anguilla`,
`America/Antigua`,
`America/Araguaina`,
`America/Argentina/Buenos_Aires`,
`America/Argentina/Catamarca`,
`America/Argentina/ComodRivadavia`,
`America/Argentina/Cordoba`,
`America/Argentina/Jujuy`,
`America/Argentina/La_Rioja`,
`America/Argentina/Mendoza`,
`America/Argentina/Rio_Gallegos`,
`America/Argentina/Salta`,
`America/Argentina/San_Juan`,
`America/Argentina/San_Luis`,
`America/Argentina/Tucuman`,
`America/Argentina/Ushuaia`,
`America/Aruba`,
`America/Asuncion`,
`America/Atikokan`,
`America/Atka`,
`America/Bahia`,
`America/Bahia_Banderas`,
`America/Barbados`,
`America/Belem`,
`America/Belize`,
`America/Blanc-Sablon`,
`America/Boa_Vista`,
`America/Bogota`,
`America/Boise`,
`America/Buenos_Aires`,
`America/Cambridge_Bay`,
`America/Campo_Grande`,
`America/Cancun`,
`America/Caracas`,
`America/Catamarca`,
`America/Cayenne`,
`America/Cayman`,
`America/Chicago`,
`America/Chihuahua`,
`America/Coral_Harbour`,
`America/Cordoba`,
`America/Costa_Rica`,
`America/Creston`,
`America/Cuiaba`,
`America/Curacao`,
`America/Danmarkshavn`,
`America/Dawson`,
`America/Dawson_Creek`,
`America/Denver`,
`America/Detroit`,
`America/Dominica`,
`America/Edmonton`,
`America/Eirunepe`,
`America/El_Salvador`,
`America/Ensenada`,
`America/Fort_Nelson`,
`America/Fort_Wayne`,
`America/Fortaleza`,
`America/Glace_Bay`,
`America/Godthab`,
`America/Goose_Bay`,
`America/Grand_Turk`,
`America/Grenada`,
`America/Guadeloupe`,
`America/Guatemala`,
`America/Guayaquil`,
`America/Guyana`,
`America/Halifax`,
`America/Havana`,
`America/Hermosillo`,
`America/Indiana/Indianapolis`,
`America/Indiana/Knox`,
`America/Indiana/Marengo`,
`America/Indiana/Petersburg`,
`America/Indiana/Tell_City`,
`America/Indiana/Vevay`,
`America/Indiana/Vincennes`,
`America/Indiana/Winamac`,
`America/Indianapolis`,
`America/Inuvik`,
`America/Iqaluit`,
`America/Jamaica`,
`America/Jujuy`,
`America/Juneau`,
`America/Kentucky/Louisville`,
`America/Kentucky/Monticello`,
`America/Knox_IN`,
`America/Kralendijk`,
`America/La_Paz`,
`America/Lima`,
`America/Los_Angeles`,
`America/Louisville`,
`America/Lower_Princes`,
`America/Maceio`,
`America/Managua`,
`America/Manaus`,
`America/Marigot`,
`America/Martinique`,
`America/Matamoros`,
`America/Mazatlan`,
`America/Mendoza`,
`America/Menominee`,
`America/Merida`,
`America/Metlakatla`,
`America/Mexico_City`,
`America/Miquelon`,
`America/Moncton`,
`America/Monterrey`,
`America/Montevideo`,
`America/Montreal`,
`America/Montserrat`,
`America/Nassau`,
`America/New_York`,
`America/Nipigon`,
`America/Nome`,
`America/Noronha`,
`America/North_Dakota/Beulah`,
`America/North_Dakota/Center`,
`America/North_Dakota/New_Salem`,
`America/Nuuk`,
`America/Ojinaga`,
`America/Panama`,
`America/Pangnirtung`,
`America/Paramaribo`,
`America/Phoenix`,
`America/Port-au-Prince`,
`America/Port_of_Spain`,
`America/Porto_Acre`,
`America/Porto_Velho`,
`America/Puerto_Rico`,
`America/Punta_Arenas`,
`America/Rainy_River`,
`America/Rankin_Inlet`,
`America/Recife`,
`America/Regina`,
`America/Resolute`,
`America/Rio_Branco`,
`America/Rosario`,
`America/Santa_Isabel`,
`America/Santarem`,
`America/Santiago`,
`America/Santo_Domingo`,
`America/Sao_Paulo`,
`America/Scoresbysund`,
`America/Shiprock`,
`America/Sitka`,
`America/St_Barthelemy`,
`America/St_Johns`,
`America/St_Kitts`,
`America/St_Lucia`,
`America/St_Thomas`,
`America/St_Vincent`,
`America/Swift_Current`,
`America/Tegucigalpa`,
`America/Thule`,
`America/Thunder_Bay`,
`America/Tijuana`,
`America/Toronto`,
`America/Tortola`,
`America/Vancouver`,
`America/Virgin`,
`America/Whitehorse`,
`America/Winnipeg`,
`America/Yakutat`,
`America/Yellowknife`,
`Antarctica/Casey`,
`Antarctica/Davis`,
`Antarctica/DumontDUrville`,
`Antarctica/Macquarie`,
`Antarctica/Mawson`,
`Antarctica/McMurdo`,
`Antarctica/Palmer`,
`Antarctica/Rothera`,
`Antarctica/South_Pole`,
`Antarctica/Syowa`,
`Antarctica/Troll`,
`Antarctica/Vostok`,
`Arctic/Longyearbyen`,
`Asia/Aden`,
`Asia/Almaty`,
`Asia/Amman`,
`Asia/Anadyr`,
`Asia/Aqtau`,
`Asia/Aqtobe`,
`Asia/Ashgabat`,
`Asia/Ashkhabad`,
`Asia/Atyrau`,
`Asia/Baghdad`,
`Asia/Bahrain`,
`Asia/Baku`,
`Asia/Bangkok`,
`Asia/Barnaul`,
`Asia/Beirut`,
`Asia/Bishkek`,
`Asia/Brunei`,
`Asia/Calcutta`,
`Asia/Chita`,
`Asia/Choibalsan`,
`Asia/Chongqing`,
`Asia/Chungking`,
`Asia/Colombo`,
`Asia/Dacca`,
`Asia/Damascus`,
`Asia/Dhaka`,
`Asia/Dili`,
`Asia/Dubai`,
`Asia/Dushanbe`,
`Asia/Famagusta`,
`Asia/Gaza`,
`Asia/Harbin`,
`Asia/Hebron`,
`Asia/Ho_Chi_Minh`,
`Asia/Hong_Kong`,
`Asia/Hovd`,
`Asia/Irkutsk`,
`Asia/Istanbul`,
`Asia/Jakarta`,
`Asia/Jayapura`,
`Asia/Jerusalem`,
`Asia/Kabul`,
`Asia/Kamchatka`,
`Asia/Karachi`,
`Asia/Kashgar`,
`Asia/Kathmandu`,
`Asia/Katmandu`,
`Asia/Khandyga`,
`Asia/Kolkata`,
`Asia/Krasnoyarsk`,
`Asia/Kuala_Lumpur`,
`Asia/Kuching`,
`Asia/Kuwait`,
`Asia/Macao`,
`Asia/Macau`,
`Asia/Magadan`,
`Asia/Makassar`,
`Asia/Manila`,
`Asia/Muscat`,
`Asia/Nicosia`,
`Asia/Novokuznetsk`,
`Asia/Novosibirsk`,
`Asia/Omsk`,
`Asia/Oral`,
`Asia/Phnom_Penh`,
`Asia/Pontianak`,
`Asia/Pyongyang`,
`Asia/Qatar`,
`Asia/Qostanay`,
`Asia/Qyzylorda`,
`Asia/Rangoon`,
`Asia/Riyadh`,
`Asia/Saigon`,
`Asia/Sakhalin`,
`Asia/Samarkand`,
`Asia/Seoul`,
`Asia/Shanghai`,
`Asia/Singapore`,
`Asia/Srednekolymsk`,
`Asia/Taipei`,
`Asia/Tashkent`,
`Asia/Tbilisi`,
`Asia/Tehran`,
`Asia/Tel_Aviv`,
`Asia/Thimbu`,
`Asia/Thimphu`,
`Asia/Tokyo`,
`Asia/Tomsk`,
`Asia/Ujung_Pandang`,
`Asia/Ulaanbaatar`,
`Asia/Ulan_Bator`,
`Asia/Urumqi`,
`Asia/Ust-Nera`,
`Asia/Vientiane`,
`Asia/Vladivostok`,
`Asia/Yakutsk`,
`Asia/Yangon`,
`Asia/Yekaterinburg`,
`Asia/Yerevan`,
`Atlantic/Azores`,
`Atlantic/Bermuda`,
`Atlantic/Canary`,
`Atlantic/Cape_Verde`,
`Atlantic/Faeroe`,
`Atlantic/Faroe`,
`Atlantic/Jan_Mayen`,
`Atlantic/Madeira`,
`Atlantic/Reykjavik`,
`Atlantic/South_Georgia`,
`Atlantic/St_Helena`,
`Atlantic/Stanley`,
`Australia/ACT`,
`Australia/Adelaide`,
`Australia/Brisbane`,
`Australia/Broken_Hill`,
`Australia/Canberra`,
`Australia/Currie`,
`Australia/Darwin`,
`Australia/Eucla`,
`Australia/Hobart`,
`Australia/LHI`,
`Australia/Lindeman`,
`Australia/Lord_Howe`,
`Australia/Melbourne`,
`Australia/NSW`,
`Australia/North`,
`Australia/Perth`,
`Australia/Queensland`,
`Australia/South`,
`Australia/Sydney`,
`Australia/Tasmania`,
`Australia/Victoria`,
`Australia/West`,
`Australia/Yancowinna`,
`Brazil/Acre`,
`Brazil/DeNoronha`,
`Brazil/East`,
`Brazil/West`,
`CET`,
`CST6CDT`,
`Canada/Atlantic`,
`Canada/Central`,
`Canada/Eastern`,
`Canada/Mountain`,
`Canada/Newfoundland`,
`Canada/Pacific`,
`Canada/Saskatchewan`,
`Canada/Yukon`,
`Chile/Continental`,
`Chile/EasterIsland`,
`Cuba`,
`EET`,
`EST5EDT`,
`Egypt`,
`Eire`,
`Etc/GMT`,
`Etc/GMT+0`,
`Etc/GMT+1`,
`Etc/GMT+10`,
`Etc/GMT+11`,
`Etc/GMT+12`,
`Etc/GMT+2`,
`Etc/GMT+3`,
`Etc/GMT+4`,
`Etc/GMT+5`,
`Etc/GMT+6`,
`Etc/GMT+7`,
`Etc/GMT+8`,
`Etc/GMT+9`,
`Etc/GMT-0`,
`Etc/GMT-1`,
`Etc/GMT-10`,
`Etc/GMT-11`,
`Etc/GMT-12`,
`Etc/GMT-13`,
`Etc/GMT-14`,
`Etc/GMT-2`,
`Etc/GMT-3`,
`Etc/GMT-4`,
`Etc/GMT-5`,
`Etc/GMT-6`,
`Etc/GMT-7`,
`Etc/GMT-8`,
`Etc/GMT-9`,
`Etc/GMT0`,
`Etc/Greenwich`,
`Etc/UCT`,
`Etc/UTC`,
`Etc/Universal`,
`Etc/Zulu`,
`Europe/Amsterdam`,
`Europe/Andorra`,
`Europe/Astrakhan`,
`Europe/Athens`,
`Europe/Belfast`,
`Europe/Belgrade`,
`Europe/Berlin`,
`Europe/Bratislava`,
`Europe/Brussels`,
`Europe/Bucharest`,
`Europe/Budapest`,
`Europe/Busingen`,
`Europe/Chisinau`,
`Europe/Copenhagen`,
`Europe/Dublin`,
`Europe/Gibraltar`,
`Europe/Guernsey`,
`Europe/Helsinki`,
`Europe/Isle_of_Man`,
`Europe/Istanbul`,
`Europe/Jersey`,
`Europe/Kaliningrad`,
`Europe/Kiev`,
`Europe/Kirov`,
`Europe/Lisbon`,
`Europe/Ljubljana`,
`Europe/London`,
`Europe/Luxembourg`,
`Europe/Madrid`,
`Europe/Malta`,
`Europe/Mariehamn`,
`Europe/Minsk`,
`Europe/Monaco`,
`Europe/Moscow`,
`Europe/Nicosia`,
`Europe/Oslo`,
`Europe/Paris`,
`Europe/Podgorica`,
`Europe/Prague`,
`Europe/Riga`,
`Europe/Rome`,
`Europe/Samara`,
`Europe/San_Marino`,
`Europe/Sarajevo`,
`Europe/Saratov`,
`Europe/Simferopol`,
`Europe/Skopje`,
`Europe/Sofia`,
`Europe/Stockholm`,
`Europe/Tallinn`,
`Europe/Tirane`,
`Europe/Tiraspol`,
`Europe/Ulyanovsk`,
`Europe/Uzhgorod`,
`Europe/Vaduz`,
`Europe/Vatican`,
`Europe/Vienna`,
`Europe/Vilnius`,
`Europe/Volgograd`,
`Europe/Warsaw`,
`Europe/Zagreb`,
`Europe/Zaporozhye`,
`Europe/Zurich`,
`GB`,
`GB-Eire`,
`GMT`,
`GMT0`,
`Greenwich`,
`Hongkong`,
`Iceland`,
`Indian/Antananarivo`,
`Indian/Chagos`,
`Indian/Christmas`,
`Indian/Cocos`,
`Indian/Comoro`,
`Indian/Kerguelen`,
`Indian/Mahe`,
`Indian/Maldives`,
`Indian/Mauritius`,
`Indian/Mayotte`,
`Indian/Reunion`,
`Iran`,
`Israel`,
`Jamaica`,
`Japan`,
`Kwajalein`,
`Libya`,
`MET`,
`MST7MDT`,
`Mexico/BajaNorte`,
`Mexico/BajaSur`,
`Mexico/General`,
`NZ`,
`NZ-CHAT`,
`Navajo`,
`PRC`,
`PST8PDT`,
`Pacific/Apia`,
`Pacific/Auckland`,
`Pacific/Bougainville`,
`Pacific/Chatham`,
`Pacific/Chuuk`,
`Pacific/Easter`,
`Pacific/Efate`,
`Pacific/Enderbury`,
`Pacific/Fakaofo`,
`Pacific/Fiji`,
`Pacific/Funafuti`,
`Pacific/Galapagos`,
`Pacific/Gambier`,
`Pacific/Guadalcanal`,
`Pacific/Guam`,
`Pacific/Honolulu`,
`Pacific/Johnston`,
`Pacific/Kanton`,
`Pacific/Kiritimati`,
`Pacific/Kosrae`,
`Pacific/Kwajalein`,
`Pacific/Majuro`,
`Pacific/Marquesas`,
`Pacific/Midway`,
`Pacific/Nauru`,
`Pacific/Niue`,
`Pacific/Norfolk`,
`Pacific/Noumea`,
`Pacific/Pago_Pago`,
`Pacific/Palau`,
`Pacific/Pitcairn`,
`Pacific/Pohnpei`,
`Pacific/Ponape`,
`Pacific/Port_Moresby`,
`Pacific/Rarotonga`,
`Pacific/Saipan`,
`Pacific/Samoa`,
`Pacific/Tahiti`,
`Pacific/Tarawa`,
`Pacific/Tongatapu`,
`Pacific/Truk`,
`Pacific/Wake`,
`Pacific/Wallis`,
`Pacific/Yap`,
`Poland`,
`Portugal`,
`ROK`,
`Singapore`,
`SystemV/AST4`,
`SystemV/AST4ADT`,
`SystemV/CST6`,
`SystemV/CST6CDT`,
`SystemV/EST5`,
`SystemV/EST5EDT`,
`SystemV/HST10`,
`SystemV/MST7`,
`SystemV/MST7MDT`,
`SystemV/PST8`,
`SystemV/PST8PDT`,
`SystemV/YST9`,
`SystemV/YST9YDT`,
`Turkey`,
`UCT`,
`US/Alaska`,
`US/Aleutian`,
`US/Arizona`,
`US/Central`,
`US/East-Indiana`,
`US/Eastern`,
`US/Hawaii`,
`US/Indiana-Starke`,
`US/Michigan`,
`US/Mountain`,
`US/Pacific`,
`US/Samoa`,
`UTC`,
`Universal`,
`W-SU`,
`WET`,
`Zulu`,
`EST`,
`HST`,
`MST`,
`ACT`,
`AET`,
`AGT`,
`ART`,
`AST`,
`BET`,
`BST`,
`CAT`,
`CNT`,
`CST`,
`CTT`,
`EAT`,
`ECT`,
`IET`,
`IST`,
`JST`,
`MIT`,
`NET`,
`NST`,
`PLT`,
`PNT`,
`PRT`,
`PST`,
`SST`,
`VST`
