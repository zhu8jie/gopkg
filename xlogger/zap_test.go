package xlogger

import (
	"fmt"
	"testing"
)

func TestEncrypt(t *testing.T) {

	DefaultWriter, err := NewRotateWriter(&LogConf{
		Filename: "tmp/console.log",

		LogLevel: "debug",

		IsStdOut: true,

		// MaxSize is the maximum size in megabytes of the logger file before it gets
		// rotated. It defaults to 100 megabytes.
		MaxSize: 500,

		// MaxAge is the maximum number of days to retain old logger files based on the
		// timestamp encoded in their filename.  Note that a day is defined as 24
		// hours and may not exactly correspond to calendar days due to daylight
		// savings, leap seconds, etc. The default_xorm is not to remove old logger files
		// based on age.
		MaxAge: 7,

		// MaxBackups is the maximum number of old logger files to retain.  The default_xorm
		// is to retain all old logger files (though MaxAge may still cause them to get
		// deleted.)
		MaxBackups: 1,

		// LocalTime determines if the time used for formatting the timestamps in
		// backup files is the computer's local time.  The default_xorm is to use UTC
		// time.
		// LocalTime bool `json:"localtime" yaml:"localtime"`

		// Compress determines if the rotated logger files should be compressed
		// using gzip.
		Compress: true,

		RotateType: "date",
	})
	if err != nil {
		fmt.Println("err: ", err)
	}

	Default, err := NewZap(&LogConf{
		Filename: "tmp/console.log",

		LogLevel: "debug",

		IsStdOut: true,

		// MaxSize is the maximum size in megabytes of the logger file before it gets
		// rotated. It defaults to 100 megabytes.
		MaxSize: 500,

		// MaxAge is the maximum number of days to retain old logger files based on the
		// timestamp encoded in their filename.  Note that a day is defined as 24
		// hours and may not exactly correspond to calendar days due to daylight
		// savings, leap seconds, etc. The default_xorm is not to remove old logger files
		// based on age.
		MaxAge: 7,

		// MaxBackups is the maximum number of old logger files to retain.  The default_xorm
		// is to retain all old logger files (though MaxAge may still cause them to get
		// deleted.)
		MaxBackups: 1,

		// LocalTime determines if the time used for formatting the timestamps in
		// backup files is the computer's local time.  The default_xorm is to use UTC
		// time.
		// LocalTime bool `json:"localtime" yaml:"localtime"`

		// Compress determines if the rotated logger files should be compressed
		// using gzip.
		Compress: true,

		RotateType: "date",
	}, 0, DefaultWriter)
	if err != nil {
		fmt.Println("err1: ", err)
	}
	if Default == nil {
		// return errors.New("default log init fail.")
		fmt.Println("default log init fail.")
	}

	Default.Debugf("abcdeeft%v", 123)
}
