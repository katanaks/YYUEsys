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

//StaffJSON 列表json格式
type StaffJSON models.PubJSON

//staffSimpJSON 简易JSON
type staffSimpJSON models.SimpJSON

//StaffController 员工数据控制器
type StaffController struct {
	beego.Controller
}

//logsStaff 显示 新增部门页面 方法
func logsStaff(c *StaffController, Method string, code int, event string, info string) {
	var log utils.LogStruct
	ip := c.Ctx.Request.RemoteAddr
	ip = ip[0:strings.LastIndex(ip, ":")]

	log.Module = "基础数据-员工"
	log.Remote = ip
	log.UID = c.GetSession("UID").(int)
	log.CID = c.GetSession("CID").(int)
	log.Method = Method
	log.Code = code
	log.Event = event
	log.Info = info
	go utils.Logger(log)
}

//PageIndex 显示 部门列表页面 方法
//============================================================
func (c *StaffController) PageIndex() {
	c.TplName = "basic/staff_index.html"
}

//PageAdd 显示 新增部门页面 方法
//============================================================
func (c *StaffController) PageAdd() {

	//处理部门列表
	DepartmentList, _ := models.GetDepartmentList(c.GetSession("CID").(int), 1, 5000) //获取列表
	c.Data["DepartmentList"] = DepartmentList                                         //指定循环列表
	c.Data["Selected"] = 1

	c.TplName = "basic/staff_edit.html"

}

//PageEdit 显示 编辑部门页面 方法
//============================================================
func (c *StaffController) PageEdit() {
	id, _ := c.GetInt("ID")

	item := models.GetStaffItem(id)
	// logs.Info(item)

	c.Data["ID"] = item.ID
	c.Data["Name"] = item.Name
	c.Data["State"] = item.State
	c.Data["Telephone"] = item.Telephone
	c.Data["Gender"] = item.Gender
	c.Data["Entrydate"] = item.Entrydate.Format("2006-01-02")

	//处理部门列表
	DepartmentList, _ := models.GetDepartmentList(c.GetSession("CID").(int), 1, 5000) //获取列表
	c.Data["DepartmentList"] = DepartmentList                                         //指定循环列表
	c.Data["Selected"] = item.Department.ID                                           //已选择复制

	c.Data["Birthday"] = item.Birthday.Format("2006-01-02")
	c.Data["Special"] = item.Special
	c.Data["School"] = item.School
	c.Data["Education"] = item.Education
	c.Data["Address"] = item.Address
	c.Data["Memo"] = item.Memo
	// logs.Info(item)
	c.TplName = "basic/staff_edit.html"

}

//GetList 部门列表
///////////////////////////////////////////////////////////////////////////////
func (c *StaffController) GetList() {
	page, _ := c.GetInt("page")
	limit, _ := c.GetInt("limit")
	rslist, count := models.GetStaffList(c.GetSession("CID").(int), page, limit)
	// println(len(Rlist))
	var code int = 200
	var msg string = "success"
	if len(rslist) == 0 {
		code = 100
		msg = "无数据。"
	}

	//返回Json
	Datastr := StaffJSON{code, msg, count, rslist}
	c.Data["json"] = Datastr
	c.ServeJSON()

	//日志输出
	go logsStaff(c, "list", code, "员工列表", "Total:"+strconv.FormatInt(count, 10))
}

//Save 保存部门方法
///////////////////////////////////////////////////////////////////////////////
func (c *StaffController) Save() {
	//准备数据
	var item models.Staff
	var err error
	id, _ := c.GetInt("id") //编辑id变量
	item.ID = id
	// logs.Info(c.GetInt("Department"))
	item.Name = strings.Trim(c.GetString("Name"), "")
	item.State = strings.Trim(c.GetString("State"), "")
	item.Telephone = strings.Trim(c.GetString("Telephone"), "")
	item.Gender = strings.Trim(c.GetString("Gender"), "")
	item.Entrydate, _ = time.ParseInLocation("2006-01-02", c.GetString("Entrydate"), time.Local)

	DepartmentID, _ := c.GetInt("Department")             //读部门ID
	DepartmentItem := models.Department{ID: DepartmentID} //赋值条目
	item.Department = &DepartmentItem

	item.Birthday, _ = time.ParseInLocation("2006-01-02", c.GetString("Birthday"), time.Local)
	item.Special = strings.Trim(c.GetString("Special"), "")
	item.School = strings.Trim(c.GetString("School"), "")
	item.Education = strings.Trim(c.GetString("Education"), "")
	item.Address = strings.Trim(c.GetString("Address"), "")
	item.Memo = c.GetString("Memo")
	item.CID = c.GetSession("CID").(int)

	// logs.Debug(item)

	//准备 Json
	var code int = 200
	var errinfo string

	//验证表单
	cvinfo := ValidStaff(item)
	errinfo += cvinfo

	//重复检查
	cdinfo := models.CheckStaffDuplicate(item)

	errinfo += cdinfo

	if errinfo != "" {
		code = 100
	} else {
		//保存
		//如果无id号新增，否则进行更新
		if id == 0 {
			err = models.AddStaff(item) //新增
		} else {
			err = models.UpdateStaff(item) //编辑
		}

		if err != nil {
			code = 100
			errinfo += "写数据失败，请刷新重试"
		}
	}

	//返回 Json
	rs := &staffSimpJSON{code, errinfo}
	c.Data["json"] = rs
	c.ServeJSON()

	//日志输出
	go logsStaff(c, "save", code, "部门保存", item.Name)
}

//Del 删除部门 方法
///////////////////////////////////////////////////////////////////////////////
func (c *StaffController) Del() {
	//准备数据
	id, _ := c.GetInt("ID")
	//执行删除
	err := models.DelStaff(id)
	//准备 Json
	code := 200
	msg := "success"
	if err != nil {
		code = 100
		msg = "error"
	}

	// 返回 Json
	rs := &staffSimpJSON{code, msg}
	c.Data["json"] = rs
	c.ServeJSON()
	//日志输出
	go logsStaff(c, "del", code, "删除部门", strconv.Itoa(id))
}

//ValidStaff 表单验证
///////////////////////////////////////////////////////////////////////////////
func ValidStaff(item models.Staff) string {
	var errinfo string
	v := validation.Validation{}
	//验证项目
	v.Required(item.Name, "姓名").Message("不能为空")
	v.MinSize(item.Name, 2, "姓名").Message("最少2个字符")
	v.MaxSize(item.Name, 4, "姓名").Message("最大32个字符")
	v.Required(item.Name, "手机").Message("不能为空")
	v.MinSize(item.Telephone, 7, "手机").Message("最少7个字符")
	v.MaxSize(item.Name, 11, "手机").Message("最大11个字符")
	v.Required(item.Entrydate, "入职日期").Message("不能为空")
	v.Required(item.Birthday, "出生年月").Message("不能为空")
	v.Required(item.Address, "当前住址").Message("不能为空")

	//准备异常信息
	for _, err := range v.Errors {
		errinfo = errinfo + err.Key + " : " + err.Message + "<br> "
	}

	//验证18岁成年   6570 18岁
	Entrydate := item.Birthday
	sub := time.Now().Sub(Entrydate).Hours() / 24
	// logs.Info(sub)
	if sub < 6570 {
		errinfo = errinfo + "员工年龄必须年满18周岁" + "<br> "
	}

	//返回异常提示
	return errinfo
}
