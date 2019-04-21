package aggregator

type Data struct {
	Kubernetes []PodMetadata `json:"kubernetes"`

}

// PodMetadata represents the Pod metadata we want to enrich
type PodMetadata struct {
	Name  string `json:"name,omitempty"`
	UID string `json:"uid,omitempty"`
}

type Transactions struct {
	Val []Values `json:"values,omitempty"`
}

type Values struct {
	MonotonicSentBytes uint64 `json:"monotonic_sent_bytes"`
	//LastSentBytes      uint64 `json:"last_sent_bytes"`

	MonotonicRecvBytes uint64 `json:"monotonic_recv_bytes"`
	//LastRecvBytes      uint64 `json:"last_recv_bytes"`
	//
	//MonotonicRetransmits uint32 `json:"monotonic_retransmits"`
	//LastRetransmits      uint32 `json:"last_retransmits"`
}
