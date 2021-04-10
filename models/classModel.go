package models

/////////////////////////////////////////////////////////////////////////////
//	员工资料 Model
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

}

//GetClassList 读取列表
func GetClassList(cid, page, limit int) ([]orm.Params, int64) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Class))

	var res []orm.Params

	count, err := qs.Filter("CID", cid).OrderBy("ID").Limit(limit, page*limit-limit).RelatedSel().
		Values(&res, "ID", "Name", "Staff__Name", "StartDates", "Memo", "Updatetime")

	if err != nil {
		logs.Error(err)
	}

	//返回 列表，计数
	return res, count
}

//DelClass 删除员工
func DelClass(id int) error {
	o := orm.NewOrm()

	_, err := o.Delete(&Class{ID: id})

	if err != nil {
		logs.Error(err)
	}

	return err
}

//AddClass 新增
func AddClass(item Class) error {
	o := orm.NewOrm()

	_, err := o.Insert(&item)
	// logs.Debug("写入", item)

	if err != nil {
		logs.Error(err.Error)
	}

	return err
}

//UpdateClass 更新
func UpdateClass(item Class) error {
	o := orm.NewOrm()
	_, err := o.Update(&item)
	return err
}

//GetClassItem 读取编辑记录
func GetClassItem(ID int) Class {
	o := orm.NewOrm()
	var item Class
	item.ID = ID
	_ = o.Read(&item, "ID")

	return item
}

//CheckClassDuplicate 检测部门名称重复
func CheckClassDuplicate(item Class) string {
	o := orm.NewOrm()
	name := item.Name
	id := item.ID
	cid := item.CID

	var errinfo string

	count, _ := o.QueryTable("Class").Filter("CID", cid).Filter("Name", name).Exclude("ID", id).Count()

	// logs.Info("item.id:", item.ID, "count:", count, "name:", name, "CID:", cid)
	if count > 0 {
		errinfo = "班级名已存在"
	}

	return errinfo
}
