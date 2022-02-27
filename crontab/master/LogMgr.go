package master

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"prepare/crontab/common"
	"time"
)

//mongodb日志管理
type LogMgr struct {
	client        *mongo.Client
	logCollection *mongo.Collection
}

var (
	G_logMgr *LogMgr
)

func InitLogMgr() (err error) {
	var (
		client     *mongo.Client
		ctx        context.Context
		cancelFunc context.CancelFunc
	)

	ctx, cancelFunc = context.WithTimeout(context.TODO(), time.Duration(G_config.EtcdDialTimeout)*time.Millisecond)
	defer cancelFunc()
	if client, err = mongo.Connect(ctx, options.Client().ApplyURI(G_config.MongodbUri)); err != nil {
		return
	}

	//选择db和collection
	G_logMgr = &LogMgr{
		client:        client,
		logCollection: client.Database("cron").Collection("log"),
	}
	return
}

//查看任务日志
func (logMgr *LogMgr) ListLog(name string, skip int64, limit int64) (logArr []*common.JobLog, err error) {
	var (
		filter  *common.JobLogFilter
		logSort *common.SortLogByStartTime
		cursor  *mongo.Cursor
		jobLog  *common.JobLog
	)

	//len(logArr)
	logArr = make([]*common.JobLog, 0)

	//过滤条件
	filter = &common.JobLogFilter{JobName: name}

	//按照任务时间倒排
	logSort = &common.SortLogByStartTime{SortOrder: -1}

	//查询
	if cursor, err = logMgr.logCollection.Find(context.TODO(), filter, options.Find().SetSkip(skip).SetLimit(limit).SetSort(logSort)); err != nil {
		return
	}
	//延迟释放游标
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		jobLog = &common.JobLog{}

		//反序列化bson
		if err = cursor.Decode(jobLog); err != nil {
			continue //有日志不合法
		}

		logArr = append(logArr, jobLog)
	}

	return
}
