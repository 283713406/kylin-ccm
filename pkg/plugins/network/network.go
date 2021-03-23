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

package network

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	versionutil "k8s.io/apimachinery/pkg/util/version"
	kubekeyapi "kylin-ccm/pkg/apis/kubekey/v1alpha1"
	"kylin-ccm/pkg/plugins/network/calico"
	"kylin-ccm/pkg/plugins/network/flannel"
	"kylin-ccm/pkg/util"
	"kylin-ccm/pkg/util/logs"
	"kylin-ccm/pkg/util/manager"
)

func DeployNetworkPlugin(mgr *manager.Manager) error {
	logs.MyLogger.Infoln("Deploying network plugin ...")

	return mgr.RunTaskOnMasterNodes(deployNetworkPlugin, true)
}

func deployNetworkPlugin(mgr *manager.Manager, _ *kubekeyapi.HostCfg) error {
	if mgr.Runner.Index == 0 {
		switch mgr.Cluster.Network.Plugin {
		case "calico":
			if err := deployCalico(mgr); err != nil {
				return err
			}
		case "flannel":
			if err := deployFlannel(mgr); err != nil {
				return err
			}
		case "macvlan":
			if err := deployMacvlan(); err != nil {
				return err
			}
		default:
			return errors.New(fmt.Sprintf("This network plugin is not supported: %s", mgr.Cluster.Network.Plugin))
		}
	}
	return nil
}

func deployCalico(mgr *manager.Manager) error {
	if !util.IsExist(fmt.Sprintf("%s/network-plugin.yaml", mgr.WorkDir)) {
		var calicoContent string
		cmp, err := versionutil.MustParseSemantic(mgr.Cluster.Kubernetes.Version).Compare("v1.16.0")
		if err != nil {
			return err
		}
		if cmp == -1 {
			calicoContentStr, err := calico.GenerateCalicoFilesOld(mgr)
			if err != nil {
				return err
			}
			calicoContent = calicoContentStr
		} else {
			calicoContentStr, err := calico.GenerateCalicoFilesNew(mgr)
			if err != nil {
				return err
			}
			calicoContent = calicoContentStr
		}

		err1 := ioutil.WriteFile(fmt.Sprintf("%s/network-plugin.yaml", mgr.WorkDir), []byte(calicoContent), 0644)
		if err1 != nil {
			return errors.Wrap(errors.WithStack(err1), fmt.Sprintf("Failed to generate network plugin manifests: %s/network-plugin.yaml", mgr.WorkDir))
		}
	}

	calicoBase64, err1 := exec.Command("/bin/bash", "-c", fmt.Sprintf("tar cfz - -C %s -T /dev/stdin <<< network-plugin.yaml | base64 --wrap=0", mgr.WorkDir)).CombinedOutput()
	if err1 != nil {
		return errors.Wrap(errors.WithStack(err1), "Failed to read network plugin manifests")
	}

	_, err2 := mgr.Runner.ExecuteCmd(fmt.Sprintf("sudo -E /bin/bash -c \"base64 -d <<< '%s' | tar xz -C %s\"", strings.TrimSpace(string(calicoBase64)), "/etc/kubernetes"), 2, false)
	if err2 != nil {
		return errors.Wrap(errors.WithStack(err2), "Failed to generate network plugin manifests")
	}

	_, err3 := mgr.Runner.ExecuteCmd("sudo -E /bin/sh -c \"/usr/local/bin/kubectl apply -f /etc/kubernetes/network-plugin.yaml --force\"", 5, true)
	if err3 != nil {
		return errors.Wrap(errors.WithStack(err3), "Failed to deploy network plugin")
	}
	return nil
}

func deployFlannel(mgr *manager.Manager) error {
	if !util.IsExist(fmt.Sprintf("%s/network-plugin.yaml", mgr.WorkDir)) {
		flannelContent, err := flannel.GenerateFlannelFiles(mgr)
		if err != nil {
			return err
		}
		err1 := ioutil.WriteFile(fmt.Sprintf("%s/network-plugin.yaml", mgr.WorkDir), []byte(flannelContent), 0644)
		if err1 != nil {
			return errors.Wrap(errors.WithStack(err1), fmt.Sprintf("Failed to generate network plugin manifests: %s/network-plugin.yaml", mgr.WorkDir))
		}
	}

	flannelBase64, err1 := exec.Command("/bin/bash", "-c", fmt.Sprintf("tar cfz - -C %s -T /dev/stdin <<< network-plugin.yaml | base64 --wrap=0", mgr.WorkDir)).CombinedOutput()
	if err1 != nil {
		return errors.Wrap(errors.WithStack(err1), "Failed to read network plugin manifests")
	}

	_, err2 := mgr.Runner.ExecuteCmd(fmt.Sprintf("sudo -E /bin/bash -c \"base64 -d <<< '%s' | tar xz -C %s\"", strings.TrimSpace(string(flannelBase64)), "/etc/kubernetes"), 2, false)
	if err2 != nil {
		return errors.Wrap(errors.WithStack(err2), "Failed to generate network plugin manifests")
	}

	_, err3 := mgr.Runner.ExecuteCmd("sudo -E /bin/sh -c \"/usr/local/bin/kubectl apply -f /etc/kubernetes/network-plugin.yaml --force\"", 5, true)
	if err3 != nil {
		return errors.Wrap(errors.WithStack(err2), "Failed to deploy network plugin")
	}
	return nil
}

func deployMacvlan() error {
	return nil
}
