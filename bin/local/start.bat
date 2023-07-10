cd %~dp0
cd ../../
start .\service\basic\idsvc\bin\start.bat %1
start .\service\basic\msgsvc\bin\start.bat %1
start .\service\basic\commonsvc\bin\start.bat %1
start .\service\platform\orgsvc\bin\start.bat %1
start .\service\platform\processsvc\bin\start.bat %1
start .\service\platform\projectsvc\bin\start.bat %1
start .\service\platform\resourcesvc\bin\start.bat %1
start .\service\platform\trendssvc\bin\start.bat %1
start .\service\platform\callsvc\bin\start.bat %1
start .\service\platform\websitesvc\bin\start.bat %1
start .\service\platform\ordersvc\bin\start.bat %1
start .\app\bin\start.bat %1
