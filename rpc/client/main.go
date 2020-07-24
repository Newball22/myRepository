package main

import (
	"fmt"
	"log"
	"net/rpc"
)

type ArithRequest struct {
	A, B int
}

type ArithResponse struct {
	Pro int
	Quo int
	Rem int
}

func main() {
	conn, err := rpc.DialHTTP("tcp", ":8080")
	if err != nil {
		log.Fatalln(err)
	}
	req := ArithRequest{
		A: 9,
		B: 2,
	}

	var res ArithResponse
	err2 := conn.Call("Arith.Multiply", req, &res)
	if err2 != nil {
		log.Fatalln(err2)
	}
	fmt.Printf("%d * %d=%d\n", req.A, req.B, res.Pro)
	err3 := conn.Call("Arith.Divide", req, &res)
	if err3 != nil {
		log.Fatal(err3)
	}
	fmt.Printf("%d / %d 商= %d，余数 = %d\n", req.A, req.B, res.Quo, res.Rem)
}
