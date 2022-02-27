package worker

import (
	"context"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"prepare/crontab/common"
	"time"
)

//任务管理器
type JobMgr struct {
	client  *clientv3.Client
	kv      clientv3.KV
	lease   clientv3.Lease
	watcher clientv3.Watcher
}

var (
	//单例
	G_jobMgr *JobMgr
)

//监听任务变化
func (jobMgr *JobMgr) watchJobs() (err error) {
	var (
		getResp            *clientv3.GetResponse
		kvPair             *mvccpb.KeyValue
		job                *common.Job
		watchStartRevision int64
		watchChan          clientv3.WatchChan
		watchResp          clientv3.WatchResponse
		watchEvent         *clientv3.Event
		jobName            string
		jobEvent           *common.JobEvent
	)

	//1、get到/cron/jobs/目录下的所有任务，并且获知当前集群的revision
	if getResp, err = jobMgr.kv.Get(context.TODO(), common.JOB_SAVE_DIR, clientv3.WithPrefix()); err != nil {
		return
	}

	//当前有哪些任务
	for _, kvPair = range getResp.Kvs {
		//反序列化得到job
		if job, err = common.UnpackJob(kvPair.Value); err == nil {
			jobEvent = common.BuildJobEvent(common.JOB_EVENT_SAVE, job)
			//把job同步给scheduler(调度协程)
			G_scheduler.PushJobEvent(jobEvent)
		}
	}

	//2、从该revision向后监听变化事件
	go func() {
		//从get时刻的后续版本开始监听变化
		watchStartRevision = getResp.Header.Revision + 1
		//启动监听/cron/jobs/目录的后续变化
		watchChan = jobMgr.watcher.Watch(context.TODO(), common.JOB_SAVE_DIR, clientv3.WithRev(watchStartRevision), clientv3.WithPrefix())
		//处理监听事件
		for watchResp = range watchChan {
			for _, watchEvent = range watchResp.Events {
				switch watchEvent.Type {
				case mvccpb.PUT: //任务保存事件
					if job, err = common.UnpackJob(watchEvent.Kv.Value); err != nil {
						continue
					}
					//构造一个更新event
					jobEvent = common.BuildJobEvent(common.JOB_EVENT_SAVE, job)
				case mvccpb.DELETE: //任务删除事件
					jobName = common.ExtractJobName(string(watchEvent.Kv.Key))

					job = &common.Job{Name: jobName}

					//构造一个删除event
					jobEvent = common.BuildJobEvent(common.JOB_EVENT_DELETE, job)
				}
				//推送事件给scheduler
				G_scheduler.PushJobEvent(jobEvent)
			}
		}
	}()

	return
}

//监听强杀任务通知
func (jobMgr *JobMgr) watchKiller() {
	var (
		watchChan  clientv3.WatchChan
		watchResp  clientv3.WatchResponse
		watchEvent *clientv3.Event
		jobEvent   *common.JobEvent
		jobName    string
		job        *common.Job
	)
	//监听/cron/killer目录
	go func() {
		//启动监听/cron/killer/目录的变化
		watchChan = jobMgr.watcher.Watch(context.TODO(), common.JOB_KILLER_DIR, clientv3.WithPrefix())
		//处理监听事件
		for watchResp = range watchChan {
			for _, watchEvent = range watchResp.Events {
				switch watchEvent.Type {
				case mvccpb.PUT: //杀死任务事件
					jobName = common.ExtractKillerName(string(watchEvent.Kv.Key))
					job = &common.Job{Name: jobName}
					jobEvent = common.BuildJobEvent(common.JOB_EVENT_KILL, job)

					//推送事件给scheduler
					G_scheduler.PushJobEvent(jobEvent)
				case mvccpb.DELETE: //任务删除事件

				}

			}
		}
	}()
}

//初始化管理器
func InitJobMgr() (err error) {
	var (
		config  clientv3.Config
		client  *clientv3.Client
		kv      clientv3.KV
		lease   clientv3.Lease
		watcher clientv3.Watcher
	)

	//初始化配置
	config = clientv3.Config{
		Endpoints:   G_config.EtcdEndPoints,                                     //集群地址
		DialTimeout: time.Duration(G_config.EtcdDialTimeout) * time.Millisecond, //连接超时
	}

	//建立连接
	if client, err = clientv3.New(config); err != nil {
		return
	}

	//得到KV和Lease的API子集
	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)
	watcher = clientv3.NewWatcher(client)

	//赋值单例
	G_jobMgr = &JobMgr{
		client:  client,
		kv:      kv,
		lease:   lease,
		watcher: watcher,
	}

	//启动任务监听
	G_jobMgr.watchJobs()

	//启动监听killer
	G_jobMgr.watchKiller()

	return
}

//创建任务执行锁
func (jobMgr *JobMgr) CreateJobLock(jobName string) (jobLock *JobLock) {
	//返回一把锁
	jobLock = InitJobLock(jobName, jobMgr.kv, jobMgr.lease)
	return
}
