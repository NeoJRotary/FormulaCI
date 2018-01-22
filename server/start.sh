#!/bin/bash
set -e

# run docker daemon
# no arguments passed
# or first arg is `-f` or `--some-option`
if [ "$#" -eq 0 -o "${1#-}" != "$1" ]; then
	# add our default arguments
	set -- dockerd \
		--host=unix:///var/run/docker.sock \
		--host=tcp://0.0.0.0:2375 \
    -g /formulaci/data/docker/ \
		"$@"
fi

if [ "$1" = 'dockerd' ]; then
	# if we're running Docker, let's pipe through dind
	# (and we'll run dind explicitly with "sh" since its shebang is /bin/bash)
	set -- sh "$(which dind)" "$@"
fi

exec "$@" &


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