package logzio

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/logzio/logzio_terraform_client/unified_projects"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnifiedProjectResourceSchema(t *testing.T) {
	r := resourceUnifiedProject()
	require.NoError(t, r.InternalValidate(nil, true))

	assert.True(t, r.Schema[unifiedProjectName].Required)
	// A rename must go through the rename endpoint: replacing the project would
	// delete every dashboard stored in it.
	assert.False(t, r.Schema[unifiedProjectName].ForceNew)
	assert.True(t, r.Schema[unifiedProjectDisplayName].Optional)
	assert.True(t, r.Schema[unifiedProjectDisplayName].Computed)
	assert.True(t, r.Schema[unifiedProjectDescription].Optional)
	assert.True(t, r.Schema[unifiedProjectDescription].Computed)
	assert.True(t, r.Schema[unifiedProjectFolderId].Computed)
	assert.NotNil(t, r.Importer)
}

func TestUnifiedProjectDisplayFields(t *testing.T) {
	tests := []struct {
		name                string
		project             *unified_projects.ProjectSummary
		expectedName        string
		expectedDescription string
	}{
		{
			name:                "values from Perses document",
			project:             unifiedProjectSummary("project-1", "metadata-name", "display name", "description"),
			expectedName:        "display name",
			expectedDescription: "description",
		},
		{
			name:                "display name falls back to API name",
			project:             &unified_projects.ProjectSummary{Id: "project-1", Name: "fallback"},
			expectedName:        "fallback",
			expectedDescription: "",
		},
		{name: "nil project", project: nil, expectedName: "", expectedDescription: ""},
		{
			name:                "malformed display block",
			project:             &unified_projects.ProjectSummary{Name: "fallback", Doc: map[string]any{"spec": map[string]any{"display": "invalid"}}},
			expectedName:        "fallback",
			expectedDescription: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedName, getUnifiedProjectDisplayName(tt.project))
			assert.Equal(t, tt.expectedDescription, getUnifiedProjectDescription(tt.project))
		})
	}
}

func TestSetUnifiedProjectSharedState(t *testing.T) {
	project := unifiedProjectSummary("project-1", "metadata-name", "Display name", "A description")

	tests := []struct {
		name   string
		schema map[string]*schema.Schema
		input  map[string]any
	}{
		{name: "resource schema", schema: resourceUnifiedProject().Schema, input: map[string]any{unifiedProjectName: "metadata-name"}},
		{name: "data source schema", schema: dataSourceUnifiedProject().Schema, input: map[string]any{unifiedProjectName: "metadata-name"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, tt.schema, tt.input)
			require.NoError(t, setUnifiedProject(d, project))
			assert.Equal(t, "metadata-name", d.Get(unifiedProjectName))
			assert.Equal(t, "Display name", d.Get(unifiedProjectDisplayName))
			assert.Equal(t, "A description", d.Get(unifiedProjectDescription))
		})
	}
}

func TestUnifiedProjectSchemaSpecificIdentityWrites(t *testing.T) {
	project := unifiedProjectSummary("project-1", "metadata-name", "Display name", "A description")

	resourceData := schema.TestResourceDataRaw(t, resourceUnifiedProject().Schema, map[string]any{unifiedProjectName: "metadata-name"})
	require.NoError(t, setUnifiedProject(resourceData, project))
	require.NoError(t, resourceData.Set(unifiedProjectFolderId, project.Id))
	resourceData.SetId(project.Id)
	assert.Equal(t, "project-1", resourceData.Get(unifiedProjectFolderId))
	assert.Equal(t, "project-1", resourceData.Id())

	dataSourceData := schema.TestResourceDataRaw(t, dataSourceUnifiedProject().Schema, map[string]any{unifiedProjectName: "metadata-name"})
	require.NoError(t, setUnifiedProject(dataSourceData, project))
	require.NoError(t, dataSourceData.Set("id", project.Id))
	dataSourceData.SetId(project.Id)
	assert.Equal(t, "project-1", dataSourceData.Get("id"))
	assert.Equal(t, "project-1", dataSourceData.Id())
}

