package logzio

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/logzio/logzio_terraform_client/log_shipping_tokens"
	"strconv"
)

func dataSourceLogShippingToken() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLogShippingTokenRead,
		Schema: map[string]*schema.Schema{
			logShippingTokenTokenId: {
				Type:     schema.TypeInt,
				Optional: true,
			},
			logShippingTokenName: {
				Type:     schema.TypeString,
				Optional: true,
			},
			logShippingTokenEnabled: {
				Type:     schema.TypeBool,
				Optional: true,
			},
			logShippingTokenToken: {
				Type:     schema.TypeString,
				Computed: true,
			},
			logShippingTokenUpdatedAt: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			logShippingTokenUpdatedBy: {
				Type:     schema.TypeString,
				Computed: true,
			},
			logShippingTokenCreatedAt: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			logShippingTokenCreatedBy: {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceLogShippingTokenRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, _ := log_shipping_tokens.New(m.(Config).apiToken, m.(Config).baseUrl)
	tokenIdString, ok := d.GetOk(logShippingTokenTokenId)

	if ok {
		id, err := strconv.Atoi(tokenIdString.(string))
		if err != nil {
			return diag.FromErr(err)
		}

		token, err := client.GetLogShippingToken(int32(id))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(fmt.Sprintf("%d", id))
		setLogShippingToken(d, token)

		return nil
	}

	// If for some reason we couldn't find the token by id,
	// looking for the token by its name
	tokenName, ok := d.GetOk(logShippingTokenName)
	if ok {
		enabledValues := []bool{true, false}
		for _, v := range enabledValues {
			token, err := findLogShippingTokenByName(tokenName.(string), v, client)
			if err != nil {
				return diag.FromErr(err)
			}

			if token != nil {
				d.SetId(fmt.Sprintf("%d", token.Id))
				setLogShippingToken(d, token)
				return nil
			}
		}
	}

	return diag.Errorf("couldn't find log shipping token with specified attributes")
}

func findLogShippingTokenByName(name string, enabled bool, client *log_shipping_tokens.LogShippingTokensClient) (*log_shipping_tokens.LogShippingToken, error) {
	retrieveRequest := log_shipping_tokens.RetrieveLogShippingTokensRequest{
		Filter: log_shipping_tokens.ShippingTokensFilterRequest{Enabled: strconv.FormatBool(enabled)},
		Pagination: log_shipping_tokens.ShippingTokensPaginationRequest{
			PageNumber: 1,
			PageSize:   25,
		},
	}

	tokenFound, total, totalRetrieved, err := findToken(name, client, retrieveRequest)
	if err != nil {
		return nil, err
	}

	if tokenFound != nil {
		return tokenFound, nil
	}

	// Pagination
	for total > totalRetrieved {
		retrieveRequest.Pagination.PageNumber += 1
		tokenFound, _, currentlyRetrieved, err := findToken(name, client, retrieveRequest)
		if err != nil {
			return nil, err
		}

		if tokenFound != nil {
			return tokenFound, nil
		}

		totalRetrieved += currentlyRetrieved
	}

	return nil, fmt.Errorf("couldn't find log shipping token with specified attributes")
}

func findTokenInResultsListByName(name string, tokens []log_shipping_tokens.LogShippingToken) *log_shipping_tokens.LogShippingToken {
	for _, token := range tokens {
		if token.Name == name {
			return &token
		}
	}

	return nil
}

func findToken(name string, client *log_shipping_tokens.LogShippingTokensClient, request log_shipping_tokens.RetrieveLogShippingTokensRequest) (*log_shipping_tokens.LogShippingToken, int, int, error) {
	tokens, err := client.RetrieveLogShippingTokens(request)
	if err != nil {
		return nil, 0, 0, err
	}

	tokenFound := findTokenInResultsListByName(name, tokens.Results)

	return tokenFound, int(tokens.Total), len(tokens.Results), nil
}
