FROM golang:alpine as golang
FROM node:alpine as nodejs
FROM google/cloud-sdk:alpine as gcloud
# docker:dind is alpine
FROM docker:dind

RUN apk --no-cache add bash git sqlite

### Golang ###
COPY --from=golang /usr/local/go/bin /usr/local/go/bin
COPY --from=golang /usr/local/go/src /usr/local/go/src
COPY --from=golang /usr/local/go/pkg /usr/local/go/pkg
COPY --from=golang /go /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
# ENV GOPATH /go
# ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
# RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

### ###

### Nodejs ###
COPY --from=nodejs /usr/local /usr/local
RUN apk --no-cache add libstdc++
### ###

### GCP Cloud SDK ### 
ENV PATH /google-cloud-sdk/bin:$PATH
COPY --from=gcloud /google-cloud-sdk /google-cloud-sdk
RUN apk --no-cache add \
  python \
  py-crcmod \
  libc6-compat \
  openssh-client \
  && ln -s /lib /lib64 && \
  gcloud config set core/disable_usage_reporting true && \
  gcloud config set component_manager/disable_update_check true && \
  gcloud config set metrics/environment github_docker_image && \
  gcloud components install --quiet kubectl && \
  gcloud --version
# VOLUME ["/root/.config"]
### ###

CMD bash
