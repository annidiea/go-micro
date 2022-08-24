package client

type Config struct {
	//记录服务的地址
	Servers map[string]ServerCfg `server`
}

type ServerCfg struct {
	Network string
	Address string
}
