package config

import (
	kafkaconsumer "todo/pkg/infra/broker/consumer/kafka"
	pkglog "todo/pkg/log"
)

type HTTPConfig struct {
	Address string `yaml:"address"`
}

type Config struct {
	HTTPServer    HTTPConfig           `yaml:"http_server"`
	Logger        pkglog.Config        `yaml:"logger"`
	KafkaConsumer kafkaconsumer.Config `yaml:"kafka"`
}
