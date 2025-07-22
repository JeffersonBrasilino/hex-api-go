// Package endpoint implements the event-driven consumer pattern for message processing systems.
//
// Intenção: Fornecer uma estrutura para consumir mensagens de forma assíncrona e escalável, utilizando múltiplos processadores e integração com gateways e canais de entrada.
// Objetivo: Facilitar o consumo, processamento e roteamento de mensagens em sistemas orientados a eventos, com suporte a timeout, dead letter e interceptadores.
//
// Este pacote é parte do sistema de mensageria, permitindo a construção de consumidores desacoplados e resilientes.
package endpoint

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/hex-api-go/pkg/core/infrastructure/message_system/container"
	"github.com/hex-api-go/pkg/core/infrastructure/message_system/message"
)

// EventDrivenConsumerBuilder é responsável por construir instâncias de EventDrivenConsumer.
// referenceName identifica o canal de entrada a ser consumido.
type EventDrivenConsumerBuilder struct {
	referenceName string
}

// EventDrivenConsumer representa um consumidor orientado a eventos.
// Gerencia múltiplos processadores, fila de processamento e integração com gateway e canal de entrada.
type EventDrivenConsumer struct {
	referenceName                 string
	processingTimeoutMilliseconds int
	gateway                       *Gateway
	inboundChannelAdapter         message.InboundChannelAdapter
	amountOfProcessors            int
	processingQueue               chan *message.Message
	processorsWaitGroup           sync.WaitGroup
	ctx                           context.Context
	close                         context.CancelFunc
}

// NewEventDrivenConsumerBuilder cria uma nova instância de EventDrivenConsumerBuilder.
//
// Parâmetros:
//   - referenceName: nome de referência do canal de entrada.
//
// Retorno: ponteiro para EventDrivenConsumerBuilder.
func NewEventDrivenConsumerBuilder(
	referenceName string,
) *EventDrivenConsumerBuilder {
	return &EventDrivenConsumerBuilder{
		referenceName: referenceName,
	}
}

// NewEventDrivenConsumer cria uma nova instância de EventDrivenConsumer.
//
// Parâmetros:
//   - referenceName: nome de referência do canal de entrada.
//   - gateway: ponteiro para o Gateway associado.
//   - inboundChannelAdapter: adaptador do canal de entrada.
//
// Retorno: ponteiro para EventDrivenConsumer.
func NewEventDrivenConsumer(
	referenceName string,
	gateway *Gateway,
	inboundChannelAdapter message.InboundChannelAdapter,
) *EventDrivenConsumer {
	ctx, cancel := context.WithCancel(context.Background())
	consumer := &EventDrivenConsumer{
		referenceName:                 referenceName,
		processingTimeoutMilliseconds: 100000,
		gateway:                       gateway,
		inboundChannelAdapter:         inboundChannelAdapter,
		amountOfProcessors:            1,
		ctx:                           ctx,
		close:                         cancel,
	}
	return consumer
}

// Build constrói um EventDrivenConsumer a partir do container de dependências.
//
// Parâmetros:
//   - container: container de dependências.
//
// Retorno: ponteiro para EventDrivenConsumer e erro, se houver.
func (b *EventDrivenConsumerBuilder) Build(
	container container.Container[any, any],
) (*EventDrivenConsumer, error) {

	anyChannel, err := container.Get(b.referenceName)
	if err != nil {
		panic(
			fmt.Sprintf(
				"[event-driven-consumer] consumer channel %s not found.",
				b.referenceName,
			),
		)
	}

	inboundChannel, ok := anyChannel.(message.InboundChannelAdapter)
	if !ok {
		panic(
			fmt.Sprintf(
				"[event-driven-consumer] consumer channel %s is not a consumer channel.",
				b.referenceName,
			),
		)
	}

	gatewayBuilder := NewGatewayBuilder(inboundChannel.ReferenceName(), "")
	if inboundChannel.DeadLetterChannelName() != "" {
		gatewayBuilder.WithDeadLetterChannel(inboundChannel.DeadLetterChannelName())
	}

	if len(inboundChannel.BeforeProcessors()) > 0 {
		gatewayBuilder.WithBeforeInterceptors(inboundChannel.BeforeProcessors()...)
	}

	if len(inboundChannel.AfterProcessors()) > 0 {
		gatewayBuilder.WithAfterInterceptors(inboundChannel.AfterProcessors()...)
	}

	gateway, _ := gatewayBuilder.Build(container)
	consumer := NewEventDrivenConsumer(
		b.referenceName,
		gateway,
		inboundChannel,
	)
	return consumer, nil
}

