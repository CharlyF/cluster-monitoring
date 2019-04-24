package controllers


type Connections struct {
	Conns []ConnectionStats `json:"connections"`
}

// ConnectionStats stores statistics for a single connection
//easyjson:json
type ConnectionStats struct {
	Pid   uint64 `json:"pid"`
	// Source & Dest represented as a string to handle both IPv4 & IPv6
	Source string `json:"source"`
	Dest   string `json:"dest"`
	DPort  uint16 `json:"dport"`

	MonotonicSentBytes uint64 `json:"monotonic_sent_bytes"`
	LastSentBytes      uint64 `json:"last_sent_bytes"`

	MonotonicRecvBytes uint64 `json:"monotonic_recv_bytes"`
	LastRecvBytes      uint64 `json:"last_recv_bytes"`

	MonotonicRetransmits uint32 `json:"monotonic_retransmits"`
	LastRetransmits      uint32 `json:"last_retransmits"`

}
