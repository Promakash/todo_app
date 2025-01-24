package config

import (
	"todo/pkg/infra"
	"todo/pkg/infra/cache/redis"
	pkglog "todo/pkg/log"
)

type Config struct {
	GRPC   GRPCConfig           `yaml:"grpc"`
	PG     infra.PostgresConfig `yaml:"postgres"`
	Redis  redis.Config         `yaml:"redis"`
	Logger pkglog.Config        `yaml:"logger"`
}

type GRPCConfig struct {
	Port int `yaml:"port"`
}
