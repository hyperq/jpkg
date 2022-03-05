package log

import (
	"github.com/hyperq/jpkg/db/mongo"
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// errorLogger
var errorLogger *zap.SugaredLogger
var errorLogger2 *zap.SugaredLogger

func Init(t int) {
	errorLogger = newlog(1, t)
	errorLogger2 = newlog(2, t)
}

func newlog(skip int, t int) *zap.SugaredLogger {
	level := zapcore.DebugLevel
	runMode := gin.Mode()
	var encoder zapcore.EncoderConfig
	if runMode == "debug" {
		encoder = zap.NewDevelopmentEncoderConfig()
	} else {
		encoder = zap.NewProductionEncoderConfig()
		encoder.EncodeTime = zapcore.EpochTimeEncoder
	}
	var eng io.Writer
	if t == 0 {
		eng = mongo.New("log")
	} else {
		eng = &lumberjack.Logger{
			Filename: "logs/backend.log",
			MaxSize:  20,
			// LocalTime: true,
			Compress: true,
		}
	}
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoder),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(eng)),
		zap.NewAtomicLevelAt(level),
	)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(skip))
	return logger.Sugar()
}

// Debug Debug
func Debug(args ...interface{}) {
	errorLogger.Debug(args...)
}

// Debugf Debugf
func Debugf(template string, args ...interface{}) {
	errorLogger.Debugf(template, args...)
}

// Info Info
func Info(args ...interface{}) {
	errorLogger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	errorLogger.Infof(template, args...)
}

func Warn(args ...interface{}) {
	errorLogger.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	errorLogger.Warnf(template, args...)
}

// Error Error
func Error(args ...interface{}) {
	errorLogger.Error(args...)
}

// Error Error
func Error2(args ...interface{}) {
	errorLogger2.Error(args...)
}

// Errorf Errorf
func Errorf(template string, args ...interface{}) {
	errorLogger.Errorf(template, args...)
}

// DPanic DPanic
func DPanic(args ...interface{}) {
	errorLogger.DPanic(args...)
}

func DPanicf(template string, args ...interface{}) {
	errorLogger.DPanicf(template, args...)
}

// Panic Panic
func Panic(args ...interface{}) {
	errorLogger.Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	errorLogger.Panicf(template, args...)
}

// Fatal Fatal
func Fatal(args ...interface{}) {
	errorLogger.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	errorLogger.Fatalf(template, args...)
}
