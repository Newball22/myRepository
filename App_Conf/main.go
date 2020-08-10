package main

import (
	"App_Conf/dao"
	"App_Conf/models"
	"App_Conf/router"
)

func main() {
	err := dao.DBInit()
	if err != nil {
		panic("数据库连接出错：" + err.Error())
	}
	defer dao.Close()
	//模型绑定
	dao.DB.AutoMigrate(&models.User{})
	r := router.SetupRouter()
	r.Run(":8080")
}
