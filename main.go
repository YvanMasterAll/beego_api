package main

import (
	_ "github.com/prime_beego/beego_api/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/prime_beego/beego_api/filters"
	"time"
	"github.com/prime_beego/beego_api/controllers"
)

func init() {
	orm.RegisterDataBase("default", "mysql", "root:19920213@tcp(127.0.0.1:3306)/prime_beego")
}

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
		orm.Debug = true
	}
	beego.BConfig.WebConfig.Session.SessionOn = true //enabled session
	//insert filter
	beego.InsertFilter("/v1/*", beego.BeforeRouter, filters.BaseFilter)
	beego.InsertFilter("/v1/paylist/*", beego.BeforeRouter, filters.PaylistFilter)
	orm.DefaultTimeLoc = time.UTC //orm time sync
	beego.ErrorController(&controllers.ErrorController{}) //bad request
	beego.BConfig.ServerName = "snail server 1.0" //set server name
	beego.Run()
}

