package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func main() {
	var (
		client     *mongo.Client
		err        error
		ctx        context.Context
		cancelFunc context.CancelFunc
		database   *mongo.Database
		collection *mongo.Collection
	)

	//1、建立连接
	ctx, cancelFunc = context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancelFunc()
	if client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017")); err != nil {
		fmt.Println(err)
		return
	}

	//2、选择数据库my_db
	database = client.Database("my_db")

	//3、选择表my_collection
	collection = database.Collection("my_collection")
	_ = collection
}
