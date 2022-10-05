package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	dbr "github.com/abiewardani/dbr/v2"
	"github.com/abiewardani/dbr/v2/dialect"
	_ "github.com/lib/pq"
)

type dbrInstance struct {
	master *dbr.Session
}

func (c *dbrInstance) Master() *dbr.Session {
	return c.master
}

// DbrDatabase abstraction
type DbrDatabase interface {
	Master() *dbr.Session
}

// InitDbr ...
func InitDbr() DbrDatabase {
	inst := new(dbrInstance)

	dsnMaster := fmt.Sprintf(fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable TimeZone=UTC+7",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
	))

	sqlMaster, err := sql.Open("postgres", dsnMaster)
	connMaster := &dbr.Connection{
		DB:            sqlMaster, // <- underlying database/sql.DB is instrumented
		EventReceiver: &dbr.NullEventReceiver{},
		Dialect:       dialect.PostgreSQL,
	}

	sqlMaster.SetConnMaxLifetime(100)
	sqlMaster.SetMaxIdleConns(10)
	sqlMaster.SetConnMaxLifetime(time.Minute)

	if err != nil {
		log.Panic(err)
	}

	if err := connMaster.Ping(); err != nil {
		log.Panic(err)
	}

	inst.master = connMaster.NewSession(nil)

	return inst
}
