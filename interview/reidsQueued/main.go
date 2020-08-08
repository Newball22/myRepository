package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const RMQ string = "mqTest"

//生产者
func producer() {
	redis_conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer redis_conn.Close()
	rand.Seed(time.Now().UnixNano())
	var i = 1
	for {
		_, err = redis_conn.Do("rPush", RMQ, strconv.Itoa(i))
		if err != nil {
			fmt.Println("producer err:", err.Error())
			continue
		}
		fmt.Println("producer element %d", i)
		time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
		i++
	}
}

//消费者
func consumer() {
	redis_conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer redis_conn.Close()

	rand.Seed(time.Now().UnixNano())

	for {
		ele, err := redis.String(redis_conn.Do("lPop", RMQ))
		if err != nil {
			fmt.Println("no msg.sleep now")
			time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
		} else {
			fmt.Println("cosume element:%s", ele)
		}
	}
}

func main() {
	list := os.Args
	fmt.Println(list)
	if list[0] == "pro" {
		go producer()
	} else if list[0] == "con" {
		go consumer()
	}
	for {
		time.Sleep(time.Duration(10000) * time.Second)
	}
}
