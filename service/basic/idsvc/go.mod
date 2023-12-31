module github.com/star-table/polaris-backend/service/basic/idsvc

go 1.13

replace github.com/star-table/polaris-backend/common => ./../../../common

require (
	github.com/DeanThompson/ginpprof v0.0.0-20190408063150-3be636683586
	github.com/gin-contrib/gzip v0.0.6
	github.com/gin-gonic/gin v1.8.1
	github.com/opentracing/opentracing-go v1.2.0
	github.com/penglongli/gin-metrics v0.1.10
	github.com/smartystreets/goconvey v1.6.4
	github.com/star-table/common v1.7.1
	github.com/star-table/polaris-backend/common v0.0.0-00010101000000-000000000000
	upper.io/db.v3 v3.7.1+incompatible
)

replace upper.io/db.v3 v3.7.1+incompatible => github.com/star-table/db v0.3.75-0.20230707012646-28b2e2303a74
