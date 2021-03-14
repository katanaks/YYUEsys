package utils

/////////////////////////////////////////////////////////////////////////////
//	公共应用包
//  1. InitDB 数据库连接
//  2. GetDBLink 返回数据库连接字串{DBLinkStr类型}
/////////////////////////////////////////////////////////////////////////////

import (
	"time"

	beego "github.com/beego/beego/v2/adapter"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
)

//DBLinkStr 数据库连接字串结构
type DBLinkStr struct {
	Host     string
	Port     string
	DBname   string
	User     string
	Password string
	Linkstr  string
}

//InitDB 数据库连接
func InitDB() {
	orm.RegisterDriver("mysql", orm.DRMySQL) //注册mysql引擎
	orm.DefaultTimeLoc = time.UTC
	s := GetDBLink()
	orm.RegisterDataBase("default", "mysql", s.Linkstr) //注册数据库
	//log
	logs.Informational(">>> Database configuration complete")

}

//GetDBLink 返回数据库连接字串
func GetDBLink() DBLinkStr {
	//生成数据库连接字串内容
	var rs DBLinkStr

	rs.Host = beego.AppConfig.String("DB::Host")
	// Configfile.String("DB::Host")
	rs.Port = beego.AppConfig.String("DB::Port")
	rs.DBname = beego.AppConfig.String("DB::DBname")
	rs.User = beego.AppConfig.String("DB::User")
	rs.Password = beego.AppConfig.String("DB::Password")

	//生成Mysql连接字符串
	linkstr := rs.User + ":" + rs.Password + "@tcp(" + rs.Host + ":" + rs.Port + ")/" + rs.DBname + "?charset=utf8&loc=Local"
	rs.Linkstr = linkstr
	// 日志
	logs.Informational(">>>", "Database Host:", rs.Host+":"+rs.Port)

	return rs
}
