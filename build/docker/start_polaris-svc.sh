#!/bin/bash

set -x

#取当前目录
BASE_PATH=`cd "$(dirname "$0")"; pwd`

cd $BASE_PATH/

./start.sh &
#./health_check.sh &
#envoy -c /etc/envoy/${SERVICE_NAME}.yaml --service-cluster ${SERVICE_NAME}
envoy -c /etc/envoy/${SERVICE_NAME}.yaml --service-cluster ${SERVICE_NAME} > /data/logs/${SERVICE_NAME}/docker.log 2>&1
