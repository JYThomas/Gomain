package controllers

import (
	"github.com/beego/beego/v2/server/web"
)

type UserController struct {
	web.Controller
}

// 获取用户列表
func (user *UserController) ListUsers() {
	user.Data["json"] = map[string]string{"status": "success", "message": "User list"}
	user.ServeJSON()
}

// 获取单个用户信息
func (user *UserController) GetUser() {
	user.Data["json"] = map[string]string{"status": "success", "message": "User details"}
	user.ServeJSON()
}