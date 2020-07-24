package main

import (
	"errors"
	"math"
	"sync"
	"time"
)

type sig struct {
}

type f func() error

type Worker struct {
	pool        *Pool
	task        chan f
	recycleTime time.Time
}

//Goroutine池类型
type Pool struct {
	//Goroutine池的容量
	capacity int

	//正在执行任务的worker数量
	running int

	//worker过期的时间段
	expiryDuration time.Duration

	//通知Goroutine工作的信号
	freeSignal chan sig

	//存放空闲的worker的切片
	workers []*Worker

	//用来通知Goroutine池关闭
	release chan sig

	//锁
	lock sync.Locker

	//只执行一次Pool关闭操作的关键字
	once sync.Once
}

func NewPool(size, expiry int) (*Pool, error) {
	if size <= 0 {
		return nil, errors.New("请输入正确的参数")
	}
	p := &Pool{
		capacity:       size,
		expiryDuration: time.Duration(expiry) * time.Second,
		freeSignal:     make(chan sig, math.MaxInt32),
		release:        make(chan sig, 1),
		once:           sync.Once{},
	}
}

func main() {

}
