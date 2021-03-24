package controllers

import (
	"encoding/json"
	"kylin-ccm/controllers/base"
	"kylin-ccm/pkg/install"

	"kylin-ccm/entity"
	"kylin-ccm/pkg/util/logs"
	"kylin-ccm/service"
)

type ClusterController struct {
	base.BaseController
}

func (c *ClusterController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}

func (c *ClusterController) Post() {
	if c.IsPost() {
		logs.MyLogger.Infof("Start create cluster")
		var singleCluster entity.SingleCluster
		data := c.Ctx.Input.RequestBody
		err := json.Unmarshal(data, &singleCluster)
		if err != nil {
			logs.MyLogger.Infof( "json.Unmarshal is err:", err.Error())
		}

		isAllInOne := false
		ksEnable := false
		if singleCluster.IsAllInOne == "true" {
			isAllInOne = true
		}

		if singleCluster.KsEnable == "true" {
			ksEnable = true
		}

		if err = install.CreateCluster(singleCluster.NodeName, singleCluster.UserName, singleCluster.K8sVersion,
			singleCluster.KsVersion, isAllInOne, ksEnable); err != nil {
			logs.MyLogger.Errorf("Failed to create cluster, error: %v", err.Error())
			return
		}

		var cluster service.Cluster
		cluster.ClusterName = singleCluster.ClusterName
		cluster.Description = "this is test cluster"
		cluster.Status = "Running"
		cluster.User = "test"

		_, err = service.ClusterService.AddCluster(&cluster);
		if err != nil {
			logs.MyLogger.Errorf("Failed to add cluster, error: %v", err.Error())
			return
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

//func getSqlCluster(sc entity.SingleCluster) service.Cluster {
//	var cluster service.Cluster
//	cluster.ClusterName = sc.ClusterName
//	cluster.Description = "this is test cluster"
//	cluster.Status = "Running"
//	cluster.User = "test"
//
//	return cluster
//}