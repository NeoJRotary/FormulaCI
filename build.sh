docker build -t formulaci-dev -f ./Dockerfile-dev .
cd web
npm run production
cd ..
docker build -t formulaci .