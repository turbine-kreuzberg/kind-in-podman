# Kind in Podman

This repository builds a Podman OCI image based on Ubuntu. Various Podman configurations have been made to use [KIND](https://kind.sigs.k8s.io/) in this image.

### Motivation

I use [ArgoCD](https://argo-cd.readthedocs.io/en/stable/) with the [App of Apps pattern](https://argo-cd.readthedocs.io/en/stable/operator-manual/cluster-bootstrapping/) in my K8S setup. Before a new version of a Helm chart is rolled out, it should first be tested in a CI test with the [e2e framework](https://github.com/kubernetes-sigs/e2e-framework).

### Usage
Starting the container
```
podman run -- privileged <image>
```
This command starts a KIND setup with the following settings
```
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
name: kind-in-podman
featureGates:
  KubeletInUserNamespace: true
networking:
  ipFamily: ipv4
  apiServerPort: 6443
  apiServerAddress: 0.0.0.0
```

### Connect to the cluster
``` 
podman run --privileged -p 6443:6443 $(pwd)/kubeconfig:/root/.kube/config <image>
```
In the mounted kubeconfig file, the server address must be adjusted to\
`server: https://0.0.0.0:6443`\
Testing the connection
```
kubectl  --kubeconfig kubeconfig get nodes
```
