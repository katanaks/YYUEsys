package controllers

/////////////////////////////////////////////////////////////////////////////
//	员工资料 控制器
/////////////////////////////////////////////////////////////////////////////

import (
	"YYUEsys/models"
	"YYUEsys/models/utils"
	"strconv"
	"strings"

	"github.com/beego/beego/v2/core/validation"
	beego "github.com/beego/beego/v2/server/web"
)

//ServiceJSON 列表json格式
type ServiceJSON models.PubJSON

//RsJSON 简单json格式
type serviceSimpJSON models.SimpJSON

//ServiceController 员工数据控制器
type ServiceController struct {
	beego.Controller
}

//logsService logs输出
func logsService(c *ServiceController, Method string, code int, event string, info string) {
	var log utils.LogStruct
	ip := c.Ctx.Request.RemoteAddr
	ip = ip[0:strings.LastIndex(ip, ":")]

	log.Module = "基础数据-服务"
	log.Remote = ip
	log.UID = c.GetSession("UID").(int)
	log.CID = c.GetSession("CID").(int)
	log.Method = Method
	log.Code = code
	log.Event = event
	log.Info = info
	go utils.Logger(log)
}

//PageIndex 显示 服务列表页面 方法
//============================================================
func (c *ServiceController) PageIndex() {
	c.TplName = "basic/service_index.html"
}

//PageAdd 显示 新增服务页面 方法
//============================================================
func (c *ServiceController) PageAdd() {

	//处理员工列表
	StaffList, _ := models.GetStaffList(c.GetSession("CID").(int), 1, 1000) //获取列表
	c.Data["StaffList"] = StaffList                                         //指定循环列表
	c.Data["Selected"] = "有效"

	c.TplName = "basic/service_edit.html"
}

//PageEdit 显示 编辑服务页面 方法
//============================================================
func (c *ServiceController) PageEdit() {
	id, _ := c.GetInt("ID")

	item := models.GetServiceItem(id) //查询编辑条目
	// logs.Info(item)

	c.Data["ID"] = item.ID
	c.Data["Name"] = item.Name
	// c.Data["StartDates"] = item.StartDates.Format("2006-01-02")

	//处理员工列表
	// StaffList, _ := models.GetStaffList(c.GetSession("CID").(int), 1, 1000) //获取列表
	// c.Data["StaffList"] = StaffList                                         //指定循环列表
	c.Data["Selected"] = item.State

	c.Data["Memo"] = item.Memo

	c.TplName = "basic/service_edit.html"

}

//GetList 服务列表
///////////////////////////////////////////////////////////////////////////////
func (c *ServiceController) GetList() {
	page, _ := c.GetInt("page")
	limit, _ := c.GetInt("limit")
	rslist, count := models.GetServiceList(c.GetSession("CID").(int), page, limit)
	// println(len(Rlist))
	var code int = 200
	var msg string = "success"
	if len(rslist) == 0 {
		code = 100
		msg = "无数据。"
	}
	//返回Json
	Datastr := ServiceJSON{code, msg, count, rslist}
	c.Data["json"] = Datastr

	c.ServeJSON()

	//日志输出
	go logsService(c, "list", code, "服务列表", "Total:"+strconv.Itoa(len(rslist)))
}

//Save 保存服务方法
///////////////////////////////////////////////////////////////////////////////
func (c *ServiceController) Save() {
	//准备数据
	var item models.Service
	var err error
	id, _ := c.GetInt("id") //编辑id变量
	item.ID = id
	item.Name = strings.Trim(c.GetString("Name"), "")

	// StaffID, _ := c.GetInt("Staff")        //读员工ID
	// StaffItem := models.Staff{ID: StaffID} //赋值条目
	// item.Staff = &StaffItem
	// item.StartDates, _ = time.ParseInLocation("2006-01-02", c.GetString("StartDates"), time.Local)

	if c.GetString("State") == "on" {
		item.State = "有效"
	} else {
		item.State = "停用"
	}
	// logs.Debug(item.State)

	// item.Memo = c.GetString("Memo")
	item.CID = c.GetSession("CID").(int)

	// logs.Debug("save:", item)

	//准备 Json
	var code int = 200
	var errinfo string

	//验证表单
	cvinfo := ValidService(item)
	errinfo += cvinfo

	//重复检查
	cdinfo := models.CheckServiceDuplicate(item)

	errinfo += cdinfo

	if errinfo != "" {
		code = 100
	} else {
		//保存
		//如果无id号新增，否则进行更新
		if id == 0 {
			err = models.AddService(item) //新增
		} else {
			err = models.UpdateService(item) //编辑
		}

		if err != nil {
			code = 100
			errinfo += "写数据失败，请刷新重试"
		}
	}

	//返回 Json
	rs := &serviceSimpJSON{code, errinfo}
	c.Data["json"] = rs
	c.ServeJSON()

	//日志输出
	go logsService(c, "save", code, "服务保存", item.Name)
}

//Del 删除服务 方法
///////////////////////////////////////////////////////////////////////////////
func (c *ServiceController) Del() {
	//准备数据
	id, _ := c.GetInt("ID")
	//执行删除
	err := models.DelService(id)
	//准备 Json
	code := 200
	msg := "success"
	if err != nil {
		code = 100
		msg = err.Error()
	}

	// 返回 Json
	rs := &serviceSimpJSON{code, msg}
	c.Data["json"] = rs
	c.ServeJSON()

	//日志输出
	go logsService(c, "del", code, "删除服务", strconv.Itoa(id))
}

//ValidService 表单验证
///////////////////////////////////////////////////////////////////////////////
func ValidService(item models.Service) string {
	var errinfo string
	v := validation.Validation{}
	//验证项目
	v.Required(item.Name, "服务名").Message("不能为空")
	v.MinSize(item.Name, 2, "服务名").Message("最少2个字符")
	v.MaxSize(item.Name, 16, "服务名").Message("最大16个字符")

	//准备异常信息
	for _, err := range v.Errors {
		errinfo = errinfo + err.Key + " : " + err.Message + "<br> "
	}

	//返回异常提示
	return errinfo
}