// WithMessageProcessingTimeout define o timeout de processamento de mensagens em milissegundos.
//
// Parâmetros:
//   - milisseconds: tempo limite em milissegundos.
//
// Retorno: ponteiro para EventDrivenConsumer.
func (b *EventDrivenConsumer) WithMessageProcessingTimeout(
	milisseconds int,
) *EventDrivenConsumer {
	if milisseconds > 0 {
		b.processingTimeoutMilliseconds = milisseconds
	}
	return b
}

// WithAmountOfProcessors define a quantidade de processadores concorrentes.
//
// Parâmetros:
//   - value: quantidade de processadores.
//
// Retorno: ponteiro para EventDrivenConsumer.
func (b *EventDrivenConsumer) WithAmountOfProcessors(
	value int,
) *EventDrivenConsumer {
	if value > 1 {
		b.amountOfProcessors = value
	}
	return b
}

// Run inicia o processamento das mensagens recebidas do canal de entrada.
//
// Parâmetros:
//   - ctx: contexto para controle de cancelamento e timeout.
//
// Retorno: erro, se houver.
func (e *EventDrivenConsumer) Run(ctx context.Context) error {
	e.processingQueue = make(chan *message.Message, e.amountOfProcessors)
	e.startProcessorsNodes(e.ctx)
	defer e.shutdown()
	for {
		select {
		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.Canceled) {
				slog.Info(
					"[event-driven consumer] Context cancelled, stopping event-driven consumer.",
					"consumerName", e.referenceName,
					"error", ctx.Err(),
				)
				return nil
			}
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				slog.Info(
					"[event-driven consumer] Deadline/Timeout exceeded, stopping event-driven consumer.",
					"consumerName", e.referenceName,
					"error", ctx.Err(),
				)
				return nil
			}
		case <-e.ctx.Done():
			slog.Info(
				"[event-driven consumer] stopping event-driven consumer",
				"consumerName", e.referenceName,
			)
			return nil
		default:
		}

		msg, err := e.inboundChannelAdapter.ReceiveMessage(e.ctx)
		if err != nil {
			return err
		}

		select {
		case e.processingQueue <- msg:
		case <-e.ctx.Done():
			slog.Info(
				"[event-driven consumer] stopping event-driven consumer",
				"consumerName", e.referenceName,
			)
			return nil
		}
	}

}

// Stop solicita a parada do consumidor, cancelando o contexto interno.
func (e *EventDrivenConsumer) Stop() {
	e.close()
}

// shutdown encerra o processamento, fecha o canal de entrada e aguarda finalização dos processadores.
func (e *EventDrivenConsumer) shutdown() {
	fmt.Println("shutdowning event-driven consumer...")
	e.inboundChannelAdapter.Close()
	close(e.processingQueue)
	e.processorsWaitGroup.Wait()
}

// startProcessorsNodes inicia os processadores concorrentes para consumir mensagens da fila.
//
// Parâmetros:
//   - ctx: contexto para controle de cancelamento e timeout.
func (e *EventDrivenConsumer) startProcessorsNodes(ctx context.Context) {
	for i := 0; i < e.amountOfProcessors; i++ {
		e.processorsWaitGroup.Add(1)
		go func(workerId int) {
			defer e.processorsWaitGroup.Done()
			for msg := range e.processingQueue {
				e.sendToGateway(ctx, msg, workerId)
			}
		}(i)
	}
}

// sendToGateway envia a mensagem para o gateway para processamento.
//
// Parâmetros:
//   - ctx: contexto para controle de timeout.
//   - msg: mensagem a ser processada.
//   - nodeId: identificador do processador.
func (e *EventDrivenConsumer) sendToGateway(
	ctx context.Context,
	msg *message.Message,
	nodeId int,
) {

	opCtx, cancel := context.WithTimeout(ctx,
		time.Duration(e.processingTimeoutMilliseconds)*time.Millisecond,
	)
	defer cancel()

	select {
	case <-opCtx.Done():
		return
	default:
	}

	slog.Info("[event-driven consumer] message processing started",
		"consumerName", e.referenceName,
		"nodeId", nodeId,
		"message", msg,
	)

	var err error
	time.Sleep(time.Second * 7)
	fmt.Println("processing OKOKOKOKOKOKOKOK")
	//_, err := e.gateway.Execute(opCtx, msg)
	if err != nil {
		slog.Error("[event-driven consumer] failed to process message",
			"error", err,
			"name", e.referenceName,
			"nodeId", nodeId,
			"message", msg,
		)
		return
	}

	select {
	case <-opCtx.Done():
		if errors.Is(opCtx.Err(), context.DeadlineExceeded) {
			slog.Info("[event-driven consumer] failed to process message",
				"consumerName", e.referenceName, "nodeId", nodeId,
				"error", opCtx.Err(),
			)
		}
		return
	default:
	}

	slog.Info("[event-driven consumer] message processing completed",
		"consumerName", e.referenceName,
		"nodeId", nodeId,
		"message", msg,
	)
}
