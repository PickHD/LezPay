package infrastructure

import (
	"github.com/PickHD/LezPay/wallet/internal/application"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	walletpb "github.com/PickHD/LezPay/wallet/pkg/proto/v1/wallet"
)

func ServeGRPC(app *application.App) *grpc.Server {
	// call register
	return register(app)
}

func register(app *application.App) *grpc.Server {
	var dep = application.SetupDependencyInjection(app)

	reflection.Register(app.GRPC)

	walletpb.RegisterWalletServiceServer(app.GRPC, dep.WalletController)

	return app.GRPC
}
