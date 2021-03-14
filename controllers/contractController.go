package controllers

/////////////////////////////////////////////////////////////////////////////
//	合同资料 控制器
/////////////////////////////////////////////////////////////////////////////

import (
	"YYUEsys/models"
	"YYUEsys/models/utils"
	"strconv"
	"strings"

	"github.com/beego/beego/v2/adapter/validation"
	beego "github.com/beego/beego/v2/server/web"
)

//ContractJSON 列表json格式
type ContractJSON models.PubJSON

//contractSimpJSON 简易JSON
type contractSimpJSON models.SimpJSON

//ContractController 合同数据控制器
type ContractController struct {
	beego.Controller
}

//logsContract 显示 新增合同页面 方法
func logsContract(c *ContractController, Method string, code int, event string, info string) {
	var log utils.LogStruct
	ip := c.Ctx.Request.RemoteAddr
	ip = ip[0:strings.LastIndex(ip, ":")]

	log.Module = "基础数据-合同"
	log.Remote = ip
	log.UID = c.GetSession("UID").(int)
	log.CID = c.GetSession("CID").(int)
	log.Method = Method
	log.Code = code
	log.Event = event
	log.Info = info
	go utils.Logger(log)
}

//PageIndex 显示 合同列表页面 方法
//============================================================
func (c *ContractController) PageIndex() {
	rs, _ := models.GetServiceList(c.GetSession("CID").(int), 1, 9999)

	c.Data["Servicelist"] = rs
	c.TplName = "basic/contract_index.html"
}

//PageAdd 显示 新增合同页面 方法
//============================================================
func (c *ContractController) PageAdd() {
	allservice, _ := models.GetServiceList(c.GetSession("CID").(int), 1, 10000) // 所有服务项目列表,

	//给服务列表赋值，从服务条款读取值
	c.Data["Serviceitemlist"] = allservice //赋值给模板显示已有值

	c.TplName = "basic/contract_edit.html"

}

//PageEdit 显示 编辑合同页面 方法
//============================================================
func (c *ContractController) PageEdit() {
	id, _ := c.GetInt("ID")

	item := models.GetContractItem(id)
	// logs.Info(item)
	//给合同主体内容赋值
	c.Data["ID"] = item.ID
	c.Data["Name"] = item.Name
	c.Data["Duration"] = item.Duration
	c.Data["Price"] = item.Price
	c.Data["State"] = item.State
	c.Data["Memo"] = item.Memo

	allservice, _ := models.GetServiceList(c.GetSession("CID").(int), 1, 9999) // 所有服务项目列表,
	items := models.GetContractItems(int64(id))                                //获取相关服务列表

	//给相关服务列表赋值，从服务条款读取值
	for k, v := range allservice {
		for _, v1 := range items {
			if v1["Service__Name"] == v["Name"] {
				allservice[k]["Quantity"] = v1["Quantity"]
			}
		}
	}

	c.Data["Serviceitemlist"] = allservice //赋值给模板循环显示变量赋值

	c.TplName = "basic/contract_edit.html"
}

//GetList 合同列表
///////////////////////////////////////////////////////////////////////////////
func (c *ContractController) GetList() {
	page, _ := c.GetInt("page")
	limit, _ := c.GetInt("limit")

	rslist, count := models.GetContractList(c.GetSession("CID").(int), page, limit)

	for k, v := range rslist {
		ID := v["ID"].(int64)
		items := models.GetContractItems(ID) //获取相关服务列表

		for _, v1 := range items { //赋值服务数量
			keystr := v1["Service__Name"].(string)
			rslist[k][keystr] = v1["Quantity"]
		}
	}

	//准备json
	var code int = 200
	var msg string = "success"
	if len(rslist) == 0 {
		code = 100
		msg = "无数据。"
	}

	//返回Json
	Datastr := ContractJSON{code, msg, count, rslist}
	c.Data["json"] = Datastr
	c.ServeJSON()

	//日志输出
	go logsContract(c, "list", code, "合同列表", "Total:"+strconv.FormatInt(count, 10))
}

//Save 保存合同方法
///////////////////////////////////////////////////////////////////////////////
func (c *ContractController) Save() {
	//准备数据
	var item models.Contract
	var err error

	//准备 合同基础内容
	id, _ := c.GetInt("id") //编辑id变量
	item.ID = id
	item.Name = strings.Trim(c.GetString("Name"), "")
	item.Duration, _ = c.GetFloat("Duration")
	item.Price, _ = c.GetFloat("Price")
	item.State = strings.Trim(c.GetString("State"), "")
	item.Memo = c.GetString("Memo")
	item.CID = c.GetSession("CID").(int)
	//准备 相关服务内容
	allservices, _ := models.GetServiceList(c.GetSession("CID").(int), 1, 9999) // 所有服务项目列表,
	for k, v := range allservices {
		quantity, _ := c.GetFloat(v["Name"].(string))
		allservices[k]["Quantity"] = quantity
	}

	//准备 Json
	var code int = 200
	var errinfo string

	//验证表单
	cvinfo := ValidContract(item)
	errinfo += cvinfo

	//重复检查
	cdinfo := models.CheckContractDuplicate(item)

	errinfo += cdinfo
	if errinfo != "" {
		code = 100
	} else {
		//保存
		//如果无id号新增，否则进行更新
		if id == 0 {
			err = models.AddContract(item, allservices) //新增
		} else {
			err = models.UpdateContract(item, allservices) ///编辑
		}
		if err != nil {
			code = 100
			errinfo += "写数据失败，请刷新重试"
		}
	}

	// //返回 Json
	rs := &contractSimpJSON{code, errinfo}
	c.Data["json"] = rs
	c.ServeJSON()

	//日志输出
	go logsContract(c, "save", code, "合同保存", item.Name)
}

//Del 删除合同 方法
///////////////////////////////////////////////////////////////////////////////
func (c *ContractController) Del() {
	//准备数据
	id, _ := c.GetInt("ID")
	//执行删除
	err := models.DelContract(id)
	//准备 Json
	code := 200
	msg := "success"
	if err != nil {
		code = 100
		msg = "error"
	}
	// 返回 Json
	rs := &contractSimpJSON{code, msg}
	c.Data["json"] = rs
	c.ServeJSON()
	//日志输出
	go logsContract(c, "del", code, "删除合同", strconv.Itoa(id))
}

//ValidContract 表单验证
///////////////////////////////////////////////////////////////////////////////
func ValidContract(item models.Contract) string {
	var errinfo string
	v := validation.Validation{}
	//验证项目
	v.Required(item.Name, "合同名称").Message("不能为空")
	v.MinSize(item.Name, 2, "合同名称").Message("最少2个字符")
	v.MaxSize(item.Name, 16, "合同名称").Message("最大16个字符")
	v.Numeric(strconv.FormatFloat(item.Duration, 'f', -1, 64), "预设期限").Message("必须为大于等于0的整数")
	//-1参数不保留末尾0，可以检测是否为整数。如有小数点不通过检查

	//准备异常信息
	for _, err := range v.Errors {
		errinfo = errinfo + err.Key + " : " + err.Message + "<br> "
	}

	//返回异常提示
	return errinfo
}
