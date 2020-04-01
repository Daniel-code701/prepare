package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)

//代表一个任务
type CronJob struct {
	expr *cronexpr.Expression
	//下一次任务执行时间
	nextTime time.Time
}

func main()  {
	//需要有一个调度协程 它定时检查所有的Cron任务 谁过期了就执行谁

	var (
		cronJob *CronJob
		expr *cronexpr.Expression
		now time.Time
		//定义一个调度表
		scheduleTable map[string]*CronJob //key:任务的名字
	)

	scheduleTable = make(map[string]*CronJob)

	//当前时间
	now = time.Now()

	//1 定义2个cronjob
	expr = cronexpr.MustParse("*/5 * * * * * *")
	cronJob = &CronJob{
		expr:     expr,
		nextTime: expr.Next(now),
	}

	//任务注册到调度表
	scheduleTable["job1"] = cronJob

	//1 定义2个cronjob
	expr = cronexpr.MustParse("*/5 * * * * * *")
	cronJob = &CronJob{
		expr:     expr,
		nextTime: expr.Next(now),
	}

	//任务注册到调度表
	scheduleTable["job2"] = cronJob

	//启动一个调度协成
	go func() {
		var (
			jobName string
			cronJob *CronJob
			now time.Time
		)
		//定时检查一下调度任务表
		for {
			now = time.Now()

			for jobName,cronJob = range scheduleTable{
				//判断是否过期 如果下次执行任务时间早于当前
				if cronJob.nextTime.Before(now) || cronJob.nextTime.Equal(now) {
					// 启动一个协成 执行这个任务
					go func(jobName string) {
						fmt.Println("执行",jobName)
					}(jobName)

					//计数下一次操作时间
					cronJob.nextTime = cronJob.expr.Next(now)
					fmt.Println("下次执行时间是",cronJob.nextTime)

				}
			}
			//睡眠100毫秒
			select {
			case <- time.NewTimer(100 *time.Millisecond).C:
			}
		}
	}()
	time.Sleep(100 *time.Second)
}
