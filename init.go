package jpkg

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/hyperq/jpkg/ali"
	"github.com/hyperq/jpkg/conf"
	"github.com/hyperq/jpkg/dao"
	"github.com/hyperq/jpkg/db/mongo"
	"github.com/hyperq/jpkg/validator"
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
}
