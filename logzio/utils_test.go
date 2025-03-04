package logzio

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"time"
)

func awaitApply(sec int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		time.Sleep(time.Duration(sec) * time.Second)
		return nil
	}
}
