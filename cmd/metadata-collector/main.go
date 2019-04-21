package main

import (
	"fmt"
	"github.com/CharlyF/cluster-monitoring/cmd/metadata-collector/app"
	"os"
)

func main() {
	//go http.ListenAndServe("0.0.0.0:1234", nil)
	if err := app.MetadataCollector.Execute(); err != nil {
		fmt.Errorf(err.Error())
		os.Exit(-1)
	}
}
