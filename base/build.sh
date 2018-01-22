if [ $1 == "base" ]
then
  docker build -t formulaci-base .
fi

if [ $1 == "dev-base" ]
then
  docker build -t formulaci-dev-base -f Dockerfile-dev .
fi