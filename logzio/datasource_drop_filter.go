package logzio

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/logzio/logzio_terraform_client/drop_filters"
	"reflect"
)

func dataSourceDropFilter() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDropFilterRead,
		Schema: map[string]*schema.Schema{
			dropFilterIdField: {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			dropFilterActive: {
				Type:     schema.TypeBool,
				Computed: true,
			},
			dropFilterLogType: {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			dropFilterFieldConditions: {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						dropFilterFieldName: {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						dropFilterValue: {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			dropFilterThresholdInGB: {
				Type:     schema.TypeFloat,
				Optional: true,
			},
		},
	}
}

func dataSourceDropFilterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, _ := drop_filters.New(m.(Config).apiToken, m.(Config).baseUrl)
	dropFilterId, ok := d.GetOk(dropFilterIdField)
	dropFilters, err := client.RetrieveDropFilters()
	if err != nil {
		return diag.FromErr(err)
	}

	if ok {
		dropFilter := findDropFilterById(dropFilterId.(string), dropFilters)
		if dropFilter != nil {
			d.SetId(dropFilterId.(string))
			setDropFilter(d, dropFilter)
			return nil
		}
	}

	// If could not find id, search for drop filter with the specified attributes
	dropFilterToSearch := createDropFilterFromSchema(d)
	for _, filter := range dropFilters {
		if isSameDropFilter(dropFilterToSearch, filter) {
			d.SetId(filter.Id)
			setDropFilter(d, &filter)
			return nil
		}
	}

	return diag.Errorf("couldn't find drop filter with specified attributes")
}

func isSameDropFilter(dropFilterToSearch drop_filters.DropFilter, dropFilter drop_filters.DropFilter) bool {
	return reflect.DeepEqual(dropFilterToSearch.FieldCondition, dropFilter.FieldCondition) &&
		reflect.DeepEqual(dropFilterToSearch.LogType, dropFilter.LogType)
}
