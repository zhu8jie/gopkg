package xutils

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func SignalHandler() {
	var (
		ch = make(chan os.Signal, 1)
	)
	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM)
	for {
		si := <-ch
		switch si {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL:
			return
		default:
			return
		}
	}
}

func Int64ToStr(i int64) string {
	return strconv.FormatInt(i, 10)
}

func StrToInt64(i string) int64 {
	ret, _ := strconv.ParseInt(i, 10, 64)
	return ret
}
