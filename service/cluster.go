package service

import (
	"errors"
	"time"
)

type clusterService struct{}

type Cluster struct {
	Id            int64      `orm:"auto" json:"id,omitempty"`
	ClusterName   string     `orm:"unique;index;size(128)" json:"name,omitempty"`
	Description   string     `orm:"null;size(512)" json:"description,omitempty"`
	CreateTime    *time.Time `orm:"auto_now_add;type(datetime)" json:"createTime,omitempty"`
	UpdateTime    *time.Time `orm:"auto_now;type(datetime)" json:"updateTime,omitempty"`
	User          string     `orm:"size(128)" json:"user,omitempty"`
	Status        string     `orm:"size(128)" json:"status"`
}

func (this *clusterService) table() string {
	return tableName("cluster")
}

// 根据集群名获取集群信息
func (this *clusterService) GetClusterByName(clusterName string) (*Cluster, error) {
	cluster := &Cluster{}
	cluster.ClusterName = clusterName
	err := o.Read(cluster, "ClusterName")
	return cluster, err
}

// 获取集群总数
func (this *clusterService) GetTotal() (int64, error) {
	return o.QueryTable(this.table()).Count()
}

// 添加集群
func (this *clusterService) AddCluster(c *Cluster) (id int64, err error) {
	if cluster, _ := this.GetClusterByName(c.ClusterName); cluster.Id > 0 {
		return 0, errors.New("集群已存在")
	}

	c.CreateTime = nil
	id, err = o.Insert(c)
	return id, nil
}

// 更新集群信息
func (this *clusterService) UpdateCluster(c *Cluster, fileds ...string) error {
	if len(fileds) < 1 {
		return errors.New("更新字段不能为空")
	}
	_, err := o.Update(c, fileds...)
	return err
}

// 删除集群
func (this *clusterService) DeleteCluster(clusterName string) error {
	cluster := &Cluster{
		ClusterName: clusterName,
	}
	_, err := o.Delete(cluster)
	return err
}

