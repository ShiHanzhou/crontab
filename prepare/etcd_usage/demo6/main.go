package main

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

func main() {
	var (
		config         clientv3.Config
		client         *clientv3.Client
		err            error
		lease          clientv3.Lease
		leaseGrantResp *clientv3.LeaseGrantResponse
		kv             clientv3.KV
		leaseId        clientv3.LeaseID
		putResp        *clientv3.PutResponse
		getResp        *clientv3.GetResponse
		keepResp       *clientv3.LeaseKeepAliveResponse
		keepRespChan   <-chan *clientv3.LeaseKeepAliveResponse
		ctx            context.Context
	)

	config = clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}

	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}

	//申请一个lease(租约)
	lease = clientv3.NewLease(client)

	//申请一个5秒的租约
	if leaseGrantResp, err = lease.Grant(context.TODO(), 5); err != nil {
		fmt.Println(err)
		return
	}

	//得到lease的ID
	leaseId = leaseGrantResp.ID

	//lease续租
	ctx, _ = context.WithTimeout(context.TODO(), 5*time.Second)
	//5秒后取消自动续租
	if keepRespChan, err = lease.KeepAlive(ctx, leaseId); err != nil {
		fmt.Println(err)
		return
	}

	//处理续租应答
	go func() {
		for {
			select {
			case keepResp = <-keepRespChan:
				if keepResp == nil {
					fmt.Println("lease已过期")
					goto END
				}
				fmt.Println("收到自动续租应答:", keepResp.ID)
			}
		}
	END:
	}()

	//获得kv API子集
	kv = clientv3.NewKV(client)

	//put一个kv，与5秒租约关联
	if putResp, err = kv.Put(context.TODO(), "/cron/jobs/job1", "", clientv3.WithLease(leaseId)); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("写入成功", putResp.Header.Revision)

	for {
		if getResp, err = kv.Get(context.TODO(), "/cron/jobs/job1"); err != nil {
			fmt.Println(err)
			return
		}

		if getResp.Count == 0 {
			fmt.Println("kv过期了")
			break
		}
		fmt.Println(getResp.Kvs, "未过期")
		time.Sleep(2 * time.Second)
	}
}
