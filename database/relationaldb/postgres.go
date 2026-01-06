package relationaldb

import (
	"fmt"
	"time"

	"github.com/arvinpaundra/private-api/config"
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
		host:     config.GetString("DB_HOST"),
		port:     config.GetString("DB_PORT"),
		user:     config.GetString("DB_USER"),
		pass:     config.GetString("DB_PASS"),
		dbname:   config.GetString("DB_DBNAME"),
		sslmode:  config.GetString("DB_SSLMODE"),
		timezone: config.GetString("DB_TIMEZONE"),
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
