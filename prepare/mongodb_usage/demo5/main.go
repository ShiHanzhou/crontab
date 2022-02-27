package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

//{"$lt":timestamp}
type TimeBeforeCondition struct {
	Before int64 `bson:"$lt"`
}

//{"timePoint.startTime":{"$lt":timestamp}}
type DeleteCond struct {
	beforeCond TimeBeforeCondition `bson:"timePoint.startTime"`
}

func main() {
	var (
		client     *mongo.Client
		err        error
		ctx        context.Context
		cancelFunc context.CancelFunc
		database   *mongo.Database
		collection *mongo.Collection
		delCond    interface{}
		delResult  *mongo.DeleteResult
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

	//4、删除开始事件早于当前时间的所有日志
	//delete({"timePoint.startTime": {"$lt":当前时间}})
	delCond = &DeleteCond{beforeCond: TimeBeforeCondition{Before: time.Now().Unix()}}
	if delResult, err = collection.DeleteMany(context.TODO(), delCond); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("删除的行数", delResult.DeletedCount)
}
