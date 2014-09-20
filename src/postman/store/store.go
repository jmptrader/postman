package store

import (
	"errors"
	"strings"
	"sync"
	"time"

	"postman/util"
)

type Store interface {
	Keys(prefix string) []string

	Get(key string) (string, bool)
	Destroy(key string)
	Set(key string, value string) error

	Members(key string) []string
	Add(key string, value string) error
	Rem(key string, value string) error
	Size(key string) int

	LPush(key string, items []string) (int, error)
	RPush(key string, items []string) (int, error)
	BLPOP(key string, timeout time.Duration) (string, error)
}

type store struct {
	RootPath  string
	SecretKey []byte
	worldLock sync.RWMutex
	keymap    map[string]bool
}

func New(rootPath string, secretKey string) Store {
	s := &store{
		RootPath:  rootPath,
		SecretKey: []byte(secretKey),
		keymap:    map[string]bool{},
	}
	for _, key := range listAllKey(s) {
		s.keymap[key] = true
	}
	return s
}

// get all keys begin with prefix.
func (st *store) Keys(prefix string) (keyArr []string) {
	st.worldLock.Lock()
	defer st.worldLock.Unlock()
	for key, _ := range st.keymap {
		if strings.HasPrefix(key, prefix) {
			keyArr = append(keyArr, key)
		}
	}
	return
}

// get value of a certain key
func (st *store) Get(key string) (res string, ok bool) {
	st.worldLock.Lock()
	defer st.worldLock.Unlock()
	if _, _ok := st.keymap[key]; !_ok {
		return
	}
	n, err := loadKV(st, key)

	if err != nil || n.Type != "string" {
		ok = false
		return
	}
	res = n.Value
	return
}

// destroy a certain key
func (st *store) Destroy(key string) {
	st.worldLock.Lock()
	defer st.worldLock.Unlock()
	delete(st.keymap, key)
	removeByKey(st, key)
}

// update/create key with value
func (st *store) Set(key string, value string) error {
	st.worldLock.Lock()
	defer st.worldLock.Unlock()
	st.keymap[key] = true
	n, err := loadKV(st, key)
	if err != nil {
		_, _err := createKV(st, key, "string", value)
		return _err
	}
	if n.Type != "string" {
		return errors.New("type error")
	}
	n.Value = value
	return n.update()
}

// get list members
func (st *store) members(key string) (items []string) {
	if _, ok := st.keymap[key]; !ok {
		return
	}
	n, err := loadKV(st, key)
	if err != nil || n.Type != "list" {
		return
	}
	util.MsgDecode([]byte(n.Value), &items)
	return
}

func (st *store) setList(key string, items []string) error {
	if len(items) == 0 {
		delete(st.keymap, key)
		removeByKey(st, key)
		return nil
	}
	newValue, err := util.MsgEncode(items)
	if err != nil {
		return err
	}
	_, err = createKV(st, key, "list", string(newValue))
	if err == nil {
		st.keymap[key] = true
	}
	return err
}

// get all members of a list-key
func (st *store) Members(key string) []string {
	st.worldLock.Lock()
	defer st.worldLock.Unlock()
	return st.members(key)
}

// add value to list
// do nothing if value exist
func (st *store) Add(key string, value string) error {
	st.worldLock.Lock()
	defer st.worldLock.Unlock()
	currentItems := st.members(key)
	for _, item := range currentItems {
		if item == value {
			return nil
		}
	}
	return st.setList(key, append(currentItems, value))
}

// remove an item form a certain list
func (st *store) Rem(key string, value string) error {
	st.worldLock.Lock()
	defer st.worldLock.Unlock()
	currentItems := st.members(key)
	resultList := []string{}
	for _, item := range currentItems {
		if item != value {
			resultList = append(resultList, item)
		}
	}
	return st.setList(key, resultList)
}

// get the length of a certain list-key
func (st *store) Size(key string) int {
	st.worldLock.Lock()
	defer st.worldLock.Unlock()
	return len(st.members(key))
}

func (st *store) LPush(key string, items []string) (left int, err error) {
	return
}

func (st *store) RPush(key string, items []string) (left int, err error) {
	return
}

func (st *store) BLPOP(key string, timeout time.Duration) (item string, err error) {
	return
}
