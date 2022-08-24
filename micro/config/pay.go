package config

type Pay struct {
	AppId        string `mapstructure:"app_id"`
	AliPublicKey string `mapstructure:"ali_public_key"`
	PrivateKey   string `mapstructure:"private_key"`
	NotifyURL    string `mapstructure:"notify_url"`
}
