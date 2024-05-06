package infrastructure

import (
	"github.com/PickHD/LezPay/merchant/internal/application"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	merchantpb "github.com/PickHD/LezPay/merchant/pkg/proto/v1/merchant"
)

func ServeGRPC(app *application.App) *grpc.Server {
	// call register
	return register(app)
}

func register(app *application.App) *grpc.Server {
	var dep = application.SetupDependencyInjection(app)

	reflection.Register(app.GRPC)

	merchantpb.RegisterMerchantServiceServer(app.GRPC, dep.MerchantController)

	return app.GRPC
}
