package main

import (
	"errors"
	"log"
	"net/http"
	"net/rpc"
)

//用于注册
type Arith struct {
}

//声明参数的结构体
type ArithRequest struct {
	A, B int
}

//返回给客户端的数据
type ArithResponse struct {
	Pro int //乘积
	Quo int //除法
	Rem int //求余
}

//乘法的方法
func (this *Arith) Multiply(req ArithRequest, res *ArithResponse) error {
	res.Pro = req.A * req.B
	return nil

}

//求商，余数
func (this *Arith) Divide(req ArithRequest, res *ArithResponse) error {
	if req.B == 0 {
		return errors.New("除数不能为零")
	}
	res.Quo = req.A / req.B
	res.Rem = req.A % req.B
	return nil
}

func main() {
	//1.0注册服务
	rect := new(Arith)
	//1.1注册一个rect的服务
	_ = rpc.Register(rect)
	//2.0将服务处理绑定到http协议上
	rpc.HandleHTTP()
	//3.0监听服务
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalln(err)
	}

}
