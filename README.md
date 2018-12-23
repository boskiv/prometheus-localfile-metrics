# prometheus-localfile-metrics
[![Build Status](https://travis-ci.org/boskiv/prometheus-localfile-metrics.svg?branch=master)](https://travis-ci.org/boskiv/prometheus-localfile-metrics)

# What is PLM

`PLM stands for Prometheus Localfile Metrics`

# How it works
- App scan directory in `PLM_STATS_PATH` recursively
- find all files and folders
- make stat name from `PLM_PREFIX` variable and relative dir path to `PLM_STATS_PATH` and filename
- get metric from file content
- listen 9102 port and endpoint /metric for prometheus request
- response with string, separated with `\n`
```
➜  ~ curl localhost:9102/metrics
myapp_ccu 100
myapp_cps 200
myapp_rps 300
```

# Kubernetes usage

## Prerequisite
kubernetes + helm

```
helm repo add coreos https://s3-eu-west-1.amazonaws.com/coreos-charts/stable/
helm install coreos/prometheus-operator --name prometheus-operator --namespace monitoring
helm install coreos/kube-prometheus --name kube-prometheus --namespace monitoring
```

## Example

Run as sidekick container with shared volume to containers metric you need

For example

Your container flush metrics to stats directory:
```
/app/stats/
  ccu
  rps
  cps
```

files contains:
```
cat ccu
100
cat rps
200
cat cps
10
```

next you run pod with:
```yml
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: sampleapp
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: app
        service: stats
    spec:
      volumes:
        - name: metrics-data
          emptyDir: {}
      containers:
        - name: yourcontainer
          image: 
          volumeMounts:
          - mountPath: /app/stats
            name: metrics-data
        - name: stats
          env:
          - name: PLM_STATS_PATH
            value: "/app/stats"
          - name: PLM_STATS_PREFIX
            value: "myapp"
          volumeMounts:
          - mountPath: /app/stats
            name: metrics-data
          ports:
          - containerPort: 9102
``` 

and Service
```
apiVersion: v1
kind: Service
metadata:
  name: sampleservice
  labels:
    app: sample-service
    prometheus: kube-prometheus
spec:
  ports:
    - port: 9102
      targetPort: 9102
      protocol: TCP
      name: stats
  selector:
    app: app
    service: stats
```

May be you need also ServiceMonitor
```
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  labels:
    app: sample-exporter
    prometheus: kube-prometheus
  name: semplae-service-monitor
  namespace: monitoring
spec:
  endpoints:
  - interval: 15s
    port: metrics
  selector:
    matchLabels:
      app: sample-service
      prometheus: kube-prometheus
 
```



After this prometheus exporter installed with prometheus operator for example, aome and get your stats with get request on `/metrics` endpoint

Example respone will be
```
➜  ~ curl localhost:9102/metrics
myapp_ccu 100
myapp_cps 200
myapp_rps 300
```