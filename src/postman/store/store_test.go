package store

import (
	"os"
	"testing"
)

func init() {
	s, _ = New(os.Getenv("POSTMAN_DB_DIR"), "Zx+qB1L0oPZubOcsL/S6tjwbugbCRagX").(*store)
}

func TestStorageUpdate(t *testing.T) {
	t.Log(s.Set("name", "vincenting"))
	if result, _ := s.Get("name"); result != "vincenting" {
		t.Log(s.Get("name"))
		t.Error("Value changed unexcepted")
	}
}

func TestLock(t *testing.T) {
	go func() {
		name, ok := s.Get("name")
		if ok && name == "vt1" {
			t.Log(name)
			t.Error("Upate happened before get")
		}
	}()
	go func() {
		s.Set("name", "vt1")
	}()
}

func TestFunc(t *testing.T) {
	s.Set("name", "vt")
	if result, _ := s.Get("name"); result != "vt" {
		t.Error("Update work unexpected")
	}
	s.Destroy("name")
	if _, ok := s.Get("name"); ok {
		t.Error("Destroy work unexpected")
	}
}

func TestFilter(t *testing.T) {
	s.Set("name:vt", "Vincent")
	s.Set("notname:at", "Alvin")
	if s.Keys("name:")[0] != "name:vt" {
		t.Log(s.Keys("name:"))
		t.Error("Keys work unexpected")
	}
}

func TestSaddAndSMembers(t *testing.T) {
	t.Log(s.Add("test", "1"))
	t.Log(s.Add("test", "1"))
	t.Log(s.Add("test", "1"))
	t.Log(s.Add("test", "2"))
	t.Log(s.Members("test"))
	if s.Members("test")[0] != "1" {
		t.Error("SAdd work unexpected")
	}
	if s.Size("test") != 2 {
		t.Log(s.Members("test"))
		t.Error("SSize work unexpected")
	}
}

func TestSRem(t *testing.T) {
	s.Add("test1", "1")
	s.Rem("test1", "1")
	if s.Size("test1") != 0 {
		t.Log(s.Members("test1"))
		t.Error("SSize work unexpected")
	}
}

func TestEndWithFlushAll1(t *testing.T) {
	removeAllKey(s)
}
