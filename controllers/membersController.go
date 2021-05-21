package controllers

/////////////////////////////////////////////////////////////////////////////
//	员工资料 控制器
/////////////////////////////////////////////////////////////////////////////

import (
	"YYUEsys/models"
	"YYUEsys/models/utils"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"strconv"
	"strings"
	"time"

	beego "github.com/beego/beego/v2/server/web"
)

//MembersJSON 列表json格式
type MembersJSON models.PubJSON

//memberSimpJSON 简单json格式
//type memberSimpJSON models.SimpJSON

//MembersController 员工数据控制器
type MembersController struct {
	beego.Controller
}

//logsService logs输出
func logsMembers(c *MembersController, Method string, code int, event string, info string) {
	var log utils.LogStruct
	ip := c.Ctx.Request.RemoteAddr
	ip = ip[0:strings.LastIndex(ip, ":")]

	log.Module = "会员管理"
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
func (c *MembersController) PageIndex() {
	rs, _ := models.GetServiceList(c.GetSession("CID").(int), 1, 9999) //服务列表

	c.Data["Servicelist"] = rs //赋值给网页生成列

	c.TplName = "members/members_index.html"
}

//PageAdd 显示 新增服务页面 方法
//============================================================
func (c *MembersController) PageAdd() {

	//处理员工列表
	StaffList, _ := models.GetStaffList(c.GetSession("CID").(int), 1, 1000) //获取列表
	c.Data["StaffList"] = StaffList                                         //指定循环列表

	//处理合同列表
	ContractList, _ := models.GetContractList(c.GetSession("CID").(int), 1, 1000) //获取列表
	c.Data["ContractList"] = ContractList

	//初始化预设合同相关值
	ContractSelected := int(ContractList[0]["ID"].(int64))
	fmt.Println(ContractSelected)
	//执行查询
	rs := models.Getpredata(ContractSelected)
	start := time.Now().Format("2006-01-02")
	end := time.Now().AddDate(0, int(rs[0]["Duration"].(float64)), -1).Format("2006-01-02")

	//预设当前日期
	c.Data["Durationstart"] = start
	c.Data["Durationend"] = end

	//预设当前价格
	c.Data["Price"] = rs[0]["Price"]
	fmt.Println(rs[0]["Price"])
	//预设其它数据
	c.Data["Selected"] = "有效"
	c.Data["Birthday"] = "2015-01-01"

	c.TplName = "members/members_edit.html"
}

// //PageEdit 显示 编辑服务页面 方法
// //============================================================
// func (c *MembersController) PageEdit() {
// 	id, _ := c.GetInt("ID")

// 	item := models.GetServiceItem(id) //查询编辑条目
// 	// logs.Info(item)

// 	c.Data["ID"] = item.ID
// 	c.Data["Name"] = item.Name
// 	// c.Data["StartDates"] = item.StartDates.Format("2006-01-02")

// 	//处理员工列表
// 	// StaffList, _ := models.GetStaffList(c.GetSession("CID").(int), 1, 1000) //获取列表
// 	// c.Data["StaffList"] = StaffList                                         //指定循环列表
// 	c.Data["Selected"] = item.State

// 	c.Data["Memo"] = item.Memo

// 	c.TplName = "basic/service_edit.html"

// }

//GetList 会员列表
///////////////////////////////////////////////////////////////////////////////
func (c *MembersController) GetList() {
	page, _ := c.GetInt("page")
	limit, _ := c.GetInt("limit")
	rslist, count := models.GetMembersList(c.GetSession("CID").(int), page, limit)
	// println(len(Rlist))

	//准备列表数据
	//查询相应字段值
	for k, v := range rslist {
		MemberID := v["ID"].(int64)
		ContractName := models.GetMemberContract(MemberID)         //获取相关合同名称
		MemberDurationend := models.GetMemberDurationend(MemberID) //获取到期日期
		rslist[k]["Contactname"] = ContractName
		rslist[k]["Durationend"] = MemberDurationend.String()[0:10]
		fmt.Println("返回的合同名", ContractName)
		//for _, v1 := range items { //赋值服务数量
		//	keystr := v1["Service__Name"].(string)
		//	rslist[k][keystr] = v1["Quantity"]
		//}
	}

	var code = 200
	var msg = "success"
	if len(rslist) == 0 {
		code = 100
		msg = "无数据。"
	}
	//返回Json
	Datastr := ServiceJSON{code, msg, count, rslist}
	c.Data["json"] = Datastr

	_ = c.ServeJSON()

	//日志输出
	go logsMembers(c, "list", code, "会员列表", "Total:"+strconv.Itoa(len(rslist)))
}

//Save 保存服务方法
///////////////////////////////////////////////////////////////////////////////
func (c *MembersController) Save() {
	//准备数据
	var item models.Members
	var itemtransaction models.Membertransaction

	var err error

	//准备会员表
	item.CID = c.GetSession("CID").(int)

	id, _ := c.GetInt("id") //编辑id变量
	item.ID = id
	item.Name = strings.Trim(c.GetString("Name"), "")
	item.Gender = c.GetString("Gender")
	item.Birthday, _ = time.ParseInLocation("2006-01-02", c.GetString("Birthday"), time.Local)
	item.Idcard = c.GetString("Idcard")
	item.Contractrelationship = c.GetString("Contractrelationship")
	item.Contactidcard = c.GetString("Contactidcard")
	item.Contacttelephone = c.GetString("Contacttelephone")
	item.State = "有效"
	item.Memo = c.GetString("Memo")

	item.Createtime = time.Now()
	item.Updatetime = time.Now()

	fmt.Println("新会员信息：", item)

	//准备交易表

	itemtransaction.CID = c.GetSession("CID").(int)
	itemtransaction.Contract, _ = c.GetInt("Contract") //对应合同
	itemtransaction.StaffID, _ = c.GetInt("Staff")
	itemtransaction.Durationstart, _ = time.ParseInLocation("2006-01-02", c.GetString("Durationstart"), time.Local)
	itemtransaction.Durationend, _ = time.ParseInLocation("2006-01-02", c.GetString("Durationend"), time.Local)

	itemtransaction.Paid, _ = c.GetFloat("Paid")
	itemtransaction.Createtime = time.Now()
	itemtransaction.Updatetime = itemtransaction.Createtime
	fmt.Println("交易内容", itemtransaction)
	println(itemtransaction.Contract)
	logs.Debug("save:", item)

	//准备 Json
	var code int = 200
	var errinfo string

	if errinfo != "" {
		code = 100
	} else {
		//保存
		//如果无id号新增，否则进行更新
		if id == 0 {
			err = models.AddMember(item, itemtransaction) //新增

		} else {
			//err = models.UpdateService(item) //编辑
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
	go logsMembers(c, "save", code, "服务保存", item.Name)
}

// //Del 删除服务 方法
// ///////////////////////////////////////////////////////////////////////////////
// func (c *MembersController) Del() {
// 	//准备数据
// 	id, _ := c.GetInt("ID")
// 	//执行删除
// 	err := models.DelService(id)
// 	//准备 Json
// 	code := 200
// 	msg := "success"
// 	if err != nil {
// 		code = 100
// 		msg = err.Error()
// 	}

// 	// 返回 Json
// 	rs := &serviceSimpJSON{code, msg}
// 	c.Data["json"] = rs
// 	c.ServeJSON()

// 	//日志输出
// 	go logsMembers(c, "del", code, "删除服务", strconv.Itoa(id))
// }

// //ValidMembers 表单验证
// ///////////////////////////////////////////////////////////////////////////////
// func ValidMembers(item models.Service) string {
// 	var errinfo string
// 	v := validation.Validation{}
// 	//验证项目
// 	v.Required(item.Name, "服务名").Message("不能为空")
// 	v.MinSize(item.Name, 2, "服务名").Message("最少2个字符")
// 	v.MaxSize(item.Name, 16, "服务名").Message("最大16个字符")

// 	//准备异常信息
// 	for _, err := range v.Errors {
// 		errinfo = errinfo + err.Key + " : " + err.Message + "<br> "
// 	}

// 	//返回异常提示
// 	return errinfo
// }

//Getpredata 输入时刷新表单
func (c *MembersController) Getpredata() {
	//准备数据
	id, _ := c.GetInt("ID")
	//执行查询
	rs := models.Getpredata(id)
	prestart := time.Now().Format("2006-01-02")
	preend := time.Now().AddDate(0, int(rs[0]["Duration"].(float64)), -1).Format("2006-01-02")

	rs[0]["prestart"] = prestart
	rs[0]["preend"] = preend

	c.Data["json"] = rs[0]
	_ = c.ServeJSON()
}
