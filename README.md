# oathkeeper-k8s-controller

This project contains a Kubernetes controller that allows users to manage Oathkeeper rules definitions.

## Prerequisites

- recent version of go with support for modules (e.g: 1.12.6)
- make
- kubectl
- access to running K8s cluster / minikube / KIND (https://github.com/kubernetes-sigs/kind)


## How to use it

- `make` to build the binary
- `make test` to run tests

Other targets require a working K8s environment: a remote cluster, a minikube or a local KIND.
Set `KUBECONFIG` environment variable to the proper value.

- `make run` to run controller locally
- `make install` to generate CRD file from go sources and install it in the cluster


Refer to the Makefile for details.

