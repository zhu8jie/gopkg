package xutils

import (
	"errors"
	"flag"
	"fmt"
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

func StrToInt(i string) int {
	ret, _ := strconv.Atoi(i)
	return ret
}

func IntToStr(i int) string {
	return strconv.Itoa(i)
}

func GetFlagPath(input ...string) (map[string]string, error) {
	ret := make(map[string]string)

	tr := make(map[string]*string)
	for _, i := range input {
		t := flag.String(i, "", "")
		tr[i] = t
	}
	flag.Parse()

	for k, v := range tr {
		if *v == "" {
			return nil, errors.New(fmt.Sprintf("Not found input flag -%s \n", k))
		} else {
			ret[k] = *v
		}
	}

	return ret, nil
}
