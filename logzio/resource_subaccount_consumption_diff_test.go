package logzio

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
