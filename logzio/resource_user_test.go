package logzio

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestAccLogzioUser_CreateUser(t *testing.T) {
	username := "test_resource_user@tfacctest.com"
	accountId, _ := strconv.ParseInt(os.Getenv(envLogzioAccountId), utils.BASE_10, utils.BITSIZE_64)
	fullName := "test test"
	fullNameUpdate := "test test update"
	resourceName := "logzio_user.test_user"
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLogzioUserConfig(username, fullName, accountId),
				Check: resource.ComposeTestCheckFunc(
					awaitApply(15),
					resource.TestCheckResourceAttr(
						resourceName, userUsername, username),
					resource.TestCheckResourceAttr(resourceName, userFullName, fullName),
				),
			},
			{
				PreConfig: func() { time.Sleep(15 * time.Second) },
				Config:    testAccCheckLogzioUserConfig(username, fullNameUpdate, accountId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, userUsername, username),
					resource.TestCheckResourceAttr(resourceName, userFullName, fullNameUpdate),
				),
			},
			{
				Config:            testAccCheckLogzioUserConfig(username, fullName, accountId),
				ResourceName:      "logzio_user.test_user",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckUserDestroy(s *terraform.State) error {
	if testAccProvider.Meta() == nil {
		return nil
	}
	client := usersClient(testAccProvider.Meta())
	for _, r := range s.RootModule().Resources {
		if r.Type != "logzio_user" {
			continue
		}
		id, err := strconv.ParseInt(r.Primary.ID, 10, 64)
		if err != nil {
			return err
		}
		_, err = client.GetUser(int32(id))
		if err == nil {
			return fmt.Errorf("user %s still exists", r.Primary.ID)
		}
	}
	return nil
}

func testAccCheckLogzioUserConfig(username string, fullname string, accountId int64) string {
	return fmt.Sprintf(`
resource "logzio_user" "test_user" {
  username = "%s"
  fullname = "%s"
  account_id = %d
  role = "USER_ROLE_READONLY"
}
`, username, fullname, accountId)
}
