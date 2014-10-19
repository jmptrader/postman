package util

import (
	"sync"
	"time"
)

// cache item
type cacheItem struct {
	Data     interface{}
	ExpireIn int64
}

type Cache struct {
	dataMap    map[string]*cacheItem
	actionLock *sync.Mutex
	expire     int64
}

// return new cache instance
func NewCache(expire int64) *Cache {
	return &Cache{
		expire:     expire,
		dataMap:    map[string]*cacheItem{},
		actionLock: new(sync.Mutex),
	}
}

func (c *Cache) Get(key string) (item interface{}, ok bool) {
	c.actionLock.Lock()
	defer c.actionLock.Unlock()
	d, ok := c.dataMap[key]
	if !ok {
		return
	}
	item = d.Data
	now := time.Now().Unix()
	if now > d.ExpireIn {
		delete(c.dataMap, key)
		ok = false
		return
	}
	return
}

func (c *Cache) Delete(key string) {
	c.actionLock.Lock()
	defer c.actionLock.Unlock()
	delete(c.dataMap, key)
}

func (c *Cache) Update(key string, value interface{}) {
	c.actionLock.Lock()
	defer c.actionLock.Unlock()
	expire := time.Now().Unix() + 7200
	c.dataMap[key] = &cacheItem{
		Data:     value,
		ExpireIn: expire,
	}
}
