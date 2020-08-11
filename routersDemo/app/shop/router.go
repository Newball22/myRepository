package shop

import "github.com/gin-gonic/gin"

func SetupShop(e *gin.Engine) {
	e.GET("/shop", showShop)
}

func showShop(c *gin.Context) {
	c.JSON(200, gin.H{
		"msg": "this is app shop router",
	})
}
