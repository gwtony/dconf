package handler
import (
	"time"
	"github.com/gwtony/gapi/log"
	ec "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
	"golang.org/x/net/context"
)

type EtcdMessage struct {
	Key []byte
	Value []byte
	Version int64
}

type EtcdHandler struct {
	user string
	pwd string
	auth bool
	addrs []string
	root string
	to time.Duration
	log log.Log
}

func InitEtcdHandler(addrs []string, to time.Duration, user, pwd string, auth bool, root string, log log.Log) *EtcdHandler {
	eh := &EtcdHandler{
		addrs: addrs,
		to: to,
		user: user,
		pwd: pwd,
		auth: auth,
		root: root,
		log: log,
	}
	return eh
}

func (eh *EtcdHandler) newClient() (*ec.Client, error) {
	var err error
	var cli *ec.Client
	if eh.auth { // Auth enabled
		cli, err = ec.New(ec.Config{
			Endpoints:   eh.addrs,
			Username: eh.user,
			Password: eh.pwd,
			DialTimeout: eh.to,
		})
	} else { // Auth disabled
		cli, err = ec.New(ec.Config{
			Endpoints:   eh.addrs,
			DialTimeout: eh.to,
		})
	}

	if err != nil {
		return nil, err
	}

	return cli, nil
}

func parseEtcdError(err error, log log.Log) {
	switch err {
	case context.Canceled:
		log.Error("ctx is canceled by another routine: %v", err)
	case context.DeadlineExceeded:
		log.Error("ctx is attached with a deadline is exceeded: %v", err)
	case rpctypes.ErrEmptyKey:
		log.Error("client-side error: %v", err)
	default:
		log.Error("bad cluster endpoints, which are not etcd servers: %v", err)
	}
}

func (eh *EtcdHandler) Set(key, value string) error {
	cli, err := eh.newClient()
	if err != nil {
		eh.log.Error("Set new etcd client failed:", err)
		return err
	}
	defer cli.Close()

	eh.log.Debug("Set key: %s, value: %s", key, value)
	ctx, cancel := context.WithTimeout(context.Background(), eh.to)
	_, err = cli.Put(ctx, key, value)
	cancel()
	if err != nil {
		parseEtcdError(err, eh.log)
		return err
	}

	return nil
}

func (eh *EtcdHandler) UnSet(key string) error {
	cli, err := eh.newClient()
	if err != nil {
		eh.log.Error("UnSet new etcd client failed:", err)
		return err
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), eh.to)
	defer cancel()

	eh.log.Debug("to unset key: %s", key)
	dresp, err := cli.Delete(ctx, key)
	if err != nil {
		parseEtcdError(err, eh.log)
		return err
	}
	eh.log.Info("Delete %d keys", dresp.Deleted)

	return nil
}

func (eh *EtcdHandler) Get(key string) (*EtcdMessage, error) {
	var em EtcdMessage
	cli, err := eh.newClient()
	if err != nil {
		eh.log.Error("Get new etcd client failed:", err)
		return nil, err
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), eh.to)
	resp, err := cli.Get(ctx, key)
	cancel()
	if err != nil {
		parseEtcdError(err, eh.log)
		return nil, err
	}
	//for _, ev := range resp.Kvs {
	//	fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	//}
	//fmt.Println(resp.Kvs)
	if len(resp.Kvs) == 0 { //Not found
		return nil, nil
	}
	em.Key = resp.Kvs[0].Key
	em.Value = resp.Kvs[0].Value
	em.Version = resp.Kvs[0].Version

	return &em, nil
}

func (eh *EtcdHandler) GetWithPrefix(key string) ([]*EtcdMessage, error) {
	var ea []*EtcdMessage

	cli, err := eh.newClient()
	if err != nil {
		eh.log.Error("GetWithPrefix new etcd client failed:", err)
		return ea, err
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), eh.to)
	resp, err := cli.Get(ctx, key, ec.WithPrefix())
	cancel()
	if err != nil {
		parseEtcdError(err, eh.log)
		return ea, err
	}
	for _, ev := range resp.Kvs {
		em := &EtcdMessage{Key: ev.Key, Value: ev.Value, Version: ev.Version}
		//eh.log.Debug("GetWithPrefix: (%s):%s\n", ev.Key, ev.Value)
		ea = append(ea, em)
	}

	return ea, nil
}

