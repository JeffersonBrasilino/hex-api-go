package endpoint

import (
	"fmt"
	"time"

	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
)

type pollingConsumerBuilder struct {
	channel                 message.ConsumerChannel
	gateway                 *Gateway
	fixedRateInMilliseconds int
	rateDelayInMilliseconds int
	stopOnError             bool
	finishWhenNoMessages    bool
}

func NewPolllingConsumerBuilder(
	channel message.ConsumerChannel,
	gateway *Gateway,
) *pollingConsumerBuilder {
	return &pollingConsumerBuilder{
		channel:                 channel,
		gateway:                 gateway,
		fixedRateInMilliseconds: 1000,
		rateDelayInMilliseconds: 0,
		stopOnError:             true,
		finishWhenNoMessages:    false,
	}
}

func (b *pollingConsumerBuilder) WithFixedRateInMilliseconds(
	value int,
) *pollingConsumerBuilder {
	b.fixedRateInMilliseconds = value
	return b
}

func (b *pollingConsumerBuilder) WithRateDelayInMilliseconds(
	value int,
) *pollingConsumerBuilder {
	b.rateDelayInMilliseconds = value
	return b
}

func (b *pollingConsumerBuilder) WithStopOnError(
	value bool,
) *pollingConsumerBuilder {
	b.stopOnError = value
	return b
}

func (b *pollingConsumerBuilder) WithFinishWhenNoMessages(
	value bool,
) *pollingConsumerBuilder {
	b.finishWhenNoMessages = value
	return b
}

func (b *pollingConsumerBuilder) Build() *pollingConsumer {
	return NewPollingConsumer(
		b.fixedRateInMilliseconds,
		b.rateDelayInMilliseconds,
		b.stopOnError,
		b.finishWhenNoMessages,
	)
}

type pollingConsumer struct {
	fixedRateInMilliseconds int
	rateDelayInMilliseconds int
	stopOnError             bool
	finishWhenNoMessages    bool
	hasRunning              bool
}

func NewPollingConsumer(
	fixedRateInMilliseconds int,
	rateDelayInMilliseconds int,
	stopOnError bool,
	finishWhenNoMessages bool,
) *pollingConsumer {
	return &pollingConsumer{
		fixedRateInMilliseconds: fixedRateInMilliseconds,
		rateDelayInMilliseconds: rateDelayInMilliseconds,
		stopOnError:             stopOnError,
		finishWhenNoMessages:    finishWhenNoMessages,
	}
}

func (c *pollingConsumer) Run() error {
	c.hasRunning = true
	if c.rateDelayInMilliseconds > 0 {
		time.Sleep(time.Millisecond * time.Duration(c.fixedRateInMilliseconds))
	}

	for c.hasRunning {
		fmt.Println("run consumer ok")
		fmt.Println("hasRunnnig", c.hasRunning)
		time.Sleep(time.Millisecond * time.Duration(c.fixedRateInMilliseconds))
	}
	return nil
}

func (c *pollingConsumer) Stop() {
	c.hasRunning = false
}

type consumerGateway struct {
	gateway *Gateway
	channel message.ConsumerChannel
}

func NewConsumer(channel message.ConsumerChannel, gateway *Gateway) *consumerGateway {
	return &consumerGateway{
		gateway: gateway,
		channel: channel,
	}
}

func (c *consumerGateway) Execute() error {
	msg, err := c.channel.Receive()
	fmt.Println("received message", msg, err)

	if err != nil {
		return err
	}

	if msg == nil {
		return nil
	}

	//c.gateway.Execute(msg)
	return nil
}
