package main

import (
	"bufio"
	"fmt"
	"github.com/nsqio/go-nsq"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

//定义一个全局生产者变量
var producer *nsq.Producer

//初始化生产者，这个是重点
func initProducer(str string) (err error) {
	config := nsq.NewConfig()
	producer, err = nsq.NewProducer(str, config)
	if err != nil {
		fmt.Println("create Producer failed,err:", err)
		return err
	}
	return nil
}

//生产者
func workProducer(i int) {
	fmt.Printf("==%d workProducer start==\n", i)
	nsqAddress := "127.0.0.1:4150"
	if err := initProducer(nsqAddress); err != nil {
		fmt.Println("connect Producer is failed,err:", err)
		return
	}

	reader := bufio.NewReader(os.Stdin) //读取标准输入的数据
	for {
		data, err := reader.ReadString('\n') //以换行键为末尾
		if err != nil {
			fmt.Println("read string from stdin failed,err:", err)
			continue
		}
		data = strings.TrimSpace(data)    //去空格
		if strings.ToUpper(data) == "Q" { //输入Q退出
			break
		}
		//向'topic_demo'publish数据
		err = producer.Publish("topic_demo", []byte(data)) //这个是重点
		if err != nil {
			fmt.Println("publish msg to nsq failed,err:", err)
			continue
		}

	}
}

//消费者类型
type MyHandler struct {
	Title string
}

//是需要实现的处理消息的方法
func (m *MyHandler) HandleMessage(msg *nsq.Message) (err error) {
	fmt.Printf("%s recived from %v,msg:%v\n", m.Title, msg.NSQDAddress, string(msg.Body))
	return
}

//初始化消费者
func initConsumer(topic string, channel string, address string) (err error) {
	config := nsq.NewConfig()
	config.LookupdPollInterval = 15 * time.Second
	c, err := nsq.NewConsumer(topic, channel, config)
	if err != nil {
		fmt.Println("create Consumer failed,err:", err)
		return
	}
	consumer := &MyHandler{
		Title: "this Queued task",
	}
	c.AddHandler(consumer)
	err = c.ConnectToNSQLookupd(address)
	if err != nil {
		return err
	}
	return nil
}

//消费者
func workConsumer(i int) {
	fmt.Printf("==%d workConsumer start==\n", i)
	err := initConsumer("topic_demo", "first", "127.0.0.1:4161")
	if err != nil {
		fmt.Println("init consumer failed,err:", err)
		return
	}
	c := make(chan os.Signal)        //定义一个信号的通道
	signal.Notify(c, syscall.SIGINT) //转发键盘中断信号到c
	<-c                              //阻塞
}

//调配
func main() {
	for i := 0; i < 3; i++ {
		go workProducer(i)
		go workConsumer(i)
	}
	time.Sleep(time.Second * 30)
	fmt.Println("over~")

}
