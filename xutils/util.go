package xutils

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"hash/crc32"
	"net/url"
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

func UrlEncode(strin string) string {
	if len(strin) > 0 {
		return url.QueryEscape(strin)
	}
	return ""
}

func Base64Encode(strin string) string {
	if len(strin) > 0 {
		return base64.RawURLEncoding.EncodeToString([]byte(strin))
	} else {
		return ""
	}
}

func Base64Decode(strin string) string {
	if len(strin) > 0 {
		tmpres, _ := base64.RawURLEncoding.DecodeString(strin)
		return string(tmpres)
	} else {
		return ""
	}
}

func Md5(strin string) string {
	if strin == "" {
		return ""
	}
	h := md5.New()
	h.Write([]byte(strin))
	return hex.EncodeToString(h.Sum(nil))
}

func Crc(strin string) uint32 {
	if len(strin) < 64 {
		var scratch [64]byte
		copy(scratch[:], strin)
		return crc32.ChecksumIEEE(scratch[:len(strin)])
	}
	return crc32.ChecksumIEEE([]byte(strin))
}
