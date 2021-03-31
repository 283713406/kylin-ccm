创建集群接口
http://172.20.188.109:8080/api/v1/clusters

例子：
curl -H "Content-Type:application/json" -X POST   -d '{"id":1, "clustername":"test","nodename":"kubesphere","username":"root","k8sversion":"v1.17.9","ksversion":"v3.0.0","IsAllInOne":"true","ksenable":"false"}' http://172.20.188.109:8080/api/v1/clusters


删除集群接口
http://172.20.188.109:8080/api/v1/clusters/test

例子：
curl -H "Content-Type:application/json" -X DELETE  http://172.20.188.109:8080/api/v1/clusters/test