package logzio

import (
	"encoding/json"
  "strings"
	//"fmt"
	//	"strconv"
	//
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/logzio/logzio_terraform_client/grafana_objects"
)

const (
	dashboardId        string = "dashboard_id"
	dashboardUid       string = "dashboard_uid"
	dashboardStarred   string = "starred"
	dashboardUrl       string = "url"
	dashboardFolderId  string = "folder_id"
	dashboardFolderUid string = "folder_uid"
	dashboardSlug      string = "slug"
	dashboardJson      string = "dashboard_json"
	dashboardMessage   string = "message"
)

/**
 * the endpoint resource schema, what terraform uses to parse and read the template
 */
func resourceDashboard() *schema.Resource {
	return &schema.Resource{
		Create: resourceDashboardCreate,
		Read:   resourceDashboardRead,
		Update: resourceDashboardUpdate,
		Delete: resourceDashboardDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			dashboardJson: {
				Type:     schema.TypeString,
				Required: true,
			},
			dashboardFolderId: {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			dashboardFolderUid: {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			dashboardMessage: {
				Type:     schema.TypeString,
				Optional: true,
			},
			dashboardUrl: {
				Type:     schema.TypeString,
				Computed: true,
			},
			dashboardId: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			dashboardUid: {
				Type:     schema.TypeString,
				Computed: true,
			},
			dashboardSlug: {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dashboardClient(m interface{}) (*grafana_objects.GrafanaObjectsClient, error) {
	var client *grafana_objects.GrafanaObjectsClient
	client, err := grafana_objects.New(m.(Config).apiToken, m.(Config).baseUrl)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func resourceDashboardRead(d *schema.ResourceData, m interface{}) error {
	client, err := dashboardClient(m)
	if err != nil {
		return err
	}

	result, err := client.Get(d.Id())
	if err != nil {
    if strings.Contains(err.Error(), "Dashboard not found") {
      d.SetId("")
      return nil
    }
		return err
	}

	d.SetId(result.Dashboard.Uid)
	d.Set(dashboardId, result.Dashboard.Id)
	d.Set(dashboardUid, result.Dashboard.Uid)
	d.Set(dashboardStarred, result.Meta.IsStarred)
	d.Set(dashboardUrl, result.Meta.Url)
	d.Set(dashboardFolderId, result.Meta.FolderId)
	d.Set(dashboardFolderUid, result.Meta.FolderUid)
	d.Set(dashboardSlug, result.Meta.Slug)

	if result.Dashboard.Timepicker.Enable == false {
		result.Dashboard.Timepicker = &grafana_objects.Timepicker{}
	}

	dashboardObject, err := json.Marshal(result.Dashboard)
	if err != nil {
		return err
	}

	d.Set(dashboardJson, string(dashboardObject))

	return nil
}

func resourceDashboardCreate(d *schema.ResourceData, m interface{}) error {
	client, err := dashboardClient(m)
	if err != nil {
		return err
	}

	var dashboardObject grafana_objects.DashboardObject
	err = json.Unmarshal([]byte(d.Get(dashboardJson).(string)), &dashboardObject)
	if err != nil {
		return err
	}

  dashboardObject.Id = 0
	jsonPayload := &grafana_objects.CreateUpdatePayload{
		FolderId:  d.Get(dashboardFolderId).(int),
		FolderUid: d.Get(dashboardFolderUid).(string),
		Message:   d.Get(dashboardMessage).(string),
		Overwrite: true,
		Dashboard: dashboardObject,
	}
  

	result, err := client.CreateUpdate(*jsonPayload)
	if err != nil {
		return err
	}
	d.SetId(result.Uid)
	d.Set(dashboardUrl, result.Url)
	d.Set(dashboardSlug, result.Slug)
	d.Set(dashboardId, result.Id)
	d.Set(dashboardUid, result.Uid)

	return nil
}

func resourceDashboardUpdate(d *schema.ResourceData, m interface{}) error {
	client, err := dashboardClient(m)
	if err != nil {
		return err
	}

	var dashboardObject grafana_objects.DashboardObject
	err = json.Unmarshal([]byte(d.Get(dashboardJson).(string)), &dashboardObject)
	if err != nil {
		return err
	}
	dashboardObject.Id = d.Get(dashboardId).(int)

	jsonPayload := &grafana_objects.CreateUpdatePayload{
		FolderId:  d.Get(dashboardFolderId).(int),
		FolderUid: d.Get(dashboardFolderUid).(string),
		Message:   d.Get(dashboardMessage).(string),
		Overwrite: true,
		Dashboard: dashboardObject,
	}

	result, err := client.CreateUpdate(*jsonPayload)
	if err != nil {
		return err
	}
	d.Set(dashboardUrl, result.Url)
	d.Set(dashboardSlug, result.Slug)
	d.Set(dashboardId, result.Id)
	d.Set(dashboardUid, result.Uid)

	return nil
}

func resourceDashboardDelete(d *schema.ResourceData, m interface{}) error {
	client, err := dashboardClient(m)
	if err != nil {
		return err
	}

	_, err = client.Delete(d.Id())
	if err != nil {
		return err
	}
	return nil
}
