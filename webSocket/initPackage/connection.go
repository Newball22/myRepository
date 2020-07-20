package initPackage

import (
	"errors"
	"github.com/gorilla/websocket"
	"sync"
)

type Connection struct {
	wsConn    *websocket.Conn
	intChan   chan []byte
	outChan   chan []byte
	closeChan chan byte
	isClosed  bool
	mutex     sync.Mutex
}

/*
读已经关闭的 chan 能一直读到东西，但是读到的内容根据通道内关闭前是否有元素而不同。
如果 chan 关闭前，buffer 内有元素还未读 , 会正确读到 chan 内的值，且返回的第二个 bool 值（是否读成功）为 true。
如果 chan 关闭前，buffer 内有元素已经被读完，chan 内无值，接下来所有接收的值都会非阻塞直接成功，返回 channel 元素的零值，
但是第二个 bool 值一直为 false。
写已经关闭的 chan 会 panic
*/

func (conn *Connection) CloseFunc() {
	conn.wsConn.Close()
	conn.mutex.Lock()
	//保证通道只能关闭一次
	if conn.isClosed {
		close(conn.closeChan) //直接把这个channel关闭，关闭的通道是可以被取数据的，如果不想被读取到东西，可以设置closeChan=nil
		conn.isClosed = true
	}
	conn.mutex.Unlock()
}

func InitConnection(wsConn *websocket.Conn) (conn *Connection, err error) {
	conn = &Connection{
		wsConn:    wsConn,
		intChan:   make(chan []byte, 1000),
		outChan:   make(chan []byte, 1000),
		closeChan: make(chan byte, 1), //使用一个channel之前，一定要对它初始化，就算不放值，要不然就是nil
	}
	//启动协程,读取传入的消息
	go conn.readMessage()
	//启动协程，返回给客户端
	go conn.writeMessage()

	return

}

func (conn *Connection) ReadMgr() (data []byte, err error) {
	select {
	case data = <-conn.intChan:
	case <-conn.closeChan:
		//如果通道没关闭则会一直阻塞，读不到数据；为什么通道关闭了就会进入这里？
		//是因为关闭的通道是可以不断读取到东西非阻塞的,但是值为channel类型的默认值
		err = errors.New("connection is closed")
	}

	return

}

func (conn *Connection) WriteMgr(data []byte) (err error) {
	select {
	case conn.outChan <- data:
	case <-conn.closeChan:
		err = errors.New("connection is closed")
	}
	return

}

func (conn *Connection) readMessage() {
	var (
		data []byte
		err  error
	)
	for {
		if _, data, err = conn.wsConn.ReadMessage(); err != nil {
			conn.CloseFunc()
		}
		//intChan操作之前先判断连接是否已经关闭,因为conn断开不会影响channel的阻塞
		select {
		case conn.intChan <- data:
		case <-conn.closeChan:
			conn.CloseFunc()
		}

	}
}

func (conn *Connection) writeMessage() {
	var (
		data []byte
		err  error
	)
	for {
		select {
		case data = <-conn.outChan:
		case <-conn.closeChan:
			conn.CloseFunc()
		}
		if err = conn.wsConn.WriteMessage(websocket.TextMessage, data); err != nil {
			conn.CloseFunc()
		}
	}

}
