package relationaldb

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type pgsql struct {
	host     string
	port     string
	user     string
	pass     string
	dbname   string
	sslmode  string
	timezone string
}

func NewPostgres() *pgsql {
	return &pgsql{
		host:     viper.GetString("DB_HOST"),
		port:     viper.GetString("DB_PORT"),
		user:     viper.GetString("DB_USER"),
		pass:     viper.GetString("DB_PASS"),
		dbname:   viper.GetString("DB_DBNAME"),
		sslmode:  viper.GetString("DB_SSLMODE"),
		timezone: viper.GetString("DB_TIMEZONE"),
	}
}

func (d *pgsql) Connect() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		d.host, d.user, d.pass, d.dbname, d.port, d.sslmode, d.timezone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NowFunc:     func() time.Time { return time.Now().UTC() },
		Logger:      logger.Default,
		PrepareStmt: true,
	})

	if err != nil {
		return nil, err
	}

	return db, nil
}
