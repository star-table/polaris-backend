package consts

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestConstAppHeaderTokenName(t *testing.T) {

	convey.Convey("Test AppHeaderTokenName", t, func() {
		convey.Convey("test", func() {
			convey.So(AppHeaderTokenName, convey.ShouldNotBeBlank)
		})
	})
}

func TestConstAppDingtalkDefaultLang(t *testing.T) {

	convey.Convey("Test AppDingtalkDefaultLang", t, func() {
		convey.Convey("test", func() {
			convey.So(AppSourceChannelDingTalkDefaultLang, convey.ShouldNotBeBlank)
		})
	})
}

func TestConstBlankTime(t *testing.T) {

	convey.Convey("Test BlankTime", t, func() {
		convey.Convey("test", func() {
			convey.So(BlankTime, convey.ShouldNotBeBlank)
		})
	})
}

func TestConstBlankDate(t *testing.T) {

	convey.Convey("Test BlankDate", t, func() {
		convey.Convey("test", func() {
			convey.So(BlankDate, convey.ShouldNotBeBlank)
		})
	})
}

func TestConstBlankTimeObject(t *testing.T) {

	convey.Convey("Test BlankTimeObject", t, func() {
		convey.Convey("test", func() {
			convey.So(BlankTimeObject, convey.ShouldNotBeNil)
		})
	})
}

//默认空字符串
func TestConstBlankString(t *testing.T) {

	convey.Convey("Test BlankString", t, func() {
		convey.Convey("test", func() {
			convey.So(BlankString, convey.ShouldBeBlank)
		})
	})
}

//是否被删除
func TestConstAppDeletedFlag(t *testing.T) {

	convey.Convey("Test AppIsDeleted", t, func() {
		convey.Convey("AppIsDeleted", func() {
			convey.So(AppIsDeleted, convey.ShouldEqual, 1)
		})

		convey.Convey("AppIsNoDelete", func() {
			convey.So(AppIsDeleted, convey.ShouldEqual, 2)
		})
	})
}

func TestConstAppUserInUserFlag(t *testing.T) {

	convey.Convey("Test AppUserIsInUse", t, func() {
		convey.Convey("AppUserIsInUse", func() {
			convey.So(AppUserIsInUse, convey.ShouldEqual, 1)
		})

		convey.Convey("AppUserIsNotInUser", func() {
			convey.So(AppUserIsNotInUser, convey.ShouldEqual, 2)
		})
	})
}

//是否流程初始化状态
func TestConstAppInitStatusFlag(t *testing.T) {

	convey.Convey("Test AppIsInitStatus", t, func() {
		convey.Convey("AppIsInitStatus", func() {
			convey.So(AppIsInitStatus, convey.ShouldEqual, 1)
		})

		convey.Convey("AppIsNotInitStatus", func() {
			convey.So(AppIsNotInitStatus, convey.ShouldEqual, 2)
		})
	})
}

//是否可用
func TestConstAppStatusEnableFlag(t *testing.T) {

	convey.Convey("Test AppStatusEnable", t, func() {
		convey.Convey("AppStatusEnable", func() {
			convey.So(AppStatusEnable, convey.ShouldEqual, 1)
		})

		convey.Convey("AppStatusDisabled", func() {
			convey.So(AppStatusDisabled, convey.ShouldEqual, 2)
		})
	})
}

//是否默认
func TestConstAPPIsDefaultFlag(t *testing.T) {

	convey.Convey("Test APPIsDefault", t, func() {
		convey.Convey("APPIsDefault", func() {
			convey.So(APPIsDefault, convey.ShouldEqual, 1)
		})

		convey.Convey("AppIsNotDefault", func() {
			convey.So(AppIsNotDefault, convey.ShouldEqual, 2)
		})
	})
}

func TestConstAppDateFormat(t *testing.T) {

	convey.Convey("Test AppDateFormat", t, func() {
		convey.Convey("AppDateFormat", func() {
			convey.So(AppDateFormat, convey.ShouldEqual, "2006-01-02")
		})
	})
}

func TestConstAppTimeFormat(t *testing.T) {

	convey.Convey("Test AppTimeFormat", t, func() {
		convey.Convey("AppTimeFormat", func() {
			convey.So(AppTimeFormat, convey.ShouldEqual, "2006-01-02 15:04:05")
		})
	})
}

func TestConstAppSystemTimeFormat(t *testing.T) {

	convey.Convey("Test AppSystemTimeFormat", t, func() {
		convey.Convey("AppSystemTimeFormat", func() {
			convey.So(AppSystemTimeFormat, convey.ShouldEqual, "2006-01-02T15:04:05Z")
		})
	})
}

func TestConstAppSystemTimeFormat8(t *testing.T) {

	convey.Convey("Test AppSystemTimeFormat8", t, func() {
		convey.Convey("AppSystemTimeFormat8", func() {
			convey.So(AppSystemTimeFormat8, convey.ShouldEqual, "2006-01-02T15:04:05+08:00")
		})
	})
}

// SAAS运行模式
func TestConstAppRunmodeSaas(t *testing.T) {

	convey.Convey("Test AppRunmodeSaas", t, func() {
		convey.Convey("AppRunmodeSaas", func() {
			convey.So(AppRunmodeSaas, convey.ShouldEqual, 1)
		})
	})
}

