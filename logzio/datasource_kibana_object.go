package logzio

import (
	"context"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/logzio/logzio_terraform_client/kibana_objects"
	"strings"
	"time"
)

const (
	kibanaObjectDatasourceIdField   = "object_id"
	kibanaObjectDatasourceTypeField = "object_type"
)

func dataSourceKibanaObject() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKibanaObjectRead,
		Schema: map[string]*schema.Schema{
			kibanaObjectDatasourceIdField: {
				Type:     schema.TypeString,
				Required: true,
			},
			kibanaObjectDatasourceTypeField: {
				Type:     schema.TypeString,
				Required: true,
			},
			kibanaObjectKibanaVersionField: {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			kibanaObjectDataField: {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: dataDiff,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(60 * time.Second),
		},
	}
}

func dataSourceKibanaObjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, _ := kibana_objects.New(m.(Config).apiToken, m.(Config).baseUrl)
	kbObjId := d.Get(kibanaObjectDatasourceIdField).(string)
	kbObjType, err := getKibanaObjectType(d)
	if err != nil {
		return diag.FromErr(err)
	}

	err = retry.Do(
		func() error {
			exportRes, err := client.ExportKibanaObject(kibana_objects.KibanaObjectExportRequest{Type: *kbObjType})
			if err != nil {
				return err
			}

			for _, res := range exportRes.Hits {
				if id, ok := res.(map[string]interface{})["_id"].(string); ok {
					if id == kbObjId {
						res.(map[string]interface{})["_index"] = "logzioCustomerIndex*"
						d.SetId(kbObjId)
						err = setKibanaObject(d, res.(map[string]interface{}), exportRes.KibanaVersion)
						return err
					}
				}
			}

			return fmt.Errorf("could not find object with id %s that matches your config\n", kbObjId)
		},
		retry.RetryIf(
			func(err error) bool {
				if err != nil {
					if strings.Contains(err.Error(), "could not find kibana object with id") {
						return true
					}
				}
				return false
			}),
	)

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func getKibanaObjectType(d *schema.ResourceData) (*kibana_objects.ExportType, error) {
	objTypeFromSchema, ok := d.GetOk(kibanaObjectDatasourceTypeField)
	if !ok {
		return nil, fmt.Errorf("%s must be set", kibanaObjectDatasourceTypeField)
	}

	validTypesMap := map[string]kibana_objects.ExportType{
		kibana_objects.ExportTypeSearch.String():        kibana_objects.ExportTypeSearch,
		kibana_objects.ExportTypeVisualization.String(): kibana_objects.ExportTypeVisualization,
		kibana_objects.ExportTypeDashboard.String():     kibana_objects.ExportTypeDashboard,
	}

	if val, ok := validTypesMap[strings.ToLower(objTypeFromSchema.(string))]; ok {
		return &val, nil
	}

	return nil, fmt.Errorf("%s is invalid type", objTypeFromSchema.(string))
}
