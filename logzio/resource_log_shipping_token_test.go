package logzio

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"io/ioutil"
	"log"
	"testing"
)

const (
	logShippingTokenResourceCreateToken string = "create_log_shipping_token"
	logShippingTokenResourceUpdateToken string = "update_log_shipping_token"
)

func TestAccLogzioLogShippingToken_CreateLogShippingToken(t *testing.T) {
	tokenName := "tf_test_create"
	resourceName := "logzio_log_shipping_token." + tokenName

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckApiToken(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: resourceTestLogShippingToken(tokenName, logShippingTokenResourceCreateToken),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, logShippingTokenName, "tf_test_create"),
					resource.TestCheckResourceAttr(resourceName, logShippingTokenEnabled, "true"),
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

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckApiToken(t) },
		Providers: testAccProviders,
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

func resourceTestLogShippingToken(name string, path string) string {
	content, err := ioutil.ReadFile(fmt.Sprintf("testdata/fixtures/%s.tf", path))
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf(fmt.Sprintf("%s", content), name)
}
