package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"hello/models"
	_ "hello/routers"
	"hello/wxpaltutil"
)

func init() {
	models.RegisterDB()
}

func main() {
	o := orm.NewOrm()
	orm.RunSyncdb("default", false, true)
	wxpaltutil.WxTimeTask(o)
	//models.Add()
	beego.Run()
}
