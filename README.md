# GPU Tracker Operator

Kubernetes Operator automatically discover GPU Nodes and record them in a Custom Resource

## Summary

- [Overview](#overview)
- [CRD Example](#crd-example)
  - [Spec Details](#spec-details)
- [Installation From Registry](#installation-from-registry)
- [Installation From GitHub Repository](#installation-from-github-repository)
- [Uninstall](#uninstall)


## Overview

GPU Tracker Operator monitors all Nodes in the Kubernetes cluster and finds Nodes that have a specific label: `node-type: gpu-node`.

It collects all matching Node names, and updates the corresponding custom resource (`GPUTracker`) under `gpu_nodes`.

It is present pipeline for building docker image and publishing to GHCR (GitHub Container Registry), and pipeline for packaging helm charts and publishing to GitHub Pages.

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

## Installation From Registry

- Docker image is already deployed to GHCR (GitHub Registry) and helm package is already deployed to GH Pages (GitHub Pages).

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
    --set image.repository=ghcr.io/leooamaral/gpu-tracker-operator \
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

## Installation From GitHub Repository

- Dockerfile is present in the repository so you can also build the operator locally.

1. Build docker image
```
docker build -t <registry-name>/<repository-name>:<tag> .
```
Or if you want to build specifying the OS platform
```
docker build -t <registry-name>/<repository-name>:<tag> --build-arg TARGETOS=linux TARGETARCH=arm64 .
```

2. Push docker image
```
docker push  <registry-name>/<repository-name>:<tag>
```

3. Install a chart 
```
helm install <operator-system-name> ./deploy/charts/gpu-tracker-operator -f <values-path-file> -n <operator-system-namespace> \
    --set image.repository=<registry-name>/<repository-name> \
    --set image.tag=<tag> \
    --wait
```

4. Verify deployment
```
kubectl get pods -n <operator-system-namespace>
```
And you should see a Pod like:
```
gpu-tracker-controller-59dc7484ff-k6qlf   1/1     Running
```

5. Apply a GPUTracker CR
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

6. Verify GPUTracker CR
```
kubectl get gputrackers
```
And you should see a Pod like:
```
suse-gpu-tracker
```

## Uninstall

1. To uninstall the Operator (but keep the CRD):
```
helm uninstall <operator-system-name> -n <operator-system-namespace>
```
- ATTENTION: Helm does not delete CRDs automatically to protect your data. If you want to manually delete the CRD, run:
```
kubectl delete crd gputrackers.suse.tests.dev
```