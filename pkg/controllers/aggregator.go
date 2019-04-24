package controllers

import (
	"fmt"
	"github.com/CharlyF/cluster-monitoring/util"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"time"
)

const (
	transactionCM = "transactions"
	metadataCM    = "metadata"
	ns            = "default"
)
type AggregatorController struct {
	clientSet kubernetes.Interface
}

func NewAggregator(client kubernetes.Interface) (*AggregatorController, error) {
	_, err := client.CoreV1().ConfigMaps(ns).Get(transactionCM, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		cm := &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      transactionCM,
				Namespace: ns,
			},
			Data : make(map[string]string),
		}
		_, err := client.CoreV1().ConfigMaps(ns).Create(cm)
		if err != nil {
			return nil, err
		}
	}
	_, err = client.CoreV1().ConfigMaps(ns).Get(metadataCM, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		cm := &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      metadataCM,
				Namespace: ns,
			},
			Data : make(map[string]string),
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
	ticker := time.NewTicker(10 * time.Second)
	// ticker to update metadata cm
	go func() {
		for {
			select {
			case <-ticker.C:
				item := util.Cache.Items()
				fmt.Printf("retrieving %d items from the cache \n", len(item))

				cm, err := a.clientSet.CoreV1().ConfigMaps(ns).Get(metadataCM, metav1.GetOptions{})
				if err != nil {
					fmt.Errorf("error will pushing metadata to the store: %s \n", err.Error())
					continue
				}
				data := make(map[string]string)
				if cm.Data != nil {
					data = cm.Data
				}
				for i, j := range item {
					val := j.Object.(string)
					data[i] = val
				}
				updatedCM := v1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      metadataCM,
						Namespace: ns,
					},
					Data: data,
				}

				_, err = a.clientSet.CoreV1().ConfigMaps(ns).Update(&updatedCM)
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

// used to send the transactions to the store.
func Save(m map[string]string) error {
	k, er := util.GetKubeClient()
	if er != nil {
		return er
	}
	cm, err := k.CoreV1().ConfigMaps(ns).Get(transactionCM, metav1.GetOptions{})
	if err != nil {
		return err
	}
	//fmt.Printf("Retrieved %d from the transaction CM \n", len(cm.Data))
	// That would add data to the cm indefinitely
	if cm.Data == nil {
		cm.Data = make(map[string]string)
	}
	for k, v := range m {
		cm.Data[k] = v
	}

	//fmt.Printf("Saving %d to the transaction CM \n", len(cm.Data))
	_, err = k.CoreV1().ConfigMaps(ns).Update(cm)
	return err
}
