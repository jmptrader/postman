package store

import (
	"os"
	"testing"
)

var s *store

func init() {
	os.Mkdir(os.Getenv("POSTMAN_DB_DIR"), 0755)
	s, _ = New(os.Getenv("POSTMAN_DB_DIR"), "Zx+qB1L0oPZubOcsL/S6tjwbugbCRagX").(*store)
}

func TestEmptyDir(t *testing.T) {
	if len(listAllKey(s)) > 0 {
		t.Fatal("no db record should be found.")
	}
}

func TestCreateAndRead(t *testing.T) {
	_, err := createKV(s, "1", "string", "2")
	if err != nil {
		t.Fatal(err)
	}
	if len(listAllKey(s)) != 1 {
		t.Fatal("insert db with error.")
	}
	n, err := loadKV(s, "1")
	if n.Type != "string" {
		t.Fatal("read db with error.")
	}
}

func TestCreateAndUpdate(t *testing.T) {
	n, err := createKV(s, "2", "string", "2")
	if err != nil {
		t.Fatal(err)
	}
	n.Value = "xx"
	n.update()
	m, err := loadKV(s, "2")
	if m.Value != "xx" {
		t.Fatal("update db with error.")
	}
}

func TestDel(t *testing.T) {
	n, err := loadKV(s, "2")
	if err != nil {
		t.Fatal(err)
	}
	n.destroy()
	m, err := loadKV(s, "2")
	if err == nil && m.Value == "xx" {
		t.Fatal(m)
	}
}

func TestEndWithFlushAll(t *testing.T) {
	removeAllKey(s)
}
