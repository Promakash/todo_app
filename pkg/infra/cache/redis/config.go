package redis

import "time"

type Config struct {
	Host     string        `yaml:"host"`
	Port     int           `yaml:"port"`
	Password string        `yaml:"password"`
	TTL      time.Duration `yaml:"TTL"`
}
