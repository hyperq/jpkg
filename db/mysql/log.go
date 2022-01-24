package mysql

import (
	"fmt"
	"github.com/hyperq/jpkg/db/mongo"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var errorLogger *zap.SugaredLogger

var levelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

func getLoggerLevel(lvl string) zapcore.Level {
	if level, ok := levelMap[lvl]; ok {
		return level
	}
	return zapcore.InfoLevel
}

// GormSQLSql log init
// var GormSQLLog *gormzap.Logger
var sqlLogger *zap.Logger

// SetSQLLog  set sql logto file
func init() {
	// fileName := "logs/sql.log"
	level := getLoggerLevel("debug")
	// syncWriter := zapcore.AddSync(&lumberjack.Logger{
	// 	Filename: fileName,
	// 	MaxSize:  20,
	// 	// LocalTime: true,
	// 	Compress: true,
	// })

	runMode := gin.Mode()
	var encoder zapcore.EncoderConfig
	if runMode == "debug" {
		encoder = zap.NewDevelopmentEncoderConfig()
	} else {
		encoder = zap.NewProductionEncoderConfig()
		encoder.EncodeTime = zapcore.EpochTimeEncoder
	}
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoder), zapcore.NewMultiWriteSyncer(zapcore.AddSync(mongo.NewLog("sql"))), zap.NewAtomicLevelAt(level),
	)
	sqlLogger = zap.New(core)
}

// debugLogQueies debug log
func debugLogQueies(query string, t time.Time, err error, args ...interface{}) {
	l := new(logs)
	// 0 1 query 2 dao 3
	_, file, line, ok := runtime.Caller(3)
	if ok {
		if strings.Index(file, "dao") == -1 {
			l.source = file + ":" + strconv.Itoa(line)
		} else {
			_, file, line, ok = runtime.Caller(4)
			if ok {
				l.source = file + ":" + strconv.Itoa(line)
			}
		}
	}
	l.duration = fmt.Sprintf("%dms", time.Now().Sub(t)/1e6)
	l.sql = strings.Replace(strings.Replace(query, "\n", "", -1), "\t", "", -1)
	l.values = getFormattedValues(args)
	if err != nil {
		l.err = fmt.Sprint(err)
		sqlLogger.Error("db", l.toZapFields()...)
	} else {
		sqlLogger.Info("db", l.toZapFields()...)
	}
}

func getFormattedValues(args []interface{}) string {
	formattedValues := make([]string, 0, len(args))
	for _, v := range args {
		str := "NULL"
		if v != nil {
			str = fmt.Sprint(v)
		}
		formattedValues = append(formattedValues, str)
	}
	return "[ " + strings.Join(formattedValues, " , ") + " ]"
}

type logs struct {
	source   string
	duration string
	sql      string
	values   string
	err      string
}

func (l *logs) toZapFields() []zapcore.Field {
	return []zapcore.Field{
		zap.String("duration", l.duration),
		zap.String("sql", l.sql),
		zap.String("values", l.values),
		zap.String("error", l.err),
		zap.String("source", l.source),
	}
}
