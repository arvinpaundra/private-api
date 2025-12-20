package relationaldb

import (
	"log"
	"sync"

	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

type Connectible interface {
	Connect() (*gorm.DB, error)
}

func NewConnection(connect Connectible) {
	once.Do(func() {
		var err error

		db, err = connect.Connect()
		if err != nil {
			log.Fatalf("failed to connect to database: %s", err.Error())
		}
	})
}

func Close() error {
	sqldb, err := db.DB()
	if err != nil {
		return err
	}

	err = sqldb.Close()
	if err != nil {
		return err
	}

	return nil
}

func GetConnection() *gorm.DB {
	return db
}
