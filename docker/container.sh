docker stop $(docker ps -q)
docker container prune
docker rmi server:latest client:latest

cd client && sudo docker build -t client -f Dockerfile ../../experiment/client
cd ../server && sudo docker build -t server -f Dockerfile ../../experiment/server0