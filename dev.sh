# build formulaci-base image and formulaci-dev-base first
# > cd images
# > ./build.sh base 
# > ./build.sh dev-server
docker build -t formulaci-dev -f ./Dockerfile-dev .
# create volume to keep data
# > docker volume create formulaci
docker stop formulaci-dev && docker rm formulaci-dev
docker run -it --privileged --name=formulaci-dev -p 8099:8099 -v formulaci:/formulaci/data formulaci-dev --storage-driver=overlay2
docker stop formulaci-dev && docker rm formulaci-dev