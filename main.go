package main

import (
	"YYUEsys/models/utils"
	_ "YYUEsys/routers"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/go-sql-driver/mysql"
)

func init() {

	// 配置数据库
	utils.InitDB()

	orm.RunSyncdb("default", false, true) //表自动生成

	//log设置
	utils.LogInit()
}

//系统主入口cnbman
func main() {
	beego.Run()
}
