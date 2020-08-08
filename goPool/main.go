package main

import (
	"fmt"
	"time"
)

//任务结构体
type Task struct {
	f func() error //是一个方法类型
}

//创建一个Task类型的任务
func NewTask(argFunc func() error) *Task {
	t := Task{
		f: argFunc,
	}
	return &t
}

//Task需要一个执行业务的方法，就是提交到入口
func (t *Task) Execute() {
	t.f()
}

//协程池类型
type Pool struct {
	EntryChannel chan *Task
	JobChannel   chan *Task
	WorkNumMax   int //协程池goroutine最大的数量
}

//创建Pool的函数
func NewPool(cap int) *Pool {
	p := Pool{
		EntryChannel: make(chan *Task),
		JobChannel:   make(chan *Task),
		WorkNumMax:   cap,
	}
	return &p

}

//协程池创建一个Worker，并且让这个Worker去工作
func (p *Pool) worker(workerID int) {
	//从JobChannel拿任务交给worker
	for task := range p.JobChannel {
		task.Execute()
		fmt.Printf("woker%d开始了工作\n", workerID)
	}

}

//
func (p *Pool) Run() {
	//创建worker
	for i := 0; i < p.WorkNumMax; i++ {
		go p.worker(i)

	}
	//不断从入口获取任务,发送给JobChannel
	for task := range p.EntryChannel {
		p.JobChannel <- task
	}

}

//任务函数
func taskTest() error {
	fmt.Println(time.Now())
	return nil
}

func main() {
	//创建任务
	task := NewTask(taskTest)
	//初始化协程池
	p := NewPool(5)
	var num = 1
	//把任务交给Pool池
	go func() {
		for {
			p.EntryChannel <- task
			num++
			fmt.Printf("总共处理了%d个任务\n", num,num)
		}
	}()
	p.Run()
}
