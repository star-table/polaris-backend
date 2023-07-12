module github.com/star-table/polaris-backend/app

go 1.13

replace github.com/star-table/polaris-backend/common => ./../common

replace github.com/star-table/polaris-backend/facade => ./../facade

require (
	github.com/99designs/gqlgen v0.12.0
	github.com/DeanThompson/ginpprof v0.0.0-20190408063150-3be636683586
	github.com/agnivade/levenshtein v1.1.1 // indirect
	github.com/dchest/captcha v0.0.0-20200903113550-03f5f0333e1f
	github.com/gin-contrib/gzip v0.0.5
	github.com/gin-gonic/gin v1.7.4
	github.com/go-kratos/kratos/contrib/config/nacos/v2 v2.0.0-20230706115902-bffc1a0989a6
	github.com/go-kratos/kratos/v2 v2.6.3
	github.com/go-openapi/jsonreference v0.19.6 // indirect
	github.com/go-openapi/spec v0.20.3 // indirect
	github.com/go-openapi/swag v0.19.15 // indirect
	github.com/jtolds/gls v4.20.0+incompatible
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mozillazg/go-pinyin v0.15.0
	github.com/nacos-group/nacos-sdk-go v1.0.9
	github.com/opentracing/opentracing-go v1.2.0
	github.com/penglongli/gin-metrics v0.1.10
	github.com/smartystreets/goconvey v1.6.4
	github.com/spf13/cast v1.5.0
	github.com/star-table/common v1.7.1
	github.com/star-table/go-common v1.0.0
	github.com/star-table/interface v0.0.0-20230707032058-aa3d85d8a825
	github.com/star-table/polaris-backend/common v0.0.0-00010101000000-000000000000
	github.com/star-table/polaris-backend/facade v0.0.0-00010101000000-000000000000
	github.com/swaggo/files v0.0.0-20190704085106-630677cd5c14
	github.com/swaggo/gin-swagger v1.3.0
	github.com/swaggo/swag v1.6.7
	github.com/vektah/gqlparser/v2 v2.0.1
)

replace upper.io/db.v3 v3.7.1+incompatible => github.com/star-table/db v0.3.75-0.20230707012646-28b2e2303a74
