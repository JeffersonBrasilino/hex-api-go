package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	_ "net/http/pprof"

	"github.com/grafana/pyroscope-go"
	"github.com/jeffersonbrasilino/gomes"
	"github.com/jeffersonbrasilino/gomes/channel/kafka"
	"github.com/jeffersonbrasilino/gomes/message"
	gomesOtelTrace "github.com/jeffersonbrasilino/gomes/otel"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
)

// Underneath the hood,
// we use the CQRS pattern to process commands, queries, and events.
// First, we need to create an action (command, query, or event) and a handler for this action.
// This is a basic example.

// cqrs action
type Command struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func CreateCommand(Username, Password string) *Command {
	return &Command{
		Username,
		Password,
	}
}

// this function is responsible for the record name and action routing
func (c *Command) Name() string {
	return "createUser"
}

// CQRS acton handler
type CommandHandler struct {
	tracer gomesOtelTrace.OtelTrace
	header map[string]string
}

// response structure
type ResultCm struct {
	Result any
}

func NewComandHandler() *CommandHandler {
	return &CommandHandler{
		tracer: gomesOtelTrace.InitTrace("command-handler"),
	}
}

// note that the link between the action and its handler is the type of the data parameter.
// This indicates that this handler is responsible for this action
func (c *CommandHandler) Handle(ctx context.Context, data *Command) (*ResultCm, error) {
	time.Sleep(time.Second * 1)
	ctx, span := c.tracer.Start(
		ctx,
		"Handle Command",
	)
	defer span.End()

	slog.Info("processing command...",
		"username", data.Username,
	)
	time.Sleep(time.Second * 5)
	slog.Info("command processed.",
		"username", data.Username,
	)

	return &ResultCm{Result: "DEU BOM"}, nil
	//return nil, fmt.Errorf("DEU RUIM AO PROCESSAR A MENSAGEM")
}

func (c *CommandHandler) SetMessageHeader(header message.Header) {
	c.header = header
}

func main() {

	initOtelTraceProvider()
	ctx, stop := context.WithCancel(context.Background()) //signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	slog.Info("start message system consumer....")

	gomes.AddChannelConnection(
		kafka.NewConnection("defaultConKafka", []string{"kafka:9092"}),
	)
	topicConsumerChannel := kafka.NewConsumerChannelAdapterBuilder(
		"defaultConKafka",
		"gomes.topic",
		"test_consumer",
	)
	//topicConsumerChannel.WithSendReplyUsingReplyTo()
	//topicConsumerChannel.WithRetryTimes(2_000, 5_000)
	//topicConsumerChannel.WithDeadLetterChannelName("gomes.dlq")
	gomes.AddConsumerChannel(topicConsumerChannel)

	//response channel
	/* responseChannel := kafka.NewPublisherChannelAdapterBuilder(
		"defaultConKafka",
		"gomes.response",
	)
	gomes.AddPublisherChannel(responseChannel)

	//DLQ channel
	dlqChannel := kafka.NewPublisherChannelAdapterBuilder(
		"defaultConKafka",
		"gomes.dlq",
	)
	gomes.AddPublisherChannel(dlqChannel) */

	// Register CQRS action and action handler.
	gomes.AddActionHandler(NewComandHandler())

	//enable otel trace for the message system
	gomes.EnableOtelTrace()
	initPyroscope()
	
	//start the message system
	gomes.Start()

	//go publishMessage()

	//For the consumer channel endpoint,
	//the advantage of having an abstraction between the consumer channel and the consumer endpoint
	//is that we can have two different endpoints for the same channel (event-driven or polling).
	//Note that the consumerName parameter of the eventDrivenConsumer method is the same as the consumer name of the consumerChannel.
	consumer, err := gomes.EventDrivenConsumer("test_consumer")
	if err != nil {
		panic(err)
	}


	go func() {
		err := consumer.WithAmountOfProcessors(50).
			WithMessageProcessingTimeout(300000).
			WithStopOnError(false).
			Run(ctx)

		fmt.Println("main.go erro no consumer", "erro", err)

		if err != nil {
			stop()
		}
	}()

	/* time.Sleep(time.Second * 9)
	stop() */

	<-ctx.Done()
	//message system graceful shutdown
	gomes.Shutdown()
	fmt.Println("CONSUMIDOR STOPPED COM SUCESSO...")
}

func publishMessage() {
	maxPublishMessages := 10
	for i := 1; i <= maxPublishMessages; i++ {
		fmt.Println("publish command message...")
		busA, _ := gomes.CommandBusByChannel("gomes.topic")
		busA.SendAsync(context.Background(), CreateCommand(fmt.Sprintf("message %d", i), "123"))
	}
}

func initOtelTraceProvider() *trace.TracerProvider {
	exporter, err := otlptracegrpc.New(context.Background())
	if err != nil {
		panic(fmt.Errorf("failed to create OTLP grpc exporter: %w", err))
	}

	batchSpanProcessor := trace.NewBatchSpanProcessor(exporter)
	provider := trace.NewTracerProvider(
		trace.WithSpanProcessor(batchSpanProcessor),
	)

	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.TraceContext{})
	return provider
}

func initPyroscope() {
	// These 2 lines are only required if you're using mutex or block profiling
	// Read the explanation below for how to set these rates:
	//runtime.SetMutexProfileFraction(5)
	//runtime.SetBlockProfileRate(5)

	pyroscope.Start(pyroscope.Config{
		ApplicationName: "event-driven-consumer",

		// replace this with the address of pyroscope server
		ServerAddress: "http://pyroscope:4040",

		// you can disable logging by setting this to nil
		Logger: pyroscope.StandardLogger,

		// you can provide static tags via a map:
		Tags: map[string]string{"hostname": os.Getenv("HOSTNAME")},

		ProfileTypes: []pyroscope.ProfileType{
			// these profile types are enabled by default:
			pyroscope.ProfileCPU,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileInuseSpace,

			// these profile types are optional:
			pyroscope.ProfileGoroutines,
			pyroscope.ProfileMutexCount,
			pyroscope.ProfileMutexDuration,
			pyroscope.ProfileBlockCount,
			pyroscope.ProfileBlockDuration,
		},
	})
}
