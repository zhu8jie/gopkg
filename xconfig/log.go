package xconfig

import (
	"os"
	"strconv"
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

// SetLogConfig 使用环境变量替换配置参数
func (conf *LogConf) SetLogConfig() {
	// 日志目录
	filename := os.Getenv("LOG_FILENAME")
	if filename != "" {
		conf.Filename = filename
	}

	// 日志级别
	level := os.Getenv("LOG_LEVEL")
	if level != "" {
		conf.LogLevel = level
	}

	// 是否在控制台输出
	isStdOut := os.Getenv("LOG_STDOUT")
	if isStdOut != "" {
		if isStdOut == "true" {
			conf.IsStdOut = true
		} else if isStdOut == "false" {
			conf.IsStdOut = false
		}
	}

	// 日志大小
	maxsize := os.Getenv("LOG_MAXSIZE")
	if maxsize != "" {
		n, err := strconv.Atoi(maxsize)
		if err == nil && n > 0 {
			conf.MaxSize = n
		}
	}

	// 日志保留时长
	maxAge := os.Getenv("LOG_MAX_AGE")
	if maxAge != "" {
		n, err := strconv.Atoi(maxAge)
		if err == nil && n > 0 {
			conf.MaxAge = n
		}
	}

	// 备份数量
	maxBackups := os.Getenv("LOG_MAX_BACKUPS")
	if maxBackups != "" {
		n, err := strconv.Atoi(maxBackups)
		if err == nil && n > 0 {
			conf.MaxBackups = n
		}
	}

	// 是否使用本地时间
	localTime := os.Getenv("LOG_LOCALTIME")
	if localTime != "" {
		if localTime == "true" {
			conf.LocalTime = true
		} else if localTime == "false" {
			conf.LocalTime = false
		}
	}

	// 是否压缩
	compress := os.Getenv("LOG_COMPRESS")
	if compress != "" {
		if compress == "true" {
			conf.Compress = true
		} else if compress == "false" {
			conf.Compress = false
		}
	}

	// 日志切割
	rotateType := os.Getenv("LOG_ROTATETYPE")
	if rotateType != "" {
		conf.RotateType = rotateType
	}
}
