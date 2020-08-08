package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.Handle("/", middleFunc(http.HandlerFunc(hello)))
	http.ListenAndServe(":8080", nil)

}

//中间件
func middleFunc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Start")
		next.ServeHTTP(w, r)
		fmt.Println("end")
	})
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("中间件测试"))
}
