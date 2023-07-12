package main

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/star-table/common/library/discovery/nacos"

	"github.com/gin-contrib/gzip"

	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/penglongli/gin-metrics/ginmetrics"
	"github.com/star-table/common/core/config"
	"github.com/star-table/common/core/logger"
	"github.com/star-table/common/core/util/network"
	"github.com/star-table/polaris-backend/common/core/buildinfo"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/extra/gin/mid"
	"github.com/star-table/polaris-backend/common/extra/gin/mvc"
	"github.com/star-table/polaris-backend/common/extra/trace/gin2micro"
	trace "github.com/star-table/polaris-backend/common/extra/trace/jaeger"
	"github.com/star-table/polaris-backend/service/basic/idsvc/api"
)

var (
	log   = logger.GetDefaultLogger()
	build = false
	env   = ""
	name  = "idsvc"

	flagconf                             string
	nacosHost, nacosPort, nacosNamespace string
)

func init() {
	//配置说明：会优先读取 -conf的配置，如果没有传入，则读nacos配置
	flag.StringVar(&env, "env", "", "eg: -env test")
	flag.StringVar(&flagconf, "conf", "", "config path, eg: -conf /dir/test/config.yaml")
	flag.StringVar(&nacosHost, "register_host", "", "eg: -register_host 127.0.0.1")
	flag.StringVar(&nacosPort, "register_port", "", " eg: -register_port 33089 ")
	flag.StringVar(&nacosNamespace, "register_namespace", "", "eg: -register_namespace lesscode")
	flag.BoolVar(&build, "build", false, "build facade")

	flag.Parse()
	err := config.LoadConfig(flagconf, nacosHost, nacosPort, nacosNamespace, name)
	if err != nil {
		panic(err)
	}
}

func main() {
	// 打印程序信息
	log.Info(buildinfo.StringifySingleLine())
	fmt.Println(buildinfo.StringifyMultiLine())

	port := config.GetServerConfig().Port
	host := config.GetServerConfig().Host

	applicationName := config.GetApplication().Name

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
	r.Use(gzip.Gzip(gzip.DefaultCompression))
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
			StorageDir: "./../../../facade/idfacade",
			Package:    "idfacade",
			VoPackage:  "idvo",
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
