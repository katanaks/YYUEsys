package controllers

/////////////////////////////////////////////////////////////////////////////
//	员工资料 控制器
/////////////////////////////////////////////////////////////////////////////

import (
	"YYUEsys/models"
	"YYUEsys/models/utils"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/validation"
	beego "github.com/beego/beego/v2/server/web"
)

//ClassJSON 列表json格式
type ClassJSON models.PubJSON

//RsJSON 简单json格式
type classSimpJSON models.SimpJSON

//ClassController 员工数据控制器
type ClassController struct {
	beego.Controller
}

//logsClass logs输出
func logsClass(c *ClassController, Method string, code int, event string, info string) {
	var log utils.LogStruct
	ip := c.Ctx.Request.RemoteAddr
	ip = ip[0:strings.LastIndex(ip, ":")]
	// logs.Info(c.Ctx.Input.IP())

	log.Module = "基础数据-班级"
	log.Remote = ip
	log.UID = c.GetSession("UID").(int)
	log.CID = c.GetSession("CID").(int)
	log.Method = Method
	log.Code = code
	log.Event = event
	log.Info = info
	go utils.Logger(log)
}

//PageIndex 显示 班级列表页面 方法
//============================================================
func (c *ClassController) PageIndex() {
	c.TplName = "basic/class_index.html"
}

//PageAdd 显示 新增班级页面 方法
//============================================================
func (c *ClassController) PageAdd() {

	//处理员工列表
	StaffList, _ := models.GetStaffList(c.GetSession("CID").(int), 1, 1000) //获取列表
	c.Data["StaffList"] = StaffList                                         //指定循环列表
	c.Data["Selected"] = 1

	c.TplName = "basic/class_edit.html"

}

//PageEdit 显示 编辑班级页面 方法
//============================================================
func (c *ClassController) PageEdit() {
	id, _ := c.GetInt("ID")

	item := models.GetClassItem(id) //查询编辑条目
	// logs.Info(item)

	c.Data["ID"] = item.ID
	c.Data["Name"] = item.Name
	c.Data["StartDates"] = item.StartDates.Format("2006-01-02")

	//处理员工列表
	StaffList, _ := models.GetStaffList(c.GetSession("CID").(int), 1, 1000) //获取列表
	c.Data["StaffList"] = StaffList                                         //指定循环列表
	c.Data["Selected"] = item.Staff.ID

	c.Data["Memo"] = item.Memo

	c.TplName = "basic/class_edit.html"

}

//GetList 班级列表
///////////////////////////////////////////////////////////////////////////////
func (c *ClassController) GetList() {
	page, _ := c.GetInt("page")
	limit, _ := c.GetInt("limit")
	rslist, count := models.GetClassList(c.GetSession("CID").(int), page, limit)
	// println(len(Rlist))
	var code int = 200
	var msg string = "success"
	if len(rslist) == 0 {
		code = 100
		msg = "无数据。"
	}

	//返回Json
	Datastr := ClassJSON{code, msg, count, rslist}
	c.Data["json"] = Datastr

	c.ServeJSON()

	//日志输出
	go logsClass(c, "list", code, "班级列表", "Total:"+strconv.Itoa(len(rslist)))
}

//Save 保存班级方法
///////////////////////////////////////////////////////////////////////////////
func (c *ClassController) Save() {
	//准备数据
	var item models.Class
	var err error
	id, _ := c.GetInt("id") //编辑id变量
	item.ID = id
	item.Name = strings.Trim(c.GetString("Name"), "")

	StaffID, _ := c.GetInt("Staff")        //读员工ID
	StaffItem := models.Staff{ID: StaffID} //赋值条目
	item.Staff = &StaffItem

	item.StartDates, _ = time.ParseInLocation("2006-01-02", c.GetString("StartDates"), time.Local)

	item.Memo = c.GetString("Memo")
	item.CID = c.GetSession("CID").(int)

	// logs.Debug("save:", item)

	//准备 Json
	var code int = 200
	var errinfo string

	//验证表单
	cvinfo := ValidClass(item)
	errinfo += cvinfo

	//重复检查
	cdinfo := models.CheckClassDuplicate(item)

	errinfo += cdinfo

	if errinfo != "" {
		code = 100
	} else {
		//保存
		//如果无id号新增，否则进行更新
		if id == 0 {
			err = models.AddClass(item) //新增
		} else {
			err = models.UpdateClass(item) //编辑
		}

		if err != nil {
			code = 100
			errinfo += "写数据失败，请刷新重试"
		}
	}

	//返回 Json
	rs := &classSimpJSON{code, errinfo}
	c.Data["json"] = rs
	c.ServeJSON()

	//日志输出
	go logsClass(c, "save", code, "班级保存", item.Name)
}

//Del 删除班级 方法
///////////////////////////////////////////////////////////////////////////////
func (c *ClassController) Del() {
	//准备数据
	id, _ := c.GetInt("ID")
	//执行删除
	err := models.DelClass(id)
	//准备 Json
	code := 200
	msg := "success"
	if err != nil {
		code = 100
		msg = err.Error()
	}

	// 返回 Json
	rs := &classSimpJSON{code, msg}
	c.Data["json"] = rs
	c.ServeJSON()

	//日志输出
	go logsClass(c, "del", code, "删除班级", strconv.Itoa(id))
}

//ValidClass 表单验证
///////////////////////////////////////////////////////////////////////////////
func ValidClass(item models.Class) string {
	var errinfo string
	v := validation.Validation{}
	//验证项目
	v.Required(item.Name, "班级名").Message("不能为空")
	v.MinSize(item.Name, 2, "班级名").Message("最少2个字符")
	v.MaxSize(item.Name, 16, "班级名").Message("最大16个字符")

	//准备异常信息
	for _, err := range v.Errors {
		errinfo = errinfo + err.Key + " : " + err.Message + "<br> "
	}

	//返回异常提示
	return errinfo
}
