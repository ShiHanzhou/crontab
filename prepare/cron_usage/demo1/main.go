package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)

func main() {
	var (
		expr     *cronexpr.Expression
		err      error
		now      time.Time
		nextTime time.Time
	)
	//linux crontab 支持五个时间粒度：
	//分钟(0-59)， 小时(0-23)，日(1-31)，月(1-12)，星期(0-7, 0和7都是周日)
	//cronexpr库比crontab多支持两个时间粒度：秒(0-59)，年(?-2099)

	//每5分钟执行一次
	if expr, err = cronexpr.Parse("*/5 * * * * * *"); err != nil {
		fmt.Println(err)
		return
	}

	now = time.Now()
	//得到下次调度时间
	nextTime = expr.Next(now)

	//计时，经过nextTime减去now时间的时间间隔，执行func
	time.AfterFunc(nextTime.Sub(now), func() {
		fmt.Println("被调度了：", nextTime)
	})

	time.Sleep(5 * time.Second)
}
