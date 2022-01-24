package conf

import (
	"flag"

	"github.com/BurntSushi/toml"
)

type config struct {
	AppName   string
	Url       string
	Port      string
	Mode      string
	Database  map[string]database
	Mongo     mongo
	Redis     redis
	Wechat    map[string]wechatInfo
	WechatPay wechatPay
	Alipay    alipay
	Oss       map[string]oss
}

// Config global config
var Config config

var (
	tomlFile = flag.String("config", "config/app.toml", "config file")
)

// Init load config
func Init() {
	Config = config{}
	flag.Parse()
	_, err := toml.DecodeFile(*tomlFile, &Config)
	if err != nil {
		panic(err)
	}
}
