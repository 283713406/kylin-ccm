package service

import (
	"errors"
	"kylin-ccm/entity"
)

type clusterService struct{}

func (this *clusterService) table() string {
	return tableName("cluster")
}

// 根据集群名获取集群信息
func (this *clusterService) GetClusterByName(clusterName string) (*entity.SingleCluster, error) {
	cluster := &entity.SingleCluster{}
	cluster.ClusterName = clusterName
	err := o.Read(cluster, "UserName")
	return cluster, err
}

// 获取集群总数
func (this *clusterService) GetTotal() (int64, error) {
	return o.QueryTable(this.table()).Count()
}

// 添加集群
func (this *clusterService) AddCluster(clusterName, nodeName, k8sVersion string) error {
	if cluster, _ := this.GetClusterByName(clusterName); cluster.Id > 0 {
		return errors.New("集群已存在")
	}

	cluster := &entity.SingleCluster{}
	cluster.UserName = clusterName
	cluster.NodeName = nodeName
	cluster.K8sVersion = k8sVersion
	_, err := o.Insert(cluster)
	return err
}

// 更新集群信息
func (this *clusterService) UpdateCluster(user *entity.SingleCluster, fileds ...string) error {
	if len(fileds) < 1 {
		return errors.New("更新字段不能为空")
	}
	_, err := o.Update(user, fileds...)
	return err
}

// 删除集群
func (this *clusterService) DeleteCluster(clusterName string) error {
	cluster := &entity.SingleCluster{
		ClusterName: clusterName,
	}
	_, err := o.Delete(cluster)
	return err
}

