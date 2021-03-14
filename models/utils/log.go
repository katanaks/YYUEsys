package utils

import (
	"strconv"

	"github.com/beego/beego/v2/core/logs"
)

/////////////////////////////////////////////////////////////////////////////
//	日志功能接口
/////////////////////////////////////////////////////////////////////////////

//LogStruct 数据库连接字串结构
type LogStruct struct {
	Module string //模块
	Method string //方法 View Add Delete Update
	Remote string //客户IP
	Code   int    //返回码
	Event  string //gk wrh
	UID    int    //用户ID
	CID    int    //公司ID
	Info   interface{}
}

//LogInit 初始化Log配置
func LogInit() {

	logs.SetLogger(logs.AdapterMultiFile, `{
		"filename":"logs/YYUEsys.log",
		"color":true,
		"maxlines":1000,
		"separate":["emergency", "alert", "critical", "error", "warning", "notice", "info", "debug"]
		}`)
	logs.EnableFuncCallDepth(false) //显示文件、行号
	logs.SetLevel(logs.LevelDebug)
	logs.Informational(">>> Logger configuration complete")
}

//Logger 日志记录
func Logger(log LogStruct) {
	logs.Info(
		"-->", log.Remote,
		"|\033[32m", log.Code, "\033[0m",
		"|\033[36m", "UID", Strlen(strconv.Itoa(log.UID), 4), "\033[0m|", "\033[36mCID", Strlen(strconv.Itoa(log.CID), 3), "\033[0m",
		"|\033[30;32;1m", Strlen(log.Module, 15), "\033[0m",
		"|\033[32m", Strlen(log.Method, 5), "\033[0m",

		"|", Strlen(log.Event, 15),
		"|", log.Info)
}

//Strlen 固定字符长度
func Strlen(str string, setlen int) string {
	oldlen := len(str)
	z := setlen - oldlen
	for i := 0; i < z; i++ {
		str = str + " "
	}
	return str
}
