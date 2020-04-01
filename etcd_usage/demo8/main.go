package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

func main() {
	var (
		config clientv3.Config
		client *clientv3.Client
		err    error
		kv     clientv3.KV
		putOp  clientv3.Op
		getOp  clientv3.Op
		opResp clientv3.OpResponse
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
		fmt.Println("test")
	}
	kv = clientv3.NewKV(client)

	//OP方法 创建OP
	putOp = clientv3.OpPut("/cron/lock/job10/", "dddd")
	//执行OP
	if opResp, err = kv.Do(context.TODO(), putOp); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Revision", opResp.Put().Header.Revision)

	//Get操作
	getOp = clientv3.OpGet("/cron/lock/job10/")
	//执行OP
	if opResp, err = kv.Do(context.TODO(), getOp); err != nil {
		fmt.Println(err)
		return
	}

	//打印
	fmt.Println("Revision", opResp.Get().Kvs[0].ModRevision)
	fmt.Println("数据value", string(opResp.Get().Kvs[0].Value))

}
