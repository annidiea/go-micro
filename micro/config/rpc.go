package config

type RpcServer struct {
	CertFile string `mapstructure:"cert_file"`
	KeyFile  string `mapstructure:"key_file"`
}

type RpcClient struct {
	Servers map[string]Server `mapstructure:"servers"`
}

type Server struct {
	CertFile      string `mapstructure:"cert_file"`
	TlsServerName string `mapstructure:"tls_server_name"`
	Network       string `mapstructure:"network"`
	Address       string `mapstructure:"address"`
}