func (eh *EtcdHandler) GetWithPrefixLimit(key string, n int64) ([]*EtcdMessage, error) {
	var ea []*EtcdMessage
	var resp *ec.GetResponse


	cli, err := eh.newClient()
	if err != nil {
		eh.log.Error("GetWithPrefix new etcd client failed:", err)
		return ea, err
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), eh.to)
	if n > 0 {
		resp, err = cli.Get(ctx, key, ec.WithPrefix(), ec.WithSort(ec.SortByKey, ec.SortDescend), ec.WithLimit(n))
	} else {
		resp, err = cli.Get(ctx, key, ec.WithPrefix(), ec.WithSort(ec.SortByKey, ec.SortDescend))
	}

	cancel()
	if err != nil {
		parseEtcdError(err, eh.log)
		return ea, err
	}
	for _, ev := range resp.Kvs {
		em := &EtcdMessage{Key: ev.Key, Value: ev.Value, Version: ev.Version}
		//eh.log.Debug("GetWithPrefix: (%s):%s\n", ev.Key, ev.Value)
		ea = append(ea, em)
	}

	return ea, nil
}

func (eh *EtcdHandler) Watch(key string, deal func(m *WatchMessage)) (error) {
	var evtype string
	cli, err := eh.newClient()
	if err != nil {
		eh.log.Error("Cas new etcd client failed:", err)
		return err
	}
	defer cli.Close()

	rch := cli.Watch(context.Background(), key)
	for wresp := range rch {
		for _, ev := range wresp.Events {
			//TODO: parse ev.Type
			eh.log.Debug("Watch type is %s", ev.Type)
			if ev.Type == mvccpb.PUT {
				evtype = ETCD_EVENT_PUT
			} else if ev.Type == mvccpb.DELETE {
				evtype = ETCD_EVENT_DELETE
			}
			//op := ""
			m, err := DecodeWatchMessage(evtype, ev.Kv.Key, ev.Kv.Value)
			if err != nil {
				eh.log.Error("Decode watch message failed")
				continue
			}
			deal(m)
		}
	}

	return nil
}

func (eh *EtcdHandler) WatchWithPrefix(key string, deal func(m *WatchMessage)) (error) {
	var evtype string

	cli, err := eh.newClient()
	if err != nil {
		eh.log.Error("WatchWithPrefix new etcd client failed:", err)
		return err
	}
	defer cli.Close()

	rch := cli.Watch(context.Background(), key, ec.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			//TODO: parse ev.Type
			eh.log.Debug("etcd action type is %s", ev.Type.String())
			if ev.Type == mvccpb.PUT {
				evtype = ETCD_EVENT_PUT
			} else if ev.Type == mvccpb.DELETE {
				evtype = ETCD_EVENT_DELETE
			}

			m, err := DecodeWatchMessage(evtype, ev.Kv.Key, ev.Kv.Value)
			if err != nil {
				eh.log.Error("Decode watch message failed")
				continue
			}
			deal(m)
		}
	}

	return nil
}

//func (eh *EtcdHandler) WatchStartWithPrefix(key string, deal func(m *WatchStartMessage)) (error) {
//	cli, err := eh.newClient()
//	if err != nil {
//		eh.log.Error("WatchWithPrefix new etcd client failed:", err)
//		return err
//	}
//	defer cli.Close()
//
//	rch := cli.Watch(context.Background(), key, ec.WithPrefix())
//	for wresp := range rch {
//		for _, ev := range wresp.Events {
//			//TODO: parse ev.Type
//			op := ""
//			m, err := DecodeWatchStartMessage(op, ev.Kv.Key, ev.Kv.Value)
//			if err != nil {
//				eh.log.Error("Decode watch start message failed")
//				continue
//			}
//			deal(m)
//		}
//	}
//
//	return nil
//}

func (eh *EtcdHandler) Cas(key, value string, version int64) (error) {
	cli, err := eh.newClient()
	if err != nil {
		eh.log.Error("Cas new etcd client failed:", err)
		return err
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), eh.to)
	kvc := ec.NewKV(cli)

	//compare key and value
	_, err = kvc.Txn(ctx).
		If(ec.Compare(ec.Version(key), "=", version)). // txn value comparisons are lexical
		Then(ec.OpPut(key, value)).
		Else().
		Commit()
	cancel()
	if err != nil {
		parseEtcdError(err, eh.log)
		return err
	}

	return nil
}

func (eh *EtcdHandler) CasLess(key, value string, version int64) (error) {
	cli, err := eh.newClient()
	if err != nil {
		eh.log.Error("Cas new etcd client failed:", err)
		return err
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), eh.to)
	kvc := ec.NewKV(cli)

	//compare key and value
	_, err = kvc.Txn(ctx).
		If(ec.Compare(ec.Version(key), ">", version)). // txn value comparisons are lexical
		//Then(ec.OpPut(key, value)).
		Else(ec.OpPut(key, value)).
		Commit()
	cancel()
	if err != nil {
		parseEtcdError(err, eh.log)
		return err
	}

	return nil
}

//TODO: get top x result may use (WithSort, WithLimit)

