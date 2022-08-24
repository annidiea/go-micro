package model

// 需要注意的是 yaml表中的内容 需要与 配置文件conf.yml中的内容对应
type Config struct {
	Dbname   string
	Host     string
	Port     int
	Username string
	Password string
	Charset  string
}

const DSN = "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True"
