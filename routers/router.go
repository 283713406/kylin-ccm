package routers

import (
	"github.com/astaxie/beego"
	"kylin-ccm/controllers"
)

func init() {
    beego.Router("/api/v1/clusters", &controllers.ClusterController{})
	beego.Router("/api/v1/clusters/?:name", &controllers.ClusterController{}, "get:Get")
	beego.Router("/api/v1/clusters/?:name", &controllers.ClusterController{}, "delete:Del")
}
