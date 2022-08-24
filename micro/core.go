package micro

import (
	"github.com/mojocn/base64Captcha"
	"github.com/spf13/viper"
	"go-micro/config"
	"go-micro/rpc/client"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"io"
)

const ConfigFile = "./conf.yml"

var (
	Config config.Config
	Viper  *viper.Viper // 后面可能会对配置文件操作，可以通过它来实现
	Logs   *zap.Logger
	DB     *gorm.DB

	Jaefer io.Closer

	RpcClient client.RpcClient
)

var CaptchaStore = base64Captcha.DefaultMemStore

func Close() {
	Jaefer.Close()
}
