package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var (
	data = Config{}
)

type Config struct {
	MysqlDB MysqlDB
	Redis   Redis
}

type MysqlDB struct {
	Dialect,
	Database,
	User,
	Password,
	Charset,
	URL,
	Host string
	Port,
	MaxIdleConns,
	MaxOpenConns int
}
type Redis struct {
	User,
	Password,
	Charset,
	Host string
	Port int
}

type JsonStruct struct {
}

//初始化项目配置文件
func InitConf() *Config {
	webConfigs()
	mysqlDB()
	return &data
}

//解析JSON格式配置文件
func webConfigs() {
	JsonParse := NewJsonStruct()
	JsonParse.Load("./conf/config.json", &data)
}

func NewJsonStruct() *JsonStruct {
	return &JsonStruct{}
}

//序列化配置文件
func (jst *JsonStruct) Load(fileName string, v interface{}) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("ReadFile failed, err:", err)
		return
	}
	err = json.Unmarshal(data, v)
	if err != nil {
		fmt.Println("Unmarshal err: ", err)
		return
	}
}

//初始化mysql数据库
func mysqlDB() {
	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		data.MysqlDB.User,
		data.MysqlDB.Password,
		data.MysqlDB.Host,
		data.MysqlDB.Port,
		data.MysqlDB.Database,
		data.MysqlDB.Charset)
	data.MysqlDB.URL = url
}
