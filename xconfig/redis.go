package config

import (
	"os"
	"strconv"
	"time"
)

// RedisConf redis配置结构
type RedisConf struct {
	Address     string        `json:"address" yaml:"address"`
	Password    string        `json:"password" yaml:"password"`
	Network     string        `json:"network" yaml:"network"`
	MaxIdle     int           `json:"max_idle" yaml:"max_idle"`
	MaxActive   int           `json:"max_active" yaml:"max_active"`
	IdleTimeout time.Duration `json:"idle_timeout" yaml:"idle_timeout"`
	DataBase    int           `json:"data_base" yaml:"data_base"`
}

// SetRedisConfig 使用环境变量替换redis配置参数
func (conf *RedisConf) SetRedisConfig() {
	address := os.Getenv("REDIS_SERVER_ADDRESS")
	if address != "" {
		conf.Address = address
	}

	password := os.Getenv("REDIS_SERVER_PASSWORD")
	if password != "" {
		conf.Password = password
	}

	network := os.Getenv("REDIS_SERVER_NETWORK")
	if network != "" {
		conf.Network = network
	}

	maxActive := os.Getenv("REDIS_SERVER_MAX_ACTIVE")
	if p, err := strconv.Atoi(maxActive); err == nil && p > 0 {
		conf.MaxActive = p
	}

	maxIdle := os.Getenv("REDIS_SERVER_MAX_IDLE")
	if p, err := strconv.Atoi(maxIdle); err == nil && p > 0 {
		conf.MaxIdle = p
	}

	idleTimeout := os.Getenv("REDIS_SERVER_IDLE_TIMEOUT")
	if p, err := strconv.Atoi(idleTimeout); err == nil && p > 0 {
		conf.IdleTimeout = time.Duration(p)
	}
}
