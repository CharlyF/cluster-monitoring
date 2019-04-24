package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/CharlyF/cluster-monitoring/pkg/aggregator"
	"github.com/CharlyF/cluster-monitoring/util"
	"github.com/ericchiang/k8s"
	"k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"os"
	"time"
)

type SvcController struct {
	clientSet kubernetes.Interface
}

func NewSvcController(client kubernetes.Interface) (*SvcController, error) {
	return &SvcController{
		clientSet:client,
	}, nil
}

func (s *SvcController) Run(stop <- chan struct{}) {
	fmt.Printf("Starting runner for Svc Controller \n")
	timer := time.NewTicker(time.Second * 5)
	// rename svc controller into something more cluster agnostic
	node, err := s.clientSet.CoreV1().Nodes().Get(os.Getenv("HOSTNAME"), v12.GetOptions{})
	if err != nil {
		fmt.Errorf("could not get the metadata of the host")
	}
	for _, v := range node.Status.Addresses {
		if v.Type == "InternalIP" {
			n := &aggregator.Data{
				Node:aggregator.NodeMeta{
					Name: node.Name,
					KernelVersion: node.Status.NodeInfo.KernelVersion,
				},
			}
			nByte, err := json.Marshal(n)
			if err != nil {
				fmt.Errorf("Could not marshal %s", v.Address)
				continue
			}
			util.Cache.Add(v.Address, string(nByte), 0)
			break
		}
	}
	for {
		select {
		case <-timer.C:
			//e := n.getPodsMeta()
			//if e != nil {
			//	continue
			//}
			svc, err := s.clientSet.CoreV1().Services(k8s.AllNamespaces).List(v12.ListOptions{})
			if err != nil {
				fmt.Errorf("error retrieving svc %v \n", err.Error())
				continue
			}
			err = s.svcCaching(svc)
			if err != nil {
				fmt.Errorf("error caching %v \n", err.Error())
				continue
			}
		case <- stop:
			fmt.Printf("Stoppring service controller \n")
			return
		}
	}
}

func addToCache(svc v1.Service) error{
	sm := &aggregator.Data{
		Service: aggregator.SvcMetadata{
			Name: svc.Name,
			Type:string(svc.Spec.Type),
			Ports: fmt.Sprintf("%d/%s", svc.Spec.Ports[0].Port,string(svc.Spec.Ports[0].Protocol)),
		},
	}
	smJson, err := json.Marshal(sm)
	if err != nil {
		return err
	}
	return util.Cache.Add(svc.Spec.ClusterIP, string(smJson), 0)
}

func (p *SvcController) svcCaching(listSvc *v1.ServiceList) error {

	for _, svc := range listSvc.Items {
		//fmt.Printf("client-go cache svc: \n- Name: %s \n- IP: %s", svc.Name, svc.Spec.ClusterIP)
		_, cached := util.Cache.Get(svc.Spec.ClusterIP)
		if !cached {
			err := addToCache(svc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
