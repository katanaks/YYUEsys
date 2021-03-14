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
	// orm.RegisterModel()
}

//GetStaffList 读取列表
func GetStaffList(cid, page, limit int) ([]orm.Params, int64) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Staff))

	var res []orm.Params

	count, err := qs.Filter("CID", cid).OrderBy("ID").Limit(limit, page*limit-limit).RelatedSel().
		Values(&res, "ID", "Name", "State", "Telephone", "Gender", "Entrydate", "Department__Name", "Birthday", "Special", "School", "Education", "Address", "Memo", "Updatetime")

	if err != nil {
		logs.Error(err)
	}

	//返回 列表，计数
	return res, count
}

//DelStaff 删除员工
func DelStaff(id int) error {
	o := orm.NewOrm()

	_, err := o.Delete(&Staff{ID: id})

	if err != nil {
		logs.Error(err.Error())
	}

	return err
}

//AddStaff 读取列表
func AddStaff(item Staff) error {
	o := orm.NewOrm()

	_, err := o.Insert(&item)
	// logs.Info(item)

	if err != nil {
		logs.Error(err.Error())
	}

	return err
}

//UpdateStaff 更新数据
func UpdateStaff(item Staff) error {
	o := orm.NewOrm()
	_, err := o.Update(&item)
	return err
}

//GetStaffItem 读取员式记录
func GetStaffItem(ID int) Staff {
	o := orm.NewOrm()
	var item Staff
	item.ID = ID
	o.Read(&item, "ID")

	return item
}

//CheckStaffDuplicate 检测部门名称重复
func CheckStaffDuplicate(item Staff) string {
	o := orm.NewOrm()
	name := item.Name
	id := item.ID
	cid := item.CID

	var errinfo string

	count, _ := o.QueryTable("staff").Filter("CID", cid).Filter("Name", name).Exclude("ID", id).Count()

	// logs.Info("item.id:", item.ID, "count:", count, "name:", name, "CID:", cid)
	if count > 0 {
		errinfo = "员工姓名存在重复"
	}

	return errinfo
}
