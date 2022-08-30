package logzio

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
	"io/ioutil"
	"log"
	"regexp"
	"testing"
)

const (
	logShippingTokenResourceCreateToken string = "create_log_shipping_token"
	logShippingTokenResourceUpdateToken string = "update_log_shipping_token"
)

func TestAccLogzioLogShippingToken_CreateLogShippingToken(t *testing.T) {
	tokenName := "tf_test_create"
	resourceName := "logzio_log_shipping_token." + tokenName
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceTestLogShippingToken(tokenName, logShippingTokenResourceCreateToken),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, logShippingTokenName, "tf_test_create"),
				),
			},
			{
				Config:            resourceTestLogShippingToken(tokenName, logShippingTokenResourceCreateToken),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLogzioLogShippingToken_UpdateLogShippingToken(t *testing.T) {
	tokenName := "tf_test_update"
	resourceName := "logzio_log_shipping_token." + tokenName
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceTestLogShippingToken(tokenName, logShippingTokenResourceCreateToken),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, logShippingTokenName, "tf_test_create"),
				),
			},
			{
				Config: resourceTestLogShippingToken(tokenName, logShippingTokenResourceUpdateToken),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, logShippingTokenName, "tf_test_update"),
					resource.TestCheckResourceAttr(resourceName, logShippingTokenEnabled, "false"),
				),
			},
			{
				Config:            resourceTestLogShippingToken(tokenName, logShippingTokenResourceUpdateToken),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLogzioLogShippingToken_CreateLogShippingTokenEmptyName(t *testing.T) {
	tokenName := "tf_test_create_fail_on_name"
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      resourceTestLogShippingToken(tokenName, "create_log_shipping_token_invalid_name"),
				ExpectError: regexp.MustCompile("name must be set"),
			},
		},
	})
}

func TestAccLogzioLogShippingToken_UpdateLogShippingTokenEmptyName(t *testing.T) {
	tokenName := "tf_test_update_fail_on_name"
	resourceName := "logzio_log_shipping_token." + tokenName
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceTestLogShippingToken(tokenName, logShippingTokenResourceCreateToken),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, logShippingTokenName, "tf_test_create"),
				),
			},
			{
				Config:      resourceTestLogShippingToken(tokenName, "update_log_shipping_token_invalid_name"),
				ExpectError: regexp.MustCompile("name must be set"),
			},
		},
	})
}

func resourceTestLogShippingToken(name string, path string) string {
	content, err := ioutil.ReadFile(fmt.Sprintf("testdata/fixtures/%s.tf", path))
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf(fmt.Sprintf("%s", content), name)
}
