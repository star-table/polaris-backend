cd "$(dirname "$0")"
cd ../../
chmod +x ./service/basic/idsvc/bin/*.sh
chmod +x ./service/basic/msgsvc/bin/*.sh
chmod +x ./service/basic/commonsvc/bin/*.sh
chmod +x ./service/platform/orgsvc/bin/*.sh
chmod +x ./service/platform/processsvc/bin/*.sh
chmod +x ./service/platform/projectsvc/bin/*.sh
chmod +x ./service/platform/resourcesvc/bin/*.sh
chmod +x ./service/platform/trendssvc/bin/*.sh
chmod +x ./service/platform/callsvc/bin/*.sh
chmod +x ./service/platform/websitesvc/bin/*.sh
chmod +x ./service/platform/ordersvc/bin/*.sh

./service/basic/idsvc/bin/start.sh $1 &
./service/basic/msgsvc/bin/start.sh $1 &
./service/basic/commonsvc/bin/start.sh $1 &
./service/platform/orgsvc/bin/start.sh $1 &
./service/platform/processsvc/bin/start.sh $1 &
./service/platform/projectsvc/bin/start.sh $1 &
./service/platform/resourcesvc/bin/start.sh $1 &
./service/platform/trendssvc/bin/start.sh $1 &
./service/platform/callsvc/bin/start.sh $1 &
./service/platform/websitesvc/bin/start.sh $1 &
./service/platform/ordersvc/bin/start.sh $1 &
./app/bin/start.sh $1 &
