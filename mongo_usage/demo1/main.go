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
		databases  *mongo.Database
		collection *mongo.Collection
	)
	//1.建立链接
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	if client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017")); err != nil {
		fmt.Println(err)
	}
	//2.选择数据库
	databases = client.Database("my_db")
	//3.选择表
	collection = databases.Collection("my_collection")
	collection = collection
}
