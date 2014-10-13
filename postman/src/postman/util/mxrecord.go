package util

import (
	"bytes"
	"errors"
	"math/rand"
	"os/exec"
	"strings"

	"postman/cache"
)

var mxCache = cache.NewCache(7200)

// use dig command to get dns record
func digCmd(addr string, t string) (string, error) {
	out := bytes.Buffer{}
	cmd := exec.Command("dig", addr, t, "+short")
	cmd.Stdout = &out
	err := cmd.Run()
	return string(out.Bytes()[:]), err
}

// get cname record form dns
func cnameRecord(addr string) (record string, err error) {
	cname, err := digCmd(addr, "CNAME")
	if err != nil {
		return
	}
	if strings.Contains(cname, ".") {
		record = strings.TrimRight(strings.Fields(cname)[0], ".")
		return
	}
	err = errors.New("no legal CNAME record: " + addr)
	return
}

// Use dig command to get mx record
// dig example.com MX +short
// 10 mx.example.com.
func mxRecords(addr string) (records []string, err error) {
	value, ok := mxCache.Get(addr)
	if ok {
		records = value.([]string)
		return
	}
	out, err := digCmd(addr, "MX")
	if err != nil {
		return
	}
	for _, mxRecord := range strings.Fields(out) {
		if !strings.Contains(mxRecord, ".") {
			continue
		}
		record := strings.TrimRight(mxRecord, ".")
		records = append(records, record)
	}
	if len(records) < 1 {
		// If MX not exist, check cname record
		cname, err := cnameRecord(addr)
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
