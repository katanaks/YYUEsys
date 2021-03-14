package controllers

import beego "github.com/beego/beego/v2/adapter"

//MainController 系统主控制器
type MainController struct {
	beego.Controller
}

//Welcome 欢迎页面
type Welcome struct {
	beego.Controller
}

//Get 显示框架
func (c *MainController) Get() {
	c.TplName = "index.html" //加载页面
}

//Get 显示欢迎内页
func (c *Welcome) Get() {
	c.TplName = "welcome.html"
}
