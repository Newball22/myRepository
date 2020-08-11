package main

import (
	"viperDemo/app/blog"
	"viperDemo/app/shop"
	"viperDemo/router"
)

func main() {
	//把app的路由添加到路由组
	router.Include(shop.SetupShop, blog.SetupBlog)

	//初始化
	r := router.Init()
	r.Run(":8080")
}
