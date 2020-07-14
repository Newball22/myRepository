package engine

import "log"

//任务调度器
type Scheduler interface {
	Submit(request Request)              //提交任务
	ConfigMasterWorkerChan(chan Request) //配置初始请求任务
	WorkerReady(w chan Request)//{第二级别}
	Run()//{第二级别}
}

//定义并发引擎结构
type ConcurrentEngine struct {
	Scheduler   Scheduler //任务调度器
	WorkerCount int       //任务并发数量
}

//多线程使用的方法
func (e *ConcurrentEngine) Run(seeds ...Request) {
	log.Println("=====ConcurrentEngine'Run Start=====")
	//in := make(chan Request)               //scheduler的输入请求{第二级别}
	out := make(chan ParseResult)          //worker的输出
	e.Scheduler.Run()
	//e.Scheduler.ConfigMasterWorkerChan(in) //把初始请求交给scheduler{第二级别}
	//创建goroutine
	for i := 0; i < e.WorkerCount; i++ {
		creatWorker(out, e.Scheduler,i)
	}
	/*注意这两个'for'的先后顺序，因为先不把请求Submit到管道，所以形成阻塞，可以先创建50个worker*/
	//engine把请求任务提交给Scheduler
	for _, request := range seeds {
		e.Scheduler.Submit(request)
	}
	itemCount := 0
	for {
		//接受Worker的解析结果
		result := <-out
		for _, item := range result.Items {
			log.Printf("Got item:#%d:%v\n", itemCount, item)
			itemCount++

		}
		//把Worker解析出的Request送给 Scheduler
		for _, request := range result.Requests {
			e.Scheduler.Submit(request)
		}
	}

}

//创建任务，调用worker，分发goroutine
func creatWorker( out chan ParseResult,s Scheduler, i int) {
	log.Printf("====creatWorker Start 这是第%d个Worker====\n", i)
	//为每一个Worker创建一个channel{第三级别}
	in:=make(chan Request)
	go func() {
		for {
			//log.Println("request没有Submit到管道in的时候，in管道没request输出，暂时阻塞，只能走到这里"){第二级别}
			s.WorkerReady(in)//告诉调度器任务空闲
			request := <-in
			//log.Println("request已经Submit到管道in，in管道有request输出，取消阻塞，可以继续往下面走了"){第二级别}
			result, err := worker(request, i) //每个worker接受一个请求
			if err != nil {
				continue
			}
			out <- result
		}
	}()
}
