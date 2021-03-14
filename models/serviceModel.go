package models

/////////////////////////////////////////////////////////////////////////////
//	服务项目资料 Mode
/////////////////////////////////////////////////////////////////////////////

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
)

//Init 初始化
func init() {
	//连接数据库
	// utils.InitDB()
	//注册数据结构
	// orm.RegisterModel()
}

//GetServiceList 读取列表
func GetServiceList(cid, page, limit int) ([]orm.Params, int64) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Service))

	var res []orm.Params

	count, err := qs.Filter("CID", cid).OrderBy("ID").Limit(limit, page*limit-limit).RelatedSel().
		Values(&res, "ID", "Name", "State", "Memo", "Updatetime")

	if err != nil {
		logs.Error(err)
	}

	//返回 列表，计数
	return res, count
}

//DelService 删除员工
func DelService(id int) error {
	o := orm.NewOrm()

	_, err := o.Delete(&Service{ID: id})

	if err != nil {
		logs.Error(err)
	}

	return err
}

//AddService 新增
func AddService(item Service) error {
	o := orm.NewOrm()

	_, err := o.Insert(&item)
	// logs.Debug("写入", item)

	if err != nil {
		logs.Error(err.Error())
	}

	return err
}

//UpdateService 更新
func UpdateService(item Service) error {
	o := orm.NewOrm()
	_, err := o.Update(&item)
	return err
}

//GetServiceItem 读取编辑记录
func GetServiceItem(ID int) Service {
	o := orm.NewOrm()
	var item Service
	item.ID = ID
	o.Read(&item, "ID")

	return item
}

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

//CheckServiceDuplicate 检测服务名称重复
func CheckServiceDuplicate(item Service) string {
	o := orm.NewOrm()
	name := item.Name
	id := item.ID
	cid := item.CID

	var errinfo string

	count, _ := o.QueryTable("Service").Filter("CID", cid).Filter("Name", name).Exclude("ID", id).Count()

	// logs.Info("item.id:", item.ID, "count:", count, "name:", name, "CID:", cid)
	if count > 0 {
		errinfo = "服务名已存在"
	}

	return errinfo
}
