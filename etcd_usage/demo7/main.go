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
		err error
		kv clientv3.KV
		watcher clientv3.Watcher
		//putResp *clientv3.PutResponse
		getResp *clientv3.GetResponse
		watchStartRevision int64
		watchRespChan <-chan clientv3.WatchResponse
		watchRsp clientv3.WatchResponse
		event *clientv3.Event
	)
	config = clientv3.Config{
		Endpoints:[]string{"127.0.0.1:2379"},
		DialTimeout:5 *time.Second,
	}
	//建立连接
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}

	//用于读写etcd键值对
	kv = clientv3.NewKV(client)

	//模拟etcd中kv的变化
	go func() {
		for {
			kv.Put(context.TODO(), "/corn/jobs/job1", "i am job1")
			kv.Delete(context.TODO(), "/corn/jobs/job1")
			time.Sleep(1 * time.Second)
		}
	}()

	//先get当前的值 再监听变化
	if getResp,err = kv.Get(context.TODO(),"/corn/jobs/job1"); err != nil{
		fmt.Println(err)
		return
	}
	//现在key值存在
	if len(getResp.Kvs) != 0 {
		fmt.Println("当前值",string(getResp.Kvs[0].Value))
	}
	//当前etcd集群事物id
	watchStartRevision = getResp.Header.Revision + 1

	//创建一个监听器
	watcher = clientv3.NewWatcher(client)

	//启动监听
	fmt.Println("从该版本向后监听",watchStartRevision)


	//模拟一个取消操作
	ctx,cancelFunc := context.WithCancel(context.TODO())
	time.AfterFunc(5 *time.Second, func() {
		cancelFunc()
	})

	watchRespChan = watcher.Watch(ctx,"/corn/jobs/job1",clientv3.WithRev(watchStartRevision))




	//处理变化事件
	for watchRsp = range watchRespChan {
		for _,event = range watchRsp.Events {
			switch event.Type {
			case mvccpb.PUT:
				fmt.Println("修改操作",string(event.Kv.Value),"Revision",event.Kv.CreateRevision,event.Kv.ModRevision)
			case mvccpb.DELETE:
				fmt.Println("删除操作","Revision",event.Kv.ModRevision)
			}
		}
	}
}
