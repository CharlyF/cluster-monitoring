package app

import (
	"encoding/json"
	"fmt"
	"github.com/CharlyF/cluster-monitoring/pkg/aggregator"
	"github.com/CharlyF/cluster-monitoring/pkg/controllers"
	"github.com/CharlyF/cluster-monitoring/pkg/render"
	"github.com/CharlyF/cluster-monitoring/util"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	//dClient, err := docker.GetDockerUtil()
	//if err != nil {
	//	return err
	//}

	stopCh := make(chan struct{})
	ctx := controllers.ControllerContext{
		KubeClient: apiCl,
		//DockerClient: *dClient,
		StopCh : stopCh,
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
	tr := []aggregator.Transactions{}
	mapper := make(map[string]aggregator.Data)
	// Start with the metadata
	cmMeta, err := apiCl.CoreV1().ConfigMaps("default").Get("metadata", metav1.GetOptions{})
	if err != nil {
		return err
	}
	//fmt.Printf("From Metadata CM %v \n", cmMeta.Data)
	for k, v := range cmMeta.Data {
		d := &aggregator.Data{}
		err := json.Unmarshal([]byte(v), d)
		if err != nil {
			fmt.Errorf("Could not unmarshall vals for %s", k)
			continue
		}
		mapper[k] = *d
	}
	//fmt.Printf("Mapper is %v \n", mapper)
	// Fetch the connections
	cm, err := apiCl.CoreV1().ConfigMaps("default").Get("transactions", metav1.GetOptions{})
	//fmt.Printf("From Transactions CM %v \n", cm.Data)

	for key, value := range cm.Data {
		v := &aggregator.Values{}
		ips := strings.Split(key, "-")
		if len(ips) == 4 {
			// this is transaction key.
			err := json.Unmarshal([]byte(value), v)
			if err != nil {
				fmt.Errorf("could not unmarshall values")
				continue
			}
			src := ips[0]
			dest := ips[1]
			t := aggregator.Transactions{
				IpSrc: src,
				IpDest: dest,
				DataSrc: mapper[src],
				DataDest: mapper[dest],
				Val: *v,
			}
			tr = append(tr, t)
		}
	}
	//fmt.Printf("Formatted %#v \n", tr)
	stats := make(map[string]interface{})
	stats["transactions"] = tr
	byTr, err := json.Marshal(stats)

	fmt.Errorf("Error marshalling list of transactions")
	formattedStatus, err := render.FormatData(byTr)
	if err != nil {
		return err
	}

	fmt.Printf(formattedStatus)
	return nil
}

/*
map[172.31.34.73:{
		"kubernetes":{
			"name":"ebpf-fxgnd","uid":"1fdd80d1-59b1-11e9-b0bc-0a4ce296e9ba"},
		"docker":null}
	172.31.39.176:{
		"kubernetes":{
			"name":"colorteller-red-6d5f5849d6-nz6d7",
			"uid":"ff79b6f6-49b3-11e9-b0bc-0a4ce296e9ba"},
		"docker":null
	}
 */
