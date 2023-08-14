package xconfig

import "os"

// HttpConf http配置信息
type HttpConf struct {
	IP           string `json:"ip" yaml:"ip"`
	Port         string `json:"port" yaml:"port"`
	ReadTimeout  int    `json:"read_timeout" yaml:"read_timeout"`
	WriteTimeout int    `json:"write_timeout" yaml:"write_timeout"`
}

// SetWebConfig 使用环境变量替换web配置参数
func (conf *HttpConf) SetWebConfig() {
	// ip
	ip := os.Getenv("WEB_IP")
	if ip != "" {
		conf.IP = ip
	}

	// port
	port := os.Getenv("WEB_PORT")
	if port != "" {
		conf.Port = port
	}
}
