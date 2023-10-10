package logzio

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/logzio/logzio_terraform_client/grafana_alerts"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
	"log"
	"reflect"
	"strings"
)

const (
	grafanaAlertRuleAnnotations               = "annotations"
	grafanaAlertRuleCondition                 = "condition"
	grafanaAlertRuleData                      = "data"
	grafanaAlertRuleDataRefId                 = "ref_id"
	grafanaAlertRuleDataDatasourceUid         = "datasource_uid"
	grafanaAlertRuleDataQueryType             = "query_type"
	grafanaAlertRuleDataModel                 = "model"
	grafanaAlertRuleDataRelativeTimeRange     = "relative_time_range"
	grafanaAlertRuleDataRelativeTimeRangeFrom = "from"
	grafanaAlertRuleDataRelativeTimeRangeTo   = "to"
	grafanaAlertRuleLabels                    = "labels"
	grafanaAlertRuleIsPaused                  = "is_paused"
	grafanaAlertRuleExecErrState              = "exec_err_state"
	grafanaAlertRuleFolderUid                 = "folder_uid"
	grafanaAlertRuleFor                       = "for"
	grafanaAlertRuleId                        = "alert_rule_id"
	grafanaAlertRuleNoDataState               = "no_data_state"
	grafanaAlertRuleOrgId                     = "org_id"
	grafanaAlertRuleRuleGroup                 = "rule_group"
	grafanaAlertRuleTitle                     = "title"
	grafanaAlertRuleUid                       = "uid"

	grafanaAlertRuleRetryAttempts = 8
)

func resourceGrafanaAlertRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGrafanaAlertRuleCreate,
		ReadContext:   resourceGrafanaAlertRuleRead,
		UpdateContext: resourceGrafanaAlertRuleUpdate,
		DeleteContext: resourceGrafanaAlertRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			grafanaAlertRuleAnnotations: {
				Type:     schema.TypeMap,
				Optional: true,
				Default:  map[string]interface{}{},
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			grafanaAlertRuleCondition: {
				Type:     schema.TypeString,
				Required: true,
			},
			grafanaAlertRuleData: {
				Type:             schema.TypeList,
				Required:         true,
				MinItems:         1,
				DiffSuppressFunc: suppressDiffJSON,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						grafanaAlertRuleDataRefId: {
							Type:     schema.TypeString,
							Required: true,
						},
						grafanaAlertRuleDataDatasourceUid: {
							Type:     schema.TypeString,
							Required: true,
						},
						grafanaAlertRuleDataQueryType: {
							Type:     schema.TypeString,
							Required: true,
						},
						grafanaAlertRuleDataModel: {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsJSON,
							StateFunc:    handleModelConfig,
						},
						grafanaAlertRuleDataRelativeTimeRange: {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									grafanaAlertRuleDataRelativeTimeRangeFrom: {
										Type:     schema.TypeInt,
										Required: true,
									},
									grafanaAlertRuleDataRelativeTimeRangeTo: {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			grafanaAlertRuleLabels: {
				Type:     schema.TypeMap,
				Optional: true,
				Default:  map[string]interface{}{},
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			grafanaAlertRuleIsPaused: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			grafanaAlertRuleExecErrState: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: utils.ValidateExecErrState,
			},
			grafanaAlertRuleFolderUid: {
				Type:     schema.TypeString,
				Required: true,
			},
			grafanaAlertRuleFor: {
				Type:     schema.TypeInt,
				Required: true,
			},
			grafanaAlertRuleId: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			grafanaAlertRuleNoDataState: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: utils.ValidateExecNoDataState,
			},
			grafanaAlertRuleOrgId: {
				Type:     schema.TypeInt,
				Required: true,
			},
			grafanaAlertRuleRuleGroup: {
				Type:     schema.TypeString,
				Required: true,
			},
			grafanaAlertRuleTitle: {
				Type:     schema.TypeString,
				Required: true,
			},
			grafanaAlertRuleUid: {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func grafanaAlertRuleClient(m interface{}) (*grafana_alerts.GrafanaAlertClient, error) {
	client, err := grafana_alerts.New(m.(Config).apiToken, m.(Config).baseUrl)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func resourceGrafanaAlertRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := grafanaAlertRuleClient(m)
	if err != nil {
		return diag.FromErr(err)
	}

	req := getCreateUpdateGrafanaAlertRuleFromSchema(d)
	result, err := client.CreateGrafanaAlertRule(req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(result.Uid)
	d.Set(grafanaAlertRuleUid, result.Uid)
	d.Set(grafanaAlertRuleId, result.Id)

	return resourceGrafanaAlertRuleRead(ctx, d, m)
}

func resourceGrafanaAlertRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := grafanaAlertRuleClient(m)
	if err != nil {
		return diag.FromErr(err)
	}

	grafanaAlertRule, err := client.GetGrafanaAlertRule(d.Id())
	if err != nil {
		tflog.Error(ctx, err.Error())
		if strings.Contains(err.Error(), "missing grafana alert") {
			// If we were not able to find the resource - delete from state
			d.SetId("")
			return diag.Diagnostics{}
		} else {
			return diag.FromErr(err)
		}
	}

	setGrafanaAlertRule(d, grafanaAlertRule)

	return nil
}

func resourceGrafanaAlertRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := grafanaAlertRuleClient(m)
	if err != nil {
		return diag.FromErr(err)
	}

	req := getCreateUpdateGrafanaAlertRuleFromSchema(d)
	err = client.UpdateGrafanaAlertRule(req)
	if err != nil {
		return diag.FromErr(err)
	}

	var diagRet diag.Diagnostics
	readErr := retry.Do(func() error {
		diagRet = resourceGrafanaAlertRuleRead(ctx, d, m)
		if diagRet.HasError() {
			return fmt.Errorf("received error from read grafana alert rule")
		}

		return nil
	},
		retry.RetryIf(
			// Retry ONLY if the resource was not updated yet
			func(err error) bool {
				if err != nil {
					return false
				} else {
					// Check if the update shows on read
					// if not updated yet - retry
					grafanaAlertRuleFromSchema := getCreateUpdateGrafanaAlertRuleFromSchema(d)
					return !reflect.DeepEqual(grafanaAlertRuleFromSchema, req)
				}
			}),
		retry.DelayType(retry.BackOffDelay),
		retry.Attempts(grafanaAlertRuleRetryAttempts),
	)

	if readErr != nil {
		tflog.Error(ctx, "could not update schema")
		return diagRet
	}

	return nil
}

func resourceGrafanaAlertRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := grafanaAlertRuleClient(m)
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteGrafanaAlertRule(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func setGrafanaAlertRule(d *schema.ResourceData, grafanaAlertRule *grafana_alerts.GrafanaAlertRule) {
	d.Set(grafanaAlertRuleAnnotations, grafanaAlertRule.Annotations)
	d.Set(grafanaAlertRuleCondition, grafanaAlertRule.Condition)
	d.Set(grafanaAlertRuleLabels, grafanaAlertRule.Labels)
	d.Set(grafanaAlertRuleIsPaused, grafanaAlertRule.IsPaused)
	d.Set(grafanaAlertRuleExecErrState, string(grafanaAlertRule.ExecErrState))
	d.Set(grafanaAlertRuleFolderUid, grafanaAlertRule.FolderUID)
	d.Set(grafanaAlertRuleFor, grafanaAlertRule.For)
	d.Set(grafanaAlertRuleNoDataState, string(grafanaAlertRule.NoDataState))
	d.Set(grafanaAlertRuleOrgId, grafanaAlertRule.OrgID)
	d.Set(grafanaAlertRuleRuleGroup, grafanaAlertRule.RuleGroup)
	d.Set(grafanaAlertRuleTitle, grafanaAlertRule.Title)
	data := getDataMapFromAlertRuleObject(grafanaAlertRule.Data)
	d.Set(grafanaAlertRuleData, data)
}

func getDataMapFromAlertRuleObject(data []*grafana_alerts.GrafanaAlertQuery) []map[string]interface{} {
	dataList := make([]map[string]interface{}, 0)
	for _, v := range data {
		dataMap := map[string]interface{}{
			grafanaAlertRuleDataRefId:         v.RefId,
			grafanaAlertRuleDataDatasourceUid: v.DatasourceUid,
			grafanaAlertRuleDataQueryType:     v.QueryType,
			grafanaAlertRuleDataModel:         handleModelConfig(v.Model),
		}

		timeRange := map[string]int{}
		timeRange[grafanaAlertRuleDataRelativeTimeRangeFrom] = int(v.RelativeTimeRange.From)
		timeRange[grafanaAlertRuleDataRelativeTimeRangeTo] = int(v.RelativeTimeRange.To)
		dataMap[grafanaAlertRuleDataRelativeTimeRange] = []interface{}{timeRange}

		dataList = append(dataList, dataMap)
	}

	return dataList
}

func getCreateUpdateGrafanaAlertRuleFromSchema(d *schema.ResourceData) grafana_alerts.GrafanaAlertRule {
	var alertRuleReq grafana_alerts.GrafanaAlertRule

	alertRuleReq.Condition = d.Get(grafanaAlertRuleCondition).(string)
	alertRuleReq.IsPaused = d.Get(grafanaAlertRuleIsPaused).(bool)
	alertRuleReq.ExecErrState = grafana_alerts.ExecErrState(d.Get(grafanaAlertRuleExecErrState).(string))
	alertRuleReq.FolderUID = d.Get(grafanaAlertRuleFolderUid).(string)
	alertRuleReq.For = int64(d.Get(grafanaAlertRuleFor).(int))
	alertRuleReq.NoDataState = grafana_alerts.NoDataState(d.Get(grafanaAlertRuleNoDataState).(string))
	alertRuleReq.OrgID = int64(d.Get(grafanaAlertRuleOrgId).(int))
	alertRuleReq.RuleGroup = d.Get(grafanaAlertRuleRuleGroup).(string)
	alertRuleReq.Title = d.Get(grafanaAlertRuleTitle).(string)
	alertRuleReq.Data = getDataObjectFromSchema(d.Get(grafanaAlertRuleData).([]interface{}))

	if uid, ok := d.GetOk(grafanaAlertRuleUid); ok {
		alertRuleReq.Uid = uid.(string)
	}

	if id, ok := d.GetOk(grafanaAlertRuleId); ok {
		alertRuleReq.Id = int64(id.(int))
	}

	if annotations, ok := d.GetOk(grafanaAlertRuleAnnotations); ok {
		alertRuleReq.Annotations = utils.InterfaceToMapOfStrings(annotations)
	}

	if labels, ok := d.GetOk(grafanaAlertRuleLabels); ok {
		alertRuleReq.Labels = utils.InterfaceToMapOfStrings(labels)
	}

	return alertRuleReq
}

func getDataObjectFromSchema(dataFromSchema []interface{}) []*grafana_alerts.GrafanaAlertQuery {

}

func suppressDiffJSON(k, old, new string, d *schema.ResourceData) bool {
	var oldInterface, newInterface interface{}
	decoder := json.NewDecoder(strings.NewReader(old))
	err := decoder.Decode(&oldInterface)
	if err != nil {
		return false
	}

	decoder = json.NewDecoder(strings.NewReader(new))
	err = decoder.Decode(&newInterface)
	if err != nil {
		return false
	}

	return reflect.DeepEqual(oldInterface, newInterface)
}

func handleModelConfig(model interface{}) string {
	// Default values reference:
	// https://github.com/grafana/grafana/blob/main/pkg/services/ngalert/models/alert_query.go#L12-L13
	const defaultMaxDataPoints float64 = 43200
	const defaultIntervalMS float64 = 1000
	const intervalMsField = "intervalMs"
	const maxDataPointsField = "maxDataPoints"
	modelJsonStr := model.(string)
	var modelObj map[string]interface{}

	err := json.Unmarshal([]byte(modelJsonStr), &modelObj)
	if err != nil {
		log.Printf("Error while unmarshaling model config %v\n", err)
		return modelJsonStr
	}

	iMaxDataPoints, ok := modelObj[maxDataPointsField]
	if ok {
		maxDataPoints, ok := iMaxDataPoints.(float64)
		if ok && maxDataPoints == defaultMaxDataPoints {
			log.Printf("Found default value for %s (%f), removing from model config", maxDataPointsField, defaultMaxDataPoints)
			delete(modelObj, maxDataPointsField)
		}
	}

	iIntervalMs, ok := modelObj[intervalMsField]
	if ok {
		intervalMs, ok := iIntervalMs.(float64)
		if ok && intervalMs == defaultIntervalMS {
			log.Printf("Found default value for %s (%f), removing from model config", intervalMsField, defaultIntervalMS)
			delete(modelObj, intervalMsField)
		}
	}

	modelJson, _ := json.Marshal(modelObj)
	jsonStr := string(modelJson)
	return jsonStr
}
