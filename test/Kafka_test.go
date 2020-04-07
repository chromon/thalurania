package main

import (
	"chalurania/service/kafka"
	"testing"
)

func TestKafkaProducer(t *testing.T) {
	kafka.Producer()
}

func TestKafkaConsumer(t *testing.T) {
	kafka.Consumer()
}
