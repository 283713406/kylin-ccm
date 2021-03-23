package controllers

import (
	"encoding/json"
	"kylin-ccm/pkg/util/logs"

	"kylin-ccm/entity"
	"kylin-ccm/pkg/install"
	"kylin-ccm/service"
)

type ClusterController struct {
	BaseController
}

func (c *ClusterController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}

func (c *ClusterController) Post() {
	if c.isPost() {
		logs.MyLogger.Infof("Start create cluster")
		var cluster entity.SingleCluster
		data := c.Ctx.Input.RequestBody
		//json数据封装到user对象中
		err := json.Unmarshal(data, &cluster)
		if err != nil {
			logs.MyLogger.Infof( "json.Unmarshal is err:", err.Error())
		}

		isAllInOne := false
		ksEnable := false
		if cluster.IsAllInOne == "true" {
			isAllInOne = true
		}

		if cluster.KsEnable == "true" {
			ksEnable = true
		}

		if err = install.CreateCluster(cluster.NodeName, cluster.UserName, cluster.K8sVersion,
			cluster.KsVersion, isAllInOne, ksEnable); err != nil {
			logs.MyLogger.Errorf("Failed to create cluster, error: %v", err.Error())
		}

		if err = service.ClusterService.AddCluster(cluster.ClusterName, cluster.NodeName,
			cluster.K8sVersion); err != nil {
			logs.MyLogger.Errorf("Failed to add cluster, error: %v", err.Error())
		}
	}
	c.Data["pageTitle"] = "添加集群"

}

//func (c *ClusterController) Del() {
//	id, _ := c.GetInt("id")
//
//	if id == 1 {
//		c.showMsg("不能删除ID为1的帐号", MSG_ERR)
//	}
//
//	err := service.UserService.DeleteUser(id)
//	c.checkError(err)
//
//	c.redirect(beego.URLFor("UserController.List"))
//}
