package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/grafana/pyroscope-go"
	gomes "github.com/jeffersonbrasilino/gomes"
	"github.com/jeffersonbrasilino/hex-api-go/internal/user"
	"github.com/jeffersonbrasilino/hex-api-go/pkg"
	httpPkg "github.com/jeffersonbrasilino/hex-api-go/pkg/http"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
)

func main() {

	slog.Info("starting api server...")
	httpServer := gin.Default()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	//middlewares
	httpServer.Use(httpPkg.BadRequestResponseParser())

	//bootstrap modules 
	modules := []pkg.Module{
		user.NewUserModule(httpServer, nil),
	}

	for _, module := range modules {
		if err := module.Register(ctx); err != nil {
			panic(err)
		}
	}

	//tp := initOtelTraceProvider()
	//initPyroscope()
	//gomes.EnableOtelTrace()
	gomes.Start()

	server := &http.Server{
		Addr: ":4000",
		Handler: httpServer,
	}

	go func() {
		slog.Info("http server listening on port 4000")
		if err := server.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	<-ctx.Done()
	gomes.Shutdown()
	//tp.Shutdown(ctx)
	if err := server.Shutdown(ctx); err != nil {
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
		Logger: nil,

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
