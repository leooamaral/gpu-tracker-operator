# GPU Tracker Operator

Kubernetes Operator automatically discover GPU Nodes and record them in a Custom Resource

## Summary



## Overview

GPU Tracker Operator monitors all Nodes in the Kubernetes cluster and finds Nodes that have a specific label: `node-type: gpu-node`.

It collects all matching Node names, and updates the corresponding custom resource (`GPUTracker`) under `gpu_nodes`.

## CRD Example
```
apiVersion: suse.tests.dev/v1
kind: GPUTracker
metadata:
  name: suse-gpu-tracker
gpu_nodes: ""
```

### Spec Details

| Field      | Type   | Description nodes
| :--------: | :----: | :----------------:
| `gpu_nodes`  | `string` | List of Kubernetes node names that have the label `node-type: gpu-node`

## Build and Deploy

### Automated

- Docker image is already deployed to GHCR (GitHub Registry) and helm package is already deployed to GH Pages (GitHub Pages).

#### Instructions

1. Add Helm Chart Repository
```
helm repo add gputracker-charts https://leooamaral.github.io/gpu-tracker-operator
```

2. Update information of available charts locally from chart repositories 
```
helm repo update
```

3. Check charts available in the repository
```
helm search repo gputracker-charts
```

4. Install a chart 
```
helm install <operator-system-name> gputracker-charts/gpu-tracker-operator -f <values-path-file> -n <operator-system-namespace> \
    --set image.repository=ghcr.io/leooamaral/gpu-tracker-operator:feat-operator \
    --set image.tag=latest \
    --wait
```

5. Verify deployment
```
kubectl get pods -n <operator-system-namespace>
```
And you should see a Pod like:
```
gpu-tracker-controller-59dc7484ff-k6qlf   1/1     Running
```

6. Apply a GPUTracker CR
```
apiVersion: suse.tests.dev/v1
kind: GPUTracker
metadata:
  name: suse-gpu-tracker
gpu_nodes: ""
```
```
kubectl apply -f config/samples/suse.tests.dev_v1_gputracker.yaml
```

7. Verify GPUTracker CR
```
kubectl get gputrackers
```
And you should see a Pod like:
```
suse-gpu-tracker
```

8. 

### Manually



docker build -t ghcr.io/leooamaral/node-tracker:latest .
docker push  ghcr.io/leooamaral/node-tracker:latest

deploy controller??



deploy crd??

missing:
- helm charts - feito
- dockerfile - feito
- remove operator-sdk files related e limpar codigo - TODO
- how to test the operator in readme.md - em andamento
- pipeline (ci/cd) - upload to ghcr and github artifacts (versioning) + helm chart - em andamento




improvements:
- operator status
- controller logs
- vuln checks
- linting
- code tests for integration and check usage
- leader election