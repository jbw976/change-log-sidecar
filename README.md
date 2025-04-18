# Change Logs Sidecar Container

## Testing published packages

Create `kind` cluster with Crossplane installed:
```
kind create cluster
helm install crossplane --namespace crossplane-system --create-namespace crossplane-stable/crossplane
kubectl -n crossplane-system get pod
```

Create a XRD/XR that can dynamically create objects:
```
kubectl apply -f test/definition.yaml
kubectl apply -f test/composition.yaml
```

Install extensions (Providers/Functions), along with `DeploymentRuntimeConfig`
and other pre-reqs:
```
kubectl apply -f drc.yaml
kubectl apply -f rbac.yaml
kubectl apply -f provider.yaml
kubectl apply -f test/functions.yaml
```

Wait for all packages to become installed and healthy:
```
kubectl get pkg
```

Create a ProviderConfig:
```
kubectl apply -f providerconfig.yaml
```

Create some `Objects` (you can update `claim.yaml` to create more or less `Objects`):
```
kubectl apply -f test/claim.yaml
```

Check `Objects` are created and examine the pod logs:
```
crossplane beta trace traceperf.trace-perf.crossplane.io/traceperf-tester
kubectl -n crossplane-system logs -l pkg.crossplane.io/provider=provider-kubernetes -c changelogs-sidecar
kubectl -n crossplane-system logs -l pkg.crossplane.io/provider=provider-kubernetes -c changelogs-sidecar | jq '.timestamp + " " + .provider + " " + .name + " " + .operation'
kubectl -n crossplane-system logs -l pkg.crossplane.io/provider=provider-kubernetes -c changelogs-sidecar --tail 1 | jq .
```

You can also check the main provider logs, which should hopefully be sparse unless something went wrong:
```
kubectl -n crossplane-system logs -l pkg.crossplane.io/provider=provider-kubernetes
```

Now update the claim in order to see an update operation in the change logs,
e.g. by changing the `spec.textData` field to a new string value, then apply it
again:
```
kubectl apply -f test/claim.yaml
```

Then delete the claim in order to clean up the `Objects` and see delete
operations in the change logs:
```
kubectl delete -f test/claim.yaml
```

## Local Development

For local development steps, please refer to the [DEVELOPMENT.md](DEVELOPMENT.md) file.