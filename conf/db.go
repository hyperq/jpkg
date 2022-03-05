package conf

type redis struct {
	Addr     string
	Password string
}

type database struct {
	DSN    []string
	Active int
	Idle   int
}

type mongo struct {
	Uri      string
	Datebase string
	Type     int // 0 mongo 1 file
}
