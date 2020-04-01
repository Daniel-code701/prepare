package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"time"
)

func main() {
	var (
		config clientv3.Config
		client *clientv3.Client
		err    error
		kv     clientv3.KV
		desRes *clientv3.DeleteResponse
		kvpair *mvccpb.KeyValue
	)
	config = clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}
	//建立连接
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println()
	}

	//用于读写etcd键值对
	kv = clientv3.NewKV(client)

	if desRes, err = kv.Delete(context.TODO(), "/corn/jobs/job1", clientv3.WithPrefix()); err != nil {
		fmt.Println(err)
		return
	}
	//被删除之前的value
	if len(desRes.PrevKvs) != 0 {
		for _, kvpair = range desRes.PrevKvs {
			fmt.Println("删除了", string(kvpair.Key), string(kvpair.Value))
		}
	}
}
