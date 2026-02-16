package main

import (
	"context"
	"fmt"
	"log/slog"
	"os/signal"
	"syscall"
	"os"
	"github.com/gofiber/fiber/v2"
	"github.com/grafana/pyroscope-go"
	"github.com/hex-api-go/internal/user"
	gomes "github.com/jeffersonbrasilino/gomes"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
)

func main() {

	slog.Info("starting api server...")
	app := fiber.New()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	initPyroscope()

	user.Bootstrap().
		WithHttpProtocol(ctx, app)

	tp := initOtelTraceProvider()
	gomes.EnableOtelTrace()
	gomes.Start()

	go func() {
		slog.Info("http server listening on port 4000")
		if err := app.Listen(":4000"); err != nil {
			panic(err)
		}
	}()

	<-ctx.Done()
	gomes.Shutdown()
	tp.Shutdown(ctx)
	if err := app.Shutdown(); err != nil {
		slog.Info("shutting down server error")
	}
	slog.Info("shutdown completed")

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
		ApplicationName: os.Getenv("APP_NAME"),

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
