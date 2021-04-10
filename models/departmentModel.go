package models

/////////////////////////////////////////////////////////////////////////////
//	部门资料 Model
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

//GetDepartmentList 读取列表
func GetDepartmentList(cid, page, limit int) ([]orm.Params, int64) {
	o := orm.NewOrm()
	var res []orm.Params //返回变量

	count, err := o.QueryTable("department").Filter("CID", cid).OrderBy("ID").Limit(limit, page*limit-limit).Values(&res, "ID", "Name", "Memo", "Createtime", "Updatetime")

	if err != nil {
		logs.Error(err.Error())
	}

	return res, count
}

//Del 读取列表
func Del(id int) error {
	o := orm.NewOrm()

	_, err := o.Delete(&Department{ID: id})

	return err
}

//Add 读取列表
func Add(item Department) error {
	o := orm.NewOrm()

	_, err := o.Insert(&item)

	return err
}

//Update 更新数据
func Update(item Department) error {
	o := orm.NewOrm()
	_, err := o.Update(&item)
	return err
}

//GetDepartmentItem 读取编辑记录
func GetDepartmentItem(ID int) Department {
	o := orm.NewOrm()
	var item Department
	item.ID = ID
	_ = o.Read(&item, "ID")

	return item
}

//CheckDuplicate 检测部门名称重复
func CheckDuplicate(item Department) string {
	o := orm.NewOrm()
	name := item.Name
	id := item.ID
	cid := item.CID

	var errinfo string

	count, _ := o.QueryTable("department").Filter("CID", cid).Filter("Name", name).Exclude("ID", id).Count()

	if count > 0 {
		errinfo = "部门名称存在重复"
	}
	return errinfo
}

//GetDepartmentSelList 获取select列表
func GetDepartmentSelList(cid int) ([]orm.Params, int64) {
	o := orm.NewOrm()
	var res []orm.Params //返回变量

	count, _ := o.QueryTable("department").Filter("CID", cid).OrderBy("ID").Values(&res, "ID", "Name")

	return res, count
}
