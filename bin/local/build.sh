#!/bin/bash

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
chmod +x ./service/platform/websitesvc/bin/*.sh
chmod +x ./service/platform/ordersvc/bin/*.sh

./service/basic/idsvc/bin/build.sh
./service/basic/msgsvc/bin/build.sh
./service/basic/commonsvc/bin/build.sh
./service/platform/orgsvc/bin/build.sh
./service/platform/processsvc/bin/build.sh
./service/platform/projectsvc/bin/build.sh
./service/platform/resourcesvc/bin/build.sh
./service/platform/trendssvc/bin/build.sh
./service/platform/websitesvc/bin/build.sh
./service/platform/ordersvc/bin/build.sh
