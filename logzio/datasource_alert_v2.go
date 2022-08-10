package logzio

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/logzio/logzio_terraform_client/alerts_v2"
	"time"
)

func dataSourceAlertV2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlertV2Read,
		Schema: map[string]*schema.Schema{
			alertV2Id: {
				Type:     schema.TypeInt,
				Optional: true,
			},
			alertV2Title: {
				Type:     schema.TypeString,
				Optional: true,
			},
			alertV2Description: {
				Type:     schema.TypeString,
				Computed: true,
			},
			alertV2Tags: {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			alertV2SearchTimeFrameMinutes: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			alertV2IsEnabled: {
				Type:     schema.TypeBool,
				Computed: true,
			},
			alertV2NotificationEmails: {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			alertV2NotificationEndpoints: {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			alertV2SuppressNotificationMinutes: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			alertV2OutputType: {
				Type:     schema.TypeString,
				Computed: true,
			},
			alertV2SubComponents: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						alertV2QueryString: {
							Type:     schema.TypeString,
							Computed: true,
						},
						alertV2FilterMust: {
							Type:     schema.TypeString,
							Computed: true,
						},
						alertV2FilterMustNot: {
							Type:     schema.TypeString,
							Computed: true,
						},
						alertV2GroupBy: {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: groupByMaxItems,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						alertV2AggregationType: {
							Type:     schema.TypeString,
							Computed: true,
						},
						alertV2AggregationField: {
							Type:     schema.TypeString,
							Computed: true,
						},
						alertV2ShouldQueryOnAllAccounts: {
							Type:     schema.TypeBool,
							Computed: true,
						},
						alertV2AccountIdsToQuery: {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						alertV2Operation: {
							Type:     schema.TypeString,
							Computed: true,
						},
						alertV2SeverityThresholdTiers: {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									alertV2Severity: {
										Type:     schema.TypeString,
										Computed: true,
									},
									alertV2Threshold: {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						alertV2Columns: {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									alertV2ColumnsFieldName: {
										Type:     schema.TypeString,
										Computed: true,
									},
									alertV2ColumnsRegex: {
										Type:     schema.TypeString,
										Computed: true,
									},
									alertV2ColumnSort: {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			alertV2CorrelationOperator: {
				Type:     schema.TypeString,
				Computed: true,
			},
			alertV2Joins: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
			alertV2CreatedAt: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			alertV2CreatedBy: {
				Type:     schema.TypeString,
				Computed: true,
			},
			alertV2UpdatedAt: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			alertV2UpdatedBy: {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func getAlertV2(client *alerts_v2.AlertsV2Client, alertId int64, retries int) (*alerts_v2.AlertType, error) {
	alert, err := client.GetAlert(alertId)
	if err != nil && retries > 0 {
		time.Sleep(time.Second * 2)
		alert, err = getAlertV2(client, alertId, retries-1)
	}
	return alert, err
}

func dataSourceAlertV2Read(d *schema.ResourceData, m interface{}) error {
	client, _ := alerts_v2.New(m.(Config).apiToken, m.(Config).baseUrl)
	alertIdString, ok := d.GetOk(alertId)

	if ok {
		id := int64(alertIdString.(int))
		alert, err := getAlertV2(client, id, 3)
		if err != nil {
			return err
		}

		d.SetId(fmt.Sprintf("%d", id))
		setValuesAlertV2(d, alert)
		setCreatedUpdatedFields(d, alert)

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
				setValuesAlertV2(d, &alert)
				setCreatedUpdatedFields(d, &alert)

				return nil
			}
		}
	}

	return fmt.Errorf("couldn't find alert with specified attributes")
}
