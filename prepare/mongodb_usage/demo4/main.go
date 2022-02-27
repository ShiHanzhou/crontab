package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type TimePoint struct {
	StartTime int64 `bson:"startTime"`
	EndTime   int64 `bson:"endTime"`
}

type LogRecord struct {
	JobName   string    `bson:"jobName"`   //任务名
	Command   string    `bson:"command"`   //shell命令
	Err       string    `bson:"err"`       //脚本错误
	Content   string    `bson:"content"`   // 脚本输出
	TimePoint TimePoint `bson:"timePoint"` //执行时间
}

type FindByJobName struct {
	JobName string `bson:"jobName"`
}

func main() {
	//mongodb读取回来的时bson，需要反序列化成LogRecord对象
	var (
		client      *mongo.Client
		err         error
		ctx         context.Context
		cancelFunc  context.CancelFunc
		database    *mongo.Database
		collection  *mongo.Collection
		cond        *FindByJobName
		cursor      *mongo.Cursor
		record      *LogRecord
		findOptions *options.FindOptions
	)

	//1、建立连接
	ctx, cancelFunc = context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancelFunc()
	if client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017")); err != nil {
		fmt.Println(err)
		return
	}

	//2、选择数据库my_db
	database = client.Database("cron")

	//3、选择表my_collection
	collection = database.Collection("log")

	//4、按照jobName字段过滤，想找出jobName=job10，找出5条
	cond = &FindByJobName{JobName: "job10"}

	//5、查询(过滤 + 翻页参数)
	findOptions = options.Find()
	findOptions.SetSkip(0)
	findOptions.SetLimit(2)

	if cursor, err = collection.Find(context.TODO(), cond, findOptions); err != nil {
		fmt.Println(err)
		return
	}

	//6、遍历结果集
	for cursor.Next(context.TODO()) {
		record = &LogRecord{}

		//反序列化bson到record中
		if err = cursor.Decode(record); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(*record)
	}

	defer cursor.Close(context.TODO())
}
