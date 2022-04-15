package wechat

import (
	"github.com/hyperq/jpkg/conf"

	"github.com/silenceper/wechat/v2/miniprogram/subscribe"

	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/miniprogram"
	mpcnf "github.com/silenceper/wechat/v2/miniprogram/config"
	"github.com/silenceper/wechat/v2/officialaccount"
	oacnf "github.com/silenceper/wechat/v2/officialaccount/config"
)

var Mp *miniprogram.MiniProgram
var Oa *officialaccount.OfficialAccount

func Init() {
	Cache := cache.NewRedis(
		&cache.RedisOpts{
			Host:     conf.Config.Redis.Addr,
			Password: conf.Config.Redis.Password,
			Database: 10,
		},
	)
	if mp, ok := conf.Config.Wechat["miniapp"]; ok {
		Mp = wechat.NewWechat().GetMiniProgram(
			&mpcnf.Config{
				AppID:     mp.AppID,
				AppSecret: mp.AppSecret,
				Cache:     Cache,
			},
		)
	}
	if oa, ok := conf.Config.Wechat["oa"]; ok {
		Oa = wechat.NewWechat().GetOfficialAccount(
			&oacnf.Config{
				AppID:          oa.AppID,
				AppSecret:      oa.AppSecret,
				Token:          oa.Token,
				EncodingAESKey: oa.EncodingAESKey,
				Cache:          Cache,
			},
		)
	}
}

func SendUMessage(template_id, openid, url, mppath string, data map[string]*subscribe.DataItem) (err error) {
	msg := new(subscribe.UniformMessage)
	msg.ToUser = openid
	msg.MpTemplateMsg.Appid = conf.Config.Wechat["oa"].AppID
	msg.MpTemplateMsg.Miniprogram.Appid = conf.Config.Wechat["miniapp"].AppID
	msg.MpTemplateMsg.TemplateID = template_id
	msg.MpTemplateMsg.URL = url
	msg.MpTemplateMsg.Miniprogram.Pagepath = mppath
	msg.MpTemplateMsg.Data = data
	err = Mp.GetSubscribe().UniformSend(msg)
	return
}
