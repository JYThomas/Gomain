package routers

import (
	"api/controllers"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	// 设置控制器下方法的命名空间，规整接口
    ns := beego.NewNamespace("v1",
		beego.NSCtrlGet("/users", (*controllers.UserController).ListUsers),
		beego.NSCtrlGet("/users/:id", (*controllers.UserController).GetUser),
	)
	// 注册命名空间
	beego.AddNamespace(ns)
}
