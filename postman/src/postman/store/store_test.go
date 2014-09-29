package store

import (
	"os"
	"testing"
	"time"
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

func TestPush(t *testing.T) {
	s.Add("test2", "1")
	s.LPush("test2", []string{"-1", "0"})
	s.RPush("test2", []string{"2", "3"})
	if s.Size("test2") != 5 {
		t.Log(s.Members("test2"))
		t.Error("push work unexpected")
	}
	if s.Members("test2")[0] != "-1" && s.Members("test2")[4] != "3" {
		t.Log(s.Members("test2"))
		t.Error("push work unexpected")
	}
}

func TestPop(t *testing.T) {
	s.LPush("test3", []string{"-1", "0"})
	s.RPush("test3", []string{"2", "3"})
	reply, _ := s.BLPOP("test3", time.Microsecond)
	if reply != "-1" {
		t.Log(s.Members("test3"))
		t.Error("pop work unexpected")
	}
}

func TestPopWithBlock(t *testing.T) {
	go func() {
		wait := time.NewTimer(time.Millisecond)
		<-wait.C
		s.LPush("test4", []string{"-1", "0"})
		s.RPush("test4", []string{"2", "3"})
	}()
	reply, err := s.BLPOP("test4", time.Second)
	if reply != "-1" {
		t.Log(reply, s.Members("test4"), err)
		t.Error("BLOCK work unexpected")
	}
}

func TestPopWithBlockAndTimeout(t *testing.T) {
	go func() {
		wait := time.NewTimer(time.Millisecond * 100)
		<-wait.C
		s.LPush("test5", []string{"-1", "0"})
		s.RPush("test5", []string{"2", "3"})
	}()
	_, err := s.BLPOP("test5", time.Millisecond)
	if err == nil {
		t.Log(s.Members("test5"), err)
		t.Error("BLOCK work unexpected")
	}
}

func TestEndWithFlushAll1(t *testing.T) {
	removeAllKey(s)
}
