package log

import (
	"github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"runtime/debug"
	"time"
)

var Logger *zap.Logger
var Sugar *zap.SugaredLogger

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	//loc, err := time.LoadLocation("Asia/Shanghai")
	//if err != nil {
	//	Errorf("time load location [Asia/Shanghai] fail %v", err)
	//	loc = time.FixedZone("CST", 8*3600)
	//}
	//enc.AppendString(t.In(loc).Format("2006-01-02 15:04:05.000"))
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func InitZapLog(filename string, maxSize int, maxBackups int, maxAge int) {
	// >>> 按文件大小滚动
	//fileWriterSyncer := zapcore.AddSync(&lumberjack.Logger{
	//	Filename:   filename,
	//	MaxSize:    maxSize, //MB
	//	LocalTime:  true,
	//	MaxBackups: maxBackups, // number of log files
	//	MaxAge:     maxAge,     // days
	//
	//})
	// >>> 按日期滚动
	rotate, err := RotateLogs(filename)
	if err != nil {
		panic(err)
	}
	fileWriterSyncer := zapcore.AddSync(rotate)
	fileEncoderConfig := zap.NewProductionEncoderConfig()
	fileEncoderConfig.EncodeTime = TimeEncoder // zapcore.ISO8601TimeEncoder
	fileEncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	devEncoderConfig := zap.NewDevelopmentEncoderConfig()
	devEncoderConfig.EncodeTime = TimeEncoder
	devEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // color

	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewConsoleEncoder(fileEncoderConfig), fileWriterSyncer, zap.NewAtomicLevel()),
		zapcore.NewCore(zapcore.NewConsoleEncoder(devEncoderConfig), zapcore.WriteSyncer(os.Stdout), zap.NewAtomicLevel()),
	)

	core.Enabled(zapcore.DebugLevel)
	Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	Sugar = Logger.Sugar()
}

//按日期切割
func RotateLogs(filePath string) (*rotatelogs.RotateLogs, error) {
	filename := filePath + ".%Y%m%d"
	retate, err := rotatelogs.New(filename, rotatelogs.WithLinkName(filePath), rotatelogs.WithMaxAge(time.Hour*24*3), rotatelogs.WithRotationTime(time.Hour*24))
	return retate, err
}

func Debug(args ...interface{}) {
	Sugar.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	Sugar.Debugf(template, args...)
}

func Info(args ...interface{}) {
	Sugar.Info(args...)
}

func Infof(template string, args ...interface{}) {
	Sugar.Infof(template, args...)
}

func Warn(args ...interface{}) {
	Sugar.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	Sugar.Warnf(template, args...)
}

func Error(args ...interface{}) {
	args = append(args, string(debug.Stack()))
	Sugar.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	args = append(args, string(debug.Stack()))
	Sugar.Errorf(template+"\n", args...)
}

func Panic(args ...interface{}) {
	args = append(args, string(debug.Stack()))
	Sugar.Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	args = append(args, string(debug.Stack()))
	Sugar.Panicf(template, args...)
}

func Fatal(args ...interface{}) {
	args = append(args, string(debug.Stack()))
	Sugar.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	args = append(args, string(debug.Stack()))
	Sugar.Fatalf(template, args...)
}