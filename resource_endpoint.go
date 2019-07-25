package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jonboydell/logzio_client/endpoints"
)

const (
	endpointId            string = "id"
	endpointType          string = "endpoint_type"
	endpointTitle         string = "title"
	endpointDescription   string = "description"
	endpointUrl           string = "url"
	endpointSlack         string = "slack"
	endpointCustom        string = "custom"
	endpointMethod        string = "method"
	endpointHeaders       string = "headers"
	endpointBodyTemplate  string = "body_template"
	endpointPagerDuty     string = "pager_duty"
	endpointServiceKey    string = "service_key"
	endpointBigPanda      string = "big_panda"
	endpointApiToken      string = "api_token"
	endpointAppKey        string = "app_key"
	endpointDataDog       string = "data_dog"
	endpointApiKey        string = "api_key"
	endpointVictorOps     string = "victorops"
	endpointRoutingKey    string = "routing_key"
	endpointMessageType   string = "message_type"
	endpointServiceApiKey string = "service_api_key"
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

		Schema: map[string]*schema.Schema{
			endpointType: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateEndpointType,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					new = strings.Replace(new, "_", "", 1)
					if strings.EqualFold(old, new) {
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
			endpointSlack: {
				Type:     schema.TypeSet,
				Optional: true,
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
			endpointCustom: {
				Type:     schema.TypeSet,
				Optional: true,
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
							Required: true,
						},
						endpointBodyTemplate: {
							Type:     schema.TypeMap,
							Required: true,
						},
					},
				},
			},
			endpointPagerDuty: {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						endpointServiceKey: {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			endpointBigPanda: {
				Type:     schema.TypeSet,
				Optional: true,
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
			endpointDataDog: {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						endpointApiKey: {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			endpointVictorOps: {
				Type:     schema.TypeSet,
				Optional: true,
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
		},
	}
}

/**
 * returns the endpoints client with the api token from the provider
 */
func endpointClient(m interface{}) *endpoints.EndpointsClient {
	var client *endpoints.EndpointsClient
	client, _ = endpoints.New(m.(Config).apiToken, m.(Config).baseUrl)
	return client
}

/*
 * returns the id from terraform, parsed to an int64
 * @todo: needs to be moved out of this file and into the commons
 */
func idFromResourceData(d *schema.ResourceData) (int64, error) {
	return strconv.ParseInt(d.Id(), BASE_10, BITSIZE_64)
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

/**
 * returns an endpoint object populated from a resource data object, used for creates and updates
 */
func endpointFromResourceData(d *schema.ResourceData) endpoints.Endpoint {
	eType := d.Get(endpointType).(string)

	endpoint := endpoints.Endpoint{
		Title:       d.Get(endpointTitle).(string),
		Description: d.Get(endpointDescription).(string),
	}

	switch eType {
	case string(endpoints.EndpointTypeSlack):
		endpoint.EndpointType = endpoints.EndpointTypeSlack
		opts, _ := mappingsFromResourceData(d, endpointSlack)
		endpoint.Url = opts[endpointUrl].(string)
	case string(endpoints.EndpointTypeCustom):
		endpoint.EndpointType = endpoints.EndpointTypeCustom
		opts, _ := mappingsFromResourceData(d, endpointCustom)
		endpoint.Url = opts[endpointUrl].(string)
		endpoint.Method = opts[endpointMethod].(string)
		endpoint.BodyTemplate = opts[endpointBodyTemplate]
		headerMap := make(map[string]string)
		for k, v := range opts[endpointHeaders].(map[string]interface{}) {
			headerMap[k] = v.(string)
		}
		endpoint.Headers = headerMap
	case string(endpoints.EndpointTypePagerDuty):
		endpoint.EndpointType = endpoints.EndpointTypePagerDuty
		opts, _ := mappingsFromResourceData(d, endpointPagerDuty)
		endpoint.ServiceKey = opts[endpointServiceKey].(string)
	case string(endpoints.EndpointTypeBigPanda):
		endpoint.EndpointType = endpoints.EndpointTypeBigPanda
		opts, _ := mappingsFromResourceData(d, endpointBigPanda)
		endpoint.ApiToken = opts[endpointApiToken].(string)
		endpoint.AppKey = opts[endpointAppKey].(string)
	case string(endpoints.EndpointTypeDataDog):
		endpoint.EndpointType = endpoints.EndpointTypeDataDog
		opts, _ := mappingsFromResourceData(d, endpointDataDog)
		endpoint.ApiKey = opts[endpointApiKey].(string)
	case string(endpoints.EndpointTypeVictorOps):
		endpoint.EndpointType = endpoints.EndpointTypeVictorOps
		opts, _ := mappingsFromResourceData(d, endpointVictorOps)
		endpoint.RoutingKey = opts[endpointRoutingKey].(string)
		endpoint.MessageType = opts[endpointMessageType].(string)
		endpoint.ServiceApiKey = opts[endpointServiceApiKey].(string)
	default:
		panic(fmt.Sprintf("unhandled endpoint type %s", eType))
	}
	return endpoint
}

/**
 * creates a new endpoint in logzio
 */
func resourceEndpointCreate(d *schema.ResourceData, m interface{}) error {
	endpoint := endpointFromResourceData(d)
	client := endpointClient(m)
	e, err := client.CreateEndpoint(endpoint)

	if err != nil {
		return err
	}

	endpointId := strconv.FormatInt(e.Id, BASE_10)
	d.SetId(endpointId)

	return nil
}

/**
 * reads an endpoint from logzio
 */
func resourceEndpointRead(d *schema.ResourceData, m interface{}) error {
	client := endpointClient(m)
	endpointId, _ := idFromResourceData(d)

	var endpoint *endpoints.Endpoint
	endpoint, err := client.GetEndpoint(endpointId)
	if err != nil {
		return err
	}

	d.Set(endpointType, endpoint.EndpointType)
	d.Set(endpointTitle, endpoint.Title)
	d.Set(endpointDescription, endpoint.Description)

	set := make([]map[string]interface{}, 1)

	switch endpoint.EndpointType {
	case endpoints.EndpointTypeSlack:
		set[0] = map[string]interface{}{
			endpointUrl: endpoint.Url,
		}
		d.Set(endpointSlack, set)
	case endpoints.EndpointTypeCustom:
		set[0] = map[string]interface{}{
			endpointUrl:          endpoint.Url,
			endpointMethod:       endpoint.Method,
			endpointHeaders:      endpoint.Headers,
			endpointBodyTemplate: endpoint.BodyTemplate,
		}
		d.Set(endpointCustom, set)
	case endpoints.EndpointTypePagerDuty:
		d.Set(endpointType, endpointPagerDuty)

		set[0] = map[string]interface{}{
			endpointServiceKey: endpoint.ServiceKey,
		}
		d.Set(endpointPagerDuty, set)
	case endpoints.EndpointTypeBigPanda:
		d.Set(endpointType, endpointBigPanda)

		set[0] = map[string]interface{}{
			endpointApiToken: endpoint.ApiToken,
			endpointAppKey:   endpoint.AppKey,
		}
		d.Set(endpointBigPanda, set)
	case endpoints.EndpointTypeDataDog:
		d.Set(endpointType, endpointDataDog)

		set[0] = map[string]interface{}{
			endpointApiKey: endpoint.ApiKey,
		}
		d.Set(endpointDataDog, set)
	case endpoints.EndpointTypeVictorOps:
		d.Set(endpointType, endpointVictorOps)

		set[0] = map[string]interface{}{
			endpointRoutingKey:    endpoint.RoutingKey,
			endpointMessageType:   endpoint.MessageType,
			endpointServiceApiKey: endpoint.ApiKey,
		}
		d.Set(endpointVictorOps, set)
	default:
		return fmt.Errorf("unhandled endpoint type: %s", endpoint.EndpointType)
	}

	return nil
}

/**
 * Updates an existing endpoint, returns an error if the endpoint can't be found
 */
func resourceEndpointUpdate(d *schema.ResourceData, m interface{}) error {
	endpoint := endpointFromResourceData(d)
	endpoint.Id, _ = idFromResourceData(d)
	client := endpointClient(m)
	_, err := client.UpdateEndpoint(endpoint.Id, endpoint)

	if err != nil {
		return err
	}

	return nil
}

/**
 * deletes an existing endpoint, returns an error if the endpoint can't be found
 */
func resourceEndpointDelete(d *schema.ResourceData, m interface{}) error {
	endpointId, _ := idFromResourceData(d)
	client := endpointClient(m)
	err := client.DeleteEndpoint(endpointId)
	if err != nil {
		return err
	}

	return nil
}

// @todo - what's this structure - 	if v, ok := d.GetOk("enable_log_file_validation"); ok { some_function() } ??

/**
 * checks that the endpoint type is something we recognize (see docs for the type of endpoints supported), is case sensitive
 */
func validateEndpointType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if !findStringInArray(value, []string{
		string(endpoints.EndpointTypeSlack),
		string(endpoints.EndpointTypeCustom),
		string(endpoints.EndpointTypePagerDuty),
		string(endpoints.EndpointTypeBigPanda),
		string(endpoints.EndpointTypeDataDog),
		string(endpoints.EndpointTypeVictorOps),
	}) {
		errors = append(errors, fmt.Errorf("value for endpoint type is unknown"))
	}

	return
}

/**
 * checks that a provided url is kind of in the right format, logzio will reject URLs that it can't resolve, and there's
 * no checking for that here
 */
func validateUrl(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	b, err := regexp.Match(VALIDATE_URL_REGEX, []byte(value))

	if !b || err != nil {
		errors = append(errors, err)
	}

	return
}

/**
 * checks that the provided HTTP method is something we recognise (GET/POST/PUT/DELETE), is case sensitive
 */
func validateHttpMethod(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if !findStringInArray(value, []string{"GET", "POST", "PUT", "DELETE"}) {
		errors = append(errors, fmt.Errorf("invalid HTTP method specified"))
	}

	return
}
