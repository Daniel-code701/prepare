package main

import (
	"context"
	//"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	//"github.com/coreos/etcd/mvcc/mvccpb"
	"time"
)

func main() {
	var (
		config         clientv3.Config
		client         *clientv3.Client
		err            error
		lease          clientv3.Lease
		leaseGrantResp *clientv3.LeaseGrantResponse
		leaseld        clientv3.LeaseID
		putResp        *clientv3.PutResponse
		getResp        *clientv3.GetResponse
		kv             clientv3.KV
		keepResp       *clientv3.LeaseKeepAliveResponse
		keepRespChan   <-chan *clientv3.LeaseKeepAliveResponse
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

	//申请一个租约
	lease = clientv3.NewLease(client)

	//申请一个10秒的租约
	if leaseGrantResp, err = lease.Grant(context.TODO(), 10); err != nil {
		fmt.Println(err)
		return
	}

	//拿到租约的Id
	leaseld = leaseGrantResp.ID

	//测试代码 5秒后取消续租
	//ctx, _ := context.WithTimeout(context.TODO(),5 *time.Second)
	//if keepRespChan,err = lease.KeepAlive(ctx,leaseld);err != nil {
	//	fmt.Println(err)
	//	return
	//}

	//如果租约续住成功了 sdk会自动启动一个协程续租
	//租约续租
	if keepRespChan, err = lease.KeepAlive(context.TODO(), leaseld); err != nil {
		fmt.Println(err)
		return
	}

	//处理续租应答的协成 续租成功后 启动一个协程来消费
	go func() {
		for {
			select {
			case keepResp = <-keepRespChan:
				if keepRespChan == nil {
					fmt.Println("租约终止")
					goto END
				} else { //每秒会续租一次 所以会收到一次应答
					fmt.Println("收到自动续租应答", keepResp.ID)
				}

			}
		}
	END:
	}()

	//获得kv对象子集
	kv = clientv3.NewKV(client)

	//申请成功put一个kv 与租约关联起来
	if putResp, err = kv.Put(context.TODO(), "/corn/jobs/job1", "hello11111", clientv3.WithLease(leaseld)); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("写入成功", putResp.Header.Revision)

	//定时看一下key过期没有
	for {
		if getResp, err = kv.Get(context.TODO(), "/corn/jobs/job1"); err != nil {
			fmt.Println(err)
		}
		if getResp.Count == 0 {
			fmt.Println("kv过期了")
			break
		}
		fmt.Println("kv没过期", getResp.Kvs)
		time.Sleep(2 * time.Second)
	}

}
