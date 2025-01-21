package config

import "todo/pkg/infra"

type Config struct {
	GRPC GRPCConfig           `yaml:"grpc"`
	PG   infra.PostgresConfig `yaml:"postgres"`
}

type GRPCConfig struct {
	Port int `yaml:"port"`
}
