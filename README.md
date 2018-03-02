# Formula CI v0.1.0
Self-Hosted CI/CD server for Kubernetes on GCP

## Still In Progess, Currently Just Prototype
  
## Support
- GCP and GKE
- Webhook from : gitlab
- Webhook to : slack
- Built-in `Cloud SDK`, `Docker`, `Git`
- Built-in `Nodejs`, `Golang`
- Built-in `npm`, `yarn`
- Auto-deploy k8s pods by configuration in `.formulaci.yaml`
- All system configurations are availible in web interface 
- "Host Mode" run CI pipeline directly in server process which you can re-use `node_modules` and `go/src`
- Docker graph in `/formulaci/data` which can keep images after restart container

## What's Next
- Validate `.formulaci.yaml` format
- Support K8S not run by GKE
- Support other cloud services (Azure, AWS) when I'm boring
- Support github
- "Docker Mode" run CI pipeline in container
- Beautify web interface and add more function
- Some other ideas in my mind ~

## Build
Prepare base image, go to `base` then run `./build.sh base`  
Prepare dev-base image, go to `base` then run `./build.sh dev-base`  
Start dev server in docker, run `./dev.sh`
Build production image, run `./build.sh`

It use dev-base for server building then COPY from it. You can check Dockerfile for more detail.

## VOLUME
All formulaci data (config files, repo, docker, etc..) are at `/formulaci/data`

## DIND
You can pass "Custom daemon flags" in CMD. By default we set `-g /formulaci/data/docker` to keep docker images in volume.

## Webhook
Gitlab is listening on `{host}/webhook/gitlab`  
Others are in progress