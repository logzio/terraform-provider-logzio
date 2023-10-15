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
	"time"
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
				Optional:     true,
				Default:      string(grafana_alerts.ErrAlerting),
				ValidateFunc: utils.ValidateExecErrState,
			},
			grafanaAlertRuleFolderUid: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			grafanaAlertRuleFor: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: utils.ValidateTimeDurationString,
				StateFunc:    handleDurationString,
			},
			grafanaAlertRuleId: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			grafanaAlertRuleNoDataState: {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      string(grafana_alerts.NoData),
				ValidateFunc: utils.ValidateExecNoDataState,
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

	req, err := getCreateUpdateGrafanaAlertRuleFromSchema(d)
	if err != nil {
		return diag.FromErr(err)
	}

	result, err := client.CreateGrafanaAlertRule(req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(result.Uid)

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

	err = setGrafanaAlertRule(d, grafanaAlertRule)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceGrafanaAlertRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := grafanaAlertRuleClient(m)
	if err != nil {
		return diag.FromErr(err)
	}

	req, err := getCreateUpdateGrafanaAlertRuleFromSchema(d)
	if err != nil {
		return diag.FromErr(err)
	}

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
					grafanaAlertRuleFromSchema, _ := getCreateUpdateGrafanaAlertRuleFromSchema(d)
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

func setGrafanaAlertRule(d *schema.ResourceData, grafanaAlertRule *grafana_alerts.GrafanaAlertRule) error {
	data, err := getDataMapFromAlertRuleObject(grafanaAlertRule.Data)
	if err != nil {
		return err
	}

	forFieldStr := parseNanosecondsToDurationString(grafanaAlertRule.For)
	d.Set(grafanaAlertRuleFor, forFieldStr)
	d.Set(grafanaAlertRuleUid, grafanaAlertRule.Uid)
	d.Set(grafanaAlertRuleId, grafanaAlertRule.Id)
	d.Set(grafanaAlertRuleAnnotations, grafanaAlertRule.Annotations)
	d.Set(grafanaAlertRuleCondition, grafanaAlertRule.Condition)
	d.Set(grafanaAlertRuleLabels, grafanaAlertRule.Labels)
	d.Set(grafanaAlertRuleIsPaused, grafanaAlertRule.IsPaused)
	d.Set(grafanaAlertRuleExecErrState, string(grafanaAlertRule.ExecErrState))
	d.Set(grafanaAlertRuleFolderUid, grafanaAlertRule.FolderUID)
	d.Set(grafanaAlertRuleNoDataState, string(grafanaAlertRule.NoDataState))
	d.Set(grafanaAlertRuleRuleGroup, grafanaAlertRule.RuleGroup)
	d.Set(grafanaAlertRuleTitle, grafanaAlertRule.Title)
	d.Set(grafanaAlertRuleData, data)

	return nil
}

func getDataMapFromAlertRuleObject(data []*grafana_alerts.GrafanaAlertQuery) ([]map[string]interface{}, error) {
	dataList := make([]map[string]interface{}, 0)
	for _, v := range data {
		model, err := json.Marshal(v.Model)
		if err != nil {
			return nil, fmt.Errorf("could not marshal model: %s", err.Error())
		}

		dataMap := map[string]interface{}{
			grafanaAlertRuleDataRefId:         v.RefId,
			grafanaAlertRuleDataDatasourceUid: v.DatasourceUid,
			grafanaAlertRuleDataQueryType:     v.QueryType,
			grafanaAlertRuleDataModel:         handleModelConfig(string(model)),
		}

		timeRange := map[string]int{}
		timeRange[grafanaAlertRuleDataRelativeTimeRangeFrom] = int(v.RelativeTimeRange.From)
		timeRange[grafanaAlertRuleDataRelativeTimeRangeTo] = int(v.RelativeTimeRange.To)
		dataMap[grafanaAlertRuleDataRelativeTimeRange] = []interface{}{timeRange}

		dataList = append(dataList, dataMap)
	}

	return dataList, nil
}

func getCreateUpdateGrafanaAlertRuleFromSchema(d *schema.ResourceData) (grafana_alerts.GrafanaAlertRule, error) {
	var alertRuleReq grafana_alerts.GrafanaAlertRule

	alertRuleReq.Condition = d.Get(grafanaAlertRuleCondition).(string)
	alertRuleReq.IsPaused = d.Get(grafanaAlertRuleIsPaused).(bool)
	alertRuleReq.ExecErrState = grafana_alerts.ExecErrState(d.Get(grafanaAlertRuleExecErrState).(string))
	alertRuleReq.FolderUID = d.Get(grafanaAlertRuleFolderUid).(string)
	alertRuleReq.NoDataState = grafana_alerts.NoDataState(d.Get(grafanaAlertRuleNoDataState).(string))
	alertRuleReq.RuleGroup = d.Get(grafanaAlertRuleRuleGroup).(string)
	alertRuleReq.Title = d.Get(grafanaAlertRuleTitle).(string)
	dataFromSchema, err := getDataObjectFromSchema(d.Get(grafanaAlertRuleData).([]interface{}))
	if err != nil {
		return grafana_alerts.GrafanaAlertRule{}, err
	}

	alertRuleReq.For = parseDurationStringToNanoseconds(d.Get(grafanaAlertRuleFor).(string))

	alertRuleReq.Data = dataFromSchema

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

	// org id is irrelevant, but must be set to a non-zero value to comply with the API
	alertRuleReq.OrgID = 1

	return alertRuleReq, nil
}

func getDataObjectFromSchema(dataFromSchema []interface{}) ([]*grafana_alerts.GrafanaAlertQuery, error) {
	dataList := make([]*grafana_alerts.GrafanaAlertQuery, 0, len(dataFromSchema))
	for _, dataInterface := range dataFromSchema {
		var alertQuery grafana_alerts.GrafanaAlertQuery
		element := dataInterface.(map[string]interface{})

		var modelJson interface{}
		err := json.Unmarshal([]byte(element[grafanaAlertRuleDataModel].(string)), &modelJson)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshel data model: %s", err.Error())
		}

		alertQuery.Model = modelJson
		alertQuery.DatasourceUid = element[grafanaAlertRuleDataDatasourceUid].(string)
		alertQuery.QueryType = element[grafanaAlertRuleDataQueryType].(string)
		alertQuery.RefId = element[grafanaAlertRuleDataRefId].(string)
		if relativeTimeRangeRaw, ok := element[grafanaAlertRuleDataRelativeTimeRange]; ok {
			timeInterfaceSlice := relativeTimeRangeRaw.([]interface{})
			// See object's definition, MaxItems is 1, so no need to iterate
			timeMap := timeInterfaceSlice[0].(map[string]interface{})
			alertQuery.RelativeTimeRange = grafana_alerts.RelativeTimeRangeObj{
				From: time.Duration(timeMap[grafanaAlertRuleDataRelativeTimeRangeFrom].(int)),
				To:   time.Duration(timeMap[grafanaAlertRuleDataRelativeTimeRangeTo].(int)),
			}
		}

		dataList = append(dataList, &alertQuery)
	}

	return dataList, nil
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
			delete(modelObj, maxDataPointsField)
		}
	}

	iIntervalMs, ok := modelObj[intervalMsField]
	if ok {
		intervalMs, ok := iIntervalMs.(float64)
		if ok && intervalMs == defaultIntervalMS {
			delete(modelObj, intervalMsField)
		}
	}

	modelJson, _ := json.Marshal(modelObj)
	jsonStr := string(modelJson)
	return jsonStr
}

func parseNanosecondsToDurationString(nanoseconds int64) string {
	duration := time.Duration(nanoseconds) * time.Nanosecond
	return duration.String()
}

func handleDurationString(durationStr interface{}) string {
	duration, _ := time.ParseDuration(durationStr.(string))
	return duration.String()
}

func parseDurationStringToNanoseconds(durationString string) int64 {
	duration, _ := time.ParseDuration(durationString)
	return duration.Nanoseconds()
}
