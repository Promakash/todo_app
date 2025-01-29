#!/bin/bash

/etc/confluent/docker/run &

echo "Waiting for Kafka to start..."
while ! nc -z localhost 9092; do
    sleep 1
done

echo "Kafka is up. Creating topics..."

kafka-topics --bootstrap-server localhost:9092 --create --topic metrics --partitions 1 --replication-factor 1

wait
