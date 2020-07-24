package config

import (
	"flag"
	"github.com/BurntSushi/toml"
)

type Mysql struct {
	UserName,
	Password,
	IpHost,
	DbName string
}

type User struct {
	Name  string
	Age   int
	IsVip bool
}

type Group struct {
	ID   int
	Name string
}

type Config struct {
	Mysql   *Mysql
	Version string
	User    *User
	Group   []*Group
}

var (
	confPath string
	Conf     = &Config{}
)

func init() {
	flag.StringVar(&confPath, "conf", "config/app.toml", "-conf path")
}

func InitDb() error {
	//将解析好的confPath存到&Conf里面
	_, err := toml.DecodeFile(confPath, &Conf)
	return err
}
