# ds-resource-injector

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


