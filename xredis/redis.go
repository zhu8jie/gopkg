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
	MaxIdle     int           `json:"max_idle"`
	MaxActive   int           `json:"max_active"`
	IdleTimeout time.Duration `json:"idle_timeout"`
	DataBase    int           `json:"data_base"`
}

func NewRedisClient(redisCfg GoRedisConf) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:        redisCfg.Address,
		Network:     redisCfg.Network,
		PoolSize:    redisCfg.MaxIdle,
		IdleTimeout: redisCfg.IdleTimeout * time.Second,
		DB:          redisCfg.DataBase,
		Password:    redisCfg.Password,
	})
	err := client.Ping(context.Background()).Err()
	return client, err
}
