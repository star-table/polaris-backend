package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/star-table/common/core/config"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/core/util/strs"
	"github.com/mozillazg/go-pinyin"
	"github.com/smartystreets/goconvey/convey"
)

//
//func TestGenerateStruct(t *testing.T) {
//	err := GenerateStruct("../polaris-task/graphql/models", "ppm_pro_project")
//	t.Log(err)
//}

func TestInitConfig(t *testing.T) {
	fmt.Println(json.ToJson(config.GetMysqlConfig()))

	fmt.Println(json.ToJson(config.GetMysqlConfig()))

	fmt.Println(os.Getenv("POL_ENV"))

}

func TestInitConfig2(t *testing.T) {
	hz := "龘"
	a := pinyin.NewArgs()
	fmt.Println(pinyin.Pinyin(hz, a))

	hz = "CallCenter"
	fmt.Println(pinyin.Pinyin(hz, a))
}

func TestStringLen(t *testing.T) {
	str := "你是一个大苹果"
	t.Log(strs.Len(str))
}

//添加单元测试使用
func TestInitEnvConfig(t *testing.T) {

	convey.Convey("测试加载env2", t, func() {

		config.LoadEnvConfig("config", "application", "")

		localJson, _ := json.ToJson(config.GetMysqlConfig())

		fmt.Println("local配置json:", localJson)

		convey.So(localJson, convey.ShouldNotBeBlank)

		convey.Convey("测试加载测试环境env2", func() {

			config.LoadEnvConfig("config", "application", "test")

			testJson, _ := json.ToJson(config.GetMysqlConfig())

			fmt.Println("test配置json:", testJson)

			convey.So(testJson, convey.ShouldNotBeBlank)
		})

	})
}

func TestPinyin(t *testing.T) {

	convey.Convey("拼音", t, func() {

		hz := "龘"

		a := pinyin.NewArgs()

		fmt.Println("第一次pinyin的结果 :", pinyin.Pinyin(hz, a))

		convey.So(a, convey.ShouldNotBeNil)

		convey.Convey("拼音处理test2", func() {

			hz = "CallCenter"
			fmt.Println("第二次pinyin的结果", pinyin.Pinyin(hz, a))
			convey.So(pinyin.Pinyin(hz, a), convey.ShouldNotBeNil)
		})

	})
}
