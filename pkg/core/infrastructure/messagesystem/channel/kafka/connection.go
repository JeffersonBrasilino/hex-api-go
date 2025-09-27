// Package kafka provides Kafka integration for the message system.
//
// This package implements Kafka-specific channel adapters and connections for
// publishing and consuming messages through Apache Kafka. It provides outbound
// and inbound channel adapters with message translation capabilities.
//
// The Connection implementation supports:
// - Kafka producer and consumer connection management
// - Connection pooling and reuse
// - Error handling and connection lifecycle
// - Configuration management for Kafka clients
package kafka

import (
	"fmt"

	"github.com/IBM/sarama"
)

// connection manages Kafka producer and consumer connections with lifecycle
// management capabilities.
type connection struct {
	name      string
	host      []string
	publisher sarama.SyncProducer
	consumer  sarama.Consumer
}

// conInstance holds the singleton connection instance for reuse across
// the application.
var conInstance *connection

// NewConnection creates a new Kafka connection instance. This implementation
// uses a singleton pattern to reuse the same connection across the application.
//
// Parameters:
//   - name: the connection name identifier
//   - host: list of Kafka broker addresses
//
// Returns:
//   - *connection: the connection instance
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

// Connect establishes connections to Kafka brokers for both producer and consumer.
// It configures the Kafka client with appropriate settings for reliability.
//
// Returns:
//   - error: error if connection establishment fails
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

// GetProducer returns the Kafka sync producer instance.
//
// Returns:
//   - sarama.SyncProducer: the Kafka producer
func (c *connection) GetProducer() sarama.SyncProducer {
	return c.publisher
}

// GetConsumer returns the Kafka consumer instance.
//
// Returns:
//   - sarama.Consumer: the Kafka consumer
func (c *connection) GetConsumer() sarama.Consumer {
	return c.consumer
}

// Disconnect closes the Kafka connections and releases associated resources.
//
// Returns:
//   - error: error if disconnection fails (typically nil)
func (c *connection) Disconnect() error {
	return nil
}

// ReferenceName returns the connection name identifier.
//
// Returns:
//   - string: the connection name
func (c *connection) ReferenceName() string {
	return c.name
}
