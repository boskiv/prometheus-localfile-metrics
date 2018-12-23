# prometheus-localfile-metrics
[![Build Status](https://travis-ci.org/boskiv/prometheus-localfile-metrics.svg?branch=master)](https://travis-ci.org/boskiv/prometheus-localfile-metrics)

Docker image: `boskiv/prometheus-localfile-metrics:0.1.0`

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

Look at example/README.md