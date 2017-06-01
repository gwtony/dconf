package handler

import (
	"sync"
	"strings"
	"github.com/gwtony/gapi/log"
)

var DManager DictManager

type DictManager struct {
	lock sync.RWMutex
	host string
	prefix string
	eh  *EtcdHandler
	dict map[string]string
	log log.Log
}

func InitDictManager(host string, eh *EtcdHandler, log log.Log) *DictManager {
	DManager.host = host
	DManager.prefix = eh.root + ETCD_HOST_VIEW + "/" + host + "/"
	DManager.eh = eh
	DManager.log = log
	DManager.dict = make(map[string]string, DEFAULT_DICT_SIZE)

	return &DManager
}

func GetConfig(key string) string {
	DManager.lock.RLock()
	defer DManager.lock.RUnlock()

	if value, ok := DManager.dict[key]; ok {
		return value
	}

	return ""
}

func (dm *DictManager) WatcherCallback(wm *WatchMessage) {
	dm.log.Debug("Got a watch message", wm)

	dm.lock.Lock()
	defer dm.lock.Unlock()

	key := strings.TrimPrefix(wm.Key, dm.prefix)
	if wm.Type == ETCD_EVENT_PUT {
		dm.log.Debug("Watch add event, key is %s", key)
		dm.dict[key] = wm.Value
	} else if wm.Type == ETCD_EVENT_DELETE {
		dm.log.Debug("Watch delete event, key is %s", key)
		if _, ok := dm.dict[key]; ok {
			delete(dm.dict, key)
		}
	} else {
		dm.log.Debug("Watch invalid status, skip", wm.Type)
	}
}

func (dm *DictManager) Run() {
	for {
		dm.eh.WatchWithPrefix(dm.prefix, dm.WatcherCallback)
		dm.log.Info("Watch interrupted")
	}
}

func (dm *DictManager) PullAll() error {
	dm.log.Debug("Pull all from %s", dm.prefix)
	da, err := dm.eh.GetWithPrefix(dm.prefix)
	if err != nil {
		dm.log.Error("Pull all from etcd failed:", err)
		return err
	}

	dm.lock.Lock()
	defer dm.lock.Unlock()

	for _, m := range da {
		key := strings.TrimPrefix(m.Key, dm.prefix)

		dm.log.Debug("Got a key: %s", key)
		dm.dict[key] = m.Value
	}

	return nil
}

