package preinstall

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/modood/table"
	kubekeyapi "kylin-ccm/pkg/apis/kubekey/v1alpha1"
	log "kylin-ccm/pkg/util/logs"
	"kylin-ccm/pkg/util/manager"
	"os"
	"strings"
)

type PrecheckResults struct {
	Name      string `table:"name"`
	Sudo      string `table:"sudo"`
	Curl      string `table:"curl"`
	Openssl   string `table:"openssl"`
	Ebtables  string `table:"ebtables"`
	Socat     string `table:"socat"`
	Ipset     string `table:"ipset"`
	Conntrack string `table:"conntrack"`
	Docker    string `table:"docker"`
	Nfs       string `table:"nfs client"`
	Ceph      string `table:"ceph client"`
	Glusterfs string `table:"glusterfs client"`
	Time      string `table:"time"`
}

var (
	CheckResults  = make(map[string]interface{})
	BaseSoftwares = []string{"sudo", "curl", "openssl", "ebtables", "socat", "ipset", "conntrack", "docker", "showmount", "rbd", "glusterfs"}
)

func Precheck(mgr *manager.Manager) error {
	log.MyLogger.Infof("Start precheck before create cluster")
	if !mgr.SkipCheck {
		if err := mgr.RunTaskOnAllNodes(PrecheckNodes, true); err != nil {
			return err
		}
		// PrecheckConfirm(mgr)
	}
	return nil
}

func PrecheckNodes(mgr *manager.Manager, node *kubekeyapi.HostCfg) error {
	log.MyLogger.Infof("Start precheck node")
	var results = make(map[string]interface{})
	results["name"] = node.Name
	for _, software := range BaseSoftwares {
		log.MyLogger.Infof("check software: %s", software)
		_, err := mgr.Runner.ExecuteCmd(fmt.Sprintf("sudo -E /bin/sh -c \"which %s\"", software), 0, false)
		switch software {
		case "showmount":
			software = "nfs"
		case "rbd":
			software = "ceph"
		case "glusterfs":
			software = "glusterfs"
		}
		if err != nil {
			results[software] = ""
		} else {
			results[software] = "y"
		}
	}
	output, err := mgr.Runner.ExecuteCmd("date +\"%Z %H:%M:%S\"", 0, false)
	log.MyLogger.Infof("output is: %s", output)
	if err != nil {
		results["time"] = ""
	} else {
		results["time"] = strings.TrimSpace(output)
	}

	CheckResults[node.Name] = results
	return nil
}

func PrecheckConfirm(mgr *manager.Manager) {

	var results []PrecheckResults
	for node := range CheckResults {
		var result PrecheckResults
		_ = mapstructure.Decode(CheckResults[node], &result)
		results = append(results, result)
	}
	table.OutputA(results)
	// reader := bufio.NewReader(os.Stdin)
	fmt.Println("")
	fmt.Println("This is a simple check of your environment.")
	fmt.Println("Before installation, you should ensure that your machines meet all requirements specified at")
	fmt.Println("https://kylin-ccm#requirements-and-recommendations")
	fmt.Println("")
Loop:
	for {
		fmt.Printf("Continue this installation? [yes/no]: ")
		/*input, err := reader.ReadString('\n')
		if err != nil {
			mgr.Logger.Fatal(err)
		}*/
		input := "yes"
		input = strings.TrimSpace(input)

		switch input {
		case "yes":
			break Loop
		case "no":
			os.Exit(0)
		default:
			continue
		}
	}
}
