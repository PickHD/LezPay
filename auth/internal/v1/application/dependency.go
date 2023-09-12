package application

import (
	"github.com/PickHD/LezPay/auth/internal/v1/controller"
	"github.com/PickHD/LezPay/auth/internal/v1/repository"
	"github.com/PickHD/LezPay/auth/internal/v1/service"

	customerpb "github.com/PickHD/LezPay/auth/pkg/proto/v1/customer"
	merchantpb "github.com/PickHD/LezPay/auth/pkg/proto/v1/merchant"
	walletpb "github.com/PickHD/LezPay/auth/pkg/proto/v1/wallet"
)

type Dependency struct {
	HealthCheckController controller.HealthCheckController
	AuthController        controller.AuthController
}

func SetupDependencyInjection(app *App) *Dependency {
	customerServiceClient := customerpb.NewCustomerServiceClient(app.CustomerGRPC)
	merchantServiceClient := merchantpb.NewMerchantServiceClient(app.MerchantGRPC)
	walletServiceClient := walletpb.NewWalletServiceClient(app.WalletGRPC)

	// repository
	healthCheckRepoImpl := repository.NewHealthCheckRepository(app.Context, app.Config, app.Logger, app.Tracer, app.DB, app.Redis)
	authRepoImpl := repository.NewAuthRepository(app.Context, app.Config, app.Logger, app.Redis, app.Tracer, app.Mailer)

	// service
	healthCheckSvcImpl := service.NewHealthCheckService(app.Context, app.Config, app.Tracer, healthCheckRepoImpl)
	authSvcImpl := service.NewAuthService(app.Context, app.Config, app.Logger, app.Tracer, authRepoImpl, customerServiceClient, merchantServiceClient, walletServiceClient)

	// controller
	healthCheckControllerImpl := controller.NewHealthCheckController(app.Context, app.Config, app.Tracer, healthCheckSvcImpl)
	authControllerImpl := controller.NewAuthController(app.Context, app.Config, app.Logger, app.Tracer, authSvcImpl)

	return &Dependency{
		HealthCheckController: healthCheckControllerImpl,
		AuthController:        authControllerImpl,
	}
}
