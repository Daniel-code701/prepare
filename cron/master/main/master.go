package main

import (
	"flag"
	"fmt"
	"prepare/cron/master"
	"runtime"
	"time"
)

var (
	confFile string //配置文件路径
)

//解析命令行参数
func initArgs() {
	//master-config ./master.json
	flag.StringVar(&confFile, "config", "./master.json", "指定master.json")
	flag.Parse()
}

//让线程数和CPU核心数量相等
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
	if err = master.InitConfig(confFile); err != nil {
		goto ERR
	}

	//启动任务管理器
	if err = master.InitJobMgr(); err != nil {
		goto ERR
	}

	// 启动Api HTTP服务
	if err = master.InitApiServer(); err != nil {
		goto ERR
	}

	//正常退出
	for {
		time.Sleep(1 * time.Second)
	}
	return

	//错误打印到屏幕上 如果启动失败 可以看到错误原因
ERR:
	fmt.Println(err)
}
