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

func TestUnifiedDashboardDataSourceSchema(t *testing.T) {
	ds := dataSourceUnifiedDashboard()
	require.NoError(t, ds.InternalValidate(nil, false))

	assert.True(t, ds.Schema[unifiedDashboardFolderId].Required)
	assert.True(t, ds.Schema[unifiedDashboardUid].Required)
	assert.True(t, ds.Schema[unifiedDashboardJson].Computed)
	assert.True(t, ds.Schema[unifiedDashboardName].Computed)
	assert.True(t, ds.Schema[unifiedDashboardVersion].Computed)
}

func TestUnifiedDashboardDataSourceRead(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/perses-public/api/v1/projects/project-1/dashboards/dashboard-1", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
  "id": "version-row-4",
  "uid": "dashboard-1",
  "projectId": "project-1",
  "version": 4,
  "doc": {
    "kind": "Dashboard",
    "metadata": {
      "name": "cpu-usage",
      "project": "project-1",
      "updatedAt": "server-value"
    },
    "spec": {
      "duration": "6h"
    }
  }
}`)
	}))
	defer server.Close()

	d := schema.TestResourceDataRaw(t, dataSourceUnifiedDashboard().Schema, map[string]any{
		unifiedDashboardFolderId: "project-1",
		unifiedDashboardUid:      "dashboard-1",
	})
	diags := dataSourceUnifiedDashboardRead(context.Background(), d, Config{apiToken: "token", baseUrl: server.URL})

	require.False(t, diags.HasError(), diagnosticsString(diags))
	assert.Equal(t, "project-1/dashboard-1", d.Id())
	assert.Equal(t, "project-1", d.Get(unifiedDashboardFolderId))
	assert.Equal(t, "dashboard-1", d.Get(unifiedDashboardUid))
	assert.Equal(t, "cpu-usage", d.Get(unifiedDashboardName))
	assert.Equal(t, 4, d.Get(unifiedDashboardVersion))
	assert.Equal(t, `{"kind":"Dashboard","metadata":{"name":"cpu-usage"},"spec":{"duration":"6h"}}`, d.Get(unifiedDashboardJson).(string))
}

func TestUnifiedDashboardDataSourceReadError(t *testing.T) {
	server := httptest.NewServer(http.NotFoundHandler())
	defer server.Close()

	d := schema.TestResourceDataRaw(t, dataSourceUnifiedDashboard().Schema, map[string]any{
		unifiedDashboardFolderId: "project-1",
		unifiedDashboardUid:      "missing",
	})
	diags := dataSourceUnifiedDashboardRead(context.Background(), d, Config{apiToken: "token", baseUrl: server.URL})

	require.True(t, diags.HasError())
	assert.Contains(t, diagnosticsString(diags), "missing unified dashboard")
	assert.Empty(t, d.Id())
}

func TestUnifiedDashboardDataSourceRejectsResponseIdentityMismatch(t *testing.T) {
	tests := []struct {
		name          string
		response      string
		expectedError string
	}{
		{
			name:          "dashboard uid",
			response:      `{"uid":"other-dashboard","projectId":"project-1","doc":{"metadata":{"name":"cpu"}}}`,
			expectedError: `response uid "other-dashboard" does not match requested dashboard_uid "dashboard-1"`,
		},
		{
			name:          "project id",
			response:      `{"uid":"dashboard-1","projectId":"other-project","doc":{"metadata":{"name":"cpu"}}}`,
			expectedError: `response project id "other-project" does not match requested folder_id "project-1"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, tt.response)
			}))
			defer server.Close()

			d := schema.TestResourceDataRaw(t, dataSourceUnifiedDashboard().Schema, map[string]any{
				unifiedDashboardFolderId: "project-1",
				unifiedDashboardUid:      "dashboard-1",
			})
			diags := dataSourceUnifiedDashboardRead(context.Background(), d, Config{apiToken: "token", baseUrl: server.URL})

			require.True(t, diags.HasError())
			assert.Contains(t, diagnosticsString(diags), tt.expectedError)
			assert.Empty(t, d.Id())
		})
	}
}

func TestProviderRegistersUnifiedDataSources(t *testing.T) {
	provider := Provider()
	dashboard, dashboardOk := provider.DataSourcesMap[resourceUnifiedDashboardType]
	project, projectOk := provider.DataSourcesMap[resourceUnifiedProjectType]

	require.True(t, dashboardOk)
	require.True(t, projectOk)
	require.NotNil(t, dashboard)
	require.NotNil(t, project)
	assert.Contains(t, dashboard.Schema, unifiedDashboardFolderId)
	assert.Contains(t, dashboard.Schema, unifiedDashboardUid)
	assert.NotContains(t, dashboard.Schema, "id")
	assert.Contains(t, project.Schema, "id")
	assert.Contains(t, project.Schema, unifiedProjectName)
	assert.NotContains(t, project.Schema, unifiedDashboardFolderId)
	assert.NotContains(t, project.Schema, unifiedDashboardUid)
}
