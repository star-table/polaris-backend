#!/bin/bash

set -x

#取当前目录
BASE_PATH=`cd "$(dirname "$0")"; pwd`

cd $BASE_PATH

APP_NAME=$1
BUILD_NUMBER=$2
POL_ENV=$3

PACKAGE_NAMES=("app idsvc msgsvc callsvc orgsvc processsvc projectsvc resourcesvc rolesvc trendssvc commonsvc websitesvc ordersvc schedule front-proxy-inside front-proxy-outside")

function restartdocker(){
    echo 'package : '$1

    local appName=$1
    local buildNumber=$2
    local polEnv=$3

    if [ -f "docker-compose-polaris-${appName}.yaml" ]; then
      docker-compose -f docker-compose-polaris-${appName}.yaml down --rmi all
      rm docker-compose-polaris-{appName}.yaml
    fi

    cp docker-compose-base-polaris-${appName}.yaml docker-compose-polaris-${appName}.yaml

    sed -i "s/#DOCKER_TAG#/${buildNumber}/g" docker-compose-polaris-${appName}.yaml
    sed -i "s/#POL_ENV#/${polEnv}/g" docker-compose-polaris-${appName}.yaml

    docker-compose -f docker-compose-polaris-${appName}.yaml up -d
}


if [[ "$APP_NAME" == "all" || "$POL_ENV" == "" ]]; then
    for svcname in ${PACKAGE_NAMES[*]}; do
      restartdocker $svcname $BUILD_NUMBER $POL_ENV
    done
else
    restartdocker $APP_NAME $BUILD_NUMBER $POL_ENV
fi







