package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

//定义一个全局的channel，用于接收从各个客户端读到的登录退出在线等消息和聊天信息，然后广播出去
var message = make(chan string)

//定义一个结构体，用于存储每位聊天室用户的信息名字+各自的通道
type UserInfo struct {
	Name      string
	Address   string
	userChan  chan string //用于用户进入或者退出当前聊天室的提醒信息
	closeChan chan int
}

//定义一个全局的map，用于存储聊天室中所有在线的用户信息
//key值为IP+端口
var onlineUsers = make(map[string]UserInfo)

//创建一个全局的读写锁，为了保护onlineUsers公共区数据
var rwMutex sync.RWMutex

func main() {
	//启动服务器
	listener, err := net.Listen("tcp", "127.0.0.1:8011")
	if err != nil {
		fmt.Println("net.Listen failed, err:", err.Error())
		return
	}

	fmt.Println("聊天室已经启动")
	//聊天室启动成功就创建全局的管理信息方法，阻塞等待着
	go Manager()

	for {
		//循环监听所有的连接
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept failed, err:", err.Error())
			continue //这里不能用return，因为一个客户端Accept错误，不能把整个程序给return掉，毕竟还有其他成功的客户端
		}
		fmt.Printf("地址为[%v]的客户端已经连接成功\n", conn.RemoteAddr())
		/*
			如果监听到连接请求并成功以后，服务器进去go线程，该线程处理服务器和客户端之间的读写或其他事件，
			同时，服务器在go主线程中回去继续监听着其他客户端的连接请求
		*/
		go HandleConnect(conn) //对每个用户都开一个goroutine去并发处理
	}

}

//全局监听message里面的数据，跟各自的每个客户端状态没有关系，所以不放在HandleConnect里面
func Manager() {
	for {
		msg := <-message //无数据的话就阻塞,不会往下走
		//加一个读锁
		rwMutex.RLock()
		//如果有数据，则把在线的用户遍历出来，然后往每个用户的channel里面
		for _, user := range onlineUsers {
			user.userChan <- msg
		}
		//解锁
		rwMutex.RUnlock()
	}
}

//处理信息返回到各自的客户端
func WriteMsgToClient(user UserInfo, conn net.Conn) {
	for {
		msg := <-user.userChan
		//将得到的数据写给自己对应的客户端
		conn.Write([]byte(msg + "\n"))
	}

}

//解析客户端发来的信息
func ParseUserMsg(user UserInfo, conn net.Conn, buf []byte) string {
	n, err := conn.Read(buf) //这个Read方法是阻塞的，直到读到东西或者出错为止
	if err != nil && err != io.EOF {
		fmt.Println("用户断开聊天室连接 err:", err.Error())
		return err.Error()
	}
	if n == 0 {
		message <- MakeMsg(user, "对方已经下线")
		return ""
	}

	return string(buf[:n-1]) //'n-1是为了去掉换行符'

}

func MakeMsg(user UserInfo, msg string) string {
	if user.Name == "" {
		return user.Address + ":" + msg
	}
	return user.Name + ":" + msg
}

//处理连接成功后的客户端业务
func HandleConnect(conn net.Conn) {
	var (
		isClosed bool
	)
	defer func() {
		conn.Write([]byte("服务器进入维护状态，大家洗洗睡吧..."))
		conn.Close()
	}()
	//用于用户超时处理
	overTime := make(chan bool)

	//获取客户端地址结构
	clientAddr := conn.RemoteAddr().String()

	//初始化客户端
	user := UserInfo{
		Name:      "",
		Address:   clientAddr,
		userChan:  make(chan string),
		closeChan: make(chan int, 1),
	}
	//启动另一个单独的go程读取自己channel里面的数据，阻塞等待着
	go WriteMsgToClient(user, conn)

	//开启一个独立的线程:设置一个定时用来防止用户既不输入昵称，又不断开连接，导致这个连接一直被Read方法阻塞着
	go func(user *UserInfo) {
		for {
			select {
			case <-user.closeChan:
				fmt.Println("用户设置昵称成功")
				return
			case <-time.NewTimer(time.Second * 30).C:
				conn.Close()
				isClosed = true
				return
			}
		}

	}(&user)

LOOP:
	conn.Write([]byte("请您输入一个昵称:"))
	bufN := make([]byte, 1024)
	MsgN := ParseUserMsg(user, conn, bufN)

	if MsgN == "" {
		conn.Write([]byte("昵称不能为空\n"))
		goto LOOP
	}
	if isClosed == true {
		return
	}
	user.Name = MsgN
	user.userChan <- "您的昵称为:" + MsgN
	user.closeChan <- 1 //说明昵称设置成功，让上面select执行第一个条件，回收goroutine

	//组织一个自己上线的广播信息【IP+PORT】Name
	msg := user.Name + "上线了"
	//将上线消息写入全局消息channel
	message <- msg

	//添加新用户之前先加写锁
	rwMutex.Lock()
	//添加到用户列表map中
	onlineUsers[clientAddr] = user
	//解锁
	rwMutex.Unlock()

	//创建一个新的go线程处理客户端发过来的聊天数据
	go func() {
		buf := make([]byte, 4096)
		//循环读取用户发送的消息
		for {
			msg := ParseUserMsg(user, conn, buf)
			if msg == "" {
				break
			}
			//组织一下用户聊天的内容
			news := MakeMsg(user, msg)
			message <- news
			overTime <- true

		}

	}()

	//判断用户是否为僵尸粉
	for {
		select {
		case <-overTime:
		case <-time.After(time.Second * 300):
			user.userChan <- MakeMsg(user, "由于您长时间不在线，聊天室已经断开连接...")
			delete(onlineUsers, clientAddr)
			return
		}
	}
}
