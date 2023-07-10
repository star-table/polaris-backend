#!/bin/bash

set -x

#取当前目录
BASE_PATH=`cd "$(dirname "$0")"; pwd`

cd $BASE_PATH


mkdir -p /data/logs/polaris-front-proxy-outside
mkdir -p /data/logs/polaris-front-proxy-inside
mkdir -p /data/logs/polaris-schedule
mkdir -p /data/logs/polaris-idsvc
mkdir -p /data/logs/polaris-msgsvc
mkdir -p /data/logs/polaris-appsvc
mkdir -p /data/logs/polaris-app
mkdir -p /data/logs/polaris-orgsvc
mkdir -p /data/logs/polaris-projectsvc
mkdir -p /data/logs/polaris-processsvc
mkdir -p /data/logs/polaris-resourcesvc
mkdir -p /data/logs/polaris-noticesvc
mkdir -p /data/logs/polaris-rolesvc
mkdir -p /data/logs/polaris-callsvc
mkdir -p /data/logs/polaris-trendssvc
mkdir -p /data/logs/polaris-websitesvc
mkdir -p /data/logs/polaris-commonsvc
mkdir -p /data/logs/polaris-ordersvc



# 停止服务，并删除镜像
# docker-compose down --rmi all

# 构建镜像
docker-compose build --parallel

# 推送镜像到仓库
docker-compose push

# 启动服务
#docker-compose up -d


