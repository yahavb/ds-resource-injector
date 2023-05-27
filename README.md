# ds-resource-injector

Daemonset log and metric forwarders (e.g., Fluentbit) require CPU and memory resources proportional to node size. Cluster autoscaler or Karpenter do not factor in the anticipated CPU or memory load from the Daemonset pods resulting in pod eviction. Daemonset pods (agents) should allocate static resources to avoid impact on service pods that run on a node. If that is not possible you can use this webhook that [interceptss](./pkg/admission/admission.go) any `pod` create call in a [namespace decorated](./specs/apps.ns.yaml) with an `admission-webhook: enabled` label and [allocates](./pkg/mutation/inject_ds_res.go) 10% CPU of the node size. 

## build the binary and the docker image
```
make docker-build
```

## push the docker image to ECR
```
make docker-push
```
## generate and deploy cluster certificates

```
cd config
./gen-certs.sh
kubectl apply -f specs/webhook.tls.secret.yaml 
```
copy the `caBundle` to `./specs/mutating.config.yaml` and `./specs/validating.config.yaml`

## deploy the MutatingWebhookConfiguration and ValidatingWebhookConfiguration

```
kubectl apply -f ./specs/mutating.config.yaml
kubectl apply -f ./specs/validating.config.yaml
```

## deploy the webhook server and service

```
kubectl apply -f ./specs/webhook.deploy.yaml
kubectl apply -f ./specs/webhook.svc.yaml
```

## test the webhook

The webhook MUST be deployed in OTHER namesapce than it is running to avoid cyclic dependencies between the webhook and the workloads it controls. Therefore, we created the name space `apps` to deploy our test pods. We deploy two type of pods, service pods (under ReplicaSet) and daemonsets pods. 

create the apps namespace:

```
kubectl apply -f ./specs/apps.ns.yaml 
```

The webhook allocate 10% of the node CPU capacity. e.g., c6a.large has 2 CPUs so every daemonset pods will use 200m CPU. 

### Test 1 - deploy service pods 

```
kubectl apply -f ./specs/alpine-deploy.yaml
```
note that the pods in namespace `apps` defined in `./specs/alpine-deploy.yaml` have no `Resources` configured because the pods are NOT daemonset pods. 

### Test 2 - deploy daemonset pods

```
kubectl apply -f ./specs/alpine-ds.yaml
```
note that the pods created in namespace `apps` the belongs to `./specs/alpine-ds.yaml` now decorated with cpu resources requests and limits. 


