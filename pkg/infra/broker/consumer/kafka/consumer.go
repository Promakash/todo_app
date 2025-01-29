package kafka

import (
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Config struct {
	Address string `yaml:"broker_address"`
	Topic   string `yaml:"topic"`
	GroupID string `yaml:"group_id"`
	Offset  string `yaml:"offset"`
}

type BrokerConsumer struct {
	Topic         string
	GroupID       string
	KafkaConsumer *kafka.Consumer
}

func NewConsumer(cfg Config) (*BrokerConsumer, error) {
	config := &kafka.ConfigMap{
		"bootstrap.servers": cfg.Address,
		"group.id":          cfg.GroupID,
		"auto.offset.reset": cfg.Offset,
	}

	conn, err := kafka.NewConsumer(config)
	if err != nil {
		return nil, err
	}
	err = conn.Subscribe(cfg.Topic, nil)
	if err != nil {
		return nil, err
	}

	return &BrokerConsumer{
		Topic:         cfg.Topic,
		GroupID:       cfg.GroupID,
		KafkaConsumer: conn,
	}, nil
}

func (c *BrokerConsumer) Consume(msg any) error {
	m, err := c.KafkaConsumer.ReadMessage(-1)
	if err != nil {
		return err
	}

	err = json.Unmarshal(m.Value, &msg)
	if err != nil {
		return err
	}

	return nil
}

func (c *BrokerConsumer) Close() {
	_ = c.KafkaConsumer.Close()
}
