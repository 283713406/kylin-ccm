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

package config

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
	kubekeyapi "kylin-ccm/pkg/apis/kubekey/v1alpha1"
	"kylin-ccm/pkg/kubesphere"
	"kylin-ccm/pkg/util"
	"kylin-ccm/pkg/util/logs"
)

func ParseClusterCfg(nodeName, userName, k8sVersion, ksVersion string, isAllInOne, ksEnabled bool) (*kubekeyapi.Cluster, error) {
	var clusterCfg *kubekeyapi.Cluster
	if isAllInOne {
		if userName != "root" {
			return nil, errors.New(fmt.Sprintf("Current user is %s. Please use root!", userName))
		}
		clusterCfg = AllinoneCfg(nodeName, userName, k8sVersion, ksVersion, ksEnabled)
	} else {
		cfg, err := ParseCfg("clusterCfgPath", k8sVersion, ksVersion, ksEnabled)
		if err != nil {
			return nil, err
		}
		clusterCfg = cfg
	}

	return clusterCfg, nil
}

func ParseCfg(clusterCfgPath, k8sVersion, ksVersion string, ksEnabled bool) (*kubekeyapi.Cluster, error) {
	clusterCfg := kubekeyapi.Cluster{}
	fp, err := filepath.Abs(clusterCfgPath)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to look up current directory")
	}
	if len(k8sVersion) != 0 {
		_ = exec.Command("/bin/sh", "-c", fmt.Sprintf("sed -i \"/version/s/\\:.*/\\: %s/g\" %s", k8sVersion, fp)).Run()
	}
	file, err := os.Open(fp)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to open the given cluster configuration file")
	}
	defer file.Close()
	b1 := bufio.NewReader(file)
	for {
		result := make(map[string]interface{})
		content, err := k8syaml.NewYAMLReader(b1).Read()
		if len(content) == 0 {
			break
		}
		if err != nil {
			return nil, errors.Wrap(err, "Unable to read the given cluster configuration file")
		}
		err = yaml.Unmarshal(content, &result)
		if err != nil {
			return nil, errors.Wrap(err, "Unable to unmarshal the given cluster configuration file")
		}
		if result["kind"] == "Cluster" {
			if err := yaml.Unmarshal(content, &clusterCfg); err != nil {
				return nil, errors.Wrap(err, "Unable to convert file to yaml")
			}
		}

		if result["kind"] == "ConfigMap" || result["kind"] == "ClusterConfiguration" {
			metadata := result["metadata"].(map[interface{}]interface{})
			labels := metadata["labels"].(map[interface{}]interface{})
			clusterCfg.Spec.KubeSphere.Enabled = true
			_, ok := labels["version"]
			if ok {
				switch labels["version"] {
				case "v3.0.0":
					clusterCfg.Spec.KubeSphere.Configurations = "---\n" + string(content)
					clusterCfg.Spec.KubeSphere.Version = "v3.0.0"
				case "v2.1.1":
					clusterCfg.Spec.KubeSphere.Configurations = "---\n" + string(content)
					clusterCfg.Spec.KubeSphere.Version = "v2.1.1"
				default:
					return nil, errors.Wrap(err, fmt.Sprintf("Unsupported versions: %s", labels["version"]))
				}
			}
		}
	}

	if ksEnabled {
		clusterCfg.Spec.KubeSphere.Enabled = true
		switch strings.TrimSpace(ksVersion) {
		case "":
			clusterCfg.Spec.KubeSphere.Version = "v3.0.0"
			clusterCfg.Spec.KubeSphere.Configurations = kubesphere.V3_0_0
		case "v3.0.0":
			clusterCfg.Spec.KubeSphere.Version = "v3.0.0"
			clusterCfg.Spec.KubeSphere.Configurations = kubesphere.V3_0_0
		case "v2.1.1":
			clusterCfg.Spec.KubeSphere.Version = "v2.1.1"
			clusterCfg.Spec.KubeSphere.Configurations = kubesphere.V2_1_1
		default:
			return nil, errors.New(fmt.Sprintf("Unsupported version: %s", strings.TrimSpace(ksVersion)))
		}
	}

	return &clusterCfg, nil
}

func AllinoneCfg(nodeName, userName, k8sVersion, ksVersion string, ksEnabled bool) *kubekeyapi.Cluster {
	allinoneCfg := kubekeyapi.Cluster{}
	if output, err := exec.Command("/bin/sh", "-c", "if [ ! -f \"$HOME/.ssh/id_rsa\" ]; then ssh-keygen -t rsa -P \"\" -f $HOME/.ssh/id_rsa && ls $HOME/.ssh;fi;").CombinedOutput(); err != nil {
		logs.MyLogger.Fatalf("Failed to generate public key: %v\n%s", err, string(output))
	}
	if output, err := exec.Command("/bin/sh", "-c", "echo \"\n$(cat $HOME/.ssh/id_rsa.pub)\" >> $HOME/.ssh/authorized_keys && awk ' !x[$0]++{print > \"'$HOME'/.ssh/authorized_keys\"}' $HOME/.ssh/authorized_keys").CombinedOutput(); err != nil {
		logs.MyLogger.Fatalf("Failed to copy public key to authorized_keys: %v\n%s", err, string(output))
	}

	allinoneCfg.Spec.Hosts = append(allinoneCfg.Spec.Hosts, kubekeyapi.HostCfg{
		Name:            nodeName,
		Address:         util.LocalIP(),
		InternalAddress: util.LocalIP(),
		Port:            "",
		User:            userName,
		Password:        "",
		PrivateKeyPath:  "/root/.ssh/id_rsa",
		Arch:            runtime.GOARCH,
	})

	allinoneCfg.Spec.RoleGroups = kubekeyapi.RoleGroups{
		Etcd:   []string{nodeName},
		Master: []string{nodeName},
		Worker: []string{nodeName},
	}
	if k8sVersion != "" {
		allinoneCfg.Spec.Kubernetes = kubekeyapi.Kubernetes{
			Version: k8sVersion,
		}
	} else {
		allinoneCfg.Spec.Kubernetes = kubekeyapi.Kubernetes{
			Version: kubekeyapi.DefaultKubeVersion,
		}
	}

	if ksEnabled {
		allinoneCfg.Spec.KubeSphere.Enabled = true
		switch strings.TrimSpace(ksVersion) {
		case "":
			allinoneCfg.Spec.KubeSphere.Version = "v3.0.0"
			allinoneCfg.Spec.KubeSphere.Configurations = kubesphere.V3_0_0
		case "v3.0.0":
			allinoneCfg.Spec.KubeSphere.Version = "v3.0.0"
			allinoneCfg.Spec.KubeSphere.Configurations = kubesphere.V3_0_0
		case "v2.1.1":
			allinoneCfg.Spec.KubeSphere.Version = "v2.1.1"
			allinoneCfg.Spec.KubeSphere.Configurations = kubesphere.V2_1_1
		default:
			logs.MyLogger.Error("Unsupported version: %s", strings.TrimSpace(ksVersion))
		}
	}

	return &allinoneCfg
}
