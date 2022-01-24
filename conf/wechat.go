package conf

type wechatInfo struct {
	AppID     string
	AppSecret string
}

type wechatPay struct {
	MchID    string
	Apikeyv2 string
	Apikeyv3 string
	Serial   string
}
