# k8sload
CLI tool to handle deployment manifests according to your Kubernetes topology


## Proof of Concept

Basically you would like to setup;

* pod-to-pod (according to your topology requirements )

* pod-to-node (pod from Cluster network and pod load on host network)

* node-to-node (two pods in host network)


Sample YAML manifest

```yaml

## manifest.yaml

topology:
  producer:
    name: producer
    namespace: producer
    spec:
      port: 8080
      replicas: 2
      antiAffinity: true
      topologyKeys:
      - app: nginx
      nodeSelector:
        node: node-two
      nodeAffinityRules:
        az: us-west-1a

  consumer:
    name: consumer
    namespace: consumer
    spec:
      port: 8080
      replicas: 1
      antiAffinity: false
      # topologyKeys:
      # - app: nginx
      nodeSelector:
        node: node-one
      nodeAffinityRules:
        az: us-west-1b
```

then 

```sh
$ kubectl load -f manifest.yaml
```

```sh
$ kubectl get po -n consumer -o wide 

producer-1  ----- located on nodes (us-west-1a)
producer-2  ----- located on nodes (us-west-1a)

```