#!/bin/bash

# docker daemon
# dockerd-entrypoint.sh &

# sleep 2s

# init ssh
mkdir ~/.ssh
ssh-keyscan gitlab.com >> ~/.ssh/known_hosts
ssh-keyscan github.com >> ~/.ssh/known_hosts

# check git ssh rsa 
if [ -e /formulaci/data/id_rsa ]
then
  cp /formulaci/data/id_rsa ~/.ssh/id_rsa
  eval $(ssh-agent -s)
  ssh-add ~/.ssh/id_rsa
fi

# check gcloud key
if [ -e /formulaci/data/gcpServiceAccountKey.json ]
then
  gcloud auth activate-service-account --key-file /formulaci/data/gcpServiceAccountKey.json
#   gcloud auth login --no-launch-browser
fi

/formulaci/server