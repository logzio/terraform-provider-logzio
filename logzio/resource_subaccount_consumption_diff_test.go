package logzio

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/logzio/logzio_terraform_client/sub_accounts"
)

func TestSuppressConsumptionMaxDailyGBDiff(t *testing.T) {
	cases := []struct {
		name      string
		old       string
		new       string
		softLimit interface{}
		expected  bool
	}{
		{"consumption sub account, API stored 0.001", "0.0010000000474974513", "1", 1.0, true},
		{"consumption sub account, a real value changed", "2", "1", 1.0, false},
		{"subscription sub account stored 0.001", "0.0010000000474974513", "1", nil, false},
		{"not a number", "", "1", 1.0, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			raw := map[string]interface{}{}
			if tc.softLimit != nil {
				raw[subAccountSoftLimitGB] = tc.softLimit
			}
			d := schema.TestResourceDataRaw(t, resourceSubAccount().Schema, raw)

			if got := suppressConsumptionMaxDailyGBDiff(subAccountMaxDailyGB, tc.old, tc.new, d); got != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func TestIsSubAccountCacheSynced(t *testing.T) {
	requested := sub_accounts.CreateOrUpdateSubAccount{
		SharingObjectsAccounts: []int32{1, 2},
		UtilizationSettings:    sub_accounts.AccountUtilizationSettingsCreateOrUpdate{FrequencyMinutes: 3, UtilizationEnabled: "true"},
	}
	synced := sub_accounts.SubAccount{
		SharingObjectsAccounts: []sub_accounts.SharingAccount{{AccountId: 2}, {AccountId: 1}},
		UtilizationSettings:    sub_accounts.AccountUtilizationSettings{FrequencyMinutes: 3, UtilizationEnabled: true},
	}
	cacheDefaults := sub_accounts.SubAccount{}

	if !isSubAccountCacheSynced(requested, &synced) {
		t.Error("expected a sub account reporting the requested values to be synced")
	}
	if isSubAccountCacheSynced(requested, &cacheDefaults) {
		t.Error("expected the API's cache defaults (no sharing, utilization off) not to be synced")
	}

	partial := synced
	partial.SharingObjectsAccounts = []sub_accounts.SharingAccount{{AccountId: 1}}
	if isSubAccountCacheSynced(requested, &partial) {
		t.Error("expected a missing sharing account not to be synced")
	}

	disabled := sub_accounts.CreateOrUpdateSubAccount{
		SharingObjectsAccounts: []int32{},
		UtilizationSettings:    sub_accounts.AccountUtilizationSettingsCreateOrUpdate{UtilizationEnabled: "false"},
	}
	if !isSubAccountCacheSynced(disabled, &cacheDefaults) {
		t.Error("expected a sub account created without sharing or utilization to be synced right away")
	}
}
