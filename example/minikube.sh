## Start our minikube and enable the local registry
minikube delete
minikube start \
    --cpus=4 \
    --memory=4096 \
    --extra-config=kubelet.authentication-token-webhook=true \
    --extra-config=kubelet.authorization-mode=Webhook \

# Install monitoring
kubectl apply -f helm-rbac-config.yml
helm init --service-account tiller
kubectl rollout status -w deployment/tiller-deploy --namespace=kube-system
helm repo add coreos https://s3-eu-west-1.amazonaws.com/coreos-charts/stable/
helm install coreos/prometheus-operator --name prometheus-operator --namespace monitoring
helm install coreos/kube-prometheus --name kube-prometheus --namespace monitoring
# Wait for grafana ready
kubectl rollout status -w deployment/kube-prometheus-grafana --namespace=monitoring

kubectl apply -f kube-example.yml
kubectl rollout status -w deployment/sampleapp --namespace=example

echo "Minikube is ready"
echo "run grafana access with kubectl port-forward --namespace monitoring svc/kube-prometheus-grafana 3000:80"
echo "login as admin:admin"
echo "and import dashboard.json as Dashboard to view sample metric"