package main

import (
	"MyProject/autostruct/dbtools"
	"MyProject/autostruct/generate"
)

func main() {
	//初始化数据库
	dbtools.Init()
	//generate.Generate() //生成所有表信息
	generate.Generate("users", "matters") //生成指定表信息，可变参数可传入多个表名
}
