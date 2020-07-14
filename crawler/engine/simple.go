package engine

import (
	"log"
)

type SimpleEngine struct {
}

//任务执行函数{单线程使用的方法}
//func Run(seeds ...Request) { //这里一定要有'...'要不然不可以range遍历
func (s SimpleEngine) Run(seeds ...Request) {
	log.Println("======SimpleEngine'Run Start=====")
	var requests []Request
	//把传入的任务遍历出来添加到任务队列
	for _, r := range seeds {
		requests = append(requests, r)
	}
	//只要任务不为空就一只爬
	for len(requests) > 0 {
		//拿第一个出来
		request := requests[0]
		//第二个开始索引移动到第一个
		requests = requests[1:]

		//根据任务请求中定义好的解析函数进行解析网页数据
		parseResult, err := worker(request, 1)
		if err != nil {
			continue
		}
		//1.0先把解析出的请求添加到任务队列中
		requests = append(requests, parseResult.Requests...)
		//2.0这里把解析出来的指定内容打印出来
		for _, item := range parseResult.Items {
			log.Printf("item:>>%v<<\n", item)
		}
	}
}
