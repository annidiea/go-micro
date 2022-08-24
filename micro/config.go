package micro

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go-micro/config"
)

func initViper(config *config.Config, configPath string) *viper.Viper {
	// 创建 viper 配置文件的处理器
	v := viper.New() // 可以自动检测配置文件的变化而更新配置信息
	// 设置配置文件
	v.SetConfigFile(configPath)
	// 设置配置文件类型
	v.SetConfigType("yml")
	// 通过viper读取配置文件
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file : %s \n", err))
	}
	// 设置配置热加载
	v.WatchConfig()
	// 定义一个用于解析配置文件信息到global.Config结构体中的方法
	cfg := func(v *viper.Viper) {
		if err := v.Unmarshal(config); err != nil {
			panic(fmt.Errorf("Fatal error config parse : %s \n", err))
		}
	}
	// 在配置发生改变的时候触发的函数
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		cfg(v)
	})

	cfg(v)

	return v
}
