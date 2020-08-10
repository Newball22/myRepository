package dao

import (
	"App_Conf/config"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	DB *gorm.DB
	//con = &config.Config{}
)

func DBInit() (err error) {
	fmt.Println("数据库持久化服务器启动...")
	con := config.InitConf()
	conf := con.MysqlDB
	//fmt.Println(conf)
	DB, err = gorm.Open(conf.Dialect, conf.URL)
	if err != nil {
		panic("数据库连接失败:" + err.Error())
	}
	DB.DB().SetMaxIdleConns(con.MysqlDB.MaxIdleConns)
	DB.DB().SetMaxOpenConns(con.MysqlDB.MaxOpenConns)
	DB.LogMode(true)
	err = DB.DB().Ping()
	return

}

//
//func InitModel() {
//	//DB.AutoMigrate(&models.User{}/&models.其他/..),可以一下子绑定N个model
//	DB.AutoMigrate(&models.User{})
//}

func Close() {
	DB.Close()
}
