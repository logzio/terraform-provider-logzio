package logzio

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func awaitApply(sec int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		time.Sleep(time.Duration(sec) * time.Second)
		return nil
	}
}

func getRandomId() string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return strconv.Itoa(rand.Intn(10000))
}
