package entity

type SingleCluster struct {
	Id            int
	ClusterName   string     `orm:"unique;size(32)"`             // 集群名
	NodeName      string     `orm:"size(32)"`                    // 节点名
	UserName      string     `orm:"size(32)"`                    // 用户名
	K8sVersion    string     `orm:"size(32)"`                    // k8s版本号
	KsVersion     string     `orm:"size(32)"`                    // kubesphere版本号
	IsAllInOne    string     `orm:"size(32)"`                    // kubesphere版本号
	KsEnable      string     `orm:"size(32)"`                    // 是否安装kubesphere
}
