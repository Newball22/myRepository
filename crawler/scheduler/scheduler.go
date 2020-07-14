package scheduler

import (
	"crawler/engine"
	"log"
)

//这个文件是给第二级并发改进用的
type SimpleScheduler struct {
	workerChan chan engine.Request
}

//把初始请求发送给Scheduler
func (s *SimpleScheduler) ConfigMasterWorkerChan(in chan engine.Request) {
	log.Println("====ConfigMasterWorkerChan Start====")
	s.workerChan = in
}

func (s *SimpleScheduler) Submit(request engine.Request) {
	log.Println("====Submit Start====")
	//为每一个Request创建goroutine
	go func() {
		s.workerChan <- request
	}()

}
