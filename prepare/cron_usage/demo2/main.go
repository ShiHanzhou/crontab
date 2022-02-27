package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)

//代表一个任务
type CronJob struct {
	expr     *cronexpr.Expression
	nextTime time.Time //通过expr.Next(now)得到下一次的调度时间
}

func main() {
	//需要一个调度协程，定时检查所有的cron任务，谁过期了就执行谁
	var (
		cronJob       *CronJob
		expr          *cronexpr.Expression
		now           time.Time
		scheduleTable map[string]*CronJob
	)

	scheduleTable = make(map[string]*CronJob)

	now = time.Now()

	expr = cronexpr.MustParse("*/5 * * * * * *")
	cronJob = &CronJob{
		expr:     expr,
		nextTime: expr.Next(now),
	}
	//任务注册到调度表
	scheduleTable["job1"] = cronJob

	expr = cronexpr.MustParse("*/5 * * * * * *")
	cronJob = &CronJob{
		expr:     expr,
		nextTime: expr.Next(now),
	}
	//任务注册到调度表
	scheduleTable["job2"] = cronJob

	go func() {
		var (
			jobName string
			cronjob *CronJob
			now     time.Time
		)
		//定时检查任务调度表
		for {
			now = time.Now()

			for jobName, cronjob = range scheduleTable {
				//判断是否过期
				if cronjob.nextTime.Before(now) || cronjob.nextTime.Equal(now) {
					//启动一个协程执行这个任务
					go func(jobName string) {
						fmt.Println("执行：", jobName)
					}(jobName)

					//计算下一次调度时间
					cronjob.nextTime = cronjob.expr.Next(now)
					fmt.Println("下次执行时间：", cronjob.nextTime)
				}
			}

			select {
			case <-time.NewTimer(100 * time.Millisecond).C: //100毫秒可读一次
			}
		}
	}()

	time.Sleep(100 * time.Second)
}
