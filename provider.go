package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

const API_TOKEN = "api_token"

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			API_TOKEN: {
				Type: schema.TypeString,
				Required: true,
				Description: descriptions[API_TOKEN],
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"logzio_alert": resourceAlert(),
		},
		ConfigureFunc: providerConfigure,
	}
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{API_TOKEN: "Your API token"}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		api_token: d.Get("api_token").(string),
	}
	return config, nil
}