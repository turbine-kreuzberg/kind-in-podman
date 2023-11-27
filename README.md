# Kind in Podman

This repository builds a Podman OCI image based on Ubuntu to start a [KIND](https://kind.sigs.k8s.io/) K8S Cluster.

### Motivation

I use [ArgoCD](https://argo-cd.readthedocs.io/en/stable/) with the [App of Apps pattern](https://argo-cd.readthedocs.io/en/stable/operator-manual/cluster-bootstrapping/) in my K8S setup. Before a new version of a Helm chart is rolled out, it should first be tested in a CI test with the [e2e framework](https://github.com/kubernetes-sigs/e2e-framework).

### Requirements

Cgroups V2
IPv6 Kernel Modules
```
ip6_tables
ip6table_nat
```

### Usage
@TODO
