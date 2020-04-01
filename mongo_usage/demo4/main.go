package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

//任务执行时间点
type TimePoint struct {
	StartTime int64 `bson:"startTime"`
	EndTime   int64 `bson:"endTime"`
}

//jobname过滤条件
type FindByJobName struct {
	JobName string `bson:"job_name"` //job_name赋值给job10
}

//小于某时间
//{"$lt":timestamp}
type TimeBeforeCond struct {
	Before int64 `bson:"$lt"`
}

//{"timePoint.startTime":{"$lt":当前时间}}
type DeleteCond struct {
	BeforeCond TimeBeforeCond `bson:"timePoint.startTime"`
}

//一条日志
type LogRecord struct {
	JobName   string    `bson:"job_name"`  //任务名
	Command   string    `bson:"command"`   //shell命令
	Err       string    `bson:"err"`       //脚本错误
	Content   string    `bson:"content"`   //脚本输出
	TimePoint TimePoint `bson:"timePoint"` //执行时间
}

func main() {
	var (
		client     *mongo.Client
		err        error
		databases  *mongo.Database
		collection *mongo.Collection
		record     *LogRecord
		cond       *FindByJobName
		cursor     *mongo.Cursor
		delCond    *DeleteCond
		delResult  *mongo.DeleteResult
	)
	//1.建立链接
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	if client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017")); err != nil {
		fmt.Println(err)
	}
	//2.选择数据库
	databases = client.Database("cron")
	//3.选择表
	collection = databases.Collection("log")

	//5.按照jobName字段过滤 想找出jobName=job10 找出五条
	cond = &FindByJobName{JobName: "job10"} //{"jobName":"job10"}

	//6.发起查询
	if cursor, err = collection.Find(context.TODO(), cond, options.Find().SetSkip(0).SetLimit(2)); err != nil {
		fmt.Println(err)
		return
	}

	//延迟释放游标 如果在查询或者其它函数中函数中断执行 使用defer释放游标
	defer cursor.Close(context.TODO())

	//7.遍历结果集
	for cursor.Next(context.TODO()) {
		//定义一个日志对象
		record = &LogRecord{}
		//反序列化bson到对象
		if err = cursor.Decode(record); err != nil {
			fmt.Println(err)
			return
		}
		//日志行打印出来
		fmt.Println(*record)

	}

	//8.要删除时间早于当前时间的日志
	//delete({"timePoint.startTime":{"$lt":当前时间}})
	delCond = &DeleteCond{BeforeCond: TimeBeforeCond{Before: time.Now().Unix()}}

	if delResult, err = collection.DeleteMany(context.TODO(), delCond); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("删除了", delResult.DeletedCount)

}
