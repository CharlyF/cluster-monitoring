package aggregator

type Data struct {
	Kubernetes PodMetadata `json:"kubernetes,omitempty"`
	//Docker     map[string]ContainerMeta `json:"docker,omitempty"`
	Service    SvcMetadata `json:"service,omitempty"`
	Node	 NodeMeta   `json:"node,omitempty"`
}
type NodeMeta struct {
	Name string `json:"name,omitempty"`
	KernelVersion string `json:"kernel,omitempty"`
}

// PodMetadata represents the Pod metadata we want to enrich
type PodMetadata struct {
	Name  string `json:"name,omitempty"`
	UID string `json:"uid,omitempty"`
	Image string `json:"image,omitempty"`
}

// SvcMetadata represents the Pod metadata we want to enrich
type SvcMetadata struct {
	Name  string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
	Ports string `json:"ports,omitempty"`
}

type ContainerMeta struct {
	Id string `json:"id,omitempty"`
	Image string `json:"image,omitempty"`
	//Ip string `json:"ip,omitempty"`
	// Add ports and pids.
}

type Transactions struct {
	Val Values `json:"values,omitempty"`
	IpSrc string `json:"ipsrc"`
	DataSrc Data  `json:"datadrc"`
	IpDest string `json:"ipdest"`
	DataDest Data  `json:"datadest"`
}

type Values struct {
	MonotonicSentBytes uint64 `json:"monotonic_sent_bytes"`
	//LastSentBytes      uint64 `json:"last_sent_bytes"`

	MonotonicRecvBytes uint64 `json:"monotonic_recv_bytes"`
	//LastRecvBytes      uint64 `json:"last_recv_bytes"`
	//
	MonotonicRetransmits uint32 `json:"monotonic_retransmits"`
	//LastRetransmits      uint32 `json:"last_retransmits"`
}
