package logzio

import (
	"math/rand"
	"strconv"
	"time"
)

func getRandomId() string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return strconv.Itoa(rand.Intn(10000))
}
