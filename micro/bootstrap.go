package micro

import (
	"fmt"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"go-micro/config"
	"go-micro/core/cache"
	"go-micro/core/debug"
	"go-micro/core/log"
	"go-micro/core/model"
	"go-micro/rpc/client"
	"strconv"
	"time"
)

func Init(cfgPath string) {
	debug.SetPrintPrefix("[shop-micro][go-micro]")

	initConfig(cfgPath)

	initLog(Config.Log)

	initCache(Config.Cache)

	initModel(Config.Mysql)

	//loadValidator()

	//initRpcClient(Config.RpcClient)

	InitJaeger(Config.App.ServerName, Config.Jaeger.Address)

}

func initConfig(cfgPath string) {
	Viper = initViper(&Config, cfgPath)
}

func initLog(cfg *log.Config) {
	Logs = log.InitLogger(cfg)
}

func initCache(cfg *cache.Config) {
	// 初始化缓存
	if Config.Cache.Default == "freecache" {
		cache.CacheManager = cache.NewCache(cache.NewFreeCache(cfg))
	}
}

func initModel(cfg *model.Config) {
	// 初识mysql
	DB = model.InitDb(cfg)
}

func InitRpcClient(cfg config.RpcClient, opts ...client.DialOption) {

	//初始化rpc
	//global.RpcClient = client.NewClient(global.Config.RpcClient)

	debug.DD("开始rpc...")
	if len(cfg.Servers) > 0 {
		for k, v := range cfg.Servers {
			debug.DD("v = %v", v)
			opts = append(opts, client.SetServer(k, &client.Server{
				CertFile:      v.CertFile,
				TlsServerName: v.TlsServerName,
				NetWork:       v.Network,
				Address:       v.Address,
			}))
		}
	}

	RpcClient = client.NewClient(opts...)

}

func InitJaeger(service, address string) {
	var err error
	if service == "" {
		service = strconv.Itoa(int(time.Now().Unix()))
	}

	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
			//将span发送给jaeger-collector的服务中
			CollectorEndpoint: address,
		},
	}

	Jaefer, err = cfg.InitGlobalTracer(service, jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("Error: connect jaeger:%v \n", err))
	}
}
