package main

import (
	"MyProject/house/engine"
	"MyProject/house/fangtianxia/parser"
	"MyProject/house/scheduler"
)

func main() {
	e := engine.ConcurrendEngine{
		Scheduler:   &scheduler.SimpleScheduler{},
		WorkerCount: 50,
	}
	e.Run(engine.Request{
		Url:       "https://tj.esf.fang.com/newsecond/esfcities.aspx",
		ParseFunc: parser.ParseCityList,
	})
	//
	//engine.Run(engine.Request{
	//	Url:       "",
	//	ParseFunc: parser.ParseCityList,
	//})

}
