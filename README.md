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

- 2 configmaps

 - `transactions` keep track of the transatctions seen on each nodes
 - `metadata` keeps the metadata for all pods on each nodes

`transactions`
Keep a key value record of all the transactions.
key is `ip_source-ip_dest-port_src-port_dest` and value contains the metrics of the datagram.
e.g. `monotonic_sent_bytes`
As the app is running on all nodes, they all add the metadata of the containers living on their nodes.

`metadata` is filled by all apps as well, the key is the ip of the pod and the value is the metadata coming from different sources.
- Docker's API
- Kubelet (in our case)
- APIServer for services metadata.

The Configmap is used as an example of bad storage implementation.
As writes are concurrent, an update can be overriden.
A proper implementation would be to use a proper store with the ability to Lock.

## Limitations

-

## Building

-

## Deploying

-
