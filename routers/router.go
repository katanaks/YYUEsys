package routers

import (
	"YYUEsys/controllers"
	"strings"

	beego "github.com/beego/beego/v2/adapter"
	"github.com/beego/beego/v2/adapter/context"
)

//FilterUser 注册过滤器
var FilterUser = func(ctx *context.Context) {
	_, ok1 := ctx.Input.Session("UID").(int)

	ok2 := strings.Contains(ctx.Request.RequestURI, "/login")
	ok3 := strings.Contains(ctx.Request.RequestURI, "/check")

	if !ok1 && !ok2 && !ok3 {
		ctx.Redirect(302, "/login")
	}
}

func init() {
	// 注册器 验证
	beego.InsertFilter("/*", beego.BeforeRouter, FilterUser)

	beego.Router("/login", &controllers.VerificationController{}, "get:PageLogin") //显示 登录
	beego.Router("/check", &controllers.VerificationController{}, "post:Check")    //显示 登录

	//首页
	beego.Router("/", &controllers.MainController{})
	beego.Router("/welcome", &controllers.Welcome{})

	//=================================================================================================
	//基础数据路由
	//==================================================================================================

	//部门
	beego.Router("/basic/department/", &controllers.DepartmentController{}, "get:PageDepartmentIndex")        //显示 列表
	beego.Router("/basic/department/pageadd", &controllers.DepartmentController{}, "get:PageDepartmentAdd")   //显示 增加
	beego.Router("/basic/department/pageedit", &controllers.DepartmentController{}, "get:PageDepartmentEdit") //显示 编辑

	beego.Router("/basic/department/getlist", &controllers.DepartmentController{}, "get:GetDepartmentList") //方法 获取列表
	beego.Router("/basic/department/save", &controllers.DepartmentController{}, "post:SaveDepartment")      //方法 增加
	beego.Router("/basic/department/del", &controllers.DepartmentController{}, "post:DelDepartment")        //方法 删除

	//员工
	//=====================================
	beego.Router("/basic/staff/", &controllers.StaffController{}, "get:PageIndex")        //显示 列表
	beego.Router("/basic/staff/pageadd", &controllers.StaffController{}, "get:PageAdd")   //显示 增加
	beego.Router("/basic/staff/pageedit", &controllers.StaffController{}, "get:PageEdit") //显示 编辑

	beego.Router("/basic/staff/getlist", &controllers.StaffController{}, "get:GetList") //方法 获取列表
	beego.Router("/basic/staff/save", &controllers.StaffController{}, "post:Save")      //方法 保存
	beego.Router("/basic/staff/del", &controllers.StaffController{}, "post:Del")        //方法 删除

	//班级
	//=====================================
	beego.Router("/basic/class/", &controllers.ClassController{}, "get:PageIndex")        //显示 列表
	beego.Router("/basic/class/pageadd", &controllers.ClassController{}, "get:PageAdd")   //显示 增加
	beego.Router("/basic/class/pageedit", &controllers.ClassController{}, "get:PageEdit") //显示 编辑

	beego.Router("/basic/class/getlist", &controllers.ClassController{}, "get:GetList") //方法 获取列表
	beego.Router("/basic/class/save", &controllers.ClassController{}, "post:Save")      //方法 保存
	beego.Router("/basic/class/del", &controllers.ClassController{}, "post:Del")        //方法 删除

	//服务项目
	//=====================================
	beego.Router("/basic/service/", &controllers.ServiceController{}, "get:PageIndex")        //显示 列表
	beego.Router("/basic/service/pageadd", &controllers.ServiceController{}, "get:PageAdd")   //显示 增加
	beego.Router("/basic/service/pageedit", &controllers.ServiceController{}, "get:PageEdit") //显示 编辑

	beego.Router("/basic/service/getlist", &controllers.ServiceController{}, "get:GetList") //方法 获取列表
	beego.Router("/basic/service/save", &controllers.ServiceController{}, "post:Save")      //方法 保存
	beego.Router("/basic/service/del", &controllers.ServiceController{}, "post:Del")        //方法 删除

	//合同
	//=====================================
	beego.Router("/basic/contract/", &controllers.ContractController{}, "get:PageIndex")        //显示 列表
	beego.Router("/basic/contract/pageadd", &controllers.ContractController{}, "get:PageAdd")   //显示 增加
	beego.Router("/basic/contract/pageedit", &controllers.ContractController{}, "get:PageEdit") //显示 编辑

	beego.Router("/basic/contract/getlist", &controllers.ContractController{}, "get:GetList") //方法 获取列表
	beego.Router("/basic/contract/save", &controllers.ContractController{}, "post:Save")      //方法 保存
	beego.Router("/basic/contract/del", &controllers.ContractController{}, "post:Del")        //方法 删除

	//会员管理
	//=====================================
	beego.Router("/members/", &controllers.MembersController{}, "get:PageIndex")      //显示 列表
	beego.Router("/members/pageadd", &controllers.MembersController{}, "get:PageAdd") //显示 增加
	// beego.Router("/basic/service/pageedit", &controllers.ServiceController{}, "get:PageEdit") //显示 编辑
	beego.Router("/members/pagerenewal", &controllers.MembersController{}, "get:PageRenewal") //显示 续费约

	beego.Router("/members/getlist", &controllers.MembersController{}, "get:GetList") //方法 获取列表
	beego.Router("/members/save", &controllers.MembersController{}, "post:Save")      //方法 保存
	// beego.Router("/members/del", &controllers.ServiceController{}, "post:Del")        //方法 删除
	beego.Router("/members/getpredata", &controllers.MembersController{}, "get:Getpredata") //方法 获取列表
}
