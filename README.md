# Change Logs Sidecar Container

## Local development

### Run Locally

run server locally:
```
sudo go run main.go
```

run client locally:
```
sudo go run client/main.go
```

### Create Control Plane with a local package cache

Create `kind` cluster with Crossplane installed:
```
kind create cluster --config=package-cache/kind.yaml
kubectl apply -f package-cache/pv.yaml
kubectl create ns crossplane-system
kubectl apply -f package-cache/pvc.yaml
helm install crossplane --namespace crossplane-system --create-namespace crossplane-stable/crossplane --set packageCache.pvc=package-cache
```

Create a XRD/XR that can dynamically create objects:
```
kubectl apply -f test/definition.yaml
kubectl apply -f test/composition.yaml
```

### Build and load components into local package cache

build and load sidecar OCI image:
```
go mod tidy
GOOS=linux GOARCH=amd64 go build -o change-log-sidecar
SIDECAR_VERSION=v0.0.3
docker build -t jbw976/change-log-sidecar:${SIDECAR_VERSION} .
kind load docker-image jbw976/change-log-sidecar:${SIDECAR_VERSION}
```

build and load provider:
```
VERSION=v0.0.1 make build.all
docker tag build-cae27abd/provider-kubernetes-amd64 xpkg.upbound.io/provider-kubernetes
kind load docker-image xpkg.upbound.io/provider-kubernetes
up alpha xpkg xp-extract --from-xpkg _output/xpkg/linux_amd64/provider-kubernetes-v0.0.1.xpkg -o ~/dev/package-cache/provider-kubernetes.gz && chmod 644 ~/dev/package-cache/provider-kubernetes.gz
```

### Set up testing scenario with providers and objects

Install extensions (Providers/Functions), along with `DeploymentRuntimeConfig` and other pre-reqs:
```
kubectl apply -f drc.yaml
kubectl apply -f rbac.yaml
kubectl apply -f package-cache/provider.yaml
kubectl apply -f test/functions.yaml
```

Create a ProviderConfig:
```
kubectl apply -f providerconfig.yaml
```

Create some Objects:
```
kubectl apply -f test/claim.yaml
```

Check objects are created and examine the pod logs:
```
crossplane beta trace traceperf.trace-perf.crossplane.io/traceperf-tester
kubectl -n crossplane-system logs -l pkg.crossplane.io/provider=provider-kubernetes --tail=500
kubectl -n crossplane-system logs -l pkg.crossplane.io/provider=provider-kubernetes -c change-log-sidecar
kubectl -n crossplane-system logs -l pkg.crossplane.io/provider=provider-kubernetes -c change-log-sidecar | jq '.timestamp + " " + .provider + " " + .name + " " + .operation'
```

#### Local dev inner loop

Clean up the objects and provider:
```
kubectl delete -f test/claim.yaml
kubectl delete -f providerconfig.yaml
kubectl delete -f provider.yaml
```

The build and set up the testing scenario again:
* build/load provider/sidecar
* install provider/providerconfig
* create objects
* check results/logs again

## Appendix

### Local package cache debugging resources

* https://github.com/crossplane/crossplane/pull/1807 
* https://github.com/crossplane-contrib/provider-aws/blob/master/cluster/local/integration_tests.sh#L64?
* troubleshooting package cache and inspecting crossplane pod filesystem: https://stackoverflow.com/a/78331043 

#### commands to examine the pacakge cache in the crossplane pod (via the kind container)
```
❯ kcs get pod -l app=crossplane -o jsonpath='{.items[0].status.containerStatuses[0].containerID}'
❯ docker exec -it 261c35a10d0d sh
# ctr -n k8s.io t ls | grep 229d50f1f29563366c0252abf5d6453fce5e46056790380a179a52bbfe7dd90d | awk '{print $2}'
# ls -al /proc/2049/root/cache
```

### Push to a registry
```
crossplane xpkg push \
  --package-files=_output/xpkg/linux_amd64/provider-kubernetes-v0.0.1.xpkg,_output/xpkg/linux_arm64/provider-kubernetes-v0.0.1.xpkg \
  xpkg.upbound.io/jaredorg/provider-kubernetes:v0.0.3
```