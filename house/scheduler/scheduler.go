package scheduler

import "MyProject/house/engine"

type SimpleScheduler struct {
	workChan chan engine.Request
}

func (s *SimpleScheduler) Submit(request engine.Request) {
	//为每一个Request创建goroutine
	go func() {
		s.workChan <- request
	}()
}

func (s *SimpleScheduler) ConfigMasterWorkerChan(in chan engine.Request) {
	s.workChan = in
}
