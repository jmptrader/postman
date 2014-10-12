package store

import (
	"errors"
	"strings"
	"sync"
	"time"

	"postman/util"
)

// TODO add cache in mem for some certain keys
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
	RootPath   string
	SecretKey  []byte
	worldLock  sync.RWMutex
	keymap     map[string]bool
	watcherHub *watcherHub
}

func New(rootPath string, secretKey string) Store {
	s := &store{
		RootPath:   rootPath,
		SecretKey:  []byte(secretKey),
		keymap:     map[string]bool{},
		watcherHub: newWatchHub(),
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
	return n.Value, true
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

// push a new item to the left of a list
func (st *store) LPush(key string, items []string) (left int, err error) {
	st.worldLock.Lock()
	defer func() {
		st.worldLock.Unlock()
		st.watcherHub.notify(key)
	}()
	currentItems := st.members(key)
	err = st.setList(key, append(items, currentItems...))
	if err != nil {
		return
	}
	left = len(st.members(key))
	return
}

// push a new item to the right of a list
func (st *store) RPush(key string, items []string) (left int, err error) {
	st.worldLock.Lock()
	defer func() {
		st.worldLock.Unlock()
		st.watcherHub.notify(key)
	}()
	currentItems := st.members(key)
	err = st.setList(key, append(currentItems, items...))
	if err != nil {
		return
	}
	left = len(st.members(key))
	return
}

// Remove and return from the left. return err if not found
func (st *store) lPop(key string) (item string, err error) {
	st.worldLock.Lock()
	defer st.worldLock.Unlock()
	currentItems := st.members(key)
	if len(currentItems) > 0 {
		item = currentItems[0]
		st.setList(key, currentItems[1:])
		return
	}
	err = errors.New("No record found")
	return
}

func (st *store) BLPOP(key string, timeout time.Duration) (item string, err error) {
	item, err = st.lPop(key)
	if err == nil {
		return
	}
	err = st.watcherHub.watch(key, timeout)
	if err != nil {
		return
	}
	return st.lPop(key)
}
