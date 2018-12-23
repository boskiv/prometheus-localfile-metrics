# examples

Run minikube.sh

Wait for app deployed

Open grafana and import dashboard.json

```
kubectl port-forward --namespace monitoring svc/kube-prometheus-grafana 3000:80
```
Open prometheus an look Targets and Service discovery settings
```
kubectl port-forward --namespace monitoring svc/kube-prometheus 9090:9090
```