package logzio

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/logzio/logzio_terraform_client/log_shipping_tokens"
	"strconv"
)

const (
	logShippingTokenId        string = "id"
	logShippingTokenName      string = "name"
	logShippingTokenEnabled   string = "enabled"
	logShippingTokenToken     string = "token"
	logShippingTokenUpdatedAt        = "updated_at"
	logShippingTokenUpdatedBy        = "updated_by"
	logShippingTokenCreatedAt        = "created_at"
	logShippingTokenCreatedBy        = "created_by"
)

// TODO:
// 1. What to do with ENABLED - required? computed? optional? instructions to always set on create to true?
// 2. Retry + timeout

func resourceLogShippingToken() *schema.Resource {
	return &schema.Resource{
		Create: resourceLogShippingTokenCreate,
		Read:   resourceLogShippingTokenRead,
		Update: resourceLogShippingTokenUpdate,
		Delete: resourceLogShippingTokenDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			logShippingTokenId: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			logShippingTokenName: {
				Type:     schema.TypeString,
				Required: true,
			},
			logShippingTokenEnabled: {
				Type:     schema.TypeBool,
				Required: true,
			},
			logShippingTokenToken: {
				Type:     schema.TypeString,
				Computed: true,
			},
			logShippingTokenUpdatedAt: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			logShippingTokenUpdatedBy: {
				Type:     schema.TypeString,
				Computed: true,
			},
			logShippingTokenCreatedAt: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			logShippingTokenCreatedBy: {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// Creates a new log shipping token in logz.io
func resourceLogShippingTokenCreate(d *schema.ResourceData, m interface{}) error {
	createToken := log_shipping_tokens.CreateLogShippingToken{Name: d.Get(logShippingTokenName).(string)}
	token, err := logShippingTokenClient(m).CreateLogShippingToken(createToken)
	if err != nil {
		return err
	}

	d.SetId(strconv.FormatInt(int64(token.Id), 10))

	return resourceLogShippingTokenRead(d, m)
}

func resourceLogShippingTokenRead(d *schema.ResourceData, m interface{}) error {
	id, err := idFromResourceData(d)
	if err != nil {
		return err
	}

	token, err := logShippingTokenClient(m).GetLogShippingToken(int32(id))
	if err != nil {
		return err
	}

	setLogShippingToken(d, token)
	return nil
}

func resourceLogShippingTokenUpdate(d *schema.ResourceData, m interface{}) error {
	id, err := idFromResourceData(d)
	if err != nil {
		return err
	}

	updateToken := log_shipping_tokens.UpdateLogShippingToken{
		Name:    d.Get(logShippingTokenName).(string),
		Enabled: strconv.FormatBool(d.Get(logShippingTokenEnabled).(bool)),
	}

	token, err := logShippingTokenClient(m).UpdateLogShippingToken(int32(id), updateToken)
	if err != nil {
		return err
	}

	return resourceLogShippingTokenRead(d, token)
}

func resourceLogShippingTokenDelete(d *schema.ResourceData, m interface{}) error {
	id, err := idFromResourceData(d)
	if err != nil {
		return err
	}

	err = logShippingTokenClient(m).DeleteLogShippingToken(int32(id))
	if err != nil {
		return err
	}

	return nil
}

func logShippingTokenClient(m interface{}) *log_shipping_tokens.LogShippingTokensClient {
	var client *log_shipping_tokens.LogShippingTokensClient
	client, _ = log_shipping_tokens.New(m.(Config).apiToken, m.(Config).baseUrl)
	return client
}

func setLogShippingToken(d *schema.ResourceData, token *log_shipping_tokens.LogShippingToken) {
	d.Set(logShippingTokenId, token.Id)
	d.Set(logShippingTokenName, token.Name)
	d.Set(logShippingTokenEnabled, token.Enabled)
	d.Set(logShippingTokenToken, token.Token)
	d.Set(logShippingTokenUpdatedAt, token.UpdatedAt)
	d.Set(logShippingTokenUpdatedBy, token.UpdatedBy)
	d.Set(logShippingTokenCreatedAt, token.CreatedAt)
	d.Set(logShippingTokenCreatedBy, token.CreatedBy)
}
