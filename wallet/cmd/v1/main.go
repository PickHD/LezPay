package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/PickHD/LezPay/wallet/internal/v1/application"
	"github.com/PickHD/LezPay/wallet/internal/v1/infrastructure"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

const (
	localServerMode = "local"
	httpServerMode  = "http"
	consumerMode    = "consumer"
	grpcMode        = "grpc"
)

// @title           LezPay API
// @version         1.0
// @description     LezPay API - wallet Services
// @contact.name    Taufik Januar
// @contact.email   taufikjanuar35@gmail.com
// @license.name    MIT
// @host            localhost:8083
// @BasePath        /v1
// @Schemes         http
func main() {
	err := godotenv.Load("./cmd/v1/.env")
	if err != nil {
		panic(err)
	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	// Checking command arguments
	var (
		args = os.Args[1:]
		mode = localServerMode
	)

	if len(args) > 0 {
		mode = os.Args[1]
	}

	// create a context with background for setup the application
	ctx := context.Background()
	app, err := application.SetupApplication(ctx)
	if err != nil {
		app.Logger.Error("Failed to initialize app. Error: ", err)
		panic(err)
	}

	switch mode {
	case localServerMode, httpServerMode:
		var (
			httpServer = infrastructure.ServeHTTP(app)
		)

		server := &http.Server{
			Addr:    fmt.Sprintf(":%d", app.Config.Server.AppPort),
			Handler: httpServer,
		}

		// Create a channel to receive OS signals
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt)

		// Start the HTTP server in a separate Goroutine
		go func() {
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				app.Logger.Error("Failed to to start server. Error: ", err)
			}
		}()

		// Wait for a SIGINT or SIGTERM signal
		<-sigCh

		// Create a context with a timeout of 5 seconds
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		app.Close(ctx)

		// Shutdown the server gracefully
		if err := server.Shutdown(ctx); err != nil {
			app.Logger.Error("Failed to shutdown server. Error: ", err)
		}

		app.Logger.Info("SERVER SHUTDOWN GRACEFULLY")
	case grpcMode:
		var (
			grpcServer = infrastructure.ServeGRPC(app)
		)

		errC := make(chan error, 1)

		ctx, stop := signal.NotifyContext(context.Background(),
			os.Interrupt,
			syscall.SIGTERM,
			syscall.SIGQUIT)

		go func() {
			addr := fmt.Sprintf("0.0.0.0:%d", app.Config.Common.GrpcPort)

			lis, err := net.Listen("tcp", addr)
			if err != nil {
				app.Logger.Error("cannot listen tcp grpc", err)
			}

			app.Logger.Info("Listening and serving GRPC server", lis.Addr().String())

			if err := grpcServer.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
				errC <- err
			}
		}()

		go func() {
			<-ctx.Done()

			app.Logger.Info("Shutdown signal received")

			defer func() {
				stop()
				close(errC)
			}()

			app.Logger.Info("Shutdown completed")
		}()

		if err := <-errC; err != nil {
			app.Logger.Error("Error received by channel", err)
		}
	case consumerMode:
		// Make a channel to receive messages into infinite loop.
		forever := make(chan bool)

		topics := []string{app.Config.Kafka.TopicTopupTransaction, app.Config.Kafka.TopicPayoutTransaction}

		for _, q := range topics {
			go infrastructure.ConsumeMessages(app, q)
		}

		<-forever
	}
}
