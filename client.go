package main

import (
	"context"
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"
)

type Etcd struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

var (
	G_etcd *Etcd
)

func Init() () {
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

	G_etcd = &Etcd{
		client: client,
		kv:     kv,
		lease:  lease,
	}
	return
}

func (etcd *Etcd) GetKeyValue(key string) error {
	getResp, err := etcd.kv.Get(context.TODO(), key, clientv3.WithPrefix())
	if err != nil {
		fmt.Println(err)
	}

	for _, kvPair := range getResp.Kvs {
		fmt.Println(kvPair)
	}
	return err
}

func (etcd *Etcd) PutKeyValue(key, value string) error {
	putResp, err := etcd.kv.Put(context.TODO(), key, value, clientv3.WithPrevKV())
	if err != nil {
		fmt.Println(err)
	}
	if putResp.PrevKv != nil {
		fmt.Println(putResp.PrevKv.Value)
	}
	return err
}

func (etcd *Etcd) DeleteKey(key string) error {
	delResp, err := etcd.kv.Delete(context.TODO(), key, clientv3.WithPrevKV())
	if err != nil {
		fmt.Println(err)
	}
	if len(delResp.PrevKvs) != 0 {
		fmt.Println(delResp.PrevKvs[0].Value)
	}
	return err
}

func (etcd *Etcd) PutKeyWithLease(key string, timeout int64) (err error) {

	leaseGrantResp, err := etcd.lease.Grant(context.TODO(), timeout)
	if err != nil {
		return
	}

	leaseId := leaseGrantResp.ID
	_, err = etcd.kv.Put(context.TODO(), key, "", clientv3.WithLease(leaseId))
	if err != nil {
		return
	}
	return
}

func main() {
	InitEtcd()

	G_etcd.PutKeyValue("/demo/A/B", "hello world")
	G_etcd.GetKeyValue("/demo/A/B")
	G_etcd.DeleteKey("/demo/A/B")
	G_etcd.PutKeyWithLease("/demo/A/B", 100)
}

