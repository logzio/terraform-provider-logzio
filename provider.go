package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

const API_TOKEN = "api_token"

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			API_TOKEN: {
				Type:        schema.TypeString,
				Required:    true,
				Description: descriptions[API_TOKEN],
				DefaultFunc: schema.EnvDefaultFunc("LOGZIO_API_TOKEN", nil),
				Sensitive:   true,
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
	apiToken, ok := d.GetOk("api_token")
	if !ok {
		return nil, fmt.Errorf("Can't find the api_token, either set it in the provider or set the LOGZIO_API_TOKEN env var")
	}
	config := Config{
		api_token: apiToken.(string),
	}
	return config, nil
}
