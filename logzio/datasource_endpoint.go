package logzio

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/logzio/logzio_terraform_client/endpoints"
)

func dataSourceEndpoint() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEndpointRead,
		Schema: map[string]*schema.Schema{
			endpointId: {
				Type:     schema.TypeInt,
				Optional: true,
			},
			endpointTitle: {
				Type:     schema.TypeString,
				Optional: true,
			},
			endpointType: {
				Type:     schema.TypeString,
				Optional: true,
			},
			endpointDescription: {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceEndpointRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var client *endpoints.EndpointsClient
	client, _ = endpoints.New(m.(Config).apiToken, m.(Config).baseUrl)

	id, ok := d.GetOk(endpointId)
	if ok {
		endpoint, err := client.GetEndpoint(int64(id.(int)))
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(fmt.Sprintf("%d", id))
		d.Set(endpointTitle, endpoint.Title)
		d.Set(endpointDescription, endpoint.Description)
		d.Set(endpointType, endpoint.Type)
		return nil
	}

	title, ok := d.GetOk(endpointTitle)
	if ok {
		list, err := client.ListEndpoints()
		if err != nil {
			return diag.FromErr(err)
		}
		for _, endpoint := range list {
			if endpoint.Title == title {
				d.SetId(fmt.Sprintf("%d", endpoint.Id))
				d.Set(endpointTitle, endpoint.Title)
				d.Set(endpointDescription, endpoint.Description)
				d.Set(endpointType, endpoint.Type)
				return nil
			}
		}
	}

	return diag.Errorf("couldn't find endpoint with specified attributes")
}
