package memorydb

import (
	"context"
	"fmt"

	"github.com/arvinpaundra/private-api/config"
	"github.com/redis/go-redis/v9"
)

type redisdb struct {
	host string
	port string
	pass string
	db   int
}

func NewRedisDB() *redisdb {
	return &redisdb{
		host: config.GetString("REDIS_HOST"),
		port: config.GetString("REDIS_PORT"),
		pass: config.GetString("REDIS_PASS"),
		db:   config.GetInt("REDIS_DB"),
	}
}

func (d *redisdb) Connect() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", d.host, d.port),
		Password: d.pass,
		DB:       d.db,
	})

	err := rdb.Ping(context.TODO()).Err()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}
