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
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	slog.Info("starting api server...")
	httpServer := gin.Default()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	dbConn := connectToDatabase()
	//tp := initOtelTraceProvider()
	//initPyroscope()
	
	//bootstrap modules
	modules := []pkg.Module{
		user.NewUserModule(httpServer, dbConn),
	}
	
	for _, module := range modules {
		if err := module.Register(ctx); err != nil {
			panic(err)
		}
	}
	
	gomes.Start()
	//gomes.EnableOtelTrace()

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", os.Getenv("APP_PORT")),
		Handler: httpServer,
	}

	go func() {
		slog.Info("http server listening", "port", os.Getenv("APP_PORT"))
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

func connectToDatabase() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASS"),
		os.Getenv("POSTGRES_DBNAME"),
		os.Getenv("POSTGRES_PORT"))

	dbConn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(fmt.Errorf("failed to connect to database: %w", err))
	}

	return dbConn
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
		ServerAddress: os.Getenv("PYROSCOPE_SERVER_ADDRESS"),

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
