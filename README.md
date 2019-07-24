# oathkeeper-k8s-controller

This project contains a Kubernetes controller that uses Custom Resources to manage Oathkeeper rules definitions.
The project is based on [Kubebuilder](https://github.com/kubernetes-sigs/kubebuilder)

## Prerequisites

- recent version of Go language with support for modules (e.g: 1.12.6)
- make
- kubectl
- kustomize
- [kind](https://github.com/kubernetes-sigs/kind) for local integration testing
- [ginkgo](https://onsi.github.io/ginkgo/) for local integration testing
- access to K8s environment: minikube or KIND (https://github.com/kubernetes-sigs/kind), or a remote K8s cluster


## How to use it

- `make` to build the binary
- `make test` to run tests
- `make test-integration` to run integration tests with local KIND environment

Other targets require a working K8s environment.
Set `KUBECONFIG` environment variable to the proper value.

- `make install` to generate CRD file from go sources and install it in the cluster
- `make run` to run controller locally

Refer to the Makefile for the details.

