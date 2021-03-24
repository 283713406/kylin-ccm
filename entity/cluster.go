package entity

type SingleCluster struct {
	ClusterName   string                 // 集群名
	NodeName      string                 // 节点名
	UserName      string                 // 用户名
	K8sVersion    string                 // k8s版本号
	KsVersion     string                 // kubesphere版本号
	IsAllInOne    string                 // kubesphere版本号
	KsEnable      string                 // 是否安装kubesphere
}
