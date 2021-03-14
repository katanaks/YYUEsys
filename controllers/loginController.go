package controllers

/////////////////////////////////////////////////////////////////////////////
//	用户验证 控制器
//  1.
//  2.
/////////////////////////////////////////////////////////////////////////////

import (
	beego "github.com/beego/beego/v2/adapter"
	_ "github.com/beego/beego/v2/adapter/session/mysql" //session引擎
)

//VerificationController 用户验证控制器
type VerificationController struct {
	beego.Controller
}

//PageLogin 显示 登录页面
func (c *VerificationController) PageLogin() {
	c.TplName = "login.html"

}

//Check 检查用户登录
func (c *VerificationController) Check() {

	c.SetSession("UID", 3)  //用户ID
	c.SetSession("CID", 15) //公司ID

	rs := "登录成功"
	c.Data["json"] = rs
	c.ServeJSON()
}
