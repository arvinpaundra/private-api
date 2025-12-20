package memorydb

import (
	"log"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	rdb  *redis.Client
	once sync.Once
)

type InMemoryConnectible interface {
	Connect() (*redis.Client, error)
}

func NewInMemoryConnection(connect InMemoryConnectible) {
	once.Do(func() {
		var err error

		rdb, err = connect.Connect()
		if err != nil {
			log.Fatalf("failed to in memory database: %s", err.Error())
		}
	})
}

func Close() error {
	return rdb.Close()
}

func GetInMemoryConnection() *redis.Client {
	return rdb
}
