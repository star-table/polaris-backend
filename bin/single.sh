#!/bin/bash

# 需定义环境变量 POL_ENV：取值范围：local：本机；dev：开发环境；sit：测试环境：uat：预发布环境；prod：生产环境

APP_NAME="#APP_NAME#"

#取当前目录
BASE_PATH=`cd "$(dirname "$0")"; pwd`
#日志路径
LOG_PATH=/data/logs/$APP_NAME
mkdir -p $LOG_PATH

exist(){
    if test $( pgrep -f "${APP_NAME}-app" | wc -l ) -eq 0
    then
        return 1
    else
        return 0
    fi
}

start(){
    if exist; then
        echo "$APP_NAME is already running."
        exit 1
    else
        cd $BASE_PATH/../
        nohup ./$APP_NAME-app >$LOG_PATH/stdout.log 2>&1 &
#        if [[ "$POL_ENV" == "local" || "$POL_ENV" == "" || "$POL_ENV" == "dev" || "$POL_ENV" == "test" || "$POL_ENV" == "unittest" ]]; then
#            nohup ./$APP_NAME-app >$LOG_PATH/stdout.log 2>&1 &
#        else
#            nohup ./$APP_NAME-app >/dev/stdout 2>&1 &
#        fi
    fi
}

stop(){
    runningPID=`pgrep -f "${APP_NAME}-app"`
    if [ "$runningPID" ]; then
        echo "$APP_NAME pid: $runningPID"
        count=0
        kwait=5
        echo "$APP_NAME is stopping, please wait..."
        kill -15 $runningPID
        until [ `ps --pid $runningPID 2> /dev/null | grep -c $runningPID 2> /dev/null` -eq '0' ] || [ $count -gt $kwait ]
        do
            sleep 1
            let count=$count+1;
        done

        if [ $count -gt $kwait ]; then
            kill -9 $runningPID
        fi
        clear
        echo "$APP_NAME is stopped."
    else
        echo "$APP_NAME has not been started."
    fi
}

check(){
   if exist; then
   	 echo "$APP_NAME is alive."
   	 exit 0
   else
   	 echo "$APP_NAME is dead."
   	 exit -1
   fi
}

restart(){
    stop
    start
}


case "$1" in

start)
    start
;;
stop)
    stop
;;
restart)
    restart
;;
check)
    check
;;
*)
    echo "available operations: [start|stop|restart|check]"
    exit 1
;;
esac