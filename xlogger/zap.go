package xlogger

import (
	"io"
	"os"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LogConf 日志信息结构体
type LogConf struct {
	Filename string `json:"filename" yaml:"filename"`

	LogLevel string `json:"level" yaml:"level"`

	IsStdOut bool `json:"stdout" yaml:"stdout"`

	// MaxSize is the maximum size in megabytes of the logger file before it gets
	// rotated. It defaults to 100 megabytes.
	MaxSize int `json:"max_size" yaml:"max_size"`

	// MaxAge is the maximum number of days to retain old logger files based on the
	// timestamp encoded in their filename.  Note that a day is defined as 24
	// hours and may not exactly correspond to calendar days due to daylight
	// savings, leap seconds, etc. The default_xorm is not to remove old logger files
	// based on age.
	MaxAge int `json:"max_age" yaml:"max_age"`

	// MaxBackups is the maximum number of old logger files to retain.  The default_xorm
	// is to retain all old logger files (though MaxAge may still cause them to get
	// deleted.)
	MaxBackups int `json:"max_backups" yaml:"max_backups"`

	// LocalTime determines if the time used for formatting the timestamps in
	// backup files is the computer's local time.  The default_xorm is to use UTC
	// time.
	LocalTime bool `json:"localtime" yaml:"localtime"`

	// Compress determines if the rotated logger files should be compressed
	// using gzip.
	Compress bool `json:"compress" yaml:"compress"`

	RotateType string `json:"rotate_type" yaml:"rotate_type"`
}

//生成时间格式
func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func NewZap(conf *LogConf, skipNum int, writers ...io.Writer) (*zap.SugaredLogger, error) {

	iow := make([]zapcore.WriteSyncer, 0, len(writers)+1)
	if len(writers) > 0 {
		for _, value := range writers {
			iow = append(iow, zapcore.AddSync(value))
		}
	}
	// var encoderConfig zapcore.EncoderConfig
	// var encoder zapcore.Encoder

	// if runEnv == "debug" {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeTime = TimeEncoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	// } else {
	// 	encoderConfig = zap.NewProductionEncoderConfig()
	// 	encoderConfig.EncodeTime = TimeEncoder
	// 	encoder = zapcore.NewJSONEncoder(encoderConfig)
	// }
	if conf.IsStdOut {
		iow = append(iow, zapcore.AddSync(os.Stdout))
	}
	writer := zapcore.NewMultiWriteSyncer(iow...)
	core := zapcore.NewCore(
		encoder,
		writer,
		parseLevel(conf.LogLevel),
	)

	//add options in here
	// options := []zap.Option{zap.AddCaller()}
	// options = append(options, zap.AddCallerSkip(skipNum+1))

	// logger := zap.New(core, options...)
	logger := zap.New(core)
	zap.RedirectStdLog(logger)
	return logger.Sugar(), nil
}

func NewRotateWriter(conf *LogConf) (io.Writer, error) {
	pathList := strings.Split(conf.Filename, "/")
	baseLogPath := ""
	len := len(pathList)
	for i := 0; i < len-1; i++ {
		os.Mkdir(strings.Join(pathList[:i], "/"), os.ModePerm)
		baseLogPath += pathList[i] + string(os.PathSeparator)
	}
	baseLogPath += pathList[len-1:][0]
	newLogName := baseLogPath[:strings.LastIndex(baseLogPath, ".")]
	fileExe := baseLogPath[strings.LastIndex(baseLogPath, ".")+1:]
	file := newLogName + `.%F.` + fileExe
	return rotatelogs.New(
		file,
		rotatelogs.WithClock(rotatelogs.Local), // 使用本地时区，走的是标准库time对象的时区
		//rotatelogs.WithLinkName(path[1]),       // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(time.Duration(conf.MaxAge)*24*time.Hour), // 文件最大保存时间
		rotatelogs.WithRotationTime(24*time.Hour),                      // 日志切割时间间隔
	)
}

//级别转换
func parseLevel(lvl string) zapcore.Level {
	switch strings.ToLower(lvl) {
	case "panic", "dpanic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	case "error":
		return zapcore.ErrorLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "info":
		return zapcore.InfoLevel
	case "debug":
		return zapcore.DebugLevel
	default:
		return zapcore.DebugLevel
	}
}
