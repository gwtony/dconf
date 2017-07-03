package handler

import (
	"os"
	"sync"
	"strings"
	"github.com/gwtony/gapi/log"
)

var DManager DictManager

type DictManager struct {
	lock     sync.RWMutex
	host     string
	store    string
	storable bool
	prefix   string
	eh       *EtcdHandler
	dict     map[string]string
	ch       chan ConfigMeta
	log      log.Log
}

func InitDictManager(host string, eh *EtcdHandler, store string, log log.Log) (*DictManager, error) {
	DManager.host = host
	DManager.prefix = eh.root + ETCD_HOST_VIEW + "/" + host + "/"
	DManager.eh = eh
	DManager.log = log
	DManager.dict = make(map[string]string, DEFAULT_DICT_SIZE)
	DManager.store = store
	if store != "" {
		DManager.storable = true
		_, err := os.Stat(store)
		if os.IsNotExist(err) {
			err = os.MkdirAll(store, 0755)
			if err != nil {
				log.Error("Mkdir in path: %s failed, ", store, err)
				return nil, err
			}
		}

		go DManager.runStorer()
	}

	DManager.ch = make(chan ConfigMeta, 100)

	return &DManager, nil
}

func GetConfig(key string) string {
	DManager.lock.RLock()
	defer DManager.lock.RUnlock()

	if value, ok := DManager.dict[key]; ok {
		return value
	}

	return ""
}

func (dm *DictManager) runStorer() {
	for {
		select {
		case cm := <-dm.ch:
			dm.log.Debug("Got message: ", cm)

			arr := strings.Split(cm.Key, "/")
			dir := dm.store + "/" + arr[0]
			_, err := os.Stat(dir)
			if os.IsNotExist(err) {
				err = os.Mkdir(dir, 0755)
				if err != nil {
					dm.log.Error("Storer mkdir %s failed", dir)
					continue //TODO: next select
				}
			}

			func () {
				var f *os.File

				name := dir + "/" + arr[1]
				dm.log.Debug("Open path: %s", name)
				f, err = os.OpenFile(name, os.O_RDWR | os.O_CREATE, 0755)
				if err != nil {
					dm.log.Error("Open file: %s failed", name)
					return
				}

				defer f.Close()

				n, err := f.Write([]byte(cm.Value))
				if err != nil {
					dm.log.Error("Write error: ", err)
					return
				}
				if n != len(cm.Value) {
					dm.log.Error("Write len error: %d", n)
					f.Close()
					return
				}
				dm.log.Debug("Write %d to file %s", n, name)
			}()
		}
	}
}

func (dm *DictManager) WatcherCallback(wm *WatchMessage) {
	dm.log.Debug("Got a watch message", wm)

	dm.lock.Lock()
	defer dm.lock.Unlock()

	key := strings.TrimPrefix(wm.Key, dm.prefix)
	if wm.Type == ETCD_EVENT_PUT {
		dm.log.Debug("Watch add event, key is %s", key)
		dm.dict[key] = wm.Value
		if dm.storable {
			cfm := &ConfigMeta{Key: key, Value: wm.Value}
			go func() {dm.ch <- *cfm}()
		}
	} else if wm.Type == ETCD_EVENT_DELETE {
		dm.log.Debug("Watch delete event, key is %s", key)
		if _, ok := dm.dict[key]; ok {
			delete(dm.dict, key)
			//TODO: delete file
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
		//TODO: load from store dir
		return err
	}

	dm.lock.Lock()
	defer dm.lock.Unlock()

	for _, m := range da {
		key := strings.TrimPrefix(m.Key, dm.prefix)

		dm.log.Debug("Got a key: %s", key)
		dm.dict[key] = m.Value

		if dm.storable {
			cfm := &ConfigMeta{Key: key, Value: m.Value}
			dm.ch <- *cfm
		}
	}

	return nil
}

