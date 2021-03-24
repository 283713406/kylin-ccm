package base

import (
	"strings"

	"github.com/astaxie/beego"
)

type BaseController struct {
	beego.Controller
	controllerName string               // 控制器名
}

// 重写GetString方法，移除前后空格
func (b *BaseController) GetString(name string, def ...string) string {
	return strings.TrimSpace(b.Controller.GetString(name, def...))
}

func (b *BaseController) GetBool(name string) bool {
	bo, _ := b.Controller.GetBool(name)
	return bo
}

func (b *BaseController) isAllInOne(name string, def ...string) bool {
	return b.Controller.GetString("is_all_in_one", def...) == "true"
}

func (b *BaseController) Prepare() {
	controllerName, _ := b.GetControllerAndAction()
	b.controllerName = strings.ToLower(controllerName[0 : len(controllerName)-10])
}

// 是否POST提交
func (b *BaseController) IsPost() bool {
	return b.Ctx.Request.Method == "POST"
}

//获取用户IP地址
func (b *BaseController) getClientIp() string {
	if p := b.Ctx.Input.Proxy(); len(p) > 0 {
		return p[0]
	}
	return b.Ctx.Input.IP()
}