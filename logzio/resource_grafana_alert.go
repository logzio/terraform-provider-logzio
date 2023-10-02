package logzio

import (
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
							Type:     schema.TypeString,
							Required: true,
							// TODO - validatefunc, statefunc
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
				Type:     schema.TypeString,
				Required: true,
				// TODO - validatefunc
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
				Type:     schema.TypeString,
				Required: true,
				// TODO - validatefunc
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

func suppressDiffJSON(k, oldValue, newValue string, data *schema.ResourceData) bool {
	var o, n interface{}
	d := json.NewDecoder(strings.NewReader(oldValue))
	if err := d.Decode(&o); err != nil {
		return false
	}

	d = json.NewDecoder(strings.NewReader(newValue))
	if err := d.Decode(&n); err != nil {
		return false
	}

	return reflect.DeepEqual(o, n)
}
