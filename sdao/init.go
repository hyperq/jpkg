package sdao

import (
	"github.com/hyperq/jpkg/conf"
	"github.com/hyperq/jpkg/db/mssql"
	"time"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
)

var Msdb *mssql.DB

// Init load DB
func Init() {
	// mysql
	pkx = make(map[string]string)
	registerTypeTimeAsString()
	var err error

	if len(conf.Config.Database["mswrite"].DSN) != 0 {
		Msdb, err = mssql.New(
			mssql.Config{
				DSN:     conf.Config.Database["mswrite"].DSN[0],
				ReadDSN: conf.Config.Database["msread"].DSN,
				Active:  conf.Config.Database["msread"].Active,
				Idle:    conf.Config.Database["msread"].Idle,
			},
		)
		if err != nil {
			panic(err)
		}
		mssql.LogInit(conf.Config.Mongo.Type)
	}

}

const (
	formatDateTime = "2006-01-02 15:04:05"
)

type Time struct{}

func (wt *Time) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	ts := iter.ReadString()
	*((*time.Time)(ptr)), _ = time.ParseInLocation(formatDateTime, ts, time.Local)
}

func (wt *Time) IsEmpty(ptr unsafe.Pointer) bool {
	ts := *((*time.Time)(ptr))
	return ts.UnixNano() == 0
}

func (wt *Time) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	ts := *((*time.Time)(ptr))
	stream.WriteString(ts.Format(formatDateTime))
}

func registerTypeTimeAsString() {
	jsoniter.RegisterTypeDecoder("time.Time", &Time{})
	jsoniter.RegisterTypeEncoder("time.Time", &Time{})
}
