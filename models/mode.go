package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

type WxAccessToken struct {
	Id          int `orm:"auto"`
	AccessToken string
}

func RegisterDB() {
	orm.RegisterModel(new(WxAccessToken))
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:695981642@/wechat?charset=utf8")
}
