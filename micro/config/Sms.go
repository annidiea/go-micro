package config

//api: http://api.smsbao.com/
//user: jinmin
//pass: www.jinmin.com
//statusStr:

type Sms struct {
	Service string `yaml:"service"`
	Long    int    `yaml:"long"`
	Overdue int    `yaml:"overdue"`

	Temp struct {
		Code string `yaml:"code"`
	} `yaml:"temp"`
}

type Smsbao struct {
	Api       string            `yaml:"api"`
	User      string            `yaml:"user"`
	Pass      string            `yaml:"pass"`
	StatusStr map[string]string `yaml:"statusStr"`
}
