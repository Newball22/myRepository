package engine

import "log"

type ConcurrendEngine struct {
	Scheduler   Scheduler //任务调度器
	WorkerCount int       //任务并发数
}

type Scheduler interface {
	Submit(request Request)              //提交任务
	   ConfigMasterWorkerChan(chan Request) //配置初始请求任务
}

func (e *ConcurrendEngine) Run(seeds ...Request) {
	in := make(chan Request)
	out := make(chan ParseResult)
	e.Scheduler.ConfigMasterWorkerChan(in)

	for i := 0; i < e.WorkerCount; i++ {
		creatWorker(in, out)
	}
	for _, request := range seeds {
		e.Scheduler.Submit(request)
	}
	itemCount := 0
	for {
		//接受Worker的解析结果
		result := <-out
		for _, item := range result.Items {
			log.Printf("Got item:#%d: %v\n", itemCount, item)
			itemCount++
		}
		for _, request := range result.Request {
			e.Scheduler.Submit(request)
		}
	}
}

//创建任务，调用worker，分发goroutine
func creatWorker(in chan Request, out chan ParseResult) {
	go func() {
		for {
			request := <-in
			result, err := worker(request)
			if err != nil {
				continue
			}
			out <- result
		}
	}()
}
