package xredis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type GoRedisConf struct {
	Address     string        `json:"address"`
	Password    string        `json:"password"`
	Network     string        `json:"network"`
	DataBase    int           `json:"data_base"`
	MaxIdle     int           `json:"max_idle"`     // 池子大小
	MinIdle     int           `json:"min_idle"`     // 最小活跃数，保证快速响应
	Timeout     time.Duration `json:"timeout"`      // 客户端等待连接的最长时间
	DialTimeout time.Duration `json:"dial_timeout"` // 建立新链接的超时时长
	IdleTimeout time.Duration `json:"idle_timeout"` // 客户端关闭空闲连接的时间间隔。该时间间隔应小于服务器的超时时间
}

func NewRedisClient(redisCfg GoRedisConf) (*redis.Client, error) {
	if redisCfg.MaxIdle == 0 {
		redisCfg.MaxIdle = 10
	}
	if redisCfg.Timeout == 0 {
		redisCfg.Timeout = 10 * time.Second
	}
	if redisCfg.MinIdle == 0 {
		redisCfg.MinIdle = 2
	}
	if redisCfg.DialTimeout == 0 {
		redisCfg.DialTimeout = 1 * time.Second
	}
	if redisCfg.IdleTimeout == 0 {
		redisCfg.IdleTimeout = 1 * time.Second
	}
	client := redis.NewClient(&redis.Options{
		Addr:         redisCfg.Address,
		Password:     redisCfg.Password,
		Network:      redisCfg.Network,
		DB:           redisCfg.DataBase,
		PoolSize:     redisCfg.MaxIdle,
		PoolTimeout:  redisCfg.Timeout,
		MinIdleConns: redisCfg.MinIdle,
		DialTimeout:  redisCfg.DialTimeout,
		IdleTimeout:  redisCfg.IdleTimeout,
	})
	err := client.Ping(context.Background()).Err()
	return client, err
}
