package util

import (
	"errors"
	"math/rand"
	"net"
	"strings"

	"postman/cache"
)

var mxCache = cache.NewCache(7200)

func mxRecords(addr string) (records []string, err error) {
	value, ok := mxCache.Get(addr)
	if ok {
		records = value.([]string)
		return
	}
	mxs, err := net.LookupMX(addr)
	if err != nil {
		return
	}
	for _, mxRecord := range mxs {
		records = append(records, strings.TrimSuffix(mxRecord.Host, "."))
	}
	if len(records) < 1 {
		// If MX not exist, check cname record
		cname, err := net.LookupCNAME(addr)
		if err == nil {
			return mxRecords(cname)
		}
	}
	mxCache.Update(addr, records)
	return
}

// return random record
func MxRecord(addr string) (record string, err error) {
	records, err := mxRecords(addr)
	if err != nil {
		return
	}
	if len(records) == 0 {
		err = errors.New("No legal MX record find for " + addr)
		return
	}
	record = records[rand.Intn(len(records))]
	return
}
