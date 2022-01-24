package log

import (
	"github.com/hyperq/jpkg/db/mongo"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// errorLogger
var errorLogger *zap.SugaredLogger
var errorLogger2 *zap.SugaredLogger

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

func init() {
	errorLogger = newlog(1)
	errorLogger2 = newlog(2)
}

func newlog(skip int) *zap.SugaredLogger {
	level := getLoggerLevel("debug")
	runMode := gin.Mode()
	var encoder zapcore.EncoderConfig
	if runMode == "debug" {
		encoder = zap.NewDevelopmentEncoderConfig()
	} else {
		encoder = zap.NewProductionEncoderConfig()
		encoder.EncodeTime = zapcore.EpochTimeEncoder
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoder), zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(mongo.NewLog("log"))),
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

// Infof Infof
func Infof(template string, args ...interface{}) {
	errorLogger.Infof(template, args...)
}

// Warn Warn
func Warn(args ...interface{}) {
	errorLogger.Warn(args...)
}

// Warnf Warnf
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

// DPanicf DPanicf
func DPanicf(template string, args ...interface{}) {
	errorLogger.DPanicf(template, args...)
}

// Panic Panic
func Panic(args ...interface{}) {
	errorLogger.Panic(args...)
}

// Panicf Panicf
func Panicf(template string, args ...interface{}) {
	errorLogger.Panicf(template, args...)
}

// Fatal Fatal
func Fatal(args ...interface{}) {
	errorLogger.Fatal(args...)
}

// Fatalf Fatalf
func Fatalf(template string, args ...interface{}) {
	errorLogger.Fatalf(template, args...)
}
