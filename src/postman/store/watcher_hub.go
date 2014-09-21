package store

import (
	"errors"
	"sync"
	"time"
)

type watcherHub struct {
	emitLock *sync.Mutex
	watchers map[string]*watcher
}

type watcher struct {
	emitChan chan bool
}

func newWatchHub() *watcherHub {
	return &watcherHub{
		emitLock: new(sync.Mutex),
		watchers: make(map[string]*watcher),
	}
}

func (wh *watcherHub) watch(key string, timeout time.Duration) error {
	if _, ok := wh.watchers[key]; ok {
		panic("key can not watch by muti watchers currently")
	}
	defer delete(wh.watchers, key)
	wh.watchers[key] = &watcher{make(chan bool, 1)}
	select {
	// return nil if new item found
	case <-wh.watchers[key].emitChan:
		return nil
	// timeout otherwise
	case <-time.After(timeout):
		close(wh.watchers[key].emitChan)
		return errors.New("Timeout")
	}
}

func (wh *watcherHub) notify(key string) {
	wh.emitLock.Lock()
	defer wh.emitLock.Unlock()
	wc, ok := wh.watchers[key]
	if !ok {
		return
	}
	wc.emitChan <- true
}
