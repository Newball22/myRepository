package main

import (
	"MyProject/webSocket/initPackage"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

var (
	//允许跨域，所以设置为true
	upGrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func webChat(w http.ResponseWriter, r *http.Request) {
	var (
		wsConn *websocket.Conn
		conn   *initPackage.Connection
		err    error
		data   []byte
	)
	if wsConn, err = upGrader.Upgrade(w, r, nil); err != nil {
		fmt.Println("upgrader.Upgrade, err:", err.Error())
		return
	}
	if conn, err = initPackage.InitConnection(wsConn); err != nil {
		conn.CloseFunc()
	}

	for {
		if data, err = conn.ReadMgr(); err != nil {
			fmt.Println("conn ReadMgr failed err:", err.Error())
			conn.CloseFunc()
		}

		if err = conn.WriteMgr(data); err != nil {
			fmt.Println("conn WriteMgr failed err:", err.Error())
			conn.CloseFunc()
		}
	}

}

func main() {
	http.HandleFunc("/", webChat)
	err := http.ListenAndServe(":8088", nil)
	if err != nil {
		fmt.Println("ListenAndServe failed, err:", err.Error())
		return
	}

}
