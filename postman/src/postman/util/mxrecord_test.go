package util

import (
	"strings"
	"testing"
)

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
	t.Log(records)
	if !asExpected {
		t.Error("Record return not as expected.")
	}
}

func TestRandMX(t *testing.T) {
	record, err := MxRecord("qq.com")
	t.Log(record)
	if err != nil {
		t.Error(err)
	}
	if !strings.Contains(record, "qq") {
		t.Error("Record return wrong")
	}
}

func TestBothCnameMX(t *testing.T) {
	record, err := MxRecord("gmail.com")
	t.Log(record)
	if err != nil {
		t.Error(err)
	}
	if !strings.Contains(record, "google.com") {
		t.Error("Record return wrong")
	}
}
