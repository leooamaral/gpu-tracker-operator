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