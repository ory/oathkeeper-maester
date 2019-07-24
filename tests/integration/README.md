# oathkeeper-k8s-controller integration tests.

This directory contains integration tests for oathkeeper-k8s-controller
The tests execute against a cluster. For local testing use either minikube or KIND environment.


## How to run in with "KIND"
- ensure KUBECONFIG is not set: `unset KUBECONFIG`
- execute `make test-integration` from project's root directory

## How to run it against a cluster

- Setup a test environment: either a K8s cluster or minikube. Install the controller.
- Export KUBECONFIG environment variable
- Execute the tests with: `ginkgo -v ./tests/integration/...`
  If you don't have ginkgo binary installed, standard `go test -v ./tests/integration/...` also works, but the output isn't formatted nicely.

