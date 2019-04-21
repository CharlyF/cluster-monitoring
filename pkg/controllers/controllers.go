package controllers

import (
	"fmt"
	"github.com/DataDog/datadog-agent/pkg/util/docker"
	"github.com/DataDog/datadog-agent/pkg/util/kubernetes/kubelet"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"os"
)

type ControllerContext struct {
	InformerFactory informers.SharedInformerFactory
	KubeClient          kubernetes.Interface
	DockerClient        *docker.DockerUtil
	KubeletClient       *kubelet.KubeUtil
	StopCh              chan struct{}
}

var controllerCatalog = map[string]controllerFuncs{
	//"cluster-metadata": {
	//	startMetadataController,
	//},
	"node-metadata": {
		start: nodeMetadataController,
	},
	"network-ingest": {
		start: startSocketPolling,
	},
	"aggregator": {
	start: startAggregator,
	},

}
// need cache for pods with containers and containers from outside k8s.

type controllerFuncs struct {
	start   func(ControllerContext) error
}

func Start(ctx ControllerContext) error {
	for name, cntrlFuncs := range controllerCatalog {
		err := cntrlFuncs.start(ctx)
		fmt.Printf("starting %s \n", name)
		if err != nil {
			fmt.Errorf("Error starting %s: %s", name, err.Error())
			return err
		}
	}
	ctx.InformerFactory.Start(ctx.StopCh)
	return nil
}

func startMetadataController(ctx ControllerContext) error {
	podCtrl, err := NewPodController(ctx.KubeClient, ctx.InformerFactory.Core().V1().Pods())
	if err != nil {
		return err
	}
	go podCtrl.Run(ctx.StopCh)
	return nil
}

func startSocketPolling(ctx ControllerContext) error {
   path := os.Getenv("SOCKET_PATH")
   fmt.Printf("socket used is %s \n", path)
   skt, err  := NewSocketController(path)
   if err != nil {
   	return err
   }
   go skt.Run(ctx.StopCh)
   return nil
}

func nodeMetadataController(ctx ControllerContext) error {
	node, err := NewNodeController(ctx)
	if err != nil {
		return err
	}
	go node.Run(ctx.StopCh)
	return nil
}

func startAggregator(ctx ControllerContext) error {
	// should reconcile the data from the network-ingest and the metadata from the metadata controller
	a, err := NewAggregator(ctx.KubeClient)
	if err != nil {
		return err
	}
	go a.Run(ctx.StopCh)

	return nil
}
