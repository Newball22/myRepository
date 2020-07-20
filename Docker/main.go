package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", sayHello)
	if err := http.ListenAndServe(":8888", nil); err != nil {
		log.Println(fmt.Println("Open ListenAndServe failed err:", err.Error()))
		return
	}

}

func sayHello(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Server is Start...")
	w.Write([]byte("I am Newball"))
}
