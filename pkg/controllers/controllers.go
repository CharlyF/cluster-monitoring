package controllers

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"os"
)

type ControllerContext struct {
	KubeClient          kubernetes.Interface
	StopCh              chan struct{}
}

var controllerCatalog = map[string]controllerFuncs{
	"cluster-metadata": {
		startMetadataController,
	},
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
		fmt.Printf("|| Starting %s ||\n", name)
		if err != nil {
			fmt.Errorf("Error starting %s: %s", name, err.Error())
			return err
		}
	}
	return nil
}

func startMetadataController(ctx ControllerContext) error {
	svcCtrl, err := NewSvcController(ctx.KubeClient)
	if err != nil {
		return err
	}
	go svcCtrl.Run(ctx.StopCh)
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
