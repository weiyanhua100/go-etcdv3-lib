package main

import (
	//"context"
	"fmt"
	"time"
	"go-etcdv3-lib/utils"
	//"github.com/coreos/etcd/clientv3"
)



func main() {
	G_etcd, err := utils.Client()
        if err != nil {
		fmt.Println("etcd open error.")
                return	
	}
	G_etcd.PutKeyValue("/demo/A/B", "hello world")
	G_etcd.GetKeyValue("/demo/A/B")
	G_etcd.DeleteKey("/demo/A/B")
	G_etcd.PutKeyWithLease("/demo/A/B", "hello world ttl-100", 100)
	go G_etcd.Watch("/demo/A/")	

	go func() {
        for {
	    G_etcd.PutKeyValue("/demo/A/B", "hello world")
		time.Sleep(time.Second)
        	G_etcd.DeleteKey("/demo/A/B")
            time.Sleep(1 * time.Second)


            G_etcd.PutKeyValue("/demo/A/C", "qinghua")
                time.Sleep(time.Second)
                G_etcd.DeleteKey("/demo/A/C")
            time.Sleep(1 * time.Second)
        }
    }()
	
	for {

	}

}

