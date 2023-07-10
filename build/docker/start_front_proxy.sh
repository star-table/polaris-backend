#!/bin/bash

set -x

#取当前目录
BASE_PATH=`cd "$(dirname "$0")"; pwd`

/usr/local/bin/envoy -c /etc/envoy/${SERVICE_NAME}.${POL_ENV}.yaml --service-cluster ${SERVICE_NAME} > /data/logs/${SERVICE_NAME}/docker.log 2>&1
#/usr/local/bin/envoy -c /etc/envoy/${SERVICE_NAME}.${POL_ENV}.yaml --service-cluster ${SERVICE_NAME}