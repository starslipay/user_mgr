#!/bin/bash
MODULE_NAME="user_mgr"
VERSION="v1.0.0"
IMAGE_NAME="${MODULE_NAME}:${VERSION}"

docker rm -f $MODULE_NAME
docker rmi -f $IMAGE_NAME
docker build -t $IMAGE_NAME .
docker run -d --name $MODULE_NAME --network dev_pay_net -p 30880:8080 $IMAGE_NAME
# docker run -d --name user_mgr --network dev_pay_net -p 30884:8080 user_mgr:v1.0.0
docker ps
docker logs $MODULE_NAME