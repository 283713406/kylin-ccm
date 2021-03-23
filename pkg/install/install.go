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

package install

import (
	"fmt"
	"kylin-ccm/pkg/util/logs"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"kylin-ccm/pkg/cluster/etcd"
	"kylin-ccm/pkg/cluster/kubernetes"
	"kylin-ccm/pkg/cluster/preinstall"
	"kylin-ccm/pkg/config"
	"kylin-ccm/pkg/container-engine/docker"
	"kylin-ccm/pkg/kubesphere"
	"kylin-ccm/pkg/plugins/network"
	"kylin-ccm/pkg/util"
	"kylin-ccm/pkg/util/executor"
	"kylin-ccm/pkg/util/manager"
)

func CreateCluster(nodeName, userName, k8sVersion, ksVersion string, isAllInOne, ksEnabled bool) error {
	currentDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return errors.Wrap(err, "Faild to get current dir")
	}
	if err := util.CreateDir(fmt.Sprintf("%s/kubekey", currentDir)); err != nil {
		return errors.Wrap(err, "Failed to create work dir")
	}

	cfg, err := config.ParseClusterCfg(nodeName, userName, k8sVersion, ksVersion, isAllInOne, ksEnabled)
	if err != nil {
		return errors.Wrap(err, "Failed to download cluster config")
	}

	return Execute(executor.NewExecutor(&cfg.Spec, true, false, false))
}

func ExecTasks(mgr *manager.Manager) error {
	createTasks := []manager.Task{
		{Task: preinstall.Precheck, ErrMsg: "Failed to precheck"},
		{Task: preinstall.DownloadBinaries, ErrMsg: "Failed to download kube binaries"},
		{Task: preinstall.InitOS, ErrMsg: "Failed to init OS"},
		{Task: docker.InstallerDocker, ErrMsg: "Failed to install docker"},
		{Task: preinstall.PrePullImages, ErrMsg: "Failed to pre-pull images"},
		{Task: etcd.GenerateEtcdCerts, ErrMsg: "Failed to generate etcd certs"},
		{Task: etcd.SyncEtcdCertsToMaster, ErrMsg: "Failed to sync etcd certs"},
		{Task: etcd.GenerateEtcdService, ErrMsg: "Failed to create etcd service"},
		{Task: etcd.SetupEtcdCluster, ErrMsg: "Failed to start etcd cluster"},
		{Task: etcd.RefreshEtcdConfig, ErrMsg: "Failed to refresh etcd configuration"},
		{Task: kubernetes.GetClusterStatus, ErrMsg: "Failed to get cluster status"},
		{Task: kubernetes.InstallKubeBinaries, ErrMsg: "Failed to install kube binaries"},
		{Task: kubernetes.InitKubernetesCluster, ErrMsg: "Failed to init kubernetes cluster"},
		{Task: network.DeployNetworkPlugin, ErrMsg: "Failed to deploy network plugin"},
		{Task: kubernetes.JoinNodesToCluster, ErrMsg: "Failed to join node"},
		{Task: kubesphere.DeployLocalVolume, ErrMsg: "Failed to deploy localVolume"},
		{Task: kubesphere.DeployKubeSphere, ErrMsg: "Failed to deploy kubesphere"},
	}

	for _, step := range createTasks {
		if err := step.Run(mgr); err != nil {
			return errors.Wrap(err, step.ErrMsg)
		}
	}
	if mgr.KsEnable {
		logs.MyLogger.Infoln(`Installation is complete.

Please check the result using the command:

       kubectl logs -n kubesphere-system $(kubectl get pod -n kubesphere-system -l app=ks-install -o jsonpath='{.items[0].metadata.name}') -f

`)
	} else {
		logs.MyLogger.Infoln("Congradulations! Installation is successful.")
	}

	return nil
}

func Execute(executor *executor.Executor) error {
	mgr, err := executor.CreateManager()
	if err != nil {
		return err
	}
	return ExecTasks(mgr)
}
