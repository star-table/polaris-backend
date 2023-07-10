#!/bin/bash

set -x

#取当前目录
BASE_PATH=`cd "$(dirname "$0")"; pwd`

cd $BASE_PATH


# 停止服务，并删除镜像
docker-compose down --rmi all


