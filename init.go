package jpkg

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/hyperq/jpkg/ali"
	"github.com/hyperq/jpkg/conf"
	"github.com/hyperq/jpkg/dao"
	"github.com/hyperq/jpkg/db/mongo"
	"github.com/hyperq/jpkg/express"
	"github.com/hyperq/jpkg/rate"
	"github.com/hyperq/jpkg/validator"
	"github.com/hyperq/jpkg/wechat"
)

func Init() {
	// binding init
	binding.Validator = new(validator.DefaultValidator)
	// config init
	conf.Init()
	// mongo init
	mongo.Init(conf.Config.Mongo.Uri)
	// aliyun
	ali.Init()
	// dao
	dao.Init()
	// 限流器
	if conf.Config.Rate {
		rate.Init()
	}
	//
	wechat.Init()
	// express
	express.Init()
}
