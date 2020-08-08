package main

import (
	"fmt"
	"github.com/go-redis/redis"
)

//初始化redis连接
func createClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	_, err := client.Ping().Result()
	if err != nil {
		fmt.Println("Ping failed err:", err)
	}
	return client
}

//判断当前用户是否点赞
func is_like_id1(client *redis.Client, userID int) {
	val, err := client.SIsMember("like_id1", userID).Result()
	if err != nil {
		fmt.Println("SIsMember failed err:", err)
	}
	if val == false {
		fmt.Println("user don't like it")
	}

}

//点赞&取消点赞
func like_id1_set(client *redis.Client, userID int) {
	val, err := client.SIsMember("like_id1", userID).Result()
	if err != nil {
		fmt.Println("SIsMember failed err:", err)
	}
	if val == false {
		_, err := client.SAdd("like_id1", userID).Result()
		if err != nil {
			fmt.Println("SAdd failed err:", err)
		}
	} else {
		_, err := client.SRem("like_id1", userID).Result()
		if err != nil {
			fmt.Println("SRem failed err:", err)
		}
	}

}

//统计获赞次数
func like_id1_count(client *redis.Client) int64 {
	val, err := client.SCard("like_id1").Result()
	if err != nil {
		fmt.Println("SCard failed err:", err)
	}
	return val
}

func main() {
	client := createClient()
	fmt.Println(client)
	like_id1_set(client, 5)
	like_id1_set(client, 6)
	data := like_id1_count(client)
	fmt.Println(data)

}
