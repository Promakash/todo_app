package config

import (
	"time"
	"todo/api-gateway/internal/clients/todo/grpc"
	pkglog "todo/pkg/log"
)

type HTTPConfig struct {
	Address      string        `yaml:"address"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

type Config struct {
	HTTPServer HTTPConfig    `yaml:"http_server"`
	GRPCClient grpc.Config   `yaml:"grpc_client"`
	Logger     pkglog.Config `yaml:"logger"`
}
