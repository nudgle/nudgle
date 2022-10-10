# nudgle
[![Go Reference](https://pkg.go.dev/badge/github.com/nudgle/nudgle.svg)](https://pkg.go.dev/github.com/nudgle/nudgle)

nudgle is a monitoring application that watches every new transaction in the Stratis blockchain.  
It uses filters to detect certain details within transactions and alerts your Discord channel about it.

## Table of contents

## Architecture
![https://i.ibb.co/W5R4Sg2/img.png](https://i.ibb.co/W5R4Sg2/img.png)

# Deploy nudgle
### Kubernetes
#### Prerequisites
In order to deploy this to kubernetes you need the following tools:
1. Helm 3
2. Golang
3. Docker
4. Kubectl

#### Steps
In order to deploy nudgle to Kubernetes, we first need to build the image, and upload it to a registry that your  
cluster has access to.

To build an image run the following:
```bash
    export REGISTRY_URI=127.0.0.1:5000
    docker build --target app -t $REGISTRY_URI/nudgle/indexer:1.0.0 -f build/indexer/Dockerfile .
    docker build --target app -t $REGISTRY_URI/nudgle/monitor:1.0.0 -f build/monitor/Dockerfile .
```
Once the docker images are built, you can push them to your registry:  
```bash
    docker push $REGISTRY_URI/nudgle/indexer:1.0.0
    docker push $REGISTRY_URI/nudgle/monitor:1.0.0
```
#### Configuration
Before we can deploy, we need to configure this application so that it can communicate with your stratis node.

Open the `deploy/nudgle/values.yaml` file and configure your docker image registry and set your stratis node connection info



Now that the images are uploaded to your registry, and we have configured it correctly you can deploy the application using helm
```bash
cd deploy/nudgle
helm -n nudgle install nudgle . --create-namespace
```

### Docker
Configure the application by setting the correct settings on the following files:
`config/indexer.yaml` & `config/monitor.yaml`  

Then execute:  
```shell
    docker-compose up -d
```

# Filters
An example filter has been added as disabled under `internal/filters/example-filter`  
In order to enable it you have to uncomment it from: `internal/filters/loadplugins.go`