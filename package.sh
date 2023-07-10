#!/bin/bash

set -x

echo "package start ..."

UNAME=`uname`

PROJECT_BRANCH=$1

BUILD_NUMBER=$2

POL_ENV=$3

#取当前目录
BASE_PATH=`cd "$(dirname "$0")"; pwd`

PROJECT_PRE="polaris-"

PACKAGE_NAMES=("service/basic/idsvc service/basic/commonsvc service/basic/msgsvc service/platform/callsvc service/platform/orgsvc service/platform/ordersvc service/platform/processsvc service/platform/projectsvc service/platform/resourcesvc service/platform/rolesvc service/platform/trendssvc service/platform/websitesvc app schedule")

# 获取源码最近一次 git commit log，包含 commit sha 值，以及 commit message
GitCommitLog=`git log --pretty=oneline -n 1`
# 将 log 原始字符串中的单引号替换成双引号
GitCommitLog=${GitCommitLog//\'/\"}
# 检查源码在 git commit 基础上，是否有本地修改，且未提交的内容
GitStatus=`git status -s`
# 获取当前时间
BuildTime=`date +'%Y-%m-%dT%H%M%S'`
# 获取 Go 的版本
BuildGoVersion=`go version`

# 将以上变量序列化至 LDFlags 变量中
LDFlags=" \
    -X 'github.com/star-table/polaris-backend/common/core/buildinfo.GitCommitLog=${GitCommitLog}' \
    -X 'github.com/star-table/polaris-backend/common/core/buildinfo.GitStatus=${GitStatus}' \
    -X 'github.com/star-table/polaris-backend/common/core/buildinfo.BuildTime=${BuildTime}' \
    -X 'github.com/star-table/polaris-backend/common/core/buildinfo.BuildGoVersion=${BuildGoVersion}' \
"

#发布目录
TARGET_PATH="$BASE_PATH/target"

#备份目录
BAK_PATH="$BASE_PATH/bak"
CUR_TIME=`date +%Y%m%d%H%m%s`
BAK_PATH=$BAK_PATH/$CUR_TIME

# 创建编译目标目录和备份目录
mkdir -p $TARGET_PATH
mkdir -p $BAK_PATH

# 清理发布目录
mv $TARGET_PATH/* $BAK_PATH

# 复制docker相关文件
cp -rf $BASE_PATH/build/docker/*  $TARGET_PATH

# 设置环境变量
#export GO111MODULE=on
#export GOPROXY=https://goproxy.io

function build(){
    echo 'package : '$1

    local package_path=$BASE_PATH/$1

    # 获取目录最后一段
    local paths=(${1//\// })
    local project_name=${PROJECT_PRE}${paths[$[${#paths[*]}-1]]}

    local package_target_path=$TARGET_PATH/${project_name}
    echo $package_target_path

    mkdir -p $package_target_path

    # 复制配置文件
    cp -rf $BASE_PATH/config  $package_target_path/
    cp -rf $package_path/config/*  $package_target_path/config/

    if [ -d "$package_path/resources" ]; then
      cp -rf $package_path/resources  $package_target_path/;
    fi

    cp -rf $BASE_PATH/init  $package_target_path/
    cp $BASE_PATH/README.md $package_target_path/
    cp $BASE_PATH/CHANGELOG.md $package_target_path/

    # 复制执行命令
    cp -rf $BASE_PATH/bin  $package_target_path/

    # 修改启动脚本中的应用名
    sed -i "s/#APP_NAME#/${project_name}/g" $package_target_path/bin/single.sh
#    if [[ "$UNAME" == "Darwin" ]]; then
#      sed -ie "s/#APP_NAME#/${project_name}/g" $package_target_path/bin/single.sh
#    else
#      sed -i "s/#APP_NAME#/${project_name}/g" $package_target_path/bin/single.sh
#    fi


    cd $package_path

    # 编译
    if [[ "$POL_ENV" == "local" || "$POL_ENV" == "" ]]; then
        go build -ldflags "$LDFlags" -o $package_target_path/${project_name}-app
    else
        GOOS=linux GOARCH=amd64 go build -ldflags "$LDFlags" -o $package_target_path/${project_name}-app
    fi

    echo '------'

    #清理文件
#    rm $package_path/$1
}


for path in ${PACKAGE_NAMES[*]}; do
  build $path
done


cd ${TARGET_PATH}

chmod +x *.sh

./build_yaml_dockerfile.sh ${BUILD_NUMBER} ${POL_ENV}

./build_docker.sh

echo "package end ..."
