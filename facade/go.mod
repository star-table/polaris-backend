module github.com/star-table/polaris-backend/facade

go 1.13

replace github.com/star-table/polaris-backend/common => ./../common

require (
	github.com/go-kratos/kratos/v2 v2.2.0
	github.com/smartystreets/assertions v1.0.0 // indirect
	github.com/spf13/cast v1.5.0
	github.com/star-table/common v1.6.9
	github.com/star-table/go-common v1.0.0
	github.com/star-table/interface v0.0.0-20230707032058-aa3d85d8a825
	github.com/star-table/polaris-backend/common v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.53.0
	google.golang.org/protobuf v1.28.1
)
