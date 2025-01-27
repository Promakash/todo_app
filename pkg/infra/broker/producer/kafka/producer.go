package kafka

import (
	"context"
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"time"
)

type Config struct {
	Address string        `yaml:"broker_address"`
	Topic   string        `yaml:"topic"`
	Retries int           `yaml:"retries"`
	Timeout time.Duration `yaml:"timeout"`
	Acks    string        `yaml:"acks"`
}

type BrokerProducer struct {
	Topic         string
	KafkaProducer *kafka.Producer
}

func NewProducer(cfg Config) (*BrokerProducer, error) {
	config := &kafka.ConfigMap{
		"bootstrap.servers":   cfg.Address,
		"acks":                cfg.Acks,
		"retries":             cfg.Retries,
		"delivery.timeout.ms": cfg.Timeout,
	}

	conn, err := kafka.NewProducer(config)
	if err != nil {
		return nil, err
	}

	return &BrokerProducer{
		Topic:         cfg.Topic,
		KafkaProducer: conn,
	}, nil
}

func (p *BrokerProducer) Produce(ctx context.Context, msg any) error {
	bytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	deliveryChan := make(chan kafka.Event, 1)
	defer close(deliveryChan)

	err = p.KafkaProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.Topic, Partition: kafka.PartitionAny},
		Value:          bytes,
	}, deliveryChan)
	if err != nil {
		return err
	}

	select {
	case e := <-deliveryChan:
		m := e.(*kafka.Message)
		if m.TopicPartition.Error != nil {
			return m.TopicPartition.Error
		}
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}

func (p *BrokerProducer) Close() {
	const FlushTimeout = 15 * 1000
	p.KafkaProducer.Flush(FlushTimeout)
	p.KafkaProducer.Close()
}
