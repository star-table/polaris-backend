package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/gin-contrib/gzip"

	"github.com/star-table/common/core/config"
	"github.com/star-table/common/core/logger"
	"github.com/star-table/common/core/util/network"
	"github.com/star-table/polaris-backend/app/server/handler"
	"github.com/star-table/polaris-backend/app/server/routes"
	"github.com/star-table/polaris-backend/common/core/buildinfo"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/extra/gin/mid"
	"github.com/star-table/polaris-backend/common/extra/trace/gin2micro"
	trace "github.com/star-table/polaris-backend/common/extra/trace/jaeger"
	"github.com/DeanThompson/ginpprof"
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/penglongli/gin-metrics/ginmetrics"
)

var log = logger.GetDefaultLogger()
var env = ""
var build = false
var registerHost, registerPort, registerNamespace = "127.0.0.1", "8848", "public"

const BaseConfigPath = "./../config"
const SelfConfigPath = "./config"

func init() {
	env = os.Getenv(consts.RunEnvKey)
	if "" == env {
		env = consts.RunEnvLocal
	}
	//配置
	flag.BoolVar(&build, "build", false, "build facade")
	flag.StringVar(&env, "env", env, "env")
	flag.StringVar(&registerHost, "registerHost", "127.0.0.1", "registerHost")
	flag.StringVar(&registerPort, "registerPort", "8848", "registerPort")
	flag.StringVar(&registerNamespace, "registerNamespace", "public", "registerNamespace")
	flag.Parse()

	if os.Getenv(consts.REGISTER_HOST) == "" {
		_ = os.Setenv(consts.REGISTER_HOST, registerHost)
	}
	if os.Getenv(consts.REGISTER_PORT) == "" {
		_ = os.Setenv(consts.REGISTER_PORT, registerPort)
	}
	if os.Getenv(consts.REGISTER_NAMESPACE) == "" {
		_ = os.Setenv(consts.REGISTER_NAMESPACE, registerNamespace)
	}

	if env == consts.RunEnvGray {
		err := config.LoadNacosConfigAutoConfiguration("app", env)
		if err != nil {
			panic(err)
		}
	} else {
		if runtime.GOOS != consts.LinuxGOOS {
			config.LoadEnvConfig(BaseConfigPath, "application.common", env)
			config.LoadEnvConfig(SelfConfigPath, "application", env)
		} else {
			if env == "test" {
				config.LoadEnvConfig(BaseConfigPath, "application.common", env)
				config.LoadEnvConfig(SelfConfigPath, "application", env)
			} else {
				config.LoadEnvConfig(SelfConfigPath, "application.common", env)
				config.LoadEnvConfig(SelfConfigPath, "application", env)
			}
		}
	}
}

// @title Polaris Apis
// @version v1.0.0
// @description 极星接口文档.
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name token
// @BasePath /
func main() {
	// 打印程序信息
	log.Info(buildinfo.StringifySingleLine())
	fmt.Println(buildinfo.StringifyMultiLine())

	serverConfig := config.GetServerConfig()
	port := strconv.Itoa(serverConfig.Port)

	// 单库模式才会执行
	//if (consts.AppRunmodePrivateSingleDb == config.GetApplication().RunMode) || (consts.AppRunmodeSingle == config.GetApplication().RunMode) {
	//	mysqlConfig := config.GetMysqlConfig()
	//	initErr := db.DbMigrations(env, mysqlConfig.Host, mysqlConfig.Port, mysqlConfig.Usr, mysqlConfig.Pwd, mysqlConfig.Database)
	//	if initErr != nil {
	//		panic(" init db fail....")
	//	}
	//}

	r := gin.New()

	// Metrics
	m := ginmetrics.GetMonitor()
	m.SetMetricPath("/prometheus")
	m.SetSlowTime(5)
	// +optional set request duration, default {0.1, 0.3, 1.2, 5, 10} used to p95, p99
	m.SetDuration([]float64{0.01, 0.05, 0.1, 0.2, 0.5, 1, 5})
	m.Use(r)

	sentryConfig := config.GetSentryConfig()
	sentryDsn := ""
	if sentryConfig != nil {
		sentryDsn = sentryConfig.Dsn
	}

	if config.GetJaegerConfig() != nil {
		t, io, err := trace.NewTracer(config.GetJaegerConfig())
		if err != nil {
			log.Infof("err %v", err)
		}
		defer func() {
			if err := io.Close(); err != nil {
				log.Errorf("err %v", err)
			}
		}()
		opentracing.SetGlobalTracer(t)

		r.Use(gin2micro.TracerWrapper)
	}
	applicationName := config.GetApplication().Name
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(mid.SentryMiddleware(applicationName, env, sentryDsn))
	r.Use(mid.StartTrace())
	r.Use(mid.GinContextToContextMiddleware())
	r.Use(mid.CorsMiddleware())
	r.Use(mid.AuthMiddleware())
	r.Use(handler.RestHandler())
	// 设定处理请求的中间件
	// Jaeger defer会close掉导致失效
	//routes.SetMiddlewares(r, env)

	captcha.SetCustomStore(&handler.RedisCache{})
	// 设定路由处理请求
	routes.SetRoutes(r)

	if env != consts.RunEnvNull {
		log.Info("开启pprof监控")
		ginpprof.Wrap(r)
	}

	log.Infof("POL_ENV:%s, connect to http://%s:%s/ for GraphQL playground", env, network.GetIntranetIp(), port)
	r.Run(":" + port)
}
