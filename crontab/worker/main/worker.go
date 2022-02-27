package main

import (
	"flag"
	"fmt"
	"prepare/crontab/worker"
	"runtime"
	"time"
)

var (
	confile string //配置文件路径
)

//解析命令行参数
func initArgs() {
	//worker -config ./worker.json
	//worker -h
	flag.StringVar(&confile, "config", "./worker.json", "指定worker.json")
	flag.Parse()
}

//初始化线程
func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var (
		err error
	)

	//初始化命令行参数
	initArgs()

	//初始化线程
	initEnv()

	//加载配置
	if err = worker.InitConfig(confile); err != nil {
		goto ERR
	}

	//服务注册
	if err = worker.InitRegister(); err != nil {
		goto ERR
	}

	//启动日志协程
	if err = worker.InitLogSink(); err != nil {
		goto ERR
	}

	//加载执行器
	if err = worker.InitExecutor(); err != nil {
		goto ERR
	}

	//启动调度器
	if err = worker.InitScheduler(); err != nil {
		goto ERR
	}

	//初始化任务管理器
	if err = worker.InitJobMgr(); err != nil {
		goto ERR
	}

	for {
		time.Sleep(1 * time.Second)
	}
	//正常退出
	return

ERR:
	fmt.Println(err)
}
