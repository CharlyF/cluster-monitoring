package controllers

import (
	"fmt"
	"github.com/CharlyF/cluster-monitoring/util"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
	"k8s.io/apimachinery/pkg/api/errors"
)

const (
	storeName = "transactions"
	ns = "default"
)
type AggregatorController struct {
	clientSet kubernetes.Interface

}

func NewAggregator(client kubernetes.Interface) (*AggregatorController, error) {
	_, err := client.CoreV1().ConfigMaps(ns).Get(storeName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		cm := &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      storeName,
				Namespace: ns,
			},
		}
		_, err := client.CoreV1().ConfigMaps(ns).Create(cm)
		if err != nil {
			return nil, err
		}
	}
	return &AggregatorController{
		clientSet: client,
	}, nil
}

func (a *AggregatorController) Run(stop chan struct{}) {
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				item := util.Cache.Items()
				data := make(map[string]string)
				for i, j := range item {
					val := j.Object.(string)
					data[i] = val
				}
				updatedCM := v1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      storeName,
						Namespace: ns,
					},
					Data: data,
				}

				_, err := a.clientSet.CoreV1().ConfigMaps(ns).Update(&updatedCM)
				if err != nil {
					fmt.Printf(err.Error())
					return
				}
				fmt.Printf("Saved to the ConfigMap: %v \n", data)
			case <-stop:
				fmt.Errorf("returning from CM storing \n")
				return
			}
		}
	}()
}
