package main

import (
	"context"
	"os"
	"runtime"

	"github.com/PickHD/LezPay/payment/internal/application"
	"github.com/PickHD/LezPay/payment/internal/infrastructure"
	"github.com/joho/godotenv"
)

const (
	localServerMode = "local"
	httpServerMode  = "http"
	consumerMode    = "consumer"
	grpcMode        = "grpc"
)

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

	switch mode {
	case localServerMode, httpServerMode:
	case grpcMode:
	case consumerMode:
		// Make a channel to receive messages into infinite loop.
		forever := make(chan bool)

		topics := []string{app.Config.Kafka.TopicRequestPayment}

		for _, q := range topics {
			go infrastructure.ConsumeMessages(app, q)
		}

		<-forever
	}
}
