package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/spf13/cast"

	"github.com/go-kratos/kratos/v2/registry"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"

	"github.com/star-table/polaris-backend/facade/tablefacade"

	"github.com/star-table/polaris-backend/service/platform/projectsvc/consume"

	"github.com/penglongli/gin-metrics/ginmetrics"

	"github.com/DeanThompson/ginpprof"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	kratosNacos "github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	"github.com/opentracing/opentracing-go"
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
)

var (
	log   = logger.GetDefaultLogger()
	build = false
	env   = ""
	name  = "projectsvc"

	flagconf                             string
	nacosHost, nacosPort, nacosNamespace string
)

func init() {
	//配置
	flag.StringVar(&env, "env", "", "eg: -env test")
	flag.StringVar(&flagconf, "conf", "", "config path, eg: -conf ../test/config.yaml")
	flag.StringVar(&nacosHost, "register_host", "", "eg: -register_host 127.0.0.1")
	flag.StringVar(&nacosPort, "register_port", "", " eg: -register_port 33089 ")
	flag.StringVar(&nacosNamespace, "register_namespace", "", "eg: -register_namespace lesscode")
	flag.BoolVar(&build, "build", false, "build facade")

	err := config.LoadConfig(flagconf, nacosHost, nacosPort, nacosNamespace, name)
	if err != nil {
		panic(err)
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
	nacos.Init(&config.NacosBaseConfig{
		AppName:   name,
		Host:      nacosHost,
		Port:      nacosPort,
		NameSpace: nacosNamespace,
		Group:     "DEFAULT_GROUP",
	})

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
			*constant.NewServerConfig(nacosHost, cast.ToUint64(nacosPort)),
		},
		constant.ClientConfig{
			NamespaceId:         config.GetConfig().Nacos.Client.NamespaceId, //namespace id
			TimeoutMs:           5000,
			NotLoadCacheAtStart: true,
			LogLevel:            "error",
		}
}
