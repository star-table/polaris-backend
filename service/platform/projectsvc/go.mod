module github.com/star-table/polaris-backend/service/platform/projectsvc

go 1.13

replace github.com/star-table/polaris-backend/common => ./../../../common

replace github.com/star-table/polaris-backend/facade => ./../../../facade

require (
	github.com/DeanThompson/ginpprof v0.0.0-20190408063150-3be636683586
	github.com/araddon/dateparse v0.0.0-20210429162001-6b43995a97de
	github.com/cznic/sortutil v0.0.0-20150617083342-4c7342852e65
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/getsentry/sentry-go v0.12.0
	github.com/gin-gonic/gin v1.8.1
	github.com/go-kratos/kratos/contrib/registry/nacos/v2 v2.0.0-20221130041748-2cf82fa4a75c
	github.com/go-kratos/kratos/v2 v2.5.3
	github.com/jtolds/gls v4.20.0+incompatible
	github.com/magiconair/properties v1.8.5
	github.com/modern-go/reflect2 v1.0.2
	github.com/nacos-group/nacos-sdk-go v1.0.9
	github.com/opentracing/opentracing-go v1.2.0
	github.com/penglongli/gin-metrics v0.1.10
	github.com/shopspring/decimal v1.2.0
	github.com/smartystreets/goconvey v1.6.4
	github.com/spf13/cast v1.5.0
	github.com/star-table/common v1.6.9
	github.com/star-table/go-common v1.0.0
	github.com/star-table/interface v0.0.0-20230707032058-aa3d85d8a825
	github.com/star-table/polaris-backend/common v0.0.0-00010101000000-000000000000
	github.com/star-table/polaris-backend/facade v0.0.0-00010101000000-000000000000
	github.com/tealeg/xlsx/v2 v2.0.1
	github.com/xuri/excelize/v2 v2.6.1
	golang.org/x/text v0.7.0
	google.golang.org/protobuf v1.28.1
	gopkg.in/fatih/set.v0 v0.2.1
	gopkg.in/go-playground/assert.v1 v1.2.1
	gotest.tools v2.2.0+incompatible
	upper.io/db.v3 v3.7.1+incompatible
)

replace upper.io/db.v3 v3.7.1+incompatible => github.com/star-table/db v0.3.75-0.20230707012646-28b2e2303a74
