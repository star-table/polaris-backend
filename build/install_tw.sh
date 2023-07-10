#!/bin/bash

echo "start install..."

source /root/.bash_profile
echo "POL_ENV:$3"

#取当前目录
BASE_PATH=`cd "$(dirname "$0")"; pwd`

PROJECT_BRANCH=$1

BUILD_NUMBER=$2

POL_ENV=$3

# 项目名
PROJECT_NAME=polaris-backend

PUB_NAME=polaris-project

# 源码目录
SRC_PATH=$BASE_PATH/src

# 发布目录
#PUB_PATH=/data/app/$PUB_NAME

# 备份目录
BAK_PATH=$BASE_PATH/bak
CUR_TIME=`date +%Y%m%d%H%m%s`
SRC_BAK_PATH=$BAK_PATH/src/$CUR_TIME
# PUB_BAK_PATH=$BAK_PATH/$PUB_NAME/$CUR_TIME


# mkdir -p $PUB_PATH
# mkdir -p $PUB_BAK_PATH
mkdir -p $SRC_BAK_PATH
mkdir -p $SRC_PATH


# 清理docker镜像
$SRC_PATH/$PROJECT_NAME/target/remove_docker.sh

# 备份
mv $SRC_PATH/* $SRC_BAK_PATH/
# mv $PUB_PATH/* $PUB_BAK_PATH/



cd $SRC_PATH
git clone https://gitea.bjx.cloud/allstar/$PROJECT_NAME.git

cd $SRC_PATH/$PROJECT_NAME

git checkout $1

echo "git checkout $1"


BRS=(${PROJECT_BRANCH//\// })

if [[ "${BRS[0]}" == "feature" || "${BRS[0]}" == "release" || "${BRS[0]}" == "master" ]]; then

  PACKAGE_TAG=${BUILD_NUMBER}.${BRS[0]}
  # 编译打包
  ./package_tw.sh ${PROJECT_BRANCH} ${PACKAGE_TAG} ${POL_ENV}

  cd $BASE_PATH

  echo $PACKAGE_TAG > tag.txt

  echo "end install..."
else
  echo "打包必须为 feature、release、master 分支"
  exit 1
fi


