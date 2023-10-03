package logzio

import (
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
	grafanaAlertRuleUpdated                   = "updated"
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
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			grafanaAlertRuleCondition: {
				Type:     schema.TypeString,
				Required: true,
			},
			grafanaAlertRuleData: {
				Type:             schema.TypeSet,
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
			grafanaAlertRuleUpdated: {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		},
	}
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
