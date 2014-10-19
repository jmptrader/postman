package processor

import (
	"time"

	"postman/util"
)

var frequencyCache = util.NewCache(86400)

func getDeliverFrequency(domain string) time.Duration {
	return time.Second * 10
}
