package storage

import (
	kubekeyapi "kylin-ccm/pkg/apis/kubekey/v1alpha1"
	"kylin-ccm/pkg/cluster/preinstall"
	"kylin-ccm/pkg/images"
	"kylin-ccm/pkg/util/manager"
)

func prePullStorageImages(mgr *manager.Manager, node *kubekeyapi.HostCfg) error {
	i := images.Images{}
	i.Images = []images.Image{
		preinstall.GetImage(mgr, "provisioner-localpv"),
		preinstall.GetImage(mgr, "node-disk-manager"),
		preinstall.GetImage(mgr, "node-disk-operator"),
		preinstall.GetImage(mgr, "linux-utils"),
		preinstall.GetImage(mgr, "rbd-provisioner"),
		preinstall.GetImage(mgr, "nfs-client-provisioner"),
	}
	if err := i.PullImages(mgr, node); err != nil {
		return err
	}
	return nil
}
