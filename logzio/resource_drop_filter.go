package logzio

import (
	"context"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/logzio/logzio_terraform_client/drop_filters"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
)

const (
	dropFilterIdField         = "drop_filter_id"
	dropFilterActive          = "active"
	dropFilterLogType         = "log_type"
	dropFilterFieldConditions = "field_conditions"
	dropFilterFieldName       = "field_name"
	dropFilterValue           = "value"

	dropFilterRetryAttempts = 8
)

// Returns the drop filters client with the api token from the provider
func dropFilterClient(m interface{}) *drop_filters.DropFiltersClient {
	var client *drop_filters.DropFiltersClient
	client, _ = drop_filters.New(m.(Config).apiToken, m.(Config).baseUrl)
	return client
}

func resourceDropFilter() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDropFilterCreate,
		ReadContext:   resourceDropFilterRead,
		UpdateContext: resourceDropFilterUpdate,
		DeleteContext: resourceDropFilterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			dropFilterIdField: {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			dropFilterActive: {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			dropFilterLogType: {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			dropFilterFieldConditions: {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						dropFilterFieldName: {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						dropFilterValue: {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
		},
	}
}

// resourceDropFilterCreate creates a new drop filter in logzio
func resourceDropFilterCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	createDropFilter := createCreatDropFilterFromSchema(d)
	active, exists := d.GetOk(dropFilterActive)
	if exists {
		tflog.Info(ctx, fmt.Sprintf("active attribute is set to %t, note that this field is ignored for creation. A drop filter will always be active after creation.\n", active))
	}
	dropFilter, err := dropFilterClient(m).CreateDropFilter(createDropFilter)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dropFilter.Id)

	return resourceDropFilterRead(ctx, d, m)
}

// resourceDropFilterRead gets drop filter by id
func resourceDropFilterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	dropFilters, err := dropFilterClient(m).RetrieveDropFilters()
	if err != nil {
		return diag.FromErr(err)
	}

	dropFilter := findDropFilterById(d.Id(), dropFilters)
	if dropFilter == nil {
		// If we were not able to find the resource - delete from state
		tflog.Error(ctx, fmt.Sprintf("could not find drop filter with id: %s", d.Id()))
		d.SetId("")
		return diag.Diagnostics{}
	}

	setDropFilter(d, dropFilter)
	return nil
}

// resourceDropFilterUpdate updates drop field by id - activate or deactivate
func resourceDropFilterUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	activate := d.Get(dropFilterActive).(bool)
	var err error
	if activate {
		_, err = dropFilterClient(m).ActivateDropFilter(d.Id())
	} else {
		_, err = dropFilterClient(m).DeactivateDropFilter(d.Id())
	}

	if err != nil {
		return diag.FromErr(err)
	}

	var diagRet diag.Diagnostics
	readErr := retry.Do(
		func() error {
			diagRet = resourceDropFilterRead(ctx, d, m)
			if diagRet.HasError() {
				return fmt.Errorf("received error from read drop filter")
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
					dropFilterFromSchema := createDropFilterFromSchema(d)
					return activate != dropFilterFromSchema.Active
				}
			}),
		retry.DelayType(retry.BackOffDelay),
		retry.Attempts(dropFilterRetryAttempts),
	)

	if readErr != nil {
		tflog.Error(ctx, "could not update schema")
		return diagRet
	}

	return nil
}

// resourceDropFilterDelete deletes drop filter by id
func resourceDropFilterDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := dropFilterClient(m).DeleteDropFilter(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func setDropFilter(d *schema.ResourceData, dropFilter *drop_filters.DropFilter) {
	d.Set(dropFilterIdField, dropFilter.Id)
	d.Set(dropFilterActive, dropFilter.Active)
	d.Set(dropFilterLogType, dropFilter.LogType)

	fieldConditions := getFieldConditionsMapping(dropFilter.FieldCondition)
	d.Set(dropFilterFieldConditions, fieldConditions)
}

func getFieldConditionsMapping(conditions []drop_filters.FieldConditionObject) []map[string]interface{} {
	var conditionsMapping []map[string]interface{}

	for _, condition := range conditions {
		mapping := map[string]interface{}{
			dropFilterFieldName: condition.FieldName,
			dropFilterValue:     utils.ParseObjectToString(condition.Value),
		}

		conditionsMapping = append(conditionsMapping, mapping)
	}

	return conditionsMapping
}

func createCreatDropFilterFromSchema(d *schema.ResourceData) drop_filters.CreateDropFilter {
	fieldConditionsFromSchema := d.Get(dropFilterFieldConditions).([]interface{})
	fieldConditions := getFieldConditionsList(fieldConditionsFromSchema)

	return drop_filters.CreateDropFilter{
		LogType:         d.Get(dropFilterLogType).(string),
		FieldConditions: fieldConditions,
	}
}

func getFieldConditionsList(conditionsFromSchemas []interface{}) []drop_filters.FieldConditionObject {
	var fieldConditions []drop_filters.FieldConditionObject
	var conditionToAppend drop_filters.FieldConditionObject
	for _, element := range conditionsFromSchemas {
		condition := element.(map[string]interface{})
		conditionToAppend.FieldName = condition[dropFilterFieldName].(string)
		conditionToAppend.Value = utils.ParseFromStringToType(condition[dropFilterValue].(string))
		fieldConditions = append(fieldConditions, conditionToAppend)
	}

	return fieldConditions
}

func createDropFilterFromSchema(d *schema.ResourceData) drop_filters.DropFilter {
	fieldConditionsFromSchema := d.Get(dropFilterFieldConditions).([]interface{})
	fieldConditions := getFieldConditionsList(fieldConditionsFromSchema)
	id, ok := d.GetOk(dropFilterIdField)
	if !ok {
		id = d.Id()
	}

	return drop_filters.DropFilter{
		Id:             id.(string),
		Active:         d.Get(dropFilterActive).(bool),
		LogType:        d.Get(dropFilterLogType).(string),
		FieldCondition: fieldConditions,
	}
}

func findDropFilterById(dropFilterId string, dropFilters []drop_filters.DropFilter) *drop_filters.DropFilter {
	for _, filter := range dropFilters {
		if filter.Id == dropFilterId {
			return &filter
		}
	}

	return nil
}
