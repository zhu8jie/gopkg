package xhttp

import (
	"net/http"
	"time"
)

// g_conn_http = &http.Client{
// Transport: &http.Transport{
// 	// DialContext: (&net.Dialer{}).DialContext,
// 	Dial: func(netw, addr string) (net.Conn, error) {
// 		conn_http, err := net.DialTimeout(netw, addr, time.Millisecond*600)
// 		if err != nil {
// 			fmt.Println("dail timeout", err)
// 			return nil, err
// 		}
// 		conn_http.SetDeadline(time.Now().Add(time.Second * 15))
// 		return conn_http, nil

// 	},
// 	// MaxIdleConns:          tools.StrToint(g_config["worker_cnt"]),
// 	MaxIdleConnsPerHost:   256,
// 	ResponseHeaderTimeout: time.Millisecond * 298,
// 	DisableKeepAlives:     false,
// },
// 	Timeout: time.Millisecond * time.Duration(tools.StrToint(g_config["dsp_timeout"])),
// }

var httpCli *http.Client

func NewXhttpClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
	}
}
