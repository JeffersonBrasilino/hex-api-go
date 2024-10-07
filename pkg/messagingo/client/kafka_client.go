package client

import "fmt"

type KafkaClient struct{}

func CreateKafkaClient() *KafkaClient {
	return &KafkaClient{}
}
func (k *KafkaClient) Publish() {
	fmt.Println("KAFKA publish MESSAGE")
}
