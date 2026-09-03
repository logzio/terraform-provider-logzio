package logzio

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnifiedProjectDataSourceSchema(t *testing.T) {
	ds := dataSourceUnifiedProject()
	require.NoError(t, ds.InternalValidate(nil, false))

	assert.True(t, ds.Schema["id"].Optional)
	assert.True(t, ds.Schema["id"].Computed)
	assert.Equal(t, []string{"id", unifiedProjectName}, ds.Schema["id"].ExactlyOneOf)
	assert.True(t, ds.Schema[unifiedProjectName].Optional)
	assert.True(t, ds.Schema[unifiedProjectName].Computed)
	assert.True(t, ds.Schema[unifiedProjectDisplayName].Computed)
	assert.True(t, ds.Schema[unifiedProjectDescription].Computed)
}

func TestUnifiedProjectDataSourceReadByID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/perses-public/api/v1/projects/project-1", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, mustMarshalJSON(t, unifiedProjectSummary("project-1", "metadata-name", "Display name", "Description")))
	}))
	defer server.Close()

	d := schema.TestResourceDataRaw(t, dataSourceUnifiedProject().Schema, map[string]any{"id": "project-1"})
	diags := dataSourceUnifiedProjectRead(context.Background(), d, Config{apiToken: "token", baseUrl: server.URL})

	require.False(t, diags.HasError(), diagnosticsString(diags))
	assert.Equal(t, "project-1", d.Id())
	assert.Equal(t, "project-1", d.Get("id"))
	assert.Equal(t, "metadata-name", d.Get(unifiedProjectName))
	assert.Equal(t, "Display name", d.Get(unifiedProjectDisplayName))
	assert.Equal(t, "Description", d.Get(unifiedProjectDescription))
}

func TestUnifiedProjectDataSourceRejectsDirectIDMismatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, mustMarshalJSON(t, unifiedProjectSummary("other-project", "metadata-name", "Display name", "Description")))
	}))
	defer server.Close()

	d := schema.TestResourceDataRaw(t, dataSourceUnifiedProject().Schema, map[string]any{"id": "project-1"})
	diags := dataSourceUnifiedProjectRead(context.Background(), d, Config{apiToken: "token", baseUrl: server.URL})

	require.True(t, diags.HasError())
	assert.Contains(t, diagnosticsString(diags), `response id "other-project" does not match requested id "project-1"`)
	assert.Empty(t, d.Id())
}

func TestUnifiedProjectDataSourceReadByName(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/perses-public/api/v1/projects", r.URL.Path)
		assert.Empty(t, r.URL.Query().Get("withDashboards"))

		response := []any{
			map[string]any{"project": unifiedProjectSummary("project-1", "other-project", "target-project", "display-name must not match")},
			map[string]any{"project": unifiedProjectSummary("project-2", "target-project", "Target", "Found")},
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, mustMarshalJSON(t, response))
	}))
	defer server.Close()

	d := schema.TestResourceDataRaw(t, dataSourceUnifiedProject().Schema, map[string]any{unifiedProjectName: "target-project"})
	diags := dataSourceUnifiedProjectRead(context.Background(), d, Config{apiToken: "token", baseUrl: server.URL})

	require.False(t, diags.HasError(), diagnosticsString(diags))
	assert.Equal(t, "project-2", d.Id())
	assert.Equal(t, "project-2", d.Get("id"))
	assert.Equal(t, "target-project", d.Get(unifiedProjectName))
	assert.Equal(t, "Target", d.Get(unifiedProjectDisplayName))
	assert.Equal(t, "Found", d.Get(unifiedProjectDescription))
}

func TestUnifiedProjectDataSourceReadByNameRejectsDuplicates(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, mustMarshalJSON(t, []any{
			map[string]any{"project": unifiedProjectSummary("project-1", "duplicate", "One", "")},
			map[string]any{"project": unifiedProjectSummary("project-2", "duplicate", "Two", "")},
		}))
	}))
	defer server.Close()

	d := schema.TestResourceDataRaw(t, dataSourceUnifiedProject().Schema, map[string]any{unifiedProjectName: "duplicate"})
	diags := dataSourceUnifiedProjectRead(context.Background(), d, Config{apiToken: "token", baseUrl: server.URL})

	require.True(t, diags.HasError())
	assert.Contains(t, diagnosticsString(diags), `multiple unified projects with metadata name "duplicate" found`)
	assert.Empty(t, d.Id())
}

func TestUnifiedProjectDataSourceReadByNameRejectsMissingProjectID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, mustMarshalJSON(t, []any{
			map[string]any{"project": unifiedProjectSummary("", "target-project", "Target", "")},
		}))
	}))
	defer server.Close()

	d := schema.TestResourceDataRaw(t, dataSourceUnifiedProject().Schema, map[string]any{unifiedProjectName: "target-project"})
	diags := dataSourceUnifiedProjectRead(context.Background(), d, Config{apiToken: "token", baseUrl: server.URL})

	require.True(t, diags.HasError())
	assert.Contains(t, diagnosticsString(diags), "response contained no project id")
	assert.Empty(t, d.Id())
}

func TestUnifiedProjectDataSourceReadRejectsMissingMetadataName(t *testing.T) {
	tests := []struct {
		name     string
		metadata string
	}{
		{name: "missing name", metadata: `{}`},
		{name: "malformed metadata", metadata: `"invalid"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintf(w, `{"id":"project-1","name":"Display name","doc":{"kind":"Project","metadata":%s,"spec":{}}}`, tt.metadata)
			}))
			defer server.Close()

			d := schema.TestResourceDataRaw(t, dataSourceUnifiedProject().Schema, map[string]any{"id": "project-1"})
			diags := dataSourceUnifiedProjectRead(context.Background(), d, Config{apiToken: "token", baseUrl: server.URL})

			require.True(t, diags.HasError())
			assert.Contains(t, diagnosticsString(diags), "response contained no Perses metadata.name")
			assert.Empty(t, d.Id())
		})
	}
}

func TestUnifiedProjectDataSourceReadByNameNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[]`)
	}))
	defer server.Close()

	d := schema.TestResourceDataRaw(t, dataSourceUnifiedProject().Schema, map[string]any{unifiedProjectName: "missing"})
	diags := dataSourceUnifiedProjectRead(context.Background(), d, Config{apiToken: "token", baseUrl: server.URL})

	require.True(t, diags.HasError())
	assert.Contains(t, diagnosticsString(diags), `unified project with metadata name "missing" not found`)
	assert.Empty(t, d.Id())
}
