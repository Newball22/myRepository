package main

import (
	"crawler/engine"
	"crawler/scheduler"
	"crawler/zhenai/parser"
)

func main() {
	e := engine.ConcurrentEngine{
		//Scheduler: &scheduler.SimpleScheduler{}, //这里是调用并发调度器{第二级别}
		Scheduler:   &scheduler.QueuedScheduler{}, //这里调用任务列表{第三级别}
		WorkerCount: 50,
	}
	e.Run(engine.Request{
		Url:       "http://www.zhenai.com/zhenghun",
		ParseFunc: parser.ParseCityList, //第一次是先爬取所有的城市列表
	})

}
