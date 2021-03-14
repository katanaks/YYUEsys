package models

/////////////////////////////////////////////////////////////////////////////
//	合同资料 Model
/////////////////////////////////////////////////////////////////////////////

import (
	"fmt"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
)

//Init 初始化
func init() {
	//连接数据库

	// utils.InitDB()
	//注册数据结构
}

//GetContractList 读取合同列表
func GetContractList(cid, page, limit int) ([]orm.Params, int64) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Contract))

	var res []orm.Params
	count, _ := qs.Filter("CID", cid).Count() //单独获取计数
	qs.Filter("CID", cid).OrderBy("-ID").Limit(limit, page*limit-limit).RelatedSel().
		Values(&res, "ID", "Name", "State", "Duration", "Price", "Memo", "Updatetime") //获取合同列表

	//返回 列表，计数
	return res, count
}

//GetContractItems 读取合同相关服务
func GetContractItems(ID int64) []orm.Params {
	var item []orm.Params

	orm.NewOrm().QueryTable("Contractitem").Filter("Contract", ID).
		Exclude("Quantity", 0).Exclude("Quantity__isnull", true).
		RelatedSel().
		Values(&item, "ID", "Service__name", "Quantity")

	return item
}

//DelContract 删除合同
func DelContract(id int) error {
	o := orm.NewOrm()
	//删除合同，同时删除相关条目
	_, err := o.Delete(&Contract{ID: id})
	//logs.Error("del id:", id)

	if err != nil {
		logs.Error(err.Error())
	}

	return err
}

//AddContract 读取列表
func AddContract(item Contract, allservices []orm.Params) error {
	o := orm.NewOrm()
	var err error
	//开始插入事务
	to, err := o.Begin()

	//插入合同主体
	InsID, _ := to.Insert(&item)
	//插入相关服务和数量
	item.ID = int(InsID)
	SaveserviceItem(o, item, allservices)

	//错误处理
	if err != nil {
		to.Rollback() //存在错误，回滚操作
		logs.Error(err.Error())
	} else {
		to.Commit() //提交插入-成功
	}
	return err
}

//UpdateContract 更新数据
func UpdateContract(item Contract, allservices []orm.Params) error {
	o := orm.NewOrm()
	//更新主体

	_, err := o.Update(&item)
	//清除原服务项目
	orm.NewOrm().QueryTable("Contractitem").Filter("Contract", item.ID).Delete()
	//插入相关服务和数量
	SaveserviceItem(o, item, allservices)

	return err
}

//SaveserviceItem 保存服务条款
func SaveserviceItem(o orm.Ormer, item Contract, allservices []orm.Params) {
	// o := orm.NewOrm()

	//插入相关服务和数量
	var serviceItem Contractitem
	for _, v := range allservices {
		Sid := v["ID"].(int64)
		serviceItem.ID = 0 //复位插入主键ID，否则错误
		serviceItem.CID = item.CID
		serviceItem.Contract = &Contract{ID: int(item.ID)}
		serviceItem.Service = &Service{ID: int(Sid)}
		serviceItem.Quantity = v["Quantity"].(float64)
		//逐条插入
		_, err := o.Insert(&serviceItem)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	return
}

//GetContractItem 读取记录
func GetContractItem(ID int) Contract {
	o := orm.NewOrm()
	var item Contract
	item.ID = ID
	o.Read(&item, "ID")

	return item
}

//CheckContractDuplicate 检测合同名称重复
func CheckContractDuplicate(item Contract) string {
	o := orm.NewOrm()
	name := item.Name
	id := item.ID
	cid := item.CID

	var errinfo string

	count, _ := o.QueryTable("Contract").Filter("CID", cid).Filter("Name", name).Exclude("ID", id).Count()

	// logs.Info("item.id:", item.ID, "count:", count, "name:", name, "CID:", cid)
	if count > 0 {
		errinfo = "合同姓名存在重复"
	}

	return errinfo
}
