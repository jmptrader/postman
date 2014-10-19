package processor

import (
	"log"
	"testing"
)

func TestCreateProcessor(t *testing.T) {
	dp, err := NewProcessor("a.com")
	if err != nil {
		t.Fatalf("Create processor with error %s", err.Error())
	}
	dp.DisAvailable()
	log.Print(dp)
	_, err = NewProcessor("a.com")
	if err == nil {
		t.Fatalf("Create processor again should face error")
	}
}

func TestGetProcessor(t *testing.T) {
	dp2, err := GetProcessor("a.com")
	if err != nil {
		t.Fatalf("GetProcessor does not work well %s", err.Error())
	}
	if dp2.CheckSender() {
		t.Log(dp2)
		t.Fatalf("Pointer err!")
	}
}

func TestAvailable(t *testing.T) {
	dp3, _ := NewProcessor("b.com")
	if !dp3.CheckSender() {
		t.Fatalf("Get sender with error")
	}
	dp3.DisAvailable()
	if dp3.CheckSender() {
		t.Log(dp3)
		t.Fatalf("Pointer err!")
	}
	dp3.SetAvailable()
	if !dp3.CheckSender() {
		t.Fatalf("Get sender with error")
	}
}
