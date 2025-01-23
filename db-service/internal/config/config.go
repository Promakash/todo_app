package config

import (
	"todo/pkg/infra"
	"todo/pkg/infra/cache/redis"
)

type Config struct {
	GRPC  GRPCConfig           `yaml:"grpc"`
	PG    infra.PostgresConfig `yaml:"postgres"`
	Redis redis.Config         `yaml:"redis"`
}

type GRPCConfig struct {
	Port int `yaml:"port"`
}
