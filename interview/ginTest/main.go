package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World",
		})
	})

	r.GET("/test", DirectFunc)
	r.Handle("POST", "/post", postFunc)
	r.Run()
}

func DirectFunc(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "https://baidu.com")
}

func postFunc(c *gin.Context) {
	c.JSON(200, "hello post")
}
