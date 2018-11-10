package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jonboydell/logzio_client"
)

const httpPOSTMethod = "POST"
const logzioAlertsUrl = "https://api.logz.io/v1/alerts"

func resourceAlert() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlertCreate,
		Read:   resourceAlertRead,
		Update: resourceAlertUpdate,
		Delete: resourceAlertDelete,

		Schema: map[string]*schema.Schema{
			"title": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"query_string": &schema.Schema{
				Type:schema.TypeString,
				Required:true,
			},
			"notification_emails": &schema.Schema{
			Type:schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"description": &schema.Schema{
				Type:schema.TypeString,
				Optional:true,
			},
		},
	}
}

func resourceAlertCreate(d *schema.ResourceData, m interface{}) error {

	title := d.Get("title").(string)
	description := d.Get("description").(string)
	queryString := d.Get("query_string").(string)
	notificationEmails := d.Get("notification_emails").([]interface{})

	api_token := m.(Config).api_token

	var t *logzio_client.Thing
	t = logzio_client.New(api_token)
	t.Sure()
	err := t.CreateAlert(title, description, queryString, notificationEmails)

	if err != nil {
		return err
	}

	return resourceAlertRead(d, m)
}

func resourceAlertRead(d *schema.ResourceData, m interface{}) error {
	return nil;
}

func resourceAlertUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceAlertRead(d, m)
}

func resourceAlertDelete(d *schema.ResourceData, m interface{}) error {
	return nil;
}