package controllers

import (
	"github.com/DataDog/datadog-agent/pkg/util/docker"
	"github.com/DataDog/datadog-agent/pkg/util/kubernetes/kubelet"
)

type NodeController struct {
	kClient *kubelet.KubeUtil
	dClient *docker.DockerUtil
}

func NewNodeController(ctx ControllerContext) (*NodeController, error){
	return &NodeController{
		kClient: ctx.KubeletClient,
		dClient: ctx.DockerClient,
	}, nil
}

func (n *NodeController) Run(ctx ControllerContext) error {
	// poll the list of pods every 10s
	//
	return nil
}
