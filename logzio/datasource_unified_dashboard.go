package logzio

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/logzio/logzio_terraform_client/unified_projects"
)

func dataSourceUnifiedProject() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUnifiedProjectRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"id", unifiedProjectName},
				Description:  "The unique identifier (GUID) of the unified project folder.",
			},
			unifiedProjectName: {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"id", unifiedProjectName},
				Description:  "The unique metadata name of the project.",
			},
			unifiedProjectDisplayName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The human-readable display title of the project.",
			},
			unifiedProjectDescription: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description text of the project.",
			},
		},
	}
}

func dataSourceUnifiedProjectRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := unifiedProjectClient(m)
	if err != nil {
		return diag.FromErr(err)
	}

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk(unifiedProjectName)

	var project *unified_projects.ProjectSummary

	if idOk {
		project, err = client.GetProject(id.(string))
		if err != nil {
			return diag.FromErr(err)
		}
	} else if nameOk {
		searchRes, err := client.SearchProjects(unified_projects.SearchProjectsRequest{
			Filter: &unified_projects.SearchProjectsFilter{
				SearchTerm: name.(string),
			},
			Pagination: &unified_projects.SearchProjectsPagination{
				PageSize: 100, // Ensure wide result coverage for search by name
			},
		})
		if err != nil {
			return diag.FromErr(err)
		}

		for i := range searchRes.Results {
			p := &searchRes.Results[i].Project
			if p.MetadataName() == name.(string) || p.Name == name.(string) {
				project = p
				break
			}
		}

		if project == nil {
			return diag.Errorf("unified project with name %q not found", name.(string))
		}
	}

	d.SetId(project.Id)
	d.Set(unifiedProjectName, project.MetadataName())
	d.Set(unifiedProjectDisplayName, getUnifiedProjectDisplayName(project))
	d.Set(unifiedProjectDescription, getUnifiedProjectDescription(project))

	return nil
}
