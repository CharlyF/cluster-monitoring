package app

import (
	"fmt"
	"github.com/CharlyF/cluster-monitoring/pkg/controllers"
	"github.com/CharlyF/cluster-monitoring/util"
	"github.com/DataDog/datadog-agent/pkg/util/docker"
	"github.com/DataDog/datadog-agent/pkg/util/kubernetes/kubelet"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var (

	MetadataCollector = &cobra.Command{
		Use: "meta [command]",
	}

	startCmd = &cobra.Command{
		Use: "start",
		RunE: start,
	}

	dataCmd = &cobra.Command{
		Use: "data",
		RunE: data,
	}
)

func init() {
	MetadataCollector.AddCommand(startCmd)
	MetadataCollector.AddCommand(dataCmd)
}

func start(cmd *cobra.Command, args []string) error {
	// Setup a channel to catch OS signals
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	apiCl, err := util.GetKubeClient()
	if err != nil {
		return err
	}
	dClient, err := docker.GetDockerUtil()
	if err != nil {
		return err
	}
	kClient, err := kubelet.GetKubeUtil()
	if err != nil {
		return err
	}
	factory := informers.NewSharedInformerFactory(apiCl, 0)
	stopCh := make(chan struct{})
	ctx := controllers.ControllerContext{
		KubeClient: apiCl,
		DockerClient: dClient,
		KubeletClient: kClient,
		StopCh : stopCh,
		InformerFactory:factory,
	}
	controllers.Start(ctx)

	// Block here until we receive the interrupt signal
	fmt.Printf("We are blocked here waitign for signal \n")
	<-signalCh
	fmt.Printf("We received a signal \n")
	stopCh <- struct {}{}
	fmt.Printf("We are sending a stop to the controller \n")

	return nil
}

func data(cmd *cobra.Command, args []string) error {
	apiCl, err := util.GetKubeClient()
	if err != nil {
		return err
	}
	cm, err := apiCl.CoreV1().ConfigMaps("default").Get("transactions", metav1.GetOptions{})
	//fmt.Printf("transactions: %v",cm.Data)
	mapper := make(map[string]string)
	for key, value := range cm.Data {
		ips := strings.Split(key, "-")
		if len(ips) == 4 {
			continue
		} else {
			// this is a metadata key
			mapper[key] = value
		}
	}

	for key, value := range cm.Data {
		ips := strings.Split(key, "-")
		if len(ips) == 4 {
			// this is transaction key.
			src := ips[0]
			dest := ips[1]
			fmt.Printf("Transaction between %s [%v] -> %s [%v]\n---\n%v\n---\n\n", src, mapper[src], dest,mapper[dest], value)
		} else {
			continue
			mapper[key] = value
			// this is a metadata key
			fmt.Printf("IP: %s has the following metadata: %v \n", key, value)
		}
	}
	return nil
}
