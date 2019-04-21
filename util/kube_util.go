package util

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

func GetKubeClient() (kubernetes.Interface, error) {
	var err error

	kubeconfig := os.Getenv("KUBECONFIG")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset, nil

	var clientConfig *rest.Config
	clientConfig, err = rest.InClusterConfig()

	if err != nil {
		return nil, err
	}
	clientConfig.ContentType = "application/vnd.kubernetes.protobuf"

	return kubernetes.NewForConfig(clientConfig)
}
