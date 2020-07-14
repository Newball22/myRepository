package scheduler

import (
	"crawler/engine"
)

//使用队列来调度任务
type QueuedScheduler struct {
	requestChan chan engine.Request      //Request channel
	workerChan  chan chan engine.Request //Worker channel.其中每一个Worker是一个chan engine.Request类型
}

//提交请求任务到requestChannel
func (q *QueuedScheduler) Submit(request engine.Request) {
	q.requestChan <- request
}

func (q *QueuedScheduler) ConfigMasterWorkerChan(chan engine.Request) {
	panic("implement me")
}

//告诉外界有一个worker可以接收request{意思是工作队列已经有一个'工作channel'进来了在等着工作}
func (q *QueuedScheduler) WorkerReady(w chan engine.Request) {
	q.workerChan <- w
}

/*
包含关系：requestChan->requestQueued->requestQueued[0]->activeRequest(engine.Request类型)
		workerChan->workerQueued->workerQueued[0]->activeWorker(chan engine.Request类型)
*/
func (q *QueuedScheduler) Run() {
	//初始化请求队列和工作队列的存放环境
	q.requestChan = make(chan engine.Request)
	q.workerChan = make(chan chan engine.Request)
	go func() {
		//请求队列
		var requestQueued []engine.Request
		//工作队列
		var workerQueued []chan engine.Request
		for {
			var activeRequest engine.Request
			var activeWorker chan engine.Request
			if len(requestQueued) > 0 && len(workerQueued) > 0 {
				activeRequest = requestQueued[0]
				activeWorker = workerQueued[0]

			}

			select {
			//当requestChan收到数据,放到队列
			case r := <-q.requestChan:
				requestQueued = append(requestQueued, r)
			//当workerChan收到数据，放到队列
			case w := <-q.workerChan:
				workerQueued = append(workerQueued, w)
			//当请求队列和工作队列都不为空时，给任务列表分配任务
			case activeWorker <- activeRequest:
				requestQueued = requestQueued[1:]
				workerQueued = workerQueued[1:]
			}
		}
	}()
}
