package kafka

import (
	"fmt"

	"github.com/IBM/sarama"
)

type connection struct {
	name      string
	host      []string
	publisher sarama.SyncProducer
}

var conInstance *connection

func ConnectionReferenceName(name string) string {
	return fmt.Sprintf("kafka:%s", name)
}

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
	c.createPublisher()
	return nil
}

func (c *connection) createPublisher() error {
	config := sarama.NewConfig()
	config.Producer.Retry.Max = 3
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(c.host, config)
	if err != nil {
		return fmt.Errorf("[kafka-connection] Error creating publisher %s", err)
	}

	c.publisher = producer
	return nil
}

func (c *connection) GetProducer() any {
	return c.publisher
}

func (c *connection) Disconnect() error {
	return nil
}

func (c *connection) ReferenceName() string {

	return ConnectionReferenceName(c.name)
}
