package config

import (
	"os"
	"strconv"
)

// Config struct
type Config struct {
	serviceName string
	environment string
	debug       bool
	port        int
	kafkaTopic  string
	kafkaUrl    string
}

// NewConfig func
func NewConfig() *Config {
	cfg := new(Config)

	return cfg
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

func (c *Config) KafkaUrl() string {
	return os.Getenv("KAFKA_URL")
}

func (c *Config) KafkaTopic() string {
	return os.Getenv("KAFKA_TOPIC")
}
