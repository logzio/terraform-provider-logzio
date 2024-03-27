package logzio

import (
	"context"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/logzio/logzio_terraform_client/metrics_accounts"
	"strconv"
	"strings"
	"time"
)

func dataSourceMetricsAccount() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMetricsAccountReadWrapper,
		Schema: map[string]*schema.Schema{
			subAccountId: {
				Type:     schema.TypeInt,
				Optional: true,
			},
			metricsAccountName: {
				Type:     schema.TypeString,
				Optional: true,
			},
			metricsAccountEmail: {
				Type:     schema.TypeString,
				Optional: true,
			},
			metricsAccountToken: {
				Type:     schema.TypeString,
				Optional: true,
			},
			metricsAccountPlanUts: {
				Type:     schema.TypeInt,
				Optional: true,
			},
			metricsAccountAuthorizedAccounts: {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Optional: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Second),
		},
	}
}

func dataSourceMetricsAccountReadWrapper(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var err error
	readErr := retry.Do(
		func() error {
			err = dataSourceMetricsAccountRead(d, m)
			if err != nil {
				return err
			}

			return nil
		},
		retry.RetryIf(
			func(err error) bool {
				if err != nil {
					if strings.Contains(err.Error(), "failed with missing metrics account") ||
						strings.Contains(err.Error(), "failed with status code 500") {
						return true
					}
				}
				return false
			}),
		retry.DelayType(retry.BackOffDelay),
		retry.Attempts(15),
	)

	if readErr != nil {
		return diag.FromErr(err)
	}

	return nil
}

func dataSourceMetricsAccountRead(d *schema.ResourceData, m interface{}) error {
	var client *metrics_accounts.MetricsAccountClient
	client, _ = metrics_accounts.New(m.(Config).apiToken, m.(Config).baseUrl)

	accountId, ok := d.GetOk(metricsAccountId)
	if ok {
		subAccount, err := client.GetMetricsAccount(int64(accountId.(int)))
		if err != nil {
			return err
		}
		d.SetId(strconv.FormatInt(int64(accountId.(int)), 10))
		setMetricsAccount(d, subAccount)
		return nil
	}

	accountName, ok := d.GetOk(metricsAccountName)
	if ok {
		subAccounts, err := client.ListMetricsAccounts()
		if err != nil {
			return err
		}

		for _, account := range subAccounts {
			if account.AccountName == accountName.(string) {
				d.SetId(strconv.FormatInt(int64(account.Id), 10))
				setMetricsAccount(d, &account)
				return nil
			}
		}
	}

	return fmt.Errorf("couldn't find metrics account with specified attributes")
}
