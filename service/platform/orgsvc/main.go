package main

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/DeanThompson/ginpprof"
	"github.com/getsentry/sentry-go"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
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
	"github.com/star-table/polaris-backend/service/platform/orgsvc/api"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/api/roleapi"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/consume"
)

var (
	log   = logger.GetDefaultLogger()
	build = false
	env   = ""
	name  = "orgsvc"

	flagconf                             string
	nacosHost, nacosPort, nacosNamespace string
)

func init() {
	//配置说明：会优先读取 -conf的配置，如果没有传入，则读nacos配置
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

	port := config.GetServerConfig().Port
	host := config.GetServerConfig().Host

	applicationName := config.GetApplication().Name
	msg := json.ToJsonIgnoreError(config.GetConfig())
	fmt.Printf("env: %s, config配置: %s", env, msg)

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

	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(mid.SentryMiddleware(applicationName, env, sentryDsn))
	r.Use(mid.StartTrace())
	r.Use(mid.GinContextToContextMiddleware())
	r.Use(mid.CorsMiddleware())
	r.Use(mid.AuthMiddleware())

	version := ""
	postGreeter := api.PostGreeter{Greeter: mvc.NewPostGreeter(applicationName, host, port, version)}
	getGreeter := api.GetGreeter{Greeter: mvc.NewGetGreeter(applicationName, host, port, version)}

	rolePostGreeter := roleapi.PostGreeter{Greeter: mvc.NewPostGreeter("rolesvc", host, port, version)}
	roleGetGreeter := roleapi.GetGreeter{Greeter: mvc.NewGetGreeter("rolesvc", host, port, version)}

	//build
	if build {
		facadeBuilder := mvc.FacadeBuilder{
			StorageDir: "./../../../facade/orgfacade",
			Package:    "orgfacade",
			VoPackage:  "orgvo",
			Greeters:   []interface{}{&postGreeter, &getGreeter},
		}
		facadeBuilder.Build()

		roleFacadeBuilder := mvc.FacadeBuilder{
			StorageDir: "./../../../facade/orgfacade",
			Package:    "orgfacade",
			VoPackage:  "rolevo",
			Greeters:   []interface{}{&rolePostGreeter, &roleGetGreeter},
		}
		roleFacadeBuilder.Build()
		return
	}

	// 多库库模式才会执行
	//if env != consts.RunEnvLocal && env != consts.RunEnvTest {
	//	if (consts.AppRunmodeSaas == config.GetApplication().RunMode) || (consts.AppRunmodePrivate == config.GetApplication().RunMode) {
	//		log.Infof("[startMigration] start:%v", time.Now().Unix())
	//		mysqlConfig := config.GetMysqlConfig()
	//		initErr := db.DbMigrations(env, mysqlConfig.Host, mysqlConfig.Port, mysqlConfig.Usr, mysqlConfig.Pwd, mysqlConfig.Database)
	//		if initErr != nil {
	//			panic(" init db fail....")
	//		}
	//		log.Infof("[startMigration] end:%v", time.Now().Unix())
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

	if env != consts.RunEnvLocal {
		go consume.OrgMemberChangeConsume()
		go consume.OrgMemberSyncConsumer()
	}

	ginHandler := mvc.NewGinHandler(r)
	ginHandler.RegisterGreeter(&postGreeter)
	ginHandler.RegisterGreeter(&getGreeter)
	ginHandler.RegisterGreeter(&rolePostGreeter)
	ginHandler.RegisterGreeter(&roleGetGreeter)

	if env != consts.RunEnvNull {
		log.Info("开启pprof监控")
		ginpprof.Wrap(r)
	}

	log.Infof("POL_ENV:%s, connect to http://%s:%d/ for %s service", env, network.GetIntranetIp(), port, applicationName)
	r.Run(":" + strconv.Itoa(port))
}
