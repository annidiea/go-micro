package config

type App struct {
	Address    string `yaml:"address"`
	ServerName string `mapstructure:"service_name"`
}
