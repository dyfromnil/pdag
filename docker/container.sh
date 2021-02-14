#!/bin/bash
if [[ $# -gt 1 ]] ; then
  echo "nums of params exceed 1!"
  exit 1
elif [[ $# -eq 0 ]] ; then
  echo $#
  echo "no param!"
  exit 1
else
  MODE=$1
fi

if [ "${MODE}" == "stopClient" ]; then
  if [ $(docker ps -qf name=client | wc -l) -ne 0 ]; then
    docker stop $(docker ps -qf name=client)
    echo "y" | docker container prune
    docker rmi client:latest
  fi
elif [ "${MODE}" == "stopServer" ]; then
  if [ $(docker ps -qf name=server | wc -l) -ne 0 ]; then
    docker stop $(docker ps -qf name=server)
    echo "y" | docker container prune
    docker rmi server:latest
  fi
elif [ "${MODE}" == "constructClientImage" ]; then
  cd client && sudo docker build -t client -f Dockerfile ../../experiment/client
elif [ "${MODE}" == "constructServerImage" ]; then
  cd server && sudo docker build -t server -f Dockerfile ../../experiment/server0
else
  echo "can not parse param!"
  exit 1
fi