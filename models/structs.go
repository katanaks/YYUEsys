package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

//PubJSON json格式
//============================================================
type PubJSON struct {
	//必须的大写开头
	Code  int          `json:"code"`
	Msg   string       `json:"msg"`
	Count int64        `json:"count"`
	Data  []orm.Params `json:"data"`
}

//SimpJSON json格式
//============================================================
type SimpJSON struct {
	//必须的大写开头
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

//Department 部门表 结构
type Department struct {
	ID         int       `orm:"auto;pk;column(i_d);size(11);description(部门ID)"`                           //部门ID
	CID        int       `orm:"index;column(c_id);size(11);description(所属公司)"`                            //所属公司
	Name       string    `orm:"index;column(name);size(32);description(部门名称)"`                            //部门名称
	Memo       string    `orm:"column(memo);size(255);description(备注)"`                                   //备注
	Createtime time.Time `orm:"auto_now_add;type(datetime);column(createtime);size(0);description(新建时间)"` //自动新建时间
	Updatetime time.Time `orm:"auto_now;type(datetime);column(updatetime);size(0);description(更新时间)"`     //自动更新时间
}

//Staff 员工表
type Staff struct {
	ID         int         `orm:"auto;pk;column(i_d);size(11);description(员工ID)"`
	CID        int         `orm:"index;column(c_id);size(11);description(所属公司)"`
	Name       string      `orm:"index;column(name);size(8);description(员工姓名)"`
	State      string      `orm:"column(state);size(8);description(状态)"`
	Telephone  string      `orm:"column(telephone);size(11);description(电话)"`
	Gender     string      `orm:"index;column(gender);size(2);description(性别)"`
	Entrydate  time.Time   `orm:"index;type(date);column(entrydate);size(0);description(入职日期)"`
	Department *Department `orm:"rel(fk);description(部门ID)"` //多对一 部门
	Birthday   time.Time   `orm:"type(datetime);column(birthday);size(0);description(出生年月)"`
	Special    string      `orm:"column(special);size(16);description(专业技能)"`
	School     string      `orm:"column(school);size(16);description(毕业院校)"`
	Education  string      `orm:"column(education);size(16);description(学历)"`
	Address    string      `orm:"column(address);size(64);description(住址)"`
	Memo       string      `orm:"column(memo);size(1024);description(备注)"`
	Createtime time.Time   `orm:"auto_now_add;type(datetime);column(createtime);size(0);description(新建时间)"` //自动新建时间
	Updatetime time.Time   `orm:"auto_now;type(datetime);column(updatetime);size(0);description(更新时间)"`     //自动更新时间
}

//Class 班级表
type Class struct {
	ID         int       `orm:"auto;pk;column(i_d);size(11);description(员工ID)"`
	CID        int       `orm:"index;column(c_id);size(11);description(所属公司)"`
	Name       string    `orm:"index;column(name);size(16);description(班级名称)"`
	Staff      *Staff    `orm:"rel(fk);column(staff_id);description(管理人员ID)"` //多对一 员工
	StartDates time.Time `orm:"index;type(date);column(startdates);size(0);description(开班日期)"`
	Memo       string    `orm:"column(memo);size(1024);description(备注)"`
	Createtime time.Time `orm:"auto_now_add;type(datetime);column(createtime);size(0);description(新建时间)"` //自动新建时间
	Updatetime time.Time `orm:"auto_now;type(datetime);column(updatetime);size(0);description(更新时间)"`     //自动更新时间
}

//Service 服务项目
type Service struct {
	ID         int       `orm:"auto;pk;column(i_d);size(11);description(员工ID)"`
	CID        int       `orm:"index;column(c_id);size(11);description(所属公司)"`
	Name       string    `orm:"index;column(name);size(16);description(班级名称)"`
	State      string    `orm:"index;;column(state);size(4);description(管理人员ID)"`
	Memo       string    `orm:"column(memo);size(1024);description(备注)"`
	Createtime time.Time `orm:"auto_now_add;type(datetime);column(createtime);size(0);description(新建时间)"` //自动新建时间
	Updatetime time.Time `orm:"auto_now;type(datetime);column(updatetime);size(0);description(更新时间)"`     //自动更新时间
}

//Contract 合同
type Contract struct {
	ID                int                  `orm:"auto;pk;column(i_d);size(11);description(合同ID)"`
	CID               int                  `orm:"index;column(c_id);size(11);description(所属公司)"`
	Name              string               `orm:"index;column(name);size(16);description(合同名称)"`
	Memo              string               `orm:"column(memo);size(1024);description(备注)"`
	Duration          float64              `orm:"column(duration);digits(12);decimals(2);description(备注)"`
	Price             float64              `orm:"column(price);digits(12);decimals(2);description(备注)"`
	State             string               `orm:"index;;column(state);size(4);description(启用状态)"`
	Contractitems     []*Contractitem      `orm:"reverse(many);description(合同条款)"`
	Membertransaction []*Membertransaction `orm:"reverse(many);description(交易表)"`                                           //对多 查询
	Createtime        time.Time            `orm:"auto_now_add;type(datetime);column(createtime);size(0);description(新建时间)"` //自动新建时间
	Updatetime        time.Time            `orm:"auto_now;type(datetime);column(updatetime);size(0);description(更新时间)"`     //自动更新时间
}

//Contractitem 合同内容
type Contractitem struct {
	ID         int       `orm:"auto;pk;column(i_d);size(11);description(合同ID)"`
	CID        int       `orm:"index;column(c_id);size(11);description(所属公司)"`
	Contract   *Contract `orm:"rel(fk);description(合同ID)"` //多对一 合同
	Service    *Service  `orm:"rel(fk);description(服务项目)"` //多对一 服务项目
	Quantity   float64   `orm:"column(quantity);digits(12);decimals(2);description(数量)"`
	Createtime time.Time `orm:"auto_now_add;type(datetime);column(createtime);size(0);description(新建时间)"` //自动新建时间
	Updatetime time.Time `orm:"auto_now;type(datetime);column(updatetime);size(0);description(更新时间)"`     //自动更新时间
}

//Members 会员
type Members struct {
	ID                   int       `orm:"auto;pk;column(i_d);size(11);description(会员ID)"`
	CID                  int       `orm:"index;column(c_id);size(11);description(所属公司)"`
	Name                 string    `orm:"index;column(name);size(8);description(姓名)"`                //姓名
	Gender               string    `orm:"index;column(gender);size(2);description(性别)"`              //性别
	Birthday             time.Time `orm:"type(datetime);column(birthday);size(0);description(出生年月)"` //出生年月
	Idcard               string    `orm:"column(idcard);size(18);description(身份证号)"`                 //身份证号
	Contractrelationship string    `orm:"column(contractrelationship);size(18);description(联系人关系)"`  //联系人关系
	Contactname          string    `orm:"column(contactname);size(11);description(联系人姓名)"`           //联系人姓名
	Contactidcard        string    `orm:"column(contactidcard);size(18);description(联系人身份证)"`        //联系人身份证
	Contacttelephone     string    `orm:"column(contacttelephone);size(11);description(联系人电话)"`      //联系人电话
	State                string    `orm:"column(state);size(18);description(状态)"`
	//Membersigning        []*Membersigning      `orm:"reverse(many);description(签约表)"`                                           //签约表
	Membertransaction []*Membertransaction `orm:"reverse(many);description(交易表)"` //交易表ID
	//Memberalbum          []*Memberalbum        `orm:"reverse(many);description(相册表)"`                                           //相册表
	//Memberserverrecord   []*Memberserverrecord `orm:"reverse(many);description(服务记录表)"`                                         //服务记录表
	Memo       string    `orm:"column(memo);size(1024);description(备注)"`                                  //备注
	Createtime time.Time `orm:"auto_now_add;type(datetime);column(createtime);size(0);description(新建时间)"` //自动新建时间
	Updatetime time.Time `orm:"auto_now;type(datetime);column(updatetime);size(0);description(更新时间)"`     //自动更新时间
}

//Membertransaction 交易详情
type Membertransaction struct {
	ID            int       `orm:"auto;pk;column(i_d);size(11);description(交易ID)"`
	CID           int       `orm:"index;column(c_id);size(11);description(所属公司)"`
	MemberID      int       `orm:"column(member_id);description(对应会员)"` //对应会员
	ContractID    int       `orm:"column(contract_id);description(对应对应合同)"`
	Durationstart time.Time `orm:"type(datetime);column(Durationstart);size(0);description(合同开始日期)"`         //合同开始日期
	Durationend   time.Time `orm:"type(datetime);column(Durationend);size(0);description(合同终止日期)"`           //合同终止日期
	Paid          float64   `orm:"column(paid);size(11);description(折扣)"`                                    //实付
	Memo          string    `orm:"column(memo);size(1024);description(备注)"`                                  //备注
	Createtime    time.Time `orm:"auto_now_add;type(datetime);column(createtime);size(0);description(新建时间)"` //自动新建时间
	Updatetime    time.Time `orm:"auto_now;type(datetime);column(updatetime);size(0);description(更新时间)"`     //自动更新时间
}

//Memberalbum 会员相册
type Memberalbum struct {
	ID         int       `orm:"auto;pk;column(i_d);size(11);description(会员ID)"`
	CID        int       `orm:"index;column(c_id);size(11);description(所属公司)"`
	Member     *Members  `orm:"rel(fk);description(对应会员)"`                                                //对应会员
	Memo       string    `orm:"column(memo);size(1024);description(备注)"`                                  //备注
	Createtime time.Time `orm:"auto_now_add;type(datetime);column(createtime);size(0);description(新建时间)"` //自动新建时间
	Updatetime time.Time `orm:"auto_now;type(datetime);column(updatetime);size(0);description(更新时间)"`     //自动更新时间
}

//Memberserverrecord 服务记录表
type Memberserverrecord struct {
	ID         int       `orm:"auto;pk;column(i_d);size(11);description(会员ID)"`
	CID        int       `orm:"index;column(c_id);size(11);description(所属公司)"`
	Member     *Members  `orm:"rel(fk);description(对应会员)"`                                                //对应会员
	Memo       string    `orm:"column(memo);size(1024);description(备注)"`                                  //备注
	Createtime time.Time `orm:"auto_now_add;type(datetime);column(createtime);size(0);description(新建时间)"` //自动新建时间
	Updatetime time.Time `orm:"auto_now;type(datetime);column(updatetime);size(0);description(更新时间)"`     //自动更新时间
}

//Membersigning 签约记录表
type Membersigning struct {
	ID         int       `orm:"auto;pk;column(i_d);size(11);description(会员ID)"`
	CID        int       `orm:"index;column(c_id);size(11);description(所属公司)"`
	Member     *Members  `orm:"rel(fk);description(对应会员)"`                                                //对应会员
	Memo       string    `orm:"column(memo);size(1024);description(备注)"`                                  //备注
	Createtime time.Time `orm:"auto_now_add;type(datetime);column(createtime);size(0);description(新建时间)"` //自动新建时间
	Updatetime time.Time `orm:"auto_now;type(datetime);column(updatetime);size(0);description(更新时间)"`     //自动更新时间
}

//表结构注册
func init() {
	orm.RegisterModel(
		new(Class), new(Department), new(Staff), new(Service), new(Contract), new(Contractitem),
		new(Members), new(Membertransaction), new(Memberalbum), new(Memberserverrecord), new(Membersigning),
	)
}
