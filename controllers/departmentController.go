package controllers

/////////////////////////////////////////////////////////////////////////////
//	部门资料 控制器
/////////////////////////////////////////////////////////////////////////////

import (
	"YYUEsys/models"
	"YYUEsys/models/utils"
	"strconv"
	"strings"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/adapter/validation"
)

//DepartmentJSON 列表json格式
type DepartmentJSON models.PubJSON

//departmentSimpJSON 简易JSON
type departmentSimpJSON models.SimpJSON

//DepartmentController 部门数据控制器
//===================================================
type DepartmentController struct {
	beego.Controller
}

//logsStaff 显示 新增部门页面 方法
func logsDepartment(c *DepartmentController, Method string, code int, event string, info string) {
	var log utils.LogStruct
	ip := c.Ctx.Request.RemoteAddr
	ip = ip[0:strings.LastIndex(ip, ":")]

	log.Module = "基础数据-部门"
	log.Remote = ip
	log.UID = c.GetSession("UID").(int)
	log.CID = c.GetSession("CID").(int)
	log.Method = Method
	log.Code = code
	log.Event = event
	log.Info = info
	go utils.Logger(log)
}

//PageDepartmentIndex 显示 部门列表页面 方法
//===================================================
func (c *DepartmentController) PageDepartmentIndex() {
	c.TplName = "basic/department_index.html"
}

//PageDepartmentAdd 显示 新增部门页面 方法
//===================================================
func (c *DepartmentController) PageDepartmentAdd() {
	c.TplName = "basic/department_add.html"
}

//PageDepartmentEdit 显示 编辑部门页面 方法
//===================================================
func (c *DepartmentController) PageDepartmentEdit() {
	id, _ := c.GetInt("ID")
	// logs.Info("编辑ID:", id)
	item := models.GetDepartmentItem(id)

	c.Data["ID"] = item.ID
	c.Data["Name"] = item.Name
	c.Data["Memo"] = item.Memo
	c.TplName = "basic/department_edit.html"
}

//GetDepartmentList 部门列表
//===================================================
func (c *DepartmentController) GetDepartmentList() {
	page, _ := c.GetInt("page")
	limit, _ := c.GetInt("limit")

	list, count := models.GetDepartmentList(c.GetSession("CID").(int), page, limit)
	var code int = 200
	var msg string = "success"
	if len(list) == 0 {
		code = 100
		msg = "无数据。"
	}

	//返回Json
	Datastr := DepartmentJSON{code, msg, count, list}
	c.Data["json"] = Datastr
	c.ServeJSON()

	//日志输出
	go logsDepartment(c, "list", code, "部门列表", "Total:"+strconv.FormatInt(count, 10))
}

//SaveDepartment 保存部门方法
//===================================================
func (c *DepartmentController) SaveDepartment() {
	//准备数据
	var item models.Department
	var err error
	id, _ := c.GetInt("id") //编辑id变量

	item.Name = strings.Trim(c.GetString("name"), "")
	item.Memo = c.GetString("memo")
	item.CID = c.GetSession("CID").(int)
	item.ID = id
	// logs.Debug(id, item.Name, item.Memo, item.CID)

	//准备 Json
	var code int = 200
	var errinfo string

	//验证表单
	cvinfo := Valid(item)
	errinfo += cvinfo

	//重复检查
	cdinfo := models.CheckDuplicate(item)
	errinfo += cdinfo

	if errinfo != "" {
		code = 100
	} else {
		//保存
		//如果无id号新增，否则进行更新
		if id == 0 {
			err = models.Add(item) //新增
		} else {
			err = models.Update(item) //编辑
		}

		if err != nil {
			code = 100
			errinfo += "写数据失败，请刷新重试"
		}
	}

	//返回 Json
	rs := &departmentSimpJSON{code, errinfo}
	c.Data["json"] = rs
	c.ServeJSON()

	//日志输出
	go logsDepartment(c, "list", code, "部门保存", item.Name)

}

//DelDepartment 删除部门 方法
//===================================================
func (c *DepartmentController) DelDepartment() {
	//准备数据
	id, _ := c.GetInt("ID")
	//执行删除
	err := models.Del(id)
	//准备 Json
	code := 200
	msg := "success"
	if err != nil {
		code = 100
		msg = "error"
	}

	// 返回 Json
	rs := &departmentSimpJSON{code, msg}
	c.Data["json"] = rs
	c.ServeJSON()
	//日志输出
	go logsDepartment(c, "del", code, "部门删除", "id:"+strconv.Itoa(id))

}

//Valid 表单验证
//===================================================
func Valid(item models.Department) string {
	v := validation.Validation{}
	//验证项目
	v.Required(item.Name, "部门名称").Message("不能为空")
	v.MinSize(item.Name, 2, "部门名称").Message("最少2个字符")
	v.MaxSize(item.Name, 32, "部门名称").Message("最大32个字符")
	var errinfo string

	//准备异常信息
	for _, err := range v.Errors {
		errinfo = errinfo + err.Key + " : " + err.Message + "<br> "
	}

	//返回异常提示
	return errinfo
}
