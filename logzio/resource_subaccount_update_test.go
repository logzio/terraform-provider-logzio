package logzio

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestGetUpdateSubAccountFromSchemaWarmRetention(t *testing.T) {
	t.Run("not in the configuration, the update sends the warm retention from the state", func(t *testing.T) {
		d := resourceSubAccount().Data(&terraform.InstanceState{
			ID:         "1",
			Attributes: map[string]string{subAccountsSnapSearchRetentionDays: "3"},
		})

		updateSubAccount := getUpdateSubAccountFromSchema(d)
		if updateSubAccount.SnapSearchRetentionDays == nil || *updateSubAccount.SnapSearchRetentionDays != 3 {
			t.Errorf("expected snapSearchRetentionDays 3, got %v", updateSubAccount.SnapSearchRetentionDays)
		}
	})

	t.Run("not in the configuration and no warm tier, the update sends no value", func(t *testing.T) {
		d := resourceSubAccount().Data(&terraform.InstanceState{
			ID:         "1",
			Attributes: map[string]string{subAccountsSnapSearchRetentionDays: "0"},
		})

		updateSubAccount := getUpdateSubAccountFromSchema(d)
		if updateSubAccount.SnapSearchRetentionDays != nil {
			t.Errorf("expected no snapSearchRetentionDays, got %d", *updateSubAccount.SnapSearchRetentionDays)
		}
	})

	t.Run("in the configuration, the update sends the configured value", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, resourceSubAccount().Schema, map[string]interface{}{
			subAccountsSnapSearchRetentionDays: 5,
		})

		updateSubAccount := getUpdateSubAccountFromSchema(d)
		if updateSubAccount.SnapSearchRetentionDays == nil || *updateSubAccount.SnapSearchRetentionDays != 5 {
			t.Errorf("expected snapSearchRetentionDays 5, got %v", updateSubAccount.SnapSearchRetentionDays)
		}
	})
}
