package rate

import (
	"github.com/hyperq/jpkg/log"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/config"
	"github.com/alibaba/sentinel-golang/core/system"
	"github.com/alibaba/sentinel-golang/logging"
)

func Init() {
	// 限流器
	conf := config.NewDefaultConfig()
	conf.Sentinel.Log.Logger = logging.NewConsoleLogger()
	conf.Sentinel.Log.Dir = "./logs/csp"
	err := sentinel.InitWithConfig(conf)
	if err != nil {
		panic(err)
	}
	_, err = system.LoadRules(
		[]*system.Rule{
			{
				MetricType:   system.Load,
				TriggerCount: 8.0,
				Strategy:     system.BBR,
			},
		},
	)
	if err != nil {
		log.Debug(err)
		return
	}
}
