package logzio

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/logzio/logzio_terraform_client/kibana_objects"
	"reflect"
	"strings"
)

const (
	kibanaObjectKibanaVersionField = "kibana_version"
	kibanaObjectDataField          = "data"
)

// kibanaObjectClient returns the kibana object client with the api token from the provider
func kibanaObjectClient(m interface{}) *kibana_objects.KibanaObjectsClient {
	var client *kibana_objects.KibanaObjectsClient
	client, _ = kibana_objects.New(m.(Config).apiToken, m.(Config).baseUrl)
	return client
}

func resourceKibanaObject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKibanaObjectCreate,
		ReadContext:   resourceKibanaObjectRead,
		UpdateContext: resourceKibanaObjectUpdate,
		DeleteContext: resourceKibanaObjectDelete,
		Schema: map[string]*schema.Schema{
			kibanaObjectKibanaVersionField: {
				Type:     schema.TypeString,
				Required: true,
			},
			kibanaObjectDataField: {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: dataDiff,
			},
		},
	}
}

// resourceKibanaObjectCreate wraps up the import API
func resourceKibanaObjectCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	importReq, err := createImportRequestFromSchema(d)
	if err != nil {
		return diag.FromErr(err)
	}

	kbObjId, err := getIdFromSchema(d)
	if err != nil {
		return diag.FromErr(err)
	}

	importRes, err := kibanaObjectClient(m).ImportKibanaObject(importReq)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(importRes.Created) == 0 {
		return diag.Errorf("error while trying to create. Got: %+v\n", *importRes)
	}

	d.SetId(kbObjId)
	return resourceKibanaObjectRead(ctx, d, m)
}

// resourceKibanaObjectRead wraps the export API
func resourceKibanaObjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	kbObjId, err := getIdFromSchema(d)

	if err != nil {
		return diag.FromErr(err)
	}

	if len(kbObjId) == 0 {
		return nil
	}

	objType, err := getObjectTypeFromData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	err = retry.Do(
		func() error {
			exportRes, err := kibanaObjectClient(m).ExportKibanaObject(kibana_objects.KibanaObjectExportRequest{Type: objType})
			if err != nil {
				return err
			}

			for _, res := range exportRes.Hits {
				if id, ok := res.(map[string]interface{})["_id"].(string); ok {
					if id == kbObjId {
						res.(map[string]interface{})["_index"] = "logzioCustomerIndex*"
						resStr, _ := json.Marshal(res)
						if compareData(d.Get(kibanaObjectDataField).(string), string(resStr)) {
							err = setKibanaObject(d, res.(map[string]interface{}), exportRes.KibanaVersion)
							return err
						}

						return fmt.Errorf("object is not updated yet\n")
					}
				}
			}

			return fmt.Errorf("could not find kibana object with id %s\n", kbObjId)
		},
		retry.RetryIf(
			func(err error) bool {
				if err != nil {
					if strings.Contains(err.Error(), "could not find kibana object with id") ||
						strings.Contains(err.Error(), "object is not updated yet") {
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

// resourceKibanaObjectUpdate wraps up the import API with override field set
func resourceKibanaObjectUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	importReq, err := createImportRequestFromSchema(d)
	if err != nil {
		return diag.FromErr(err)
	}

	importReq.Override = new(bool)
	*importReq.Override = true

	importRes, err := kibanaObjectClient(m).ImportKibanaObject(importReq)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(importRes.Updated) == 0 {
		return diag.Errorf("error while trying to update. Got: %+v", *importRes)
	}

	return resourceKibanaObjectRead(ctx, d, m)
}

// resourceKibanaObjectDelete just remove object from state, user has to delete manually from the app
func resourceKibanaObjectDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.SetId("")

	fmt.Printf("[INFO] Delete object in not supported - just removing from state")
	return nil
}

func setKibanaObject(d *schema.ResourceData, object map[string]interface{}, kibanaVersion string) error {
	d.Set(kibanaObjectKibanaVersionField, kibanaVersion)
	objectStr, err := json.Marshal(object)
	if err != nil {
		return fmt.Errorf("failed to marshal object: %s", err.Error())
	}

	d.Set(kibanaObjectDataField, string(objectStr))
	return nil
}

func getIdFromSchema(d *schema.ResourceData) (string, error) {
	dataStr := d.Get(kibanaObjectDataField).(string)
	var dataObj map[string]interface{}
	err := json.Unmarshal([]byte(dataStr), &dataObj)
	if err != nil {
		return "", err
	}

	if id, ok := dataObj["_id"].(string); ok {
		return id, nil
	}

	return "", fmt.Errorf("could not find id within the data field provided\n")
}

func getObjectTypeFromData(d *schema.ResourceData) (kibana_objects.ExportType, error) {
	dataStr := d.Get(kibanaObjectDataField).(string)
	var dataObj map[string]interface{}
	err := json.Unmarshal([]byte(dataStr), &dataObj)
	if err != nil {
		return "", err
	}
	typeFromData := strings.ToLower(dataObj["_source"].(map[string]interface{})["type"].(string))
	typesMap := []kibana_objects.ExportType{
		kibana_objects.ExportTypeSearch,
		kibana_objects.ExportTypeDashboard,
		kibana_objects.ExportTypeVisualization,
	}

	for _, validType := range typesMap {
		if strings.ToLower(validType.String()) == typeFromData {
			return validType, nil
		}
	}

	return "", fmt.Errorf("could not find valid type within the data field provided\n")
}

func createImportRequestFromSchema(d *schema.ResourceData) (kibana_objects.KibanaObjectImportRequest, error) {
	var importRequest kibana_objects.KibanaObjectImportRequest
	importRequest.KibanaVersion = d.Get(kibanaObjectKibanaVersionField).(string)
	var dataJson map[string]interface{}
	err := json.Unmarshal([]byte(d.Get(kibanaObjectDataField).(string)), &dataJson)
	if err != nil {
		return importRequest, err
	}

	importRequest.Hits = []map[string]interface{}{dataJson}

	return importRequest, nil
}

func dataDiff(k, old, new string, d *schema.ResourceData) bool {
	return compareData(old, new)
}

func compareData(old, new string) bool {
	var oldDataObj, newDataObj map[string]interface{}
	err := json.Unmarshal([]byte(old), &oldDataObj)
	if err != nil {
		if len(old) > 0 {
			fmt.Printf("error while trying to check diff: %s\n", err.Error())
		}
		return false
	}

	err = json.Unmarshal([]byte(new), &newDataObj)
	if err != nil {
		fmt.Printf("error while trying to check diff: %s\n", err.Error())
		return false
	}

	// Fields that we want to ignore their difference
	oldDataObj["_score"] = 0
	newDataObj["_score"] = 0
	oldDataObj["_source"].(map[string]interface{})["updated_at"] = 0
	newDataObj["_source"].(map[string]interface{})["updated_at"] = 0

	res := reflect.DeepEqual(oldDataObj, newDataObj)
	return res
}
