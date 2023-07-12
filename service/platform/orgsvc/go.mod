module github.com/star-table/polaris-backend/service/platform/orgsvc

go 1.13

replace github.com/star-table/polaris-backend/common => ./../../../common

replace github.com/star-table/polaris-backend/facade => ./../../../facade

require (
	github.com/DeanThompson/ginpprof v0.0.0-20190408063150-3be636683586
	github.com/dchest/captcha v0.0.0-20200903113550-03f5f0333e1f
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/getsentry/sentry-go v0.12.0
	github.com/gin-contrib/gzip v0.0.6
	github.com/gin-gonic/gin v1.8.1
	github.com/magiconair/properties v1.8.5
	github.com/nyaruka/phonenumbers v1.0.43
	github.com/opentracing/opentracing-go v1.2.0
	github.com/penglongli/gin-metrics v0.1.10
	github.com/smartystreets/goconvey v1.6.4
	github.com/spf13/cast v1.5.0
	github.com/star-table/common v1.7.1
	github.com/star-table/interface v0.0.0-20230707032058-aa3d85d8a825
	github.com/star-table/polaris-backend/common v0.0.0-00010101000000-000000000000
	github.com/star-table/polaris-backend/facade v0.0.0-00010101000000-000000000000
	github.com/tealeg/xlsx/v2 v2.0.1
	gotest.tools v2.2.0+incompatible
	upper.io/db.v3 v3.7.1+incompatible
)

replace upper.io/db.v3 v3.7.1+incompatible => github.com/star-table/db v0.3.75-0.20230707012646-28b2e2303a74
