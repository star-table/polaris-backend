#!/bin/bash
ROOT=$(PWD)
cd ${ROOT}/common && sh refresh_go_mod.sh
cd ${ROOT}/app && sh refresh_go_mod.sh
cd ${ROOT}/facade && sh refresh_go_mod.sh
cd ${ROOT}/service/basic/commonsvc && sh refresh_go_mod.sh
cd ${ROOT}/service/basic/idsvc && sh refresh_go_mod.sh
cd ${ROOT}/service/basic/msgsvc && sh refresh_go_mod.sh
cd ${ROOT}/service/platform/callsvc && sh refresh_go_mod.sh
cd ${ROOT}/service/platform/ordersvc && sh refresh_go_mod.sh
cd ${ROOT}/service/platform/orgsvc && sh refresh_go_mod.sh
cd ${ROOT}/service/platform/projectsvc && sh refresh_go_mod.sh
cd ${ROOT}/service/platform/resourcesvc && sh refresh_go_mod.sh
cd ${ROOT}/service/platform/trendssvc && sh refresh_go_mod.sh
