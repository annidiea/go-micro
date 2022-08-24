package config

import (
	"go-micro/core/cache"
	"go-micro/core/log"
	"go-micro/core/model"
)

type Config struct {
	App
	Sms
	Smsbao
	Captche `mapstructure:"captcha"`
	Pay

	Jaeger
	RpcClient `mapstructure:"rpc_client"`
	RpcServer `mapstructure:"rpc_server"`

	Mysql *model.Config `mapstructure:"mysql"`
	Cache *cache.Config `mapstructure:"cache"`
	Log   *log.Config   `mapstructure:"log"`
}
