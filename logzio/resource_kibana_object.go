package logzio

import (
	"encoding/json"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/logzio/logzio_terraform_client/kibana_objects"
	"strings"
	"time"
)

const (
	kibanaObjectKibanaVersionField = "kibana_version"
	kibanaObjectDataField          = "data"
)

// Returns the kibana object client with the api token from the provider
func kibanaObjectClient(m interface{}) *kibana_objects.KibanaObjectsClient {
	var client *kibana_objects.KibanaObjectsClient
	client, _ = kibana_objects.New(m.(Config).apiToken, m.(Config).baseUrl)
	return client
}

func resourceKibanaObject() *schema.Resource {
	return &schema.Resource{
		Create: resourceKibanaObjectCreate,
		Read:   resourceKibanaObjectRead,
		Update: resourceKibanaObjectUpdate,
		Delete: resourceKibanaObjectDelete,
		Schema: map[string]*schema.Schema{
			kibanaObjectKibanaVersionField: {
				Type:     schema.TypeString,
				Required: true,
			},
			kibanaObjectDataField: {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Second),
			Read:   schema.DefaultTimeout(5 * time.Second),
			Update: schema.DefaultTimeout(5 * time.Second),
			Delete: schema.DefaultTimeout(5 * time.Second),
		},
	}
}

// resourceKibanaObjectCreate wraps up the import API
func resourceKibanaObjectCreate(d *schema.ResourceData, m interface{}) error {
	importReq, err := createImportRequestFromSchema(d)
	if err != nil {
		return err
	}

	kbObjId, err := getIdFromSchema(d)
	if err != nil {
		return err
	}

	importRes, err := kibanaObjectClient(m).ImportKibanaObject(importReq)
	if err != nil {
		return err
	}

	if len(importRes.Created) == 0 {
		return fmt.Errorf("error while trying to create. Got: %+v", importRes)
	}

	d.SetId(kbObjId)

	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		err = resourceKibanaObjectRead(d, m)
		if err != nil {
			if strings.Contains(err.Error(), "could not find kibana object with id") {
				return resource.RetryableError(err)
			}
		}

		return resource.NonRetryableError(err)
	})
}

// resourceKibanaObjectRead wraps the export API
func resourceKibanaObjectRead(d *schema.ResourceData, m interface{}) error {
	kbObjId, err := getIdFromSchema(d)
	if err != nil {
		return err
	}

	objType, err := getObjectTypeFromData(d)
	if err != nil {
		return err
	}

	err = retry.Do(
		func() error {
			exportRes, err := kibanaObjectClient(m).ExportKibanaObject(kibana_objects.KibanaObjectExportRequest{Type: objType})
			if err != nil {
				return err
			}

			for _, res := range exportRes.Hits {
				if id, ok := res.(map[string]interface{})["_source"].(map[string]interface{})["id"].(string); ok {
					if id == kbObjId {
						err = setKibanaObject(d, res.(map[string]interface{}), exportRes.KibanaVersion)
						return err
					}
				}
			}

			return fmt.Errorf("could not find kibana object with id %s", kbObjId)
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

	return err
}

// resourceKibanaObjectUpdate wraps up the import API with override field set
func resourceKibanaObjectUpdate(d *schema.ResourceData, m interface{}) error {
	importReq, err := createImportRequestFromSchema(d)
	if err != nil {
		return err
	}

	*importReq.Override = true

	importRes, err := kibanaObjectClient(m).ImportKibanaObject(importReq)
	if err != nil {
		return err
	}

	if len(importRes.Updated) == 0 {
		return fmt.Errorf("error while trying to update. Got: %+v", importRes)
	}

	return resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		err = resourceKibanaObjectRead(d, m)
		if err != nil {
			if strings.Contains(err.Error(), "could not find kibana object with id") {
				return resource.RetryableError(err)
			}
		}

		return resource.NonRetryableError(err)
	})
}

// resourceKibanaObjectDelete just remove object from state, user has to delete manually from the app
func resourceKibanaObjectDelete(d *schema.ResourceData, m interface{}) error {
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

	if id, ok := dataObj["_source"].(map[string]interface{})["id"].(string); ok {
		return id, nil
	}

	objType, err := getObjectTypeFromData(d)
	if err != nil {
		return "", err
	}

	if id, ok := dataObj["_source"].(map[string]interface{})[objType.String()].(map[string]interface{})["id"].(string); ok {
		return id, nil
	}

	return "", fmt.Errorf("could not find id within the data field provided")
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

	return "", fmt.Errorf("could not find valid type within the data field provided")
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
