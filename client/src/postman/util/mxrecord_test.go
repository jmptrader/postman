package util

import (
	"strings"
	"testing"
)

func TestCacheFind(t *testing.T) {
	data := make(map[string]mxRecordData)
	cache := MXCache{data: data}
	l := []string{"name"}
	cache.update("qq.com", l)
	records, ok := cache.find("qq.com")
	if !(len(records) == 1 && ok) {
		t.Error("Can not find elem")
	}
}

func TestCacheExpire(t *testing.T) {
	data := make(map[string]mxRecordData)
	data["qq.com"] = mxRecordData{records: []string{"mx1.qq.com"}, expireIn: 1231231}
	cache := MXCache{data: data}
	records, ok := cache.find("qq.com")
	if len(records) != 1 || ok {
		t.Error("Can not find elem")
	}
}

func TestMXRecord(t *testing.T) {
	records, err := mxRecords("qq.com")
	if err != nil {
		t.Error(err)
	}
	asExpected := false
	for _, record := range records {
		if record == "mx2.qq.com" {
			asExpected = true
		}
	}
	if !asExpected {
		t.Error("Record return not as expected.")
	}
}

func TestRandMX(t *testing.T) {
	record, err := MxRecord("qq.com")
	if err != nil {
		t.Error(err)
	}
	t.Log(record)
	if !strings.Contains(record, "qq") {
		t.Error("Record return wrong")
	}
}

func TestBothCnameMX(t *testing.T) {
	record, err := MxRecord("sogou.com")
	if err != nil {
		t.Error(err)
	}
	t.Log(record)
	if !strings.Contains(record, "sogou.com") {
		t.Error("Record return wrong")
	}
}
