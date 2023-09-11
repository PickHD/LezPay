package infrastructure

import (
	"github.com/PickHD/LezPay/customer/internal/v1/application"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func ServeGRPC(app *application.App) *grpc.Server {
	// call register
	return register(app)
}

func register(app *application.App) *grpc.Server {
	// var dep = application.SetupDependencyInjection(app)

	reflection.Register(app.GRPC)

	return app.GRPC
}
