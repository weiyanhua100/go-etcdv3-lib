package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
)

type Etcd struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
	watcher clientv3.Watcher
}

var (
	G_etcd *Etcd
)

func Client() (*Etcd, error) {
	config := clientv3.Config{
		Endpoints:   []string{"172.16.4.36:2379", "localhost:22379"}, 
		DialTimeout: time.Duration(5) * time.Millisecond,
	}

	client, err := clientv3.New(config)
	if err != nil {
		fmt.Println(err)
	}

	kv := clientv3.NewKV(client)
	lease := clientv3.NewLease(client)
        watcher := clientv3.NewWatcher(client)

	G_etcd = &Etcd{
		client: client,
		kv:     kv,
		lease:  lease,
		watcher: watcher,	
	}
	return G_etcd, err
}

func (etcd *Etcd) GetKeyValue(key string) error {
	getResp, err := etcd.kv.Get(context.TODO(), key, clientv3.WithPrefix())
	if err != nil {
		fmt.Println(err)
		return err
	}

	for _, kvPair := range getResp.Kvs {
		fmt.Printf("key is %s \nvalue is %s \n", kvPair.Key, kvPair.Value)
	}
	return err
}

func (etcd *Etcd) PutKeyValue(key, value string) error {
	putResp, err := etcd.kv.Put(context.TODO(), key, value, clientv3.WithPrevKV())
	if err != nil {
		fmt.Println(err)
	}
	if putResp.PrevKv != nil {
		fmt.Printf("Pre value:%s\n", string(putResp.PrevKv.Value))
	}
	return err
}

func (etcd *Etcd) DeleteKey(key string) error {
	delResp, err := etcd.kv.Delete(context.TODO(), key, clientv3.WithPrevKV())
	if err != nil {
		fmt.Println(err)
	}
	if len(delResp.PrevKvs) != 0 {
		fmt.Printf("Delete value:%s\n", delResp.PrevKvs[0].Value)
	}
	return err
}

func (etcd *Etcd) PutKeyWithLease(key string, value string, timeout int64) (err error) {

	leaseGrantResp, err := etcd.lease.Grant(context.TODO(), timeout)
	if err != nil {
		return
	}

	leaseId := leaseGrantResp.ID
	_, err = etcd.kv.Put(context.TODO(), key, value, clientv3.WithLease(leaseId))
	if err != nil {
		return
	}
	return
}

func (etcd * Etcd) Watch(key string) (err error) {
	var event *clientv3.Event
	
	watchRespChan := etcd.watcher.Watch(context.TODO(), key, clientv3.WithPrefix())

        /*
	for wresp := range watchRespChan {
	for _, ev := range wresp.Events {
		fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
	}
	}*/

	for watchResp := range watchRespChan {
		for _, event = range watchResp.Events {
			switch event.Type {
			case mvccpb.PUT:
				fmt.Println("Watch PUT:", string(event.Kv.Value), "Revision:", event.Kv.CreateRevision, event.Kv.ModRevision)
			case mvccpb.DELETE:
			fmt.Println("Watch DELETE:", string(event.Kv.Value), "Revision:", event.Kv.CreateRevision, event.Kv.ModRevision)
			}
		}
	}
	return nil
}