func TestConstAppRunmodeSingle(t *testing.T) {

	convey.Convey("Test AppRunmodeSingle", t, func() {
		convey.Convey("AppRunmodeSingle", func() {
			convey.So(AppRunmodeSingle, convey.ShouldEqual, 2)
		})
	})
}

func TestConstAppRunmodePrivate(t *testing.T) {

	convey.Convey("Test AppRunmodePrivate", t, func() {
		convey.Convey("AppRunmodePrivate", func() {
			convey.So(AppRunmodePrivate, convey.ShouldEqual, 3)
		})
	})
}

//初始化时的一些常量定义
func TestConstInitDefaultTeamName(t *testing.T) {

	convey.Convey("Test InitDefaultTeamName", t, func() {
		convey.Convey("InitDefaultTeamName", func() {
			convey.So(InitDefaultTeamName, convey.ShouldEqual, "默认团队")
		})
	})
}

func TestConstInitDefaultTeamNickname(t *testing.T) {

	convey.Convey("Test InitDefaultTeamNickname", t, func() {
		convey.Convey("InitDefaultTeamNickname", func() {
			convey.So(InitDefaultTeamNickname, convey.ShouldEqual, "默认团队昵称")
		})
	})
}

// context key
func TestConstTraceIdKey(t *testing.T) {

	convey.Convey("Test TraceIdKey", t, func() {
		convey.Convey("TraceIdKey", func() {
			//convey.So(TraceIdKey, convey.ShouldEqual, "_traceId")
		})
	})
}

func TestConstHttpContextKey(t *testing.T) {

	convey.Convey("Test HttpContextKey", t, func() {
		convey.Convey("HttpContextKey", func() {
			//convey.So(HttpContextKey, convey.ShouldEqual, "_httpContext")
		})
	})
}

// 默认对象id步长
func TestConstDefaultObjectIdStep(t *testing.T) {

	convey.Convey("Test DefaultObjectIdStep", t, func() {
		convey.Convey("DefaultObjectIdStep", func() {
			convey.So(DefaultObjectIdStep, convey.ShouldEqual, 200)
		})
	})
}

// 系统缓存模式
func TestConstCacheModeRedis(t *testing.T) {

	convey.Convey("Test CacheModeRedis", t, func() {
		convey.Convey("CacheModeRedis", func() {
			convey.So(CacheModeRedis, convey.ShouldEqual, "Redis")
		})
	})
}

func TestConstCacheModeInside(t *testing.T) {

	convey.Convey("Test CacheModeInside", t, func() {
		convey.Convey("CacheModeInside", func() {
			convey.So(CacheModeInside, convey.ShouldEqual, "Inside")
		})
	})
}

// 系统消息队列模式
func TestConstMQModeRocketMQ(t *testing.T) {

	convey.Convey("Test MQModeRocketMQ", t, func() {
		convey.Convey("MQModeRocketMQ", func() {
			convey.So(MQModeRocketMQ, convey.ShouldEqual, "RocketMQ")
		})
	})
}

func TestConstMQModeDB(t *testing.T) {

	convey.Convey("Test MQModeDB", t, func() {
		convey.Convey("MQModeDB", func() {
			convey.So(MQModeDB, convey.ShouldEqual, "DB")
		})
	})
}

func TestConstMQModeKafka(t *testing.T) {

	convey.Convey("Test MQModeKafka", t, func() {
		convey.Convey("MQModeKafka", func() {
			convey.So(MQModeKafka, convey.ShouldEqual, "Kafka")
		})
	})
}

// 发送消息状态
func TestConstSendMQStatusSuccess(t *testing.T) {

	convey.Convey("Test SendMQStatusSuccess", t, func() {
		convey.Convey("SendMQStatusSuccess", func() {
			convey.So(SendMQStatusSuccess, convey.ShouldEqual, 1)
		})
	})
}

func TestConstSendMQStatusFail(t *testing.T) {

	convey.Convey("Test SendMQStatusFail", t, func() {
		convey.Convey("SendMQStatusFail", func() {
			convey.So(SendMQStatusFail, convey.ShouldEqual, 2)
		})
	})
}

// 消息处理状态
//待处理
func TestConstMQStatusWait(t *testing.T) {

	convey.Convey("Test MQStatusWait", t, func() {
		convey.Convey("MQStatusWait", func() {
			convey.So(MQStatusWait, convey.ShouldEqual, 1)
		})
	})
}

//处理中
func TestConstMQStatusHandle(t *testing.T) {

	convey.Convey("Test MQStatusHandle", t, func() {
		convey.Convey("MQStatusHandle", func() {
			convey.So(MQStatusHandle, convey.ShouldEqual, 2)
		})
	})
}

//处理成功
func TestConstMQStatusSuccess(t *testing.T) {

	convey.Convey("Test MQStatusSuccess", t, func() {
		convey.Convey("MQStatusSuccess", func() {
			convey.So(MQStatusSuccess, convey.ShouldEqual, 3)
		})
	})
}

//处理失败
func TestConstMQStatusFail(t *testing.T) {

	convey.Convey("Test MQStatusFail", t, func() {
		convey.Convey("MQStatusFail", func() {
			convey.So(MQStatusFail, convey.ShouldEqual, 4)
		})
	})
}
