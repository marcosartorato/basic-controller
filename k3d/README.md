# Local Environment

This folder contains everything you need to run a local k3s (via k3d) cluster,
a local image registry, deploy the controller, and test it from inside the cluster.

## Prerequisites

- Docker
- k3d (`brew install k3d` or see https://k3d.io)
- kubectl
- k9s (optional)

## Commands

The following commands can be run from the repository root directory.
This setup uses the ctrl cluster (as defined in `./k3d/cluster.yaml`).

```sh
# Create the cluster (with registry)
k3d cluster create --config ./k3d/cluster.yaml
# Start Tilt to build, deploy, and watch the controller.
# Tilt automatically rebuilds and redeploys on source code changes.
tilt up
```

The cluster is up! The cluster can be explored using `k9s` or `kubectl`.

The controller is there but neither the `configmap` or the `greeting` object.

About the `configmap`, run:

```sh
# Create the test-cm configmap.
kubectl create cm test-cm --from-literal=message=hello
# Check the configmap deployment.
kubectl get cm
# Check the configmap message field.
kubectl get cm/test-cm -o jsonpath='{.data.message}'
```

About the `greeting` object, use the following manifest:

```yaml
# greet.yaml
apiVersion: operator.example.com/v1alpha1
kind: Greeting
metadata:
  name: test
spec:
  message: "Hej!"
```

and run:

```sh
# Create the test Greeting.
kubectl apply -f greet.yaml
# Check the Greeting.
kubectl get greeting
# Check the greeting message field.
kubectl get greeting/test -o jsonpath='{.spec.message}'
```

Note the commands work because Kubernetes pluralizes. Hence, `greetings` is equivalent to `greetings.operator.example.com`.

To see the controller in action, edit the `message` field of the `greeting` object through its manifest:

```sh
kubectl edit greeting/test
```

Check that the `greeting` object has the updated value and that the controller updated the `configmap`.

To clean the local environment:

- Press `Ctrl`+`C` to stop tilt.
- Run `k3d cluster stop ctrl` to stop the cluster.
- Run `k3d cluster delete ctrl` to delete the cluster.
