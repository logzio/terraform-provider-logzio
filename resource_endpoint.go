package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jonboydell/logzio_client/alerts"
	"github.com/jonboydell/logzio_client/endpoints"
	"log"
	"regexp"
	"strconv"
)

const BASE_10 int = 10
const BITSIZE_64 int = 64

const (
	endpoint_type   string = "endpoint_type"
	title                     string = "title"
	description                    string = "description"
	url                         string = "url"
)

func resourceEndpointRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceEndpointCreate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceEndpointUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceEndpointDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func validateEndpointType(v interface{}, k string) (ws []string, errors []error) {
	return
}


func resourceEndpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceEndpointCreate,
		Read:   resourceEndpointRead,
		Update: resourceEndpointUpdate,
		Delete: resourceEndpointDelete,

		Schema: map[string]*schema.Schema{
			endpoint_type: &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				ValidateFunc: validateEndpointType,
			},
			title: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			description: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			url: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validateUrl,
			},
		},
	}
}


func validateUrl(v interface{}, k string) (ws []string, errors []error) {
	regex := "^http(s):\\/\\/"
	value := v.(string)
	b, err := regexp.Match(regex, []byte(value))

	if !b || err != nil {
		errors = append(errors, err)
	}

	return
}


func resourceEndpointDeleteDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("%s::%s", "resourceEndpointDelete", d.Id())
	api_token := m.(Config).api_token

	var client *endpoints.Endpoints
	client, _ = endpoints.New(api_token)

	alertId, _ := strconv.ParseInt(d.Id(), BASE_10, BITSIZE_64)
	err := client.DeleteEndpoint(alertId)
	return err
}
