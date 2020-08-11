package blog

import "github.com/gin-gonic/gin"

func SetupBlog(e *gin.Engine) {
	e.GET("/blog", showBlog)
}

func showBlog(c *gin.Context) {
	c.JSON(200, gin.H{
		"msg": "this is app blog router",
	})
}
