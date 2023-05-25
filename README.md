# ds-resource-injector

## build the binary and the docker image
```
make docker-build
```

## push the docker image to ECR
```
make docker-push
```
## generate cluster certidicates

```
cd config
./gen-certs.sh
kubectl apply -f specs/webhook/webhook.tls.secret.yaml 
```
copy the `caBundle` to `./config/specs/webhook-config/mutating.config.yaml` and `./config/specs/webhook-config/validating.config.yaml`
