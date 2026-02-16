package http

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/hex-api-go/internal/user/application/command/createuser"
	gomes "github.com/jeffersonbrasilino/gomes"
	"github.com/jeffersonbrasilino/gomes/channel/kafka"
	"github.com/jeffersonbrasilino/gomes/otel"
)

var otelTrace = otel.InitTrace("HTTP POST")

type Request struct {
	Username string `validate:"gte=4"`
	Password string `validate:"required"`
}

func CreateUser(ctx context.Context, fiberApp fiber.Router) {

	gomes.AddChannelConnection(
		kafka.NewConnection("defaultConKafka", []string{"kafka:9092"}),
	)

	/* publisherResponseChannel := kafka.NewPublisherChannelAdapterBuilder(
		"defaultConKafka",
		"gomes.response",
	)
	gomes.AddPublisherChannel(publisherResponseChannel) */

	a := kafka.NewPublisherChannelAdapterBuilder(
		"defaultConKafka",
		"gomes.topic",
	)
	a.WithReplyChannelName("gomes.response")
	gomes.AddPublisherChannel(a)

	fiberApp.Post("/create/kafka", func(c *fiber.Ctx) error {

		ctx, span := otelTrace.Start(
			c.Context(),
			"post /users/create/kafka",
			otel.WithSpanKind(otel.SpanKindServer),
		)
		defer span.End()

		request := new(Request)
		if err := c.BodyParser(request); err != nil {
			return c.SendStatus(400)
		}

		for i := 1; i <= 500; i++ {
			busA, _ := gomes.CommandBusByChannel("gomes.topic")
			// err := busA.SendAsync(ctx, createuser.CreateCommand(fmt.Sprintf("teste ID: %v", i), "123"))
			err := busA.SendRawAsync(
				ctx,
				"createUser",
				createuser.CreateCommand(fmt.Sprintf("teste ID: %v", i), "123"),
				map[string]string{"header1": "val 1", "header2": "val 2"},
			)
			fmt.Println("[controller] ASYNC COMMAND SEND ERROR ", err)
		}

		return c.JSON("foi OK")
	})
}

func CreateUserRabbit(ctx context.Context, fiberApp fiber.Router) {

	/* gomes.AddChannelConnection(
		rabbitmq.NewConnection("rabbit-test", "admin:admin@rabbitmq:5672"),
	)
	pubChan := rabbitmq.NewPublisherChannelAdapterBuilder("rabbit-test", "gomes-exchange").
		WithChannelType(rabbitmq.ProducerExchange).
		WithExchangeType(rabbitmq.ExchangeDirect).
		WithExchangeRoutingKeys("rota-fila-1")
	gomes.AddPublisherChannel(pubChan)


	fiberApp.Post("/create/rabbit", func(c *fiber.Ctx) error {

		ctx, span := otelTrace.Start(
			c.Context(),
			"post /users/create/rabbit",
			otel.WithSpanKind(otel.SpanKindServer),
		)
		defer span.End()

		request := new(Request)
		if err := c.BodyParser(request); err != nil {
			return c.SendStatus(400)
		}

		for i := 1; i <= 1; i++ {
			busA, _ := gomes.CommandBusByChannel("gomes-exchange")
			err := busA.SendAsync(ctx, createuser.CreateCommand(fmt.Sprintf("teste ID: %v", i), "123"))
			fmt.Println("[controller] ASYNC COMMAND SEND ERROR ", err)
		}

		return c.JSON("foi OK")
	}) */
}
func CreateUserSync(ctx context.Context, fiberApp fiber.Router) {
	fiberApp.Post("/create/sync", func(c *fiber.Ctx) error {

		ctx, span := otelTrace.Start(
			c.Context(),
			"post /users/create",
			otel.WithSpanKind(otel.SpanKindServer),
		)
		defer span.End()

		request := new(Request)
		if err := c.BodyParser(request); err != nil {
			return c.SendStatus(400)
		}

		//coreHttp.ValidateRequest(request)

		bus, _ := gomes.CommandBus()
		res, err := bus.SendRaw(
			ctx,
			"createUser",
			createuser.CreateCommand(fmt.Sprintf("teste ID: %v", 1), "123"),
			map[string]string{"header1": "val 1", "header2": "val 2"},
		)
		fmt.Println("result", res, err)
		//res, err := bus.Send(ctx, createuser.CreateCommand("teste", "123"))
		if err != nil {
			return c.SendStatus(500)
		}

		return c.JSON("foi OK")
	})
}
