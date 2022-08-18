package logzio

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
	"os"
	"regexp"
	"testing"
	"time"
)

func TestAccLogzioRestoreLogs_InitiateRestore(t *testing.T) {
	path := os.Getenv(envLogzioS3Path)
	arn := os.Getenv(envLogzioAwsArn)
	archiveName := "archive_for_restore_initiate"
	restoreName := "tf_test_restore_initiate"
	fullRestoreName := "logzio_restore_logs." + restoreName
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:  getConfigTestArchiveS3Iam(archiveName, path, arn),
				Destroy: false,
			},
			{
				Config: getConfigTestArchiveS3Iam(archiveName, path, arn) +
					getConfigTestRestore(restoreName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(fullRestoreName, restoreLogsId),
					resource.TestCheckResourceAttrSet(fullRestoreName, restoreLogsStatus),
				),
			},
			{
				SkipFunc: func() (bool, error) {
					// waiting for status to change
					time.Sleep(20 * time.Second)
					return false, nil
				},
				Config: getConfigTestArchiveS3Iam(archiveName, path, arn) +
					getConfigTestRestore(restoreName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(fullRestoreName, restoreLogsId),
					resource.TestCheckResourceAttrSet(fullRestoreName, restoreLogsStatus),
				),
			},
			{
				ResourceName:      fullRestoreName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLogzioRestoreLogs_InitiateRestoreEmptyStartTime(t *testing.T) {
	restoreName := "tf_test_empty_start_time"
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      getConfigTestRestoreEmptyStartTime(restoreName),
				ExpectError: regexp.MustCompile("StartTime must be set"),
			},
		},
	})
}

func TestAccLogzioRestoreLogs_InitiateRestoreEmptyEndTime(t *testing.T) {
	restoreName := "tf_test_empty_start_time"
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      getConfigTestRestoreEmptyEndTime(restoreName),
				ExpectError: regexp.MustCompile("EndTime must be set"),
			},
		},
	})
}

func getConfigTestRestore(name string) string {
	now := time.Now()
	accountName := fmt.Sprintf("tf-test-%s", now.Format("2006-01-02,15:04:05"))
	hourAgo := now.Add(-time.Hour)
	return fmt.Sprintf(`resource "logzio_restore_logs" "%s" {
 account_name = "%s"
 start_time = %d
 end_time = %d
}
`, name, accountName, hourAgo.Unix(), now.Unix())
}

func getConfigTestRestoreEmptyStartTime(name string) string {
	now := time.Now()
	accountName := fmt.Sprintf("tf-test-%s", now.Format("2006-01-02,15:04:05"))
	return fmt.Sprintf(`resource "logzio_restore_logs" "%s" {
 account_name = "%s"
 start_time = 0
 end_time = %d
}
`, name, accountName, now.Unix())
}

func getConfigTestRestoreEmptyEndTime(name string) string {
	now := time.Now()
	accountName := fmt.Sprintf("tf-test-%s", now.Format("2006-01-02,15:04:05"))
	hourAgo := now.Add(-time.Hour)
	return fmt.Sprintf(`resource "logzio_restore_logs" "%s" {
 account_name = "%s"
 start_time = %d
 end_time = 0
}
`, name, accountName, hourAgo.Unix())
}

func getConfigTestRestoreNewAccountName(name string) string {
	now := time.Now()
	accountName := fmt.Sprintf("tf-test-%s", now.Format("2006-01-02,15:04:05"))
	hourAgo := now.Add(-time.Hour)
	return fmt.Sprintf(`resource "logzio_restore_logs" "%s" {
 account_name = "%s"
 start_time = %d
 end_time = %d
}
`, name, accountName+"_new", hourAgo.Unix(), now.Unix())
}
