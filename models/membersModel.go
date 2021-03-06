package models

/////////////////////////////////////////////////////////////////////////////
//	服务项目资料 Mode
/////////////////////////////////////////////////////////////////////////////

import (
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"time"
)

//Init 初始化
func init() {
	//连接数据库
	// utils.InitDB()
	//注册数据结构
	// orm.RegisterModel()
}

//GetMembersList 读取列表
func GetMembersList(cid, page, limit int) ([]orm.Params, int64) {
	var res []orm.Params

	count, err := orm.NewOrm().QueryTable("Members").
		Filter("CID", cid).
		OrderBy("ID").
		Limit(limit, page*limit-limit).RelatedSel().
		Values(&res, "ID", "Name", "Gender", "Birthday", "Contacttelephone", "State", "Updatetime")

	if err != nil {
		logs.Error(err)
	}

	//返回 列表，计数
	return res, count
}

// //DelService 删除员工
// func DelService(id int) error {
// 	o := orm.NewOrm()

// 	_, err := o.Delete(&Service{ID: id})

// 	if err != nil {
// 		logs.Error(err)
// 	}

// 	return err
// }

//AddMember 新增
func AddMember(item Members, itemtransaction Membertransaction) error {
	o := orm.NewOrm()
	var err error
	//开始插入事务
	to, err := o.Begin()

	//插入会员
	//====================
	to.Insert(&item)
	//插入相关服务和数量
	itemtransaction.Members = &item

	_, err = o.Insert(&itemtransaction)

	//错误处理
	if err != nil {
		to.Rollback() //存在错误，回滚操作
		logs.Error(err.Error())
	} else {
		to.Commit() //提交插入-成功
	}
	return err
}

// //UpdateService 更新
// func UpdateService(item Service) error {
// 	o := orm.NewOrm()
// 	_, err := o.Update(&item)
// 	return err
// }

// //GetServiceItem 读取编辑记录
// func GetServiceItem(ID int) Service {
// 	o := orm.NewOrm()
// 	var item Service
// 	item.ID = ID
// 	o.Read(&item, "ID")

// 	return item
// }

//GetServiceItemList 读取服务列表
// func GetServiceItemList(cid int) []orm.Params {
// 	o := orm.NewOrm()
// 	qs := o.QueryTable(new(Service))

// 	var res []orm.Params

// 	qs.Filter("CID", cid).OrderBy("-ID").Limit(500).RelatedSel().
// 		Values(&res, "ID", "Name") //获取合同列表

// 	//返回
// 	return res
// }

// //CheckServiceDuplicate 检测服务名称重复
// func CheckServiceDuplicate(item Service) string {
// 	o := orm.NewOrm()
// 	name := item.Name
// 	id := item.ID
// 	cid := item.CID

// 	var errinfo string

// 	count, _ := o.QueryTable("Service").Filter("CID", cid).Filter("Name", name).Exclude("ID", id).Count()

// 	// logs.Info("item.id:", item.ID, "count:", count, "name:", name, "CID:", cid)
// 	if count > 0 {
// 		errinfo = "服务名已存在"
// 	}

// 	return errinfo
// }

//Getpredata 获取签约表单合同刷新数据表单
func Getpredata(ID int) []orm.Params {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Contract))

	var res []orm.Params
	qs.Filter("ID", ID).Values(&res, "ID", "Price", "Duration")
	return res
}

//GetMembertransactionItem 读取会员交易列表
func GetMembertransactionItem(MemberID int64) []orm.Params {
	//var ContractName string
	var item []orm.Params
	orm.NewOrm().QueryTable("Membertransaction").
		Filter("MemberID", MemberID).
		OrderBy("Durationend").
		Limit(1).
		RelatedSel().
		Values(&item)

	item[0]["aaaaaaaaaaaaddd"] = "kdkdkdk"

	for k, v := range item {
		Sid := v["ContractID"].(int64)
		var itemContract []orm.Params
		orm.NewOrm().QueryTable("Contract").
			Filter("ID", Sid).
			OrderBy("ID").
			Limit(1).
			RelatedSel().
			Values(&itemContract)

		item[k]["ContractName"] = itemContract[0]["Name"]
	}
	return item

}

//GetMemberContract 获取会员合同名称
func GetMemberContract(MemberID int64) string {
	//对应会员交易
	var Membertransactions []orm.Params
	//查询会员交易记录
	orm.NewOrm().QueryTable("Membertransaction").
		Filter("Members", MemberID).
		//Filter("Durationstart__lt", time.Now()).
		//Filter("Durationend__gt", time.Now()).
		OrderBy("-Durationend").
		RelatedSel().
		Values(&Membertransactions)

	fmt.Println("相关交易", Membertransactions)

	var ContractItem []orm.Params
	//var ContractName string
	//查询对应合同名称
	if Membertransactions != nil {
		orm.NewOrm().QueryTable("Contract").
			Filter("ID", Membertransactions[0]["Contract"]).
			Values(&ContractItem)

		ContractName := ContractItem[0]["Name"].(string)

		return ContractName
	} else {
		return ""
	}

}

//GetMemberDurationend 获取会员到期日期
func GetMemberDurationend(MemberID int64) time.Time {
	var Membertransactions []orm.Params

	orm.NewOrm().QueryTable("Membertransaction").
		Filter("Members", MemberID).
		OrderBy("-Durationend").
		RelatedSel().
		Values(&Membertransactions)

	MemberDurationend := Membertransactions[0]["Durationend"].(time.Time)
	return MemberDurationend
}
