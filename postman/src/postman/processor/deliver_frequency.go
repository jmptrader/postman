package processor

import (
	"strconv"
	"time"

	"postman/client"
)

const (
	FREQUENCY_PREFIX  = "domainFrequency:"
	MAX_SEND_INTERVAL = 300
)

// get deliver interval for a certain domain
func GetDeliverInterval(domain string) time.Duration {
	f, err := GetDeliverFrequency(domain)
	if err != nil || f == 0 {
		return time.Duration(MAX_SEND_INTERVAL) * time.Second
	}
	return time.Duration(60000/f) * time.Millisecond
}

func GetDeliverFrequency(domain string) (int, error) {
	var frequency string
	frequency, ok := store.Get(FREQUENCY_PREFIX + domain)
	if !ok {
		frequency, ok = store.Get(FREQUENCY_PREFIX + "default")
	}
	if !ok {
		frequency, err := client.Postman.Tunnel.RequestBlock("frequency", map[string]string{"domain": "default"})
		if err != nil {
			return 0, err
		}
		SaveDeliverFrequency(domain, frequency)
	}
	return strconv.Atoi(frequency)
}

func SaveDeliverFrequency(domain string, frequency string) error {
	return store.Set(FREQUENCY_PREFIX+domain, frequency)
}

func DeleteDliverFrequency(domain string) {
	store.Destroy(FREQUENCY_PREFIX + domain)
}
