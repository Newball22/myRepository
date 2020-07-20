package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	fmt.Println("正在连接聊天室服务器...")
	conn, err := net.Dial("tcp", "127.0.0.1:8011")
	if err != nil {
		fmt.Println("net.Dial failed, err:", err.Error())
		return
	}
	defer conn.Close()

	fmt.Println("聊天室连接成功👏")
	fmt.Println("请您起个喜欢的昵称吧💗")
	var userName string
	fmt.Scan(&userName)
	conn.Write([]byte(userName))

	bufMsg1 := make([]byte, 4096)
	n1, err := conn.Read(bufMsg1)
	if err != nil {
		fmt.Println("conn.Read failed, err:", err.Error())
		return
	}

	//客户端收到服务器的反馈
	fmt.Println(string(bufMsg1[:n1]))
	fmt.Printf("温馨提示：尊敬的%s,长时间没有发送消息，会自动退出聊天室哦")

	//
	go func() {
		for {
			bufMsg2 := make([]byte, 4096)
			n2, err := os.Stdin.Read(bufMsg2)
			if err != nil {
				fmt.Println("os.Stdin.Read failed, err:", err.Error())
				continue
			}
			conn.Write(bufMsg2[:n2])
		}
	}()

	//接收服务器发送来的数据
	for {
		buf := make([]byte, 4096)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("conn.Read err:", err.Error())
			return
		}
		if n == 0 {
			fmt.Println("服务器已关闭当前连接，正在退出,期待您下次的光临...")
			return
		}
		fmt.Println(string(bufMsg1[:n1]))
	}

}
