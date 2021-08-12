package logzio

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/logzio/logzio_terraform_client/endpoints"
)

const (
	endpointId            string = "id"
	endpointIdField       string = "endpoint_id"
	endpointType          string = "endpoint_type"
	endpointTitle         string = "title"
	endpointDescription   string = "description"
	endpointUrl           string = "url"
	endpointMethod        string = "method"
	endpointHeaders       string = "headers"
	endpointBodyTemplate  string = "body_template"
	endpointServiceKey    string = "service_key"
	endpointApiToken      string = "api_token"
	endpointAppKey        string = "app_key"
	endpointApiKey        string = "api_key"
	endpointRoutingKey    string = "routing_key"
	endpointMessageType   string = "message_type"
	endpointServiceApiKey string = "service_api_key"
	endpointUsername      string = "username"
	endpointPassword      string = "password"

	endpointTypeMicrosoftTeamsFromApi = "microsoft teams"
)

/**
 * the endpoint resource schema, what terraform uses to parse and read the template
 */
func resourceEndpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceEndpointCreate,
		Read:   resourceEndpointRead,
		Update: resourceEndpointUpdate,
		Delete: resourceEndpointDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			endpointIdField: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			endpointType: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateEndpointType,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					new = strings.ToLower(new)
					newUnderScore := strings.Replace(new, "_", "", 1)
					if strings.EqualFold(old, newUnderScore) {
						return true
					}
					newSpace := strings.Replace(new, " ", "", 1)
					if strings.EqualFold(old, newSpace) {
						return true
					}
					return false
				},
			},
			endpointTitle: {
				Type:     schema.TypeString,
				Required: true,
			},
			endpointDescription: {
				Type:     schema.TypeString,
				Optional: true,
			},
			endpoints.EndpointTypeSlack: {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						endpointUrl: {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateUrl,
						},
					},
				},
			},
			endpoints.EndpointTypeCustom: {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						endpointUrl: {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateUrl,
						},
						endpointMethod: {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateHttpMethod,
						},
						endpointHeaders: {
							Type:     schema.TypeMap,
							Optional: true,
						},
						endpointBodyTemplate: {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			endpoints.EndpointTypePagerDuty: {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						endpointServiceKey: {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			endpoints.EndpointTypeBigPanda: {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						endpointApiToken: {
							Type:     schema.TypeString,
							Required: true,
						},
						endpointAppKey: {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			endpoints.EndpointTypeDataDog: {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						endpointApiKey: {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			endpoints.EndpointTypeVictorOps: {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						endpointRoutingKey: {
							Type:     schema.TypeString,
							Required: true,
						},
						endpointMessageType: {
							Type:     schema.TypeString,
							Required: true,
						},
						endpointServiceApiKey: {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			endpoints.EndpointTypeOpsGenie: {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						endpointApiKey: {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			endpoints.EndpointTypeServiceNow: {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						endpointUsername: {
							Type:     schema.TypeString,
							Required: true,
						},
						endpointPassword: {
							Type:     schema.TypeString,
							Required: true,
						},
						endpointUrl: {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			endpoints.EndpointTypeMicrosoftTeams: {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						endpointUrl: {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Second),
			Update: schema.DefaultTimeout(5 * time.Second),
			Delete: schema.DefaultTimeout(5 * time.Second),
		},
	}
}

// returns the endpoints client with the api token from the provider
func endpointClient(m interface{}) *endpoints.EndpointsClient {
	var client *endpoints.EndpointsClient
	client, _ = endpoints.New(m.(Config).apiToken, m.(Config).baseUrl)
	return client
}

func resourceEndpointCreate(d *schema.ResourceData, m interface{}) error {
	createEndpoint := getCreateOrUpdateEndpointFromSchema(d)
	endpoint, err := endpointClient(m).CreateEndpoint(createEndpoint)
	if err != nil {
		return err
	}

	d.SetId(strconv.FormatInt(int64(endpoint.Id), 10))

	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		err = resourceEndpointRead(d, m)
		if err != nil {
			if strings.Contains(err.Error(), "failed with missing endpoint") {
				return resource.RetryableError(err)
			}
		}

		return resource.NonRetryableError(err)
	})
}

func resourceEndpointRead(d *schema.ResourceData, m interface{}) error {
	id, err := idFromResourceData(d)
	if err != nil {
		return nil
	}

	endpoint, err := endpointClient(m).GetEndpoint(id)
	if err != nil {
		return err
	}

	setEndpoint(d, endpoint)
	return nil
}

func resourceEndpointUpdate(d *schema.ResourceData, m interface{}) error {
	id, _ := idFromResourceData(d)
	updateEndpoint := getCreateOrUpdateEndpointFromSchema(d)
	_, err := endpointClient(m).UpdateEndpoint(id, updateEndpoint)
	if err != nil {
		return err
	}

	return resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		err = resourceEndpointRead(d, m)
		if err != nil {
			endpointFromSchema := getCreateOrUpdateEndpointFromSchema(d)
			if strings.Contains(err.Error(), "failed with missing endpoint") &&
				!reflect.DeepEqual(updateEndpoint, endpointFromSchema) {
				return resource.RetryableError(fmt.Errorf("endpoint is not updated yet: %s", err.Error()))
			}
		}

		return resource.NonRetryableError(err)
	})
}

func resourceEndpointDelete(d *schema.ResourceData, m interface{}) error {
	endpointId, _ := idFromResourceData(d)

	return resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		err := endpointClient(m).DeleteEndpoint(endpointId)
		if err != nil {
			return resource.RetryableError(err)
		}

		return resource.NonRetryableError(err)
	})
}

func getCreateOrUpdateEndpointFromSchema(d *schema.ResourceData) endpoints.CreateOrUpdateEndpoint {
	createEndpoint := endpoints.CreateOrUpdateEndpoint{
		Title:       d.Get(endpointTitle).(string),
		Description: d.Get(endpointDescription).(string),
		Type:        d.Get(endpointType).(string),
	}

	opts, _ := mappingsFromResourceData(d, createEndpoint.Type)
	switch createEndpoint.Type {
	case endpoints.EndpointTypeSlack:
		createEndpoint.Url = opts[endpointUrl].(string)
	case endpoints.EndpointTypeCustom:
		createEndpoint.Url = opts[endpointUrl].(string)
		createEndpoint.Method = opts[endpointMethod].(string)
		headerMap := make(map[string]string)
		for k, v := range opts[endpointHeaders].(map[string]interface{}) {
			headerMap[k] = v.(string)
		}
		createEndpoint.Headers = parseObjectToString(headerMap)
		createEndpoint.BodyTemplate = parseFromStringToType(opts[endpointBodyTemplate].(string))
	case endpoints.EndpointTypePagerDuty:
		createEndpoint.ServiceKey = opts[endpointServiceKey].(string)
	case endpoints.EndpointTypeBigPanda:
		createEndpoint.ApiToken = opts[endpointApiToken].(string)
		createEndpoint.AppKey = opts[endpointAppKey].(string)
	case endpoints.EndpointTypeDataDog:
		createEndpoint.ApiKey = opts[endpointApiKey].(string)
	case endpoints.EndpointTypeVictorOps:
		createEndpoint.RoutingKey = opts[endpointRoutingKey].(string)
		createEndpoint.MessageType = opts[endpointMessageType].(string)
		createEndpoint.ServiceApiKey = opts[endpointServiceApiKey].(string)
	case endpoints.EndpointTypeOpsGenie:
		createEndpoint.ApiKey = opts[endpointApiKey].(string)
	case endpoints.EndpointTypeServiceNow:
		createEndpoint.Username = opts[endpointUsername].(string)
		createEndpoint.Password = opts[endpointPassword].(string)
		createEndpoint.Url = opts[endpointUrl].(string)
	case endpoints.EndpointTypeMicrosoftTeams:
		createEndpoint.Url = opts[endpointUrl].(string)
	default:
		panic(fmt.Sprintf("unhandled endpoint type %s", createEndpoint.Type))
	}

	return createEndpoint
}

/*
* returns the mapping stored in a terraform map value - who knows why this is not "just" a map, but instead a map
* wrapped in an array
 */
func mappingsFromResourceData(d *schema.ResourceData, key string) (map[string]interface{}, error) {
	if v, ok := d.GetOk(key); ok {
		rawMappings := v.(*schema.Set).List()
		for i := 0; i < len(rawMappings); i++ {
			x := rawMappings[i]
			y := x.(map[string]interface{})
			return y, nil
		}
	}

	return nil, fmt.Errorf("can't load mapping for key %s", key)
}

func setEndpoint(d *schema.ResourceData, endpoint *endpoints.Endpoint) {
	d.Set(endpointIdField, endpoint.Id)
	d.Set(endpointTitle, endpoint.Title)
	d.Set(endpointDescription, endpoint.Description)
	typeLowerCase := strings.ToLower(endpoint.Type)
	if typeLowerCase == endpointTypeMicrosoftTeamsFromApi {
		// microsoft teams is the only type that's being returned with space.
		typeLowerCase = strings.Replace(typeLowerCase, " ", "", 1)
	}
	d.Set(endpointType, typeLowerCase)
	set := make([]map[string]interface{}, 1)
	switch typeLowerCase {
	case endpoints.EndpointTypeSlack:
		set[0] = map[string]interface{}{
			endpointUrl: endpoint.Url,
		}
	case endpoints.EndpointTypeCustom:
		set[0] = map[string]interface{}{
			endpointUrl:          endpoint.Url,
			endpointMethod:       endpoint.Method,
			endpointHeaders:      parseFromStringToType(endpoint.Headers),
			endpointBodyTemplate: parseObjectToString(endpoint.BodyTemplate),
		}
	case endpoints.EndpointTypePagerDuty:
		set[0] = map[string]interface{}{
			endpointServiceKey: endpoint.ServiceKey,
		}
	case endpoints.EndpointTypeBigPanda:
		set[0] = map[string]interface{}{
			endpointApiToken: endpoint.ApiToken,
			endpointAppKey:   endpoint.AppKey,
		}
	case endpoints.EndpointTypeDataDog:
		set[0] = map[string]interface{}{
			endpointApiKey: endpoint.ApiKey,
		}
	case endpoints.EndpointTypeVictorOps:
		set[0] = map[string]interface{}{
			endpointRoutingKey:    endpoint.RoutingKey,
			endpointMessageType:   endpoint.MessageType,
			endpointServiceApiKey: endpoint.ServiceApiKey,
		}
	case endpoints.EndpointTypeOpsGenie:
		set[0] = map[string]interface{}{
			endpointApiKey: endpoint.ApiKey,
		}
	case endpoints.EndpointTypeServiceNow:
		set[0] = map[string]interface{}{
			endpointUsername: endpoint.Username,
			endpointPassword: endpoint.Password,
			endpointUrl:      endpoint.Url,
		}
	case endpoints.EndpointTypeMicrosoftTeams:
		set[0] = map[string]interface{}{
			endpointUrl: endpoint.Url,
		}
	default:
		panic(fmt.Sprintf("unhandled endpoint type %s", typeLowerCase))
	}

	d.Set(typeLowerCase, set)
}