func TestSetUnifiedProjectPropagatesStateErrors(t *testing.T) {
	skipIfSetErrorsPanic(t)

	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		unifiedProjectName: {
			Type:     schema.TypeInt,
			Optional: true,
		},
		unifiedProjectDisplayName: {
			Type:     schema.TypeString,
			Computed: true,
		},
		unifiedProjectDescription: {
			Type:     schema.TypeString,
			Computed: true,
		},
	}, map[string]any{})

	err := setUnifiedProject(d, unifiedProjectSummary("project-1", "metadata-name", "Display name", "A description"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to set name")
}

// A changed metadata.name is carried by the ordinary PUT, which rewrites it in
// place and keeps the folder id — verified live against api.logz.io on
// 2026-09-02, where the dashboards in the folder survived the rename. The
// API's /rename endpoint is deliberately not used: it sets spec.display.name
// and leaves metadata.name untouched.
func TestUnifiedProjectUpdateRenamesThroughPut(t *testing.T) {
	var calls []string
	var putBody map[string]any
	project := unifiedProjectSummary("project-1", "old-name", "Display name", "A description")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, r.Method+" "+r.URL.Path)
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodPut {
			require.NoError(t, json.NewDecoder(r.Body).Decode(&putBody))
			metadata, _ := putBody["metadata"].(map[string]any)
			spec, _ := putBody["spec"].(map[string]any)
			display, _ := spec["display"].(map[string]any)
			project = unifiedProjectSummary("project-1", stringValue(metadata["name"]), stringValue(display["name"]), stringValue(display["description"]))
		}
		fmt.Fprint(w, mustMarshalJSON(t, project))
	}))
	defer server.Close()

	d, err := schema.InternalMap(resourceUnifiedProject().Schema).Data(
		&terraform.InstanceState{
			ID: "project-1",
			Attributes: map[string]string{
				unifiedProjectName:        "old-name",
				unifiedProjectDisplayName: "Display name",
				unifiedProjectDescription: "A description",
			},
		},
		&terraform.InstanceDiff{
			Attributes: map[string]*terraform.ResourceAttrDiff{
				unifiedProjectName: {Old: "old-name", New: "new-name"},
			},
		},
	)
	require.NoError(t, err)

	diags := resourceUnifiedProjectUpdate(context.Background(), d, Config{apiToken: "token", baseUrl: server.URL})
	require.False(t, diags.HasError(), diagnosticsString(diags))

	metadata, _ := putBody["metadata"].(map[string]any)
	assert.Equal(t, "new-name", metadata["name"], "the PUT carries the new metadata.name")
	assert.NotContains(t, calls, "PUT /perses-public/api/v1/projects/project-1/rename",
		"/rename sets the display name, not metadata.name, so an identity rename must not use it")
	assert.Equal(t, "new-name", d.Get(unifiedProjectName))
}

func TestUnifiedProjectResourceLifecycle(t *testing.T) {
	var methods []string
	project := unifiedProjectSummary("project-1", "metadata-name", "metadata-name", "Initial description")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		methods = append(methods, r.Method)
		if r.Method == http.MethodPost {
			assert.Equal(t, "/perses-public/api/v1/projects", r.URL.Path)
		} else {
			assert.Equal(t, "/perses-public/api/v1/projects/project-1", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		case http.MethodPost, http.MethodPut:
			var request map[string]any
			assert.NoError(t, json.NewDecoder(r.Body).Decode(&request))
			metadata, _ := request["metadata"].(map[string]any)
			spec, _ := request["spec"].(map[string]any)
			display, _ := spec["display"].(map[string]any)
			project = unifiedProjectSummary("project-1", stringValue(metadata["name"]), stringValue(display["name"]), stringValue(display["description"]))
			fmt.Fprint(w, mustMarshalJSON(t, project))
		case http.MethodGet:
			fmt.Fprint(w, mustMarshalJSON(t, project))
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
	defer server.Close()

	d := schema.TestResourceDataRaw(t, resourceUnifiedProject().Schema, map[string]any{
		unifiedProjectName:        "metadata-name",
		unifiedProjectDescription: "Initial description",
	})
	config := Config{apiToken: "token", baseUrl: server.URL}

	diags := resourceUnifiedProjectCreate(context.Background(), d, config)
	require.False(t, diags.HasError(), diagnosticsString(diags))
	assert.Equal(t, "project-1", d.Id())
	assert.Equal(t, "project-1", d.Get(unifiedProjectFolderId))
	assert.Equal(t, "metadata-name", d.Get(unifiedProjectDisplayName))

	require.NoError(t, d.Set(unifiedProjectDisplayName, "Updated display"))
	require.NoError(t, d.Set(unifiedProjectDescription, "Updated description"))
	diags = resourceUnifiedProjectUpdate(context.Background(), d, config)
	require.False(t, diags.HasError(), diagnosticsString(diags))
	assert.Equal(t, "Updated display", d.Get(unifiedProjectDisplayName))
	assert.Equal(t, "Updated description", d.Get(unifiedProjectDescription))

	diags = resourceUnifiedProjectDelete(context.Background(), d, config)
	require.False(t, diags.HasError(), diagnosticsString(diags))
	assert.Empty(t, d.Id())
	assert.Equal(t, []string{http.MethodPost, http.MethodGet, http.MethodPut, http.MethodGet, http.MethodDelete}, methods)
}

func TestUnifiedProjectReadClearsMissingResource(t *testing.T) {
	server := httptest.NewServer(http.NotFoundHandler())
	defer server.Close()

	d := schema.TestResourceDataRaw(t, resourceUnifiedProject().Schema, map[string]any{})
	d.SetId("missing")
	diags := resourceUnifiedProjectRead(context.Background(), d, Config{apiToken: "token", baseUrl: server.URL})

	assert.False(t, diags.HasError(), diagnosticsString(diags))
	assert.Empty(t, d.Id())
}

func unifiedProjectSummary(id, metadataName, displayName, description string) *unified_projects.ProjectSummary {
	return &unified_projects.ProjectSummary{
		Id:   id,
		Name: displayName,
		Doc: map[string]any{
			"kind":     "Project",
			"metadata": map[string]any{"name": metadataName},
			"spec": map[string]any{
				"display": map[string]any{"name": displayName, "description": description},
			},
		},
	}
}

func mustMarshalJSON(t *testing.T, value any) string {
	t.Helper()
	encoded, err := json.Marshal(value)
	require.NoError(t, err)
	return string(encoded)
}

func stringValue(value any) string {
	result, _ := value.(string)
	return result
}

func diagnosticsString(diags diag.Diagnostics) string {
	return fmt.Sprint(diags)
}
