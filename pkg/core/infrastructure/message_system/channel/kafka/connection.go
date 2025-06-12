package kafka

import (
	"fmt"

	"github.com/IBM/sarama"
)

type connection struct {
	name      string
	host      []string
	publisher sarama.SyncProducer
	consumer  sarama.Consumer
}

var conInstance *connection

func NewConnection(name string, host []string) *connection {
	if conInstance != nil {
		return conInstance
	}
	conInstance = &connection{
		name: name,
		host: host,
	}
	return conInstance
}

func (c *connection) Connect() error {
	config := sarama.NewConfig()
	config.Producer.Retry.Max = 3
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(c.host, config)
	if err != nil {
		return fmt.Errorf("[kafka-connection] Error creating publisher %s", err)
	}
	c.publisher = producer

	consumer, err := sarama.NewConsumer(c.host, config)
	if err != nil {
		return fmt.Errorf("[kafka-connection] Error creating consumer %s", err)
	}
	c.consumer = consumer
	return nil
}

func (c *connection) GetProducer() any {
	return c.publisher
}

func (c *connection) GetConsumer() any {
	return c.consumer
}

func (c *connection) Disconnect() error {
	return nil
}

func (c *connection) ReferenceName() string {
	return c.name
}
