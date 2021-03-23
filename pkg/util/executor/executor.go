/*
Copyright 2020 The KubeSphere Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package executor

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	kubekeyapi "kylin-ccm/pkg/apis/kubekey/v1alpha1"
	"kylin-ccm/pkg/util/logs"
	"kylin-ccm/pkg/util/manager"
	"kylin-ccm/pkg/util/ssh"
)

type Executor struct {
	Cluster        *kubekeyapi.ClusterSpec
	Debug          bool
	SkipCheck      bool
	SkipPullImages bool
}

func NewExecutor(cluster *kubekeyapi.ClusterSpec, debug, skipCheck, skipPullImages bool) *Executor {
	return &Executor{
		Cluster:        cluster,
		Debug:          debug,
		SkipCheck:      skipCheck,
		SkipPullImages: skipPullImages,
	}
}

func (executor *Executor) CreateManager() (*manager.Manager, error) {
	mgr := &manager.Manager{}
	defaultCluster, hostGroups := executor.Cluster.SetDefaultClusterSpec()
	mgr.AllNodes = hostGroups.All
	mgr.EtcdNodes = hostGroups.Etcd
	mgr.MasterNodes = hostGroups.Master
	mgr.WorkerNodes = hostGroups.Worker
	mgr.K8sNodes = hostGroups.K8s
	mgr.ClientNode = hostGroups.Client
	mgr.Cluster = defaultCluster
	mgr.ClusterHosts = GenerateHosts(hostGroups, defaultCluster)
	mgr.Connector = ssh.NewDialer()
	mgr.WorkDir = GenerateWorkDir()
	mgr.KsEnable = executor.Cluster.KubeSphere.Enabled
	mgr.KsVersion = executor.Cluster.KubeSphere.Version
	mgr.Debug = executor.Debug
	mgr.SkipCheck = executor.SkipCheck
	mgr.SkipPullImages = executor.SkipPullImages
	return mgr, nil
}

func GenerateHosts(hostGroups *kubekeyapi.HostGroups, cfg *kubekeyapi.ClusterSpec) []string {
	var lbHost string
	hostsList := []string{}

	if cfg.ControlPlaneEndpoint.Address != "" {
		lbHost = fmt.Sprintf("%s  %s", cfg.ControlPlaneEndpoint.Address, cfg.ControlPlaneEndpoint.Domain)
	} else {
		lbHost = fmt.Sprintf("%s  %s", hostGroups.Master[0].InternalAddress, cfg.ControlPlaneEndpoint.Domain)
	}

	for _, host := range cfg.Hosts {
		if host.Name != "" {
			hostsList = append(hostsList, fmt.Sprintf("%s  %s.%s %s", host.InternalAddress, host.Name, cfg.Kubernetes.ClusterName, host.Name))
		}
	}

	hostsList = append(hostsList, lbHost)
	return hostsList
}

func GenerateWorkDir() string {
	currentDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logs.MyLogger.Fatal(errors.Wrap(err, "Faild to get current dir"))
	}
	return fmt.Sprintf("%s/%s", currentDir, kubekeyapi.DefaultPreDir)
}
