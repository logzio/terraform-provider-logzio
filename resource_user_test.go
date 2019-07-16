package main

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/jonboydell/logzio_client/users"
	"os"
	"strconv"
	"testing"
)

func TestAccLogzioUser_CreateUser(t *testing.T) {

	username := "test@sometest.com"
	accountId, _ := strconv.ParseInt(os.Getenv(envLogzioAccountId), BASE_10, BITSIZE_64)
	terraformPlan := testAccCheckLogzioUserConfig(username, "test test", accountId)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccLogzioUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: terraformPlan,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogzioUserExists("logzio_user.test_user"),
					resource.TestCheckResourceAttr(
						"logzio_user.test_user", "username", username),
				),
			},
		},
	})
}

func testAccCheckLogzioUserExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no user ID is set")
		}

		id, err := strconv.ParseInt(rs.Primary.ID, BASE_10, BITSIZE_64)

		var client *users.UsersClient
		client, _ = users.New(os.Getenv(envLogzioApiToken))

		_, err = client.GetUser(int64(id))

		if err != nil {
			return fmt.Errorf("user doesn't exist")
		}

		return nil
	}
}

func testAccLogzioUserDestroy(s *terraform.State) error {
	for _, r := range s.RootModule().Resources {
		id, err := strconv.ParseInt(r.Primary.ID, BASE_10, BITSIZE_64)
		if err != nil {
			return err
		}

		var client *users.UsersClient
		client, _ = users.New(os.Getenv(envLogzioApiToken))

		err = client.DeleteUser(int64(id))
		if err == nil {
			return fmt.Errorf("endpoint still exists")
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
  roles = [2]
}
`, username, fullname, accountId)
}
