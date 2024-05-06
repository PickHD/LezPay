package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/PickHD/LezPay/customer/internal/application"
	"github.com/PickHD/LezPay/customer/internal/infrastructure"
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
// @description     LezPay API - Customer Services
// @contact.name    Taufik Januar
// @contact.email   taufikjanuar35@gmail.com
// @license.name    MIT
// @host            localhost:8081
// @BasePath        /v1
// @Schemes         http
func main() {
	err := godotenv.Load("./cmd/.env")
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

	//create a channel for listening to OS signals and connecting OS interrupts to the channel
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	serverShutdown := make(chan struct{})

	go func() {
		_ = <-c
		app.Logger.Info("SERVER SHUTDOWN GRACEFULLY")
		app.Close(app.Context)
		_ = app.Application.Shutdown()
		serverShutdown <- struct{}{}
	}()

	switch mode {
	case localServerMode, httpServerMode:
		var (
			httpServer = infrastructure.ServeHTTP(app)
		)

		if err := httpServer.Listen(fmt.Sprintf(":%d", app.Config.Server.AppPort)); err != nil {
			app.Logger.Error("Failed to to start server. Error: ", err)
			panic(err)
		}
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
