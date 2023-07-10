#!/bin/bash
#等待初始化
sleep 120

#获取健康检查路径中的字段
if [ ${SERVICE_NAME}"bjx" = "polaris-front-proxy-insidebjx" ];then
echo "skip check"
exit 0

elif [ ${SERVICE_NAME}"bjx" = "polaris-front-proxy-outsidebjx" ];then
echo "skip check"
exit 0

elif [ ${SERVICE_NAME}"bjx" = "polaris-appbjx" ];then
APP="task"

else
APP=$(echo ${SERVICE_NAME} |awk -F '-' '{print $2}')

fi

#拼接健康检查的url
url=http://localhost:${API_PORT}/api/${APP}/health
echo $url

##开始无限循环
while true
do
result="health"

#检查服务是否健康,超过5次(25秒)则为失败
for i in {1..5}
do
statu=$(curl -s -m 10 $url)
if [ ${statu}"bjx" = "okbjx" ];then
echo "health"
break
fi
result="unhealth"
echo "unhealth"
sleep 5
done

#如果健康检查失败,杀掉envoy服务
if [ ${result} = "unhealth" ];then
ps -ef |grep envoy |grep -v grep |awk '{print $1}' |xargs kill -9
fi

sleep 5
done