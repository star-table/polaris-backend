package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/spf13/cast"

	"github.com/go-kratos/kratos/v2/registry"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"

	"github.com/star-table/polaris-backend/facade/tablefacade"

	"github.com/star-table/common/library/mqtt/emt"

	"github.com/star-table/polaris-backend/service/platform/projectsvc/consume"

	"github.com/penglongli/gin-metrics/ginmetrics"

	"github.com/star-table/common/core/config"
	"github.com/star-table/common/core/logger"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/core/util/network"
	"github.com/star-table/common/library/discovery/nacos"
	"github.com/star-table/polaris-backend/common/core/buildinfo"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/extra/gin/mid"
	"github.com/star-table/polaris-backend/common/extra/gin/mvc"
	"github.com/star-table/polaris-backend/common/extra/trace/gin2micro"
	trace "github.com/star-table/polaris-backend/common/extra/trace/jaeger"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/api"
	"github.com/DeanThompson/ginpprof"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	kratosNacos "github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	"github.com/opentracing/opentracing-go"
)

var log = logger.GetDefaultLogger()
var env = ""
var build = false
var registerHost, registerPort, registerNamespace = "127.0.0.1", "8848", "public"

const BaseConfigPath = "./../../../config"
const SelfConfigPath = "./config"

func init() {
	env = os.Getenv(consts.RunEnvKey)
	if "" == env {
		env = consts.RunEnvLocal
	}
	//配置
	flag.BoolVar(&build, "build", false, "build facade")
	flag.StringVar(&env, "env", env, "env")
	flag.StringVar(&registerHost, "registerHost", "172.19.98.19", "registerHost")
	flag.StringVar(&registerPort, "registerPort", "30048", "registerPort")
	flag.StringVar(&registerNamespace, "registerNamespace", "public", "registerNamespace")
	flag.Parse()

	if os.Getenv(consts.REGISTER_HOST) == "" {
		_ = os.Setenv(consts.REGISTER_HOST, registerHost)
	} else {
		registerHost = os.Getenv(consts.REGISTER_HOST)
	}
	if os.Getenv(consts.REGISTER_PORT) == "" {
		_ = os.Setenv(consts.REGISTER_PORT, registerPort)
	} else {
		registerPort = os.Getenv(consts.REGISTER_PORT)
	}
	if os.Getenv(consts.REGISTER_NAMESPACE) == "" {
		_ = os.Setenv(consts.REGISTER_NAMESPACE, registerNamespace)
	} else {
		registerNamespace = os.Getenv(consts.REGISTER_NAMESPACE)
	}

	//配置文件
	if env != consts.RunEnvLocal && env != consts.RunEnvTest && env != "fuse_k8s" {
		err := config.LoadNacosConfigAutoConfiguration("project", env)
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

func main() {
	// 打印程序信息
	fmt.Println(buildinfo.StringifySingleLine())
	fmt.Println(buildinfo.StringifyMultiLine())

	rand.Seed(time.Now().UnixNano())

	port := config.GetServerConfig().Port
	host := config.GetServerConfig().Host

	applicationName := config.GetApplication().Name

	msg := json.ToJsonIgnoreError(config.GetConfig())

	fmt.Println("config配置:" + msg)

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
	// 初始化 sentry
	// 必须要在调用 log.xxx() 之前进行初始化 sentry。因为后续要将 sentry 实例作为一个配置放到 log 包中。
	oriErr := sentry.Init(sentry.ClientOptions{
		Dsn:         sentryDsn,
		ServerName:  applicationName,
		Environment: env,
	})
	if oriErr != nil {
		// sentry 不是必要的服务，因此异常时，业务服务还是会启动
		log.Error(oriErr)
	} else {
		log.SetExtraLoggerOption("sentryClient", sentry.CurrentHub().Client())
		log.Info("init sentryClient ok")
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
	//r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(mid.SentryMiddleware(applicationName, env, sentryDsn))
	r.Use(mid.StartTrace())
	r.Use(mid.GinContextToContextMiddleware())
	r.Use(mid.CorsMiddleware())
	r.Use(mid.AuthMiddleware())

	version := ""
	postGreeter := api.PostGreeter{Greeter: mvc.NewPostGreeter(applicationName, host, port, version)}
	getGreeter := api.GetGreeter{Greeter: mvc.NewGetGreeter(applicationName, host, port, version)}

	//build
	if build {
		facadeBuilder := mvc.FacadeBuilder{
			StorageDir: "./../../../facade/projectfacade",
			Package:    "projectfacade",
			VoPackage:  "projectvo",
			Greeters:   []interface{}{&postGreeter, &getGreeter},
		}
		facadeBuilder.Build()
		return
	}

	// 多库库模式才会执行
	//if env != consts.RunEnvLocal && env != consts.RunEnvTest {
	//	if (consts.AppRunmodeSaas == config.GetApplication().RunMode) || (consts.AppRunmodePrivate == config.GetApplication().RunMode) {
	//		mysqlConfig := config.GetMysqlConfig()
	//		initErr := db.DbMigrations(env, mysqlConfig.Host, mysqlConfig.Port, mysqlConfig.Usr, mysqlConfig.Pwd, mysqlConfig.Database)
	//		if initErr != nil {
	//			panic(" init db fail....")
	//		}
	//	}
	//}

	//启动nacos
	nacos.Init()

	//启动MQTT
	emt.Init()

	discover := newDiscovery()
	tablefacade.InitGrpcClient(discover)

	//启动mq消费者
	if env != consts.RunEnvLocal {
		go consume.IssueTrendsAndNoticeConsume()
		go consume.BatchCreateIssueConsume()
		go consume.DailyProjectReportMsgConsumer() // 项目日报
		go consume.DailyIssueReportMsgConsumer()   // 个人日报
		go consume.IssueRemindConsumer()           // 任务截止日期提醒
	}

	ginHandler := mvc.NewGinHandler(r)
	ginHandler.RegisterGreeter(&postGreeter)
	ginHandler.RegisterGreeter(&getGreeter)

	if env != consts.RunEnvNull {
		log.Info("开启pprof监控")
		ginpprof.Wrap(r)
	}

	log.Infof("POL_ENV:%s, connect to http://%s:%d/ for %s service", env, network.GetIntranetIp(), port, applicationName)
	r.Run(":" + strconv.Itoa(port))
}

func newDiscovery() registry.Discovery {
	sc, cc := getNacosServerAndClientConfig()
	client, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ServerConfigs: sc,
			ClientConfig:  &cc,
		},
	)
	if err != nil {
		log.Error(err)
		panic(err)
		return nil
	}

	r := kratosNacos.New(client)

	return r
}

func getNacosServerAndClientConfig() ([]constant.ServerConfig, constant.ClientConfig) {
	return []constant.ServerConfig{
			*constant.NewServerConfig(registerHost, cast.ToUint64(registerPort)),
		},
		constant.ClientConfig{
			NamespaceId:         config.GetConfig().Nacos.Client.NamespaceId, //namespace id
			TimeoutMs:           5000,
			NotLoadCacheAtStart: true,
			LogLevel:            "error",
		}
}
