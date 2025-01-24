package grpc

import "time"

type Config struct {
	Address string        `yaml:"address"`
	Timeout time.Duration `yaml:"timeout"`
	Retries int           `yaml:"retries"`
}
