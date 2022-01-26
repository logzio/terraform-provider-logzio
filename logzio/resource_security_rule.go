package logzio

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"time"
)

const (
	securityRuleId          = "security_rule_id"
	securityRuleTitle       = "title"
	securityRuleDescription = "description"
	securityRuleTags        = "tags"
	SecurityRuleEmails      = "emails"
)

func resourceSecurityRule() *schema.Resource {
	return &schema.Resource{
		//Create: resourceEndpointCreate,
		//Read:   resourceEndpointRead,
		//Update: resourceEndpointUpdate,
		Delete: resourceEndpointDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			securityRuleId: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			securityRuleTitle: {
				Type:     schema.TypeString,
				Required: true,
			},
			securityRuleDescription: {
				Type:     schema.TypeString,
				Optional: true,
			},
			securityRuleTags: {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			SecurityRuleEmails: {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Second),
			Update: schema.DefaultTimeout(5 * time.Second),
			Delete: schema.DefaultTimeout(5 * time.Second),
		},
	}
}
