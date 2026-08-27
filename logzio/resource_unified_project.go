package logzio

import (
	"context"
	"fmt"
	"strings"

	"github.com/avast/retry-go"
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
	displayName := unifiedProjectDisplayNameForWrite(d)
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
		if strings.Contains(err.Error(), "missing unified project") {
			d.SetId("") // Unsets state if the remote object was deleted
			return diag.Diagnostics{}
		}
		return diag.FromErr(err)
	}

	if err = validateUnifiedProjectIdentity(project, id); err != nil {
		return diag.FromErr(err)
	}
	if err = setUnifiedProject(d, project); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set(unifiedProjectFolderId, project.Id); err != nil {
		return diag.Errorf("failed to set %s: %s", unifiedProjectFolderId, err)
	}
	d.SetId(project.Id)
	return nil
}

func resourceUnifiedProjectUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := unifiedProjectClient(m)
	if err != nil {
		return diag.FromErr(err)
	}

	projectName := d.Get(unifiedProjectName).(string)
	expectedDisplayName := unifiedProjectDisplayNameForWrite(d)
	expectedDescription := d.Get(unifiedProjectDescription).(string)

	_, err = client.UpdateProject(d.Id(), unified_projects.UpdateProjectRequest{
		Name:        projectName,
		DisplayName: expectedDisplayName,
		Description: expectedDescription,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	var diagRet diag.Diagnostics
	readErr := retry.Do(func() error {
		diagRet = resourceUnifiedProjectRead(ctx, d, m)
		if diagRet.HasError() {
			return fmt.Errorf("received error from read unified project")
		}

		currentDisplayName := d.Get(unifiedProjectDisplayName).(string)
		currentDescription := d.Get(unifiedProjectDescription).(string)

		if currentDisplayName != expectedDisplayName || currentDescription != expectedDescription {
			return fmt.Errorf("unified project has not finished updating yet")
		}

		return nil
	},
		retry.RetryIf(func(err error) bool {
			return err != nil
		}),
		retry.DelayType(retry.BackOffDelay),
		retry.Attempts(unifiedProjectRetryAttempts),
	)

	if readErr != nil {
		tflog.Error(ctx, "could not update unified project schema")
		return diagRet
	}

	return nil
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

func setUnifiedProject(d *schema.ResourceData, result *unified_projects.ProjectSummary) error {
	if result == nil {
		return fmt.Errorf("cannot set unified project state from an empty response")
	}
	if err := d.Set(unifiedProjectName, result.MetadataName()); err != nil {
		return fmt.Errorf("failed to set %s: %w", unifiedProjectName, err)
	}
	if err := d.Set(unifiedProjectDisplayName, getUnifiedProjectDisplayName(result)); err != nil {
		return fmt.Errorf("failed to set %s: %w", unifiedProjectDisplayName, err)
	}
	if err := d.Set(unifiedProjectDescription, getUnifiedProjectDescription(result)); err != nil {
		return fmt.Errorf("failed to set %s: %w", unifiedProjectDescription, err)
	}

	return nil
}

func validateUnifiedProject(result *unified_projects.ProjectSummary) error {
	if result == nil {
		return fmt.Errorf("unified project response was empty")
	}
	if result.Id == "" {
		return fmt.Errorf("unified project response contained no project id")
	}
	if result.MetadataName() == "" {
		return fmt.Errorf("unified project response contained no Perses metadata.name")
	}
	return nil
}

func validateUnifiedProjectIdentity(result *unified_projects.ProjectSummary, requestedId string) error {
	if err := validateUnifiedProject(result); err != nil {
		return err
	}
	if result.Id != requestedId {
		return fmt.Errorf("unified project response id %q does not match requested id %q", result.Id, requestedId)
	}
	return nil
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

func unifiedProjectDisplayNameForWrite(d *schema.ResourceData) string {
	displayName := d.Get(unifiedProjectDisplayName).(string)
	if displayName != "" {
		return displayName
	}
	return d.Get(unifiedProjectName).(string)
}
