package logzio

import (
	"context"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/logzio/logzio_terraform_client/endpoints"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
	"reflect"
	"strconv"
	"strings"
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
		CreateContext: resourceEndpointCreate,
		ReadContext:   resourceEndpointRead,
		UpdateContext: resourceEndpointUpdate,
		DeleteContext: resourceEndpointDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
				ValidateFunc: utils.ValidateEndpointType,
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
							ValidateFunc: utils.ValidateUrl,
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
							ValidateFunc: utils.ValidateUrl,
						},
						endpointMethod: {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: utils.ValidateHttpMethod,
						},
						endpointHeaders: {
							Type:     schema.TypeString,
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
	}
}

// returns the endpoints client with the api token from the provider
func endpointClient(m interface{}) *endpoints.EndpointsClient {
	var client *endpoints.EndpointsClient
	client, _ = endpoints.New(m.(Config).apiToken, m.(Config).baseUrl)
	return client
}

func resourceEndpointCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	createEndpoint := getCreateOrUpdateEndpointFromSchema(d)
	endpoint, err := endpointClient(m).CreateEndpoint(createEndpoint)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(int64(endpoint.Id), 10))
	return resourceEndpointRead(ctx, d, m)
}

func resourceEndpointRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id, err := utils.IdFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	var endpoint *endpoints.Endpoint
	readErr := retry.Do(
		func() error {
			endpoint, err = endpointClient(m).GetEndpoint(id)
			if err != nil {
				return err
			}

			return nil
		},
		retry.RetryIf(
			func(err error) bool {
				if err != nil {
					if strings.Contains(err.Error(), "failed with missing endpoint") {
						return true
					}
				}
				return false
			}),
		retry.DelayType(retry.BackOffDelay),
		retry.Attempts(15),
	)

	if readErr != nil {
		// If we were not able to find the resource - delete from state
		d.SetId("")
		return diag.FromErr(err)
	}

	setEndpoint(d, endpoint)
	return nil
}

func resourceEndpointUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id, _ := utils.IdFromResourceData(d)
	updateEndpoint := getCreateOrUpdateEndpointFromSchema(d)
	_, err := endpointClient(m).UpdateEndpoint(id, updateEndpoint)
	if err != nil {
		return diag.FromErr(err)
	}

	var diagRet diag.Diagnostics
	readErr := retry.Do(
		func() error {
			diagRet = resourceEndpointRead(ctx, d, m)
			if diagRet.HasError() {
				return fmt.Errorf("received error from read endpoint")
			}

			return nil
		},
		retry.RetryIf(
			func(err error) bool {
				if err != nil {
					return true
				} else {
					// Check if the update shows on read
					// if not updated yet - retry
					endpointFromSchema := getCreateOrUpdateEndpointFromSchema(d)
					return !reflect.DeepEqual(updateEndpoint, endpointFromSchema)
				}
			}),
		retry.DelayType(retry.BackOffDelay),
		retry.Attempts(15),
	)

	if readErr != nil {
		tflog.Error(ctx, "could not update schema")
		return diagRet
	}

	return nil
}

func resourceEndpointDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id, err := utils.IdFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	deleteErr := retry.Do(
		func() error {
			return endpointClient(m).DeleteEndpoint(id)
		},
		retry.RetryIf(
			func(err error) bool {
				return err != nil
			}),
		retry.DelayType(retry.BackOffDelay),
		retry.Attempts(15),
	)

	if deleteErr != nil {
		return diag.FromErr(deleteErr)
	}

	d.SetId("")
	return nil
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
		createEndpoint.Headers = new(string)
		*createEndpoint.Headers = opts[endpointHeaders].(string)
		createEndpoint.BodyTemplate = utils.ParseFromStringToType(opts[endpointBodyTemplate].(string))
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
* mappingsFromResourceData returns the mapping stored in terraform map value - who knows why this is not "just" a map, but instead a map
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
		// Microsoft Teams is the only type that's being returned with space.
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
			endpointHeaders:      endpoint.Headers,
			endpointBodyTemplate: utils.ParseObjectToString(endpoint.BodyTemplate),
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
