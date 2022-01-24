package dao

import (
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/hyperq/jpkg/cache"
	"github.com/hyperq/jpkg/conf"
	"github.com/hyperq/jpkg/log"
	"time"
	"unsafe"

	"github.com/hyperq/jpkg/db/mysql"

	"github.com/didi/gendry/scanner"

	jsoniter "github.com/json-iterator/go"
)

var db *mysql.DB

// Init load DB
func Init() {
	// mysql
	scanner.SetTagName("json")
	registerTypeTimeAsString()
	if len(conf.Config.Database["write"].DSN) == 0 {
		panic(errors.New("请配置数据库"))
	}
	var err error
	db, err = mysql.New(
		mysql.Config{
			DSN:     conf.Config.Database["write"].DSN[0],
			ReadDSN: conf.Config.Database["read"].DSN,
			Active:  conf.Config.Database["read"].Active,
			Idle:    conf.Config.Database["read"].Idle,
		},
	)
	if err != nil {
		panic(err)
	}
	// redis
	RC, err = cache.New(
		&redis.Options{
			Addr:     conf.Config.Redis.Addr,
			Password: conf.Config.Redis.Password,
			DB:       0,
		},
	)
	if err != nil {
		panic(err)
	}
	keys, err := RC.KEYS(conf.Config.AppName + "*")
	if err != nil {
		log.Error(err)
	}
	err = RC.DEL(keys...)
	if len(keys) > 0 {
		if err != nil {
			log.Error(err)
		}
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
