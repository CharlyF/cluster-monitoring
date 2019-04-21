package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/orisano/uds"
)

const (
	sockPath = "/Users/charly.fontaine/var/run/net.sock"
)

func main() {
	http.HandleFunc("/connections", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{
		"connections": [{
		"pid": 14931,
		"type": 1,
		"family": 0,
		"net_ns": 4026532528,
		"source": "172.17.0.2",
		"dest": "172.17.0.5",
		"sport": 41515,
		"dport": 53,
		"direction": 2,
		"monotonic_sent_bytes": 3079,
		"last_sent_bytes": 0,
		"monotonic_recv_bytes": 5849,
		"last_recv_bytes": 0,
		"monotonic_retransmits": 0,
		"last_retransmits": 0,
		"last_update_epoch": 11457291007926789
		}, {
		"pid": 14931,
		"type": 1,
		"family": 0,
		"net_ns": 4026532528,
		"source": "172.17.0.2",
		"dest": "172.17.0.5",
		"sport": 34681,
		"dport": 53,
		"direction": 2,
		"monotonic_sent_bytes": 2049,
		"last_sent_bytes": 0,
		"monotonic_recv_bytes": 3937,
		"last_recv_bytes": 0,
		"monotonic_retransmits": 0,
		"last_retransmits": 0,
		"last_update_epoch": 11457305024024228
		}, {
		"pid": 15096,
		"type": 0,
		"family": 0,
		"net_ns": 4026532603,
		"source": "172.31.14.185",
		"dest": "172.17.0.2",
		"sport": 15000,
		"dport": 57768,
		"direction": 2,
		"monotonic_sent_bytes": 7324757,
		"last_sent_bytes": 0,
		"monotonic_recv_bytes": 8328870,
		"last_recv_bytes": 0,
		"monotonic_retransmits": 0,
		"last_retransmits": 0,
		"last_update_epoch": 11457311874323512
		}, {
		"pid": 14931,
		"type": 1,
		"family": 0,
		"net_ns": 4026532528,
		"source": "172.17.0.2",
		"dest": "172.17.0.4",
		"sport": 47018,
		"dport": 53,
		"direction": 2,
		"monotonic_sent_bytes": 2391,
		"last_sent_bytes": 0,
		"monotonic_recv_bytes": 4469,
		"last_recv_bytes": 0,
		"monotonic_retransmits": 0,
		"last_retransmits": 0,
		"last_update_epoch": 11457313010005317
		}, {
		"pid": 14931,
		"type": 1,
		"family": 0,
		"net_ns": 4026532528,
		"source": "172.17.0.2",
		"dest": "172.17.0.5",
		"sport": 57574,
		"dport": 53,
		"direction": 2,
		"monotonic_sent_bytes": 2843,
		"last_sent_bytes": 0,
		"monotonic_recv_bytes": 5256,
		"last_recv_bytes": 0,
		"monotonic_retransmits": 0,
		"last_retransmits": 0,
		"last_update_epoch": 11457285009700386
		}, {
		"pid": 14931,
		"type": 1,
		"family": 0,
		"net_ns": 4026532528,
		"source": "172.17.0.2",
		"dest": "172.17.0.4",
		"sport": 60839,
		"dport": 53,
		"direction": 2,
		"monotonic_sent_bytes": 2166,
		"last_sent_bytes": 0,
		"monotonic_recv_bytes": 4097,
		"last_recv_bytes": 0,
		"monotonic_retransmits": 0,
		"last_retransmits": 0,
		"last_update_epoch": 11457309008093123
		}]
	}`)})
	os.Remove(sockPath)
	log.Fatal(uds.ListenAndServe(sockPath, nil))
}
