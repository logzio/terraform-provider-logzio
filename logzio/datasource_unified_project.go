package logzio

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/go-cty/cty/gocty"
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

	id, idOk, err := configuredUnifiedProjectSelector(d, "id")
	if err != nil {
		return diag.FromErr(err)
	}
	name, nameOk, err := configuredUnifiedProjectSelector(d, unifiedProjectName)
	if err != nil {
		return diag.FromErr(err)
	}

	var project *unified_projects.ProjectSummary

	if idOk {
		requestedId := id
		project, err = client.GetProject(requestedId)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = validateUnifiedProjectIdentity(project, requestedId); err != nil {
			return diag.FromErr(err)
		}
	} else if nameOk {
		requestedName := name
		projects, err := client.ListProjects(false)
		if err != nil {
			return diag.FromErr(err)
		}

		for i := range projects {
			candidate := &projects[i].Project
			if candidate.MetadataName() != requestedName {
				continue
			}
			if project != nil {
				return diag.Errorf("multiple unified projects with metadata name %q found", requestedName)
			}
			project = candidate
		}

		if project == nil {
			return diag.Errorf("unified project with metadata name %q not found", requestedName)
		}
		if err = validateUnifiedProject(project); err != nil {
			return diag.Errorf("unified project with metadata name %q returned an invalid response: %s", requestedName, err)
		}
	} else {
		return diag.Errorf("exactly one of id or name must be configured")
	}

	if err = setUnifiedProject(d, project); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("id", project.Id); err != nil {
		return diag.Errorf("failed to set id: %s", err)
	}
	d.SetId(project.Id)

	return nil
}

func configuredUnifiedProjectSelector(d *schema.ResourceData, key string) (string, bool, error) {
	if d.GetRawConfig().IsNull() {
		selector, ok := d.GetOk(key)
		if !ok {
			return "", false, nil
		}
		return selector.(string), true, nil
	}

	value, diags := d.GetRawConfigAt(cty.GetAttrPath(key))
	if diags.HasError() {
		return "", false, fmt.Errorf("failed to read configured %s: %v", key, diags)
	}
	if value.IsNull() || !value.IsKnown() {
		return "", false, nil
	}

	var selector string
	if err := gocty.FromCtyValue(value, &selector); err != nil {
		return "", false, fmt.Errorf("invalid configured %s: %w", key, err)
	}
	return selector, selector != "", nil
}
