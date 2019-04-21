package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/CharlyF/cluster-monitoring/pkg/aggregator"
	"github.com/CharlyF/cluster-monitoring/util"
	"github.com/orisano/uds"
	"io/ioutil"
	"net/http"
	"time"
)

type SocketController struct {
	stop chan struct {}
	conn  *http.Client
	received chan([]byte)
}

var received chan []byte

func NewSocketController(path string) (*SocketController, error) {
	stop := make(chan struct{})
	client := uds.NewClient(path)
	received = make(chan []byte, 2400)
	return &SocketController{
		conn:client,
		stop:stop,
		received: received,
	}, nil
}

func (s *SocketController) Run(stop chan struct{}) {
	go s.processConnections(stop)
	tick := time.NewTicker(2 * time.Second)
	for {
		select {
		case <- tick.C:
			//buf := make([]byte, 1024)
			n, err := s.conn.Get("http://foo/connections")
			if err != nil {
				fmt.Errorf(err.Error())
				s.stop <- struct {}{}
				return
			}
			resp, err := ioutil.ReadAll(n.Body)
			if err != nil {
				fmt.Errorf("Could not read Body: %s", err.Error())
				continue
			}
			fmt.Printf("received %d bytes \n", len(resp))
			s.received <- resp
			n.Body.Close()
		case <- stop:
			fmt.Printf("Stopped Running \n")
			return
		}
	}
}

func (s *SocketController) processConnections(stop chan struct{}) {
	fmt.Printf("Starting processConnections \n")
	for {
		select {

		case msg := <-s.received:
			var c Connections
			err := json.Unmarshal(msg, &c)
			if err != nil {
				fmt.Errorf("err unmarshaling %s \n", err.Error())
			}
			transactionSaver(c.Conns)
		case <- stop:
			fmt.Printf("stopped processing connections... \n")
			return
		}
	}
}

// func that creates a key made of the src.dest
// which values holds the transaction details.
// Assumes that the IP is unique accross the cluster.
func transactionSaver( co []ConnectionStats){
	// Marshal Format
	for _, c := range co {
		key := fmt.Sprintf("%s-%s-%d-%d", c.Source, c.Dest, c.SPort,c.DPort)
		a := aggregator.Values{
			MonotonicSentBytes:c.MonotonicSentBytes,
			MonotonicRecvBytes: c.MonotonicRecvBytes,
		}
		byteA, err :=json.Marshal(a)
		if err != nil{
			fmt.Errorf(err.Error())
			continue
		}
		fmt.Printf("Caching: %s: %s \n", key, string(byteA))

		util.Cache.Add(key, string(byteA), 0)
	}
}
