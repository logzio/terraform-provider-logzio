package logzio

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/logzio/logzio_terraform_client/endpoints"
)

func dataSourceEndpoint() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceEndpointRead,
		Schema: map[string]*schema.Schema{
			endpointId: {
				Type:     schema.TypeInt,
				Optional: true,
			},
			endpointTitle: {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceEndpointRead(d *schema.ResourceData, m interface{}) error {
	var client *endpoints.EndpointsClient
	client, _ = endpoints.New(m.(Config).apiToken, m.(Config).baseUrl)

	endpointId, ok := d.GetOk(endpointId)
	if ok {
		endpoint, err := client.GetEndpoint(int64(endpointId.(int)))
		if err != nil {
			return err
		}
		d.SetId(fmt.Sprintf("%d", endpointId))
		d.Set(endpointTitle, endpoint.Title)
		d.Set(endpointDescription, endpoint.Description)
		d.Set(endpointType, endpoint.EndpointType)
		return nil
	}

	title, ok := d.GetOk(endpointTitle)
	if ok {
		list, err := client.ListEndpoints()
		if err != nil {
			return err
		}
		for i := 0; i < len(list); i++ {
			endpoint := list[i]
			if endpoint.Title == title {
				d.SetId(fmt.Sprintf("%d", endpoint.Id))
				d.Set(endpointTitle, endpoint.Title)
				d.Set(endpointDescription, endpoint.Description)
				d.Set(endpointType, endpoint.EndpointType)
				return nil
			}
		}
	}

	return fmt.Errorf("couldn't find endpoint with specified attributes")
}
