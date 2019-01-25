package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jonboydell/logzio_client/endpoints"
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
	apiToken := m.(Config).apiToken

	var client *endpoints.Endpoints
	client, _ = endpoints.New(apiToken)

	endpointId, ok := d.GetOk(endpointId)
	if ok {
		endpoint, err := client.GetEndpoint(endpointId.(int64))
		if err != nil {
			return err
		}
		d.SetId(fmt.Sprintf("%d", endpointId.(int64)))
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
