package main

import (
	"MyProject/toml_conf/config"
	"flag"
	"fmt"
)

func main() {
	flag.Parse()
	if err := config.InitDb(); err != nil {
		fmt.Println(err.Error())
	}
	mysqlConf := config.Conf.Mysql
	fmt.Println(mysqlConf.DbName)
	fmt.Println(mysqlConf.UserName, mysqlConf.Password, mysqlConf.IpHost)
	versionConf := config.Conf.Version
	fmt.Println(versionConf)
	userConf := config.Conf.User
	fmt.Println(userConf.Name, userConf.Age, userConf.IsVip)
	groupConf := config.Conf.Group
	for _, val := range groupConf {
		fmt.Println(val.ID)
		fmt.Println(val.Name)
	}

}
