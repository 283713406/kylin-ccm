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

package manager

import (
	kubekeyapi "kylin-ccm/pkg/apis/kubekey/v1alpha1"
	"kylin-ccm/pkg/util/runner"
	"kylin-ccm/pkg/util/ssh"
)

type Manager struct {
	Cluster        *kubekeyapi.ClusterSpec
	Connector      *ssh.Dialer
	Runner         *runner.Runner
	AllNodes       []kubekeyapi.HostCfg
	EtcdNodes      []kubekeyapi.HostCfg
	MasterNodes    []kubekeyapi.HostCfg
	WorkerNodes    []kubekeyapi.HostCfg
	K8sNodes       []kubekeyapi.HostCfg
	ClientNode     []kubekeyapi.HostCfg
	ClusterHosts   []string
	WorkDir        string
	KsEnable       bool
	KsVersion      string
	Debug          bool
	SkipCheck      bool
	SkipPullImages bool
}

func (mgr *Manager) Copy() *Manager {
	newManager := *mgr
	return &newManager
}
