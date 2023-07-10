package test

import (
	"context"
	"fmt"

	"github.com/star-table/common/core/config"
	"github.com/star-table/common/core/consts"
	"github.com/gin-gonic/gin"
)

var env = ""

const BaseConfigPath = "./../../../../config"
const SelfConfigPath = "./../config"

// 用于 service 目录下的单测的配置文件寻址。
var SubSvcBaseConfigPath = "./../../../../../config"
var SubSvcSelfConfigPath = "./../../config"

const User1001 = int64(1001)
const Org1001 = int64(1001)

const IssueId100 = int64(100)
const IssueId1083 = int64(1083)

// 用户 service 目录下的包的单元测试
func StartUpForSubSvc(f func(ctx context.Context)) func() {
	return func() {
		env = "test" // test
		//配置文件
		err := config.LoadEnvConfig(SubSvcBaseConfigPath, "application.common", env)
		if err != nil {
			fmt.Printf("err:%s\n", err)
		}
		err = config.LoadEnvConfig(SubSvcSelfConfigPath, "application", env)
		if err != nil {
			fmt.Printf("err:%s\n", err)
		}

		//添加token操作
		ginCtx := gin.Context{}
		ginCtx.Set(consts.AppHeaderTokenName, "o2571u1368436t2a471ab7a375459e9f3c6cfe656e5ffc")
		//获得一个顶级上下文
		ctx := context.Background()
		//返回父上下文  在父上下文中设置 key/value  这边是 GinContextKey:ginCtx
		ctx = context.WithValue(ctx, "GinContextKey", &ginCtx)

		//丢进来的方法立刻执行
		f(ctx)
	}
}

func StartUp(f func(ctx context.Context)) func() {
	return func() {
		env = "test"
		//配置文件
		err := config.LoadEnvConfig(BaseConfigPath, "application.common", env)

		if err != nil {
			fmt.Printf("err:%s\n", err)
		}

		err = config.LoadEnvConfig(SelfConfigPath, "application", env)

		if err != nil {
			fmt.Printf("err:%s\n", err)
		}

		//添加token操作
		ginCtx := gin.Context{}
		ginCtx.Set(consts.AppHeaderTokenName, "o2699u1368436tdd3e978e7dd74ae9958694a8b2b23625")
		//获得一个顶级上下文
		ctx := context.Background()
		//返回父上下文  在父上下文中设置 key/value  这边是 GinContextKey:ginCtx
		ctx = context.WithValue(ctx, "GinContextKey", &ginCtx)

		// 注册飞书平台 sdk 套件
		// fsConfig := config.GetConfig().FeiShu
		// platform_sdk.RegisterPlatformInfo("fs", fsConfig.AppId, fsConfig.AppSecret, psdk_consts.RunTypeCorp, nil)

		//丢进来的方法立刻执行
		f(ctx)
	}
}

func StartUpWithUserInfo(f func(userId, orgId int64)) func() {
	return func() {
		env = "local"
		//配置文件
		err := config.LoadEnvConfig(BaseConfigPath, "application.common", env)

		if err != nil {
			fmt.Printf("err:%s\n", err)
		}

		err = config.LoadEnvConfig(SelfConfigPath, "application", env)

		if err != nil {
			fmt.Printf("err:%s\n", err)
		}
		f(User1001, Org1001)
	}
}
