<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Ory Oathkeeper Maester](#ory-oathkeeper-maester)
  - [Prerequisites](#prerequisites)
  - [How to use it](#how-to-use-it)
  - [Command-line parameters](#command-line-parameters)
    - [Mode options](#mode-options)
    - [Global flags](#global-flags)
    - [Controller mode flags](#controller-mode-flags)
    - [Sidecar mode flags](#sidecar-mode-flags)
    - [Environment variables](#environment-variables)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Ory Oathkeeper Maester

⚠️ ⚠️ ⚠️

> Ory Oathkeeper Maester is developed by the Ory community and is not actively
> maintained by Ory core maintainers due to lack of resources, time, and
> knolwedge. As such please be aware that there might be issues with the system.
> If you have ideas for better testing and development principles please open an
> issue or PR!

⚠️ ⚠️ ⚠️

ORY Maester is a Kubernetes controller that watches for instances of
`rules.oathkeeper.ory.sh/v1alpha1` custom resource (CR) and creates or updates
the Oathkeeper ConfigMap with Access Rules found in the CRs. The controller
passes the Access Rules as an array in a format recognized by the Oathkeeper.

The project is based on
[Kubebuilder](https://github.com/kubernetes-sigs/kubebuilder)

## Prerequisites

- recent version of Go language with support for modules (e.g: 1.12.6)
- make
- kubectl
- kustomize
- [kind](https://github.com/kubernetes-sigs/kind) for local integration testing
- [ginkgo](https://onsi.github.io/ginkgo/) for local integration testing
- access to K8s environment: minikube or KIND
  (https://github.com/kubernetes-sigs/kind), or a remote K8s cluster

## How to use it

- `make` to build the binary
- `make test` to run tests
- `make test-integration` to run integration tests with local KIND environment

Other targets require a working K8s environment. Set `KUBECONFIG` environment
variable to the proper value.

- `make install` to generate CRD file from go sources and install it in the
  cluster
- `make run` to run controller locally

Refer to the Makefile for the details.

## Command-line parameters

Usage example: `./manager [--global-flags] mode [--mode-flags]`

### Mode options

| Name           | Description                                                                                                                                                                                           |
| :------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **controller** | This is the **default** mode of operation, in which `oathkeeper-maester` is expected to be deployed as a separate deployment. It uses the kubernetes api-server and ConfigMaps to store data.         |
| **sidecar**    | Alternative mode of operation, in which the `oathkeeper-maester` is expected to be deployed as a sidecar container to the main application. It uses local filesystem to create the access rules file. |

### Global flags

| Name                       | Description                                                                                                           | Default values |
| :------------------------- | :-------------------------------------------------------------------------------------------------------------------- | :------------: |
| **metrics-addr**           | The address the metric endpoint binds to                                                                              |     `8080`     |
| **enable-leader-election** | Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager. |    `false`     |
| **kubeconfig**             | Paths to a kubeconfig. Only required if out-of-cluster.                                                               | `$KUBECONFIG`  |

### Controller mode flags

| Name                        | Description                                              |       Default values        |
| :-------------------------- | :------------------------------------------------------- | :-------------------------: |
| **rulesConfigmapName**      | Name of the Configmap that stores Oathkeeper rules.      |     `oathkeeper-rules`      |
| **rulesConfigmapNamespace** | Namespace of the Configmap that stores Oathkeeper rules. | `oathkeeper-maester-system` |
| **rulesFileName**           | Name of the key in ConfigMap containing the rules.json   |     `access-rules.json`     |

### Sidecar mode flags

| Name              | Description                                      |         Default values          |
| :---------------- | :----------------------------------------------- | :-----------------------------: |
| **rulesFilePath** | Path to the file with converted Oathkeeper rules | `/etc/config/access-rules.json` |

### Environment variables

| Name          | Description                                                                                                                                                                            | Default values |
| :------------ | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :------------: |
| **NAMESPACE** | Namespace option to scope Oathkeeper maester to one namespace only - useful for running several instances in one cluster. Defaults to "" which means that there is no namespace scope. |       ``       |
