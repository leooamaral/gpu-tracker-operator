# GPU Node Operator Helm Charts

## Overview
- This branch holds GPU Node Operator helm charts used to deploy the operator on K8s.
- The charts are packaged and delivered by GitHub Actions to GitHub Pages.

## How to use it

1. Add Helm Chart Repository
```
helm repo add  https://
```

2. Update information of available charts locally from chart repositories 
```
helm repo update
```

3. Check charts available in the repository
```
helm search repo 
```

4. Install a chart 
```
helm install <app-name> charts/<chart-name> -f <values-path-file> -n <app-namespace>
```

5. Upgrade a chart checking if a release by this name doesn't already exist, so it will run an install
```
helm upgrade -i <app-name> -charts/<chart-name> -f <values-path-file> -n <app-namespace>
```

