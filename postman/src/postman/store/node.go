package store

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"postman/util"
)

const dbSuffix = ".db"

type node struct {
	Path  string
	Key   string
	Type  string
	Value string
	Store *store
}

type dbFile struct {
	Type  string `codec:"t"`
	Value string `codec:"v"`
}

// initKV creates a Key-Value pair
func initKV(store *store, pt string, key string, tp string, value string) *node {
	return &node{
		Path:  pt,
		Key:   key,
		Type:  tp,
		Store: store,
		Value: value,
	}
}

// list all keys exists
func listAllKey(store *store) []string {
	var filename string
	keys := make([]string, 0)
	filepath.Walk(store.RootPath, func(pt string, fi os.FileInfo, err error) error {
		if nil == fi {
			return err
		}
		if fi.IsDir() {
			return nil
		}
		if strings.HasSuffix(pt, dbSuffix) {
			filename = path.Base(strings.TrimSuffix(pt, dbSuffix))
			keys = append(keys, string(util.DecodeBase64(filename)))
		}
		return err
	})
	return keys
}

func removeByKey(store *store, key string) error {
	pt := pathByKey(store, key)
	return os.Remove(pt)
}

// destroy all db files
func removeAllKey(store *store) error {
	os.RemoveAll(store.RootPath)
	return os.Mkdir(store.RootPath, 0755)
}

// loadKV load a Key-Value pair form file by key
func loadKV(store *store, key string) (n *node, err error) {
	pt := pathByKey(store, key)
	data, err := ioutil.ReadFile(pt)
	if err != nil {
		return
	}
	if len(data) < 1 {
		err = errors.New("empty file found")
		return
	}
	value := util.Decrypt(store.SecretKey, data)
	df := new(dbFile)
	err = util.MsgDecode(value, df)
	if err != nil {
		return
	}
	n = initKV(store, pt, key, df.Type, df.Value)
	return
}

// createKV create a Key-Value pair and write to file
func createKV(store *store, key string, tp string, value string) (n *node, err error) {
	pt := pathByKey(store, key)
	n = initKV(store, pt, key, tp, value)
	err = n.update()
	return
}

// get filename by key
func pathByKey(store *store, key string) string {
	filename := string(util.EncodeBase64([]byte(key))) + dbSuffix
	return path.Join(store.RootPath, filename)
}

// update node
func (n *node) update() error {
	df := dbFile{n.Type, n.Value}
	msgData, err := util.MsgEncode(df)
	if err != nil {
		return err
	}
	data := util.Encrypt(n.Store.SecretKey, msgData)
	return ioutil.WriteFile(n.Path, data, 0644)
}

// destroy node
func (n *node) destroy() error {
	return os.Remove(n.Path)
}
