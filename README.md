# Cluster Monitoring

## Data sources

### Node level

- Network tracer

We are fetching data about the network traffic from a uds socket.
The format should be

```
{
	"connections": [
	{
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
		}
	}]
}
```

The implementation and the handling of the values is left to the user.

- Node Metadata Controller

Every 10 seconds, we collect the metadata of the pods and containers running on the node.
We also collect the node metadata.

    - Node IP and Node name
    - Pod:
        - UID
        - Name
        - IP
        - OwnerRef
    - Container:
        - PID
        - Ports
        - Image name

The list of metadata collected for each type is extensible in the
`controllers/types.go`

### Cluster level

- Cluster Metadata Controller

We watch the services

## Processing Pipeline

### Aggregator

- 

### Storage

-


## Limitations

-

## Building

-

## Deploying

-
