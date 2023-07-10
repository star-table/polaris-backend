package domain

import (
	"fmt"
	"github.com/star-table/common/core/config"
	"os"
)

var env = ""

const BaseConfigPath = "./../../../../config"
const SelfConfigPath = "./../config"

func init() {
	env = os.Getenv("POL_ENV")
	if "" == env {
		env = "unittest"
	}
	//配置文件
	err := config.LoadEnvConfig(BaseConfigPath, "application.common", env)

	if err != nil {
		fmt.Printf("err:%s\n", err)
	}

	err = config.LoadEnvConfig(SelfConfigPath, "application", "")

	if err != nil {
		fmt.Printf("err:%s\n", err)
	}
}
