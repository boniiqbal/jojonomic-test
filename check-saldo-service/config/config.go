package config

import (
	"os"
	"strconv"

	"github.com/check-saldo-service/config/database"
)

// Config struct
type Config struct {
	serviceName string
	environment string
	debug       bool
	port        int
	db          database.DbrDatabase
}

// NewConfig func
func NewConfig() *Config {
	cfg := new(Config)

	cfg.ConnectDB()
	return cfg
}

// ConnectDB func
func (c *Config) ConnectDB() {
	c.db = database.InitDbr()
}

// DB func
func (c *Config) DB() database.DbrDatabase {
	return c.db
}

// ServiceName ..
func (c *Config) ServiceName() string {
	return os.Getenv(`SERVICE_NAME`)
}

// Port func
func (c *Config) Port() int {
	v := os.Getenv("PORT")
	c.port, _ = strconv.Atoi(v)

	return c.port
}
