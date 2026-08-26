package logzio

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/logzio/logzio_terraform_client/unified_projects"
)

const (
	unifiedProjectName        = "name"
	unifiedProjectDisplayName = "display_name"
	unifiedProjectDescription = "description"
	unifiedProjectFolderId    = "folder_id"

	unifiedProjectRetryAttempts = 8
)

func resourceUnifiedProject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUnifiedProjectCreate,
		ReadContext:   resourceUnifiedProjectRead,
		UpdateContext: resourceUnifiedProjectUpdate,
		DeleteContext: resourceUnifiedProjectDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			unifiedProjectName: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			unifiedProjectDescription: {
				Type:     schema.TypeString,
				Optional: true, // Changed from Computed to Optional
				Computed: true,
			},
			unifiedProjectDisplayName: {
				Type:     schema.TypeString,
				Optional: true, // Changed from Computed to Optional
				Computed: true,
			},
			unifiedProjectFolderId: {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func unifiedProjectClient(m any) (*unified_projects.ProjectsClient, error) {
	client, err := unified_projects.New(m.(Config).apiToken, m.(Config).baseUrl)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func resourceUnifiedProjectCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := unifiedProjectClient(m)
	if err != nil {
		return diag.FromErr(err)
	}

	projectName := d.Get(unifiedProjectName).(string)
	displayName := d.Get(unifiedProjectDisplayName).(string)
	description := d.Get(unifiedProjectDescription).(string)

	result, err := client.CreateProject(unified_projects.CreateProjectRequest{Name: projectName, DisplayName: displayName, Description: description})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(result.Id)
	return resourceUnifiedProjectRead(ctx, d, m)
}

func resourceUnifiedProjectRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := unifiedProjectClient(m)
	if err != nil {
		return diag.FromErr(err)
	}

	// Fetch using the primary resource ID tracked in Terraform state
	id := d.Id()
	if id == "" {
		return nil
	}

	project, err := client.GetProject(id)
	if err != nil {
		tflog.Error(ctx, err.Error())
		if strings.Contains(err.Error(), "missing unified folder") {
			d.SetId("") // Unsets state if the remote object was deleted
			return diag.Diagnostics{}
		}
		return diag.FromErr(err)
	}

	if err = setUnifiedProject(d, project); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceUnifiedProjectUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := unifiedProjectClient(m)
	if err != nil {
		return diag.FromErr(err)
	}

	projectName := d.Get(unifiedProjectName).(string)
	displayName := d.Get(unifiedProjectDisplayName).(string)
	description := d.Get(unifiedProjectDescription).(string)

	// Send update payload to the API
	_, err = client.UpdateProject(d.Id(), unified_projects.UpdateProjectRequest{
		Name:        projectName,
		DisplayName: displayName,
		Description: description,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	// Refresh Terraform state with the updated values
	return resourceUnifiedProjectRead(ctx, d, m)
}

func resourceUnifiedProjectDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := unifiedProjectClient(m)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = client.DeleteProject(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func composeUnifiedProjectId(folderId, uid string) string {
	return folderId + "/" + uid
}

func parseUnifiedProjectId(id string) (string, string, error) {
	parts := strings.SplitN(id, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("unexpected id format %q, expected folder_id/dashboard_uid", id)
	}
	return parts[0], parts[1], nil
}

// func unifiedProjectDocFromSchema(d *schema.ResourceData) (map[string]any, error) {
// 	var doc map[string]any
// 	if err := json.Unmarshal([]byte(d.Get(unifiedProjectJson).(string)), &doc); err != nil {
// 		return nil, err
// 	}
// 	return doc, nil
// }

func setUnifiedProject(d *schema.ResourceData, result *unified_projects.ProjectSummary) error {
	d.Set(unifiedProjectFolderId, result.Id)
	d.Set(unifiedProjectName, result.MetadataName())
	d.Set(unifiedProjectDisplayName, getUnifiedProjectDisplayName(result))
	d.Set(unifiedProjectDescription, getUnifiedProjectDescription(result))

	return nil
}

// func validateUnifiedProjectJson(config any, k string) ([]string, []error) {
// 	var doc map[string]any
// 	if err := json.Unmarshal([]byte(config.(string)), &doc); err != nil {
// 		return nil, []error{err}
// 	}
// 	if kind, _ := doc["kind"].(string); kind != "Project" {
// 		return nil, []error{fmt.Errorf("%s must be a Perses Project document with kind %q", k, "Project")}
// 	}
// 	if unifiedProjectNameFromDoc(doc) == "" {
// 		return nil, []error{fmt.Errorf("%s metadata.name is required", k)}
// 	}
// 	if _, ok := doc["spec"].(map[string]any); !ok {
// 		return nil, []error{fmt.Errorf("%s spec is required", k)}
// 	}
// 	return nil, nil
// }

func unifiedProjectNameFromDoc(doc map[string]any) string {
	metadata, ok := doc["metadata"].(map[string]any)
	if !ok {
		return ""
	}
	name, _ := metadata["name"].(string)
	return name
}

func unifiedProjectNameChanged(currentName, dashboardJson string) bool {
	var doc map[string]any
	if err := json.Unmarshal([]byte(dashboardJson), &doc); err != nil {
		return true
	}
	name := unifiedProjectNameFromDoc(doc)
	return name != "" && currentName != "" && name != currentName
}

func getUnifiedProjectDescription(result *unified_projects.ProjectSummary) string {
	if result == nil || result.Doc == nil {
		return ""
	}
	if spec, ok := result.Doc["spec"].(map[string]any); ok {
		if display, ok := spec["display"].(map[string]any); ok {
			if desc, ok := display["description"].(string); ok {
				return desc
			}
		}
	}
	return ""
}

func getUnifiedProjectDisplayName(result *unified_projects.ProjectSummary) string {
	if result == nil {
		return ""
	}
	if result.Doc != nil {
		if spec, ok := result.Doc["spec"].(map[string]any); ok {
			if display, ok := spec["display"].(map[string]any); ok {
				if name, ok := display["name"].(string); ok {
					return name
				}
			}
		}
	}
	return result.Name
}
