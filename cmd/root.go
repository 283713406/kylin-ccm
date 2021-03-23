package cmd

import (
	"github.com/astaxie/beego"
	_ "kylin-ccm/routers"
	"kylin-ccm/service"
)

func Run() {

    // mysql
	service.Init()

	beego.Run()
}

