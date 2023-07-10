#!/bin/bash

set -x

#取当前目录
BASE_PATH=`cd "$(dirname "$0")"; pwd`

cd $BASE_PATH

TARGET_PATH=$BASE_PATH

BUILD_NUMBER=$1
POL_ENV=$2


cp $TARGET_PATH/docker-compose-template-hub.yaml $TARGET_PATH/docker-compose.yaml
sed -i "s/#DOCKER_TAG#/${BUILD_NUMBER}/g" $TARGET_PATH/docker-compose.yaml
sed -i "s/#POL_ENV#/${POL_ENV}/g" $TARGET_PATH/docker-compose.yaml

cp -rf $TARGET_PATH/resources/* $TARGET_PATH/

#APP_NAME=front-proxy-outside
#cp $TARGET_PATH/Dockerfile-polaris-front-proxy $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#sed -i "s/#SERVICE_NAME_ENV#/polaris-${APP_NAME}/g" $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#
#cp $TARGET_PATH/docker-compose-front-proxy-base.yaml $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
#sed -i "s/#SERVICE_NAME_ENV#/${APP_NAME}/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
## sed -i "s/#DOCKER_TAG#/${BUILD_NUMBER}/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
##sed -i "s/#POL_ENV#/${POL_ENV}/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
#sed -i "s/#DOCKER_SVC_PORT#/8001/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
#sed -i "s/#HOST_PORT#/8181/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
#
#
#cp $TARGET_PATH/polaris-front-proxy-outside.yaml $TARGET_PATH/polaris-front-proxy-outside.dev.yaml
#cp $TARGET_PATH/polaris-front-proxy-outside.yaml $TARGET_PATH/polaris-front-proxy-outside.test.yaml
#cp $TARGET_PATH/polaris-front-proxy-outside.yaml $TARGET_PATH/polaris-front-proxy-outside.stag.yaml
#cp $TARGET_PATH/polaris-front-proxy-outside.yaml $TARGET_PATH/polaris-front-proxy-outside.prod.yaml
#cp $TARGET_PATH/polaris-front-proxy-outside.yaml $TARGET_PATH/polaris-front-proxy-outside.unittest.yaml
#
#sed -i "s/#k8ssuffix#//g" $TARGET_PATH/polaris-front-proxy-outside.dev.yaml
#sed -i "s/#k8ssuffix#//g" $TARGET_PATH/polaris-front-proxy-outside.test.yaml
#sed -i "s/#k8ssuffix#/.stag.svc.cluster.local/g" $TARGET_PATH/polaris-front-proxy-outside.stag.yaml
#sed -i "s/#k8ssuffix#/.prod.svc.cluster.local/g" $TARGET_PATH/polaris-front-proxy-outside.prod.yaml
#sed -i "s/#k8ssuffix#//g" $TARGET_PATH/polaris-front-proxy-outside.unittest.yaml
#
#
#APP_NAME=front-proxy-inside
#cp $TARGET_PATH/Dockerfile-polaris-front-proxy $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#sed -i "s/#SERVICE_NAME_ENV#/polaris-${APP_NAME}/g" $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#
#cp $TARGET_PATH/docker-compose-front-proxy-base.yaml $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
#sed -i "s/#SERVICE_NAME_ENV#/${APP_NAME}/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
## sed -i "s/#DOCKER_TAG#/${BUILD_NUMBER}/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
##sed -i "s/#POL_ENV#/${POL_ENV}/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
#sed -i "s/#DOCKER_SVC_PORT#/8001/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
#sed -i "s/#HOST_PORT#/8182/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
#
#
#cp $TARGET_PATH/polaris-front-proxy-inside.yaml $TARGET_PATH/polaris-front-proxy-inside.dev.yaml
#cp $TARGET_PATH/polaris-front-proxy-inside.yaml $TARGET_PATH/polaris-front-proxy-inside.test.yaml
#cp $TARGET_PATH/polaris-front-proxy-inside.yaml $TARGET_PATH/polaris-front-proxy-inside.stag.yaml
#cp $TARGET_PATH/polaris-front-proxy-inside.yaml $TARGET_PATH/polaris-front-proxy-inside.prod.yaml
#cp $TARGET_PATH/polaris-front-proxy-inside.yaml $TARGET_PATH/polaris-front-proxy-inside.unittest.yaml
#
#sed -i "s/#k8ssuffix#//g" $TARGET_PATH/polaris-front-proxy-inside.dev.yaml
#sed -i "s/#k8ssuffix#//g" $TARGET_PATH/polaris-front-proxy-inside.test.yaml
#sed -i "s/#k8ssuffix#/.stag.svc.cluster.local/g" $TARGET_PATH/polaris-front-proxy-inside.stag.yaml
#sed -i "s/#k8ssuffix#/.prod.svc.cluster.local/g" $TARGET_PATH/polaris-front-proxy-inside.prod.yaml
#sed -i "s/#k8ssuffix#//g" $TARGET_PATH/polaris-front-proxy-inside.unittest.yaml
#
#
#APP_NAME=idsvc
#cp $TARGET_PATH/Dockerfile-polaris-svc $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#sed -i "s/#SERVICE_NAME_ENV#/polaris-${APP_NAME}/g" $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#sed -i "s/#API_PATH_ENV#/${APP_NAME}/g" $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#sed -i "s/#API_PORT_ENV#/10002/g" $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#
#cp $TARGET_PATH/docker-compose-base.yaml $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
#sed -i "s/#SERVICE_NAME_ENV#/${APP_NAME}/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
## sed -i "s/#DOCKER_TAG#/${BUILD_NUMBER}/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
##sed -i "s/#POL_ENV#/${POL_ENV}/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
#sed -i "s/#DOCKER_SVC_PORT#/10002/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
#
#
#APP_NAME=msgsvc
#cp $TARGET_PATH/Dockerfile-polaris-svc $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#sed -i "s/#SERVICE_NAME_ENV#/polaris-${APP_NAME}/g" $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#sed -i "s/#API_PATH_ENV#/${APP_NAME}/g" $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#sed -i "s/#API_PORT_ENV#/10003/g" $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#
#cp $TARGET_PATH/docker-compose-base.yaml $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
#sed -i "s/#SERVICE_NAME_ENV#/${APP_NAME}/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
## sed -i "s/#DOCKER_TAG#/${BUILD_NUMBER}/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
##sed -i "s/#POL_ENV#/${POL_ENV}/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
#sed -i "s/#DOCKER_SVC_PORT#/10003/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
#
#
#APP_NAME=appsvc
#cp $TARGET_PATH/Dockerfile-polaris-svc $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#sed -i "s/#SERVICE_NAME_ENV#/polaris-${APP_NAME}/g" $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#sed -i "s/#API_PATH_ENV#/${APP_NAME}/g" $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#sed -i "s/#API_PORT_ENV#/10001/g" $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#
#cp $TARGET_PATH/docker-compose-base.yaml $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
#sed -i "s/#SERVICE_NAME_ENV#/${APP_NAME}/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
## sed -i "s/#DOCKER_TAG#/${BUILD_NUMBER}/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
##sed -i "s/#POL_ENV#/${POL_ENV}/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
#sed -i "s/#DOCKER_SVC_PORT#/10001/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
#
#
#APP_NAME=app
#cp $TARGET_PATH/Dockerfile-polaris-svc $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#sed -i "s/#SERVICE_NAME_ENV#/polaris-${APP_NAME}/g" $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#sed -i "s/#API_PATH_ENV#/${APP_NAME}/g" $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#sed -i "s/#API_PORT_ENV#/12000/g" $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#
#cp $TARGET_PATH/docker-compose-base.yaml $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
#sed -i "s/#SERVICE_NAME_ENV#/${APP_NAME}/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
## sed -i "s/#DOCKER_TAG#/${BUILD_NUMBER}/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
##sed -i "s/#POL_ENV#/${POL_ENV}/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
#sed -i "s/#DOCKER_SVC_PORT#/12000/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
#
#
#APP_NAME=orgsvc
#cp $TARGET_PATH/Dockerfile-polaris-svc $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#sed -i "s/#SERVICE_NAME_ENV#/polaris-${APP_NAME}/g" $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#sed -i "s/#API_PATH_ENV#/${APP_NAME}/g" $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#sed -i "s/#API_PORT_ENV#/12001/g" $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#
#cp $TARGET_PATH/docker-compose-base.yaml $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
#sed -i "s/#SERVICE_NAME_ENV#/${APP_NAME}/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
## sed -i "s/#DOCKER_TAG#/${BUILD_NUMBER}/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
##sed -i "s/#POL_ENV#/${POL_ENV}/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
#sed -i "s/#DOCKER_SVC_PORT#/12001/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
#
#
#
#APP_NAME=projectsvc
#cp $TARGET_PATH/Dockerfile-polaris-svc $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#sed -i "s/#SERVICE_NAME_ENV#/polaris-${APP_NAME}/g" $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#sed -i "s/#API_PATH_ENV#/${APP_NAME}/g" $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#sed -i "s/#API_PORT_ENV#/12002/g" $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#
#cp $TARGET_PATH/docker-compose-base.yaml $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
#sed -i "s/#SERVICE_NAME_ENV#/${APP_NAME}/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
## sed -i "s/#DOCKER_TAG#/${BUILD_NUMBER}/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
##sed -i "s/#POL_ENV#/${POL_ENV}/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
#sed -i "s/#DOCKER_SVC_PORT#/12002/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
#
#
#
#APP_NAME=processsvc
#cp $TARGET_PATH/Dockerfile-polaris-svc $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#sed -i "s/#SERVICE_NAME_ENV#/polaris-${APP_NAME}/g" $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#sed -i "s/#API_PATH_ENV#/${APP_NAME}/g" $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#sed -i "s/#API_PORT_ENV#/12003/g" $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#
#cp $TARGET_PATH/docker-compose-base.yaml $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
#sed -i "s/#SERVICE_NAME_ENV#/${APP_NAME}/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
## sed -i "s/#DOCKER_TAG#/${BUILD_NUMBER}/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
##sed -i "s/#POL_ENV#/${POL_ENV}/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
#sed -i "s/#DOCKER_SVC_PORT#/12003/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
#
#
#
#APP_NAME=resourcesvc
#cp $TARGET_PATH/Dockerfile-polaris-svc $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#sed -i "s/#SERVICE_NAME_ENV#/polaris-${APP_NAME}/g" $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#sed -i "s/#API_PATH_ENV#/${APP_NAME}/g" $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#sed -i "s/#API_PORT_ENV#/12004/g" $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#
#cp $TARGET_PATH/docker-compose-base.yaml $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
#sed -i "s/#SERVICE_NAME_ENV#/${APP_NAME}/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
## sed -i "s/#DOCKER_TAG#/${BUILD_NUMBER}/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
##sed -i "s/#POL_ENV#/${POL_ENV}/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
#sed -i "s/#DOCKER_SVC_PORT#/12004/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
#
#
#
#APP_NAME=noticesvc
#cp $TARGET_PATH/Dockerfile-polaris-svc $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#sed -i "s/#SERVICE_NAME_ENV#/polaris-${APP_NAME}/g" $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#sed -i "s/#API_PATH_ENV#/${APP_NAME}/g" $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#sed -i "s/#API_PORT_ENV#/12005/g" $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#
#cp $TARGET_PATH/docker-compose-base.yaml $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
#sed -i "s/#SERVICE_NAME_ENV#/${APP_NAME}/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
## sed -i "s/#DOCKER_TAG#/${BUILD_NUMBER}/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
##sed -i "s/#POL_ENV#/${POL_ENV}/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
#sed -i "s/#DOCKER_SVC_PORT#/12005/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
#
#
#
#APP_NAME=rolesvc
#cp $TARGET_PATH/Dockerfile-polaris-svc $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#sed -i "s/#SERVICE_NAME_ENV#/polaris-${APP_NAME}/g" $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#sed -i "s/#API_PATH_ENV#/${APP_NAME}/g" $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#sed -i "s/#API_PORT_ENV#/12006/g" $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#
#cp $TARGET_PATH/docker-compose-base.yaml $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
#sed -i "s/#SERVICE_NAME_ENV#/${APP_NAME}/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
## sed -i "s/#DOCKER_TAG#/${BUILD_NUMBER}/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
##sed -i "s/#POL_ENV#/${POL_ENV}/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
#sed -i "s/#DOCKER_SVC_PORT#/12006/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
#
#
#
#
#APP_NAME=callsvc
#cp $TARGET_PATH/Dockerfile-polaris-svc $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#sed -i "s/#SERVICE_NAME_ENV#/polaris-${APP_NAME}/g" $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#sed -i "s/#API_PATH_ENV#/${APP_NAME}/g" $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#sed -i "s/#API_PORT_ENV#/12007/g" $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#
#cp $TARGET_PATH/docker-compose-base.yaml $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
#sed -i "s/#SERVICE_NAME_ENV#/${APP_NAME}/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
## sed -i "s/#DOCKER_TAG#/${BUILD_NUMBER}/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
##sed -i "s/#POL_ENV#/${POL_ENV}/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
#sed -i "s/#DOCKER_SVC_PORT#/12007/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
#
#
#
#APP_NAME=trendssvc
#cp $TARGET_PATH/Dockerfile-polaris-svc $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#sed -i "s/#SERVICE_NAME_ENV#/polaris-${APP_NAME}/g" $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#sed -i "s/#API_PATH_ENV#/${APP_NAME}/g" $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#sed -i "s/#API_PORT_ENV#/12008/g" $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#
#cp $TARGET_PATH/docker-compose-base.yaml $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
#sed -i "s/#SERVICE_NAME_ENV#/${APP_NAME}/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
## sed -i "s/#DOCKER_TAG#/${BUILD_NUMBER}/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
##sed -i "s/#POL_ENV#/${POL_ENV}/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
#sed -i "s/#DOCKER_SVC_PORT#/12008/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
#
#
#
#APP_NAME=schedule
#cp $TARGET_PATH/Dockerfile-polaris-svc $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#sed -i "s/#SERVICE_NAME_ENV#/polaris-${APP_NAME}/g" $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#
#sed -i "s/#API_PATH_ENV#/${APP_NAME}/g" $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#sed -i "s/#API_PORT_ENV#/12009/g" $TARGET_PATH/Dockerfile-polaris-${APP_NAME}
#
#cp $TARGET_PATH/docker-compose-base.yaml $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
#sed -i "s/#SERVICE_NAME_ENV#/${APP_NAME}/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
## sed -i "s/#DOCKER_TAG#/${BUILD_NUMBER}/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
##sed -i "s/#POL_ENV#/${POL_ENV}/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml
#sed -i "s/#DOCKER_SVC_PORT#/12009/g" $TARGET_PATH/docker-compose-base-polaris-${APP_NAME}.yaml