package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/CharlyF/cluster-monitoring/pkg/aggregator"
	"github.com/CharlyF/cluster-monitoring/util"
	"github.com/ericchiang/k8s"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"os"
	"time"
)

type NodeController struct {
	kClient kubernetes.Interface
}

func NewNodeController(ctx ControllerContext) (*NodeController, error){
	return &NodeController{
		kClient: ctx.KubeClient,
	}, nil
}

func (n *NodeController) Run(stop <- chan struct{}) {
	// poll the list of pods every 10s

	timer := time.NewTicker(time.Second * 5)
	for {
		select {
		case <- timer.C:
			e := n.getPodsMeta()
			if e != nil {
				continue
			}
		case <- stop:
			fmt.Printf("Stopping metadata controller \n")
			return
		}
	}
}

func (n *NodeController) getPodsMeta() error {

	hostname := os.Getenv("HOSTNAME")
	fmt.Printf("Using %s \n", hostname)
	pods, err := n.kClient.CoreV1().Pods(k8s.AllNamespaces).List(metav1.ListOptions{FieldSelector: fmt.Sprintf("spec.nodeName=%s", hostname)})

	if err != nil {
		return err
	}

	for _, p := range pods.Items {
		d := &aggregator.Data{
			Kubernetes:aggregator.PodMetadata{
				Name: p.Name,
				UID: string(p.UID),
				Image: p.Spec.Containers[0].Image,
			},
		}
		fmt.Printf("Processing pod %s \n", p.Name)
		podJson, err  := json.Marshal(d)
		if err != nil {
			fmt.Errorf("error marshalling pod %s", err.Error())
			continue
		}
		util.Cache.Add(p.Status.PodIP, string(podJson), 0)
	}
	return nil
}

//

// let's check if two containers in the same pod have the same ip.
// It's actually a third container pause that has the ip
//
//func (n *NodeController) getDockerCntr() error {
//	ctn, err := n.dClient.ListContainers(&docker.ContainerListConfig{
//		IncludeExited: false,
//		FlagExcluded:  false,
//	})
//	if err != nil {
//		return  err
//	}
//	mapPodCont := make(map[string]aggregator.ContainerMeta)
//
//	for _, c := range ctn {
//		cInspected, err := n.dClient.Inspect(c.ID, false)
//		if err != nil {
//			return err
//		}
//		ip := cInspected.NetworkSettings.IPAddress
//		id := c.ID
//		im := c.Image
//		m := aggregator.ContainerMeta{
//			Image: im,
//			Id: id,
//		}
//
//		mapPodCont[ip] = m
//		po, err := n.kClient.GetPodForContainerID(m.Id)
//		if err != nil {
//			// container is not in a pod
//			continue
//		}
//		p := aggregator.PodMetadata{
//			Name: po.Metadata.Name,
//			UID: po.Metadata.UID,
//		}
//		obj := aggregator.Data{
//			Docker: mapPodCont,
//			Kubernetes: p,
//		}
//	}
//
//}
