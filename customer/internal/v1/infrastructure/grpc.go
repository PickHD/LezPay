package infrastructure

import (
	"github.com/PickHD/LezPay/customer/internal/v1/application"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	customerpb "github.com/PickHD/LezPay/customer/pkg/proto/v1/customer"
)

func ServeGRPC(app *application.App) *grpc.Server {
	// call register
	return register(app)
}

func register(app *application.App) *grpc.Server {
	var dep = application.SetupDependencyInjection(app)

	reflection.Register(app.GRPC)

	customerpb.RegisterCustomerServiceServer(app.GRPC, dep.CustomerController)

	return app.GRPC
}
