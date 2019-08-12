package logzio

import "github.com/hashicorp/terraform/helper/schema"

const (
	accountId              string = "accountId"   //required
	email                  string = "email"       //required
	accountName            string = "accountName" //required
	maxDailyGB             string = "maxDailyGB"
	retentionDays          string = "retentionDays" //required
	accessible             string = "accessible"
	searchable             string = "searchable"
	sharingObjectsAccounts string = "sharingObjectsAccounts" //required
	docSizeSetting         string = "docSizeSetting"
	utilizationSettings    string = "utilizationSettings"
	frequencyMinutes       string = "frequencyMinutes"
	utilizationEnabled     string = "utilizationEnabled"
	accountToken           string = "accountToken"
	dailyUsagesList        string = "dailyUsagesList"
)

func resourceSubAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceSubAccountCreate,
		Read:   resourceSubAccountRead,
		Update: resourceSubAccountUpdate,
		Delete: resourceSubAccountDelete,

		Schema: map[string]*schema.Schema{
			email: {
				Type:     schema.TypeString,
				Required: true,
			},
			accountName: {
				Type:     schema.TypeString,
				Required: true,
			},
			maxDailyGB: {
				Type:     schema.TypeFloat,
				Optional: true,
				Default:  0.0,
			},
			retentionDays: {
				Type:     schema.TypeInt,
				Required: true,
			},
			accessible: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			searchable: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			sharingObjectsAccounts: {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Required: true,
			},
			docSizeSetting: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			utilizationSettings: {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						frequencyMinutes: {
							Type:     schema.TypeInt,
							Required: true,
						},
						utilizationEnabled: {
							Type:     schema.TypeBool,
							Required: true,
						},
					},
				},
				Optional: true,
			},
		},
	}
}

func resourceSubAccountCreate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceSubAccountRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceSubAccountUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceSubAccountDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
