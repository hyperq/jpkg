package sf

import "github.com/hyperq/jpkg/conf"

type SFM struct {
	PartnerID string
	CheckWord string
	ApiUrl    string
	Debug     bool
}

var SF *SFM

func Init() *SFM {
	return &SFM{
		PartnerID: conf.Config.Sf.PartnerID,
		CheckWord: conf.Config.Sf.CheckWord,
		ApiUrl:    conf.Config.Sf.ApiUrl,
		Debug:     conf.Config.Sf.Debug,
	}
}

const (
	NULL     = ""
	SUCCESS  = "SUCCESS"
	FAIL     = "FAIL"
	OK       = "OK"
	DebugOff = 0
	DebugOn  = 1
	Version  = "1.5.58"
)
