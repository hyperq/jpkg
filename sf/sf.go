package sf

type SF struct {
	PartnerID string
	CheckWord string
	ApiUrl    string
	Debug     bool
}

func New(sfs SF) *SF {
	return &sfs
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
