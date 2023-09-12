package application

import (
	"github.com/PickHD/LezPay/wallet/internal/v1/controller"
	"github.com/PickHD/LezPay/wallet/internal/v1/repository"
	"github.com/PickHD/LezPay/wallet/internal/v1/service"
)

type Dependency struct {
	HealthCheckController controller.HealthCheckController
	WalletController      *controller.WalletControllerImpl
}

func SetupDependencyInjection(app *App) *Dependency {
	// repository
	healthCheckRepoImpl := repository.NewHealthCheckRepository(app.Context, app.Config, app.Logger, app.Tracer, app.DB, app.Redis)
	walletRepoImpl := repository.NewWalletRepository(app.Context, app.Config, app.Logger, app.Tracer, app.DB, app.Redis)

	// service
	healthCheckSvcImpl := service.NewHealthCheckService(app.Context, app.Config, app.Tracer, healthCheckRepoImpl)
	walletSvcImpl := service.NewWalletService(app.Context, app.Config, app.Tracer, walletRepoImpl)

	// controller
	healthCheckControllerImpl := controller.NewHealthCheckController(app.Context, app.Config, app.Tracer, healthCheckSvcImpl)
	walletControllerImpl := controller.NewWalletController(app.Context, app.Config, app.Tracer, walletSvcImpl)

	return &Dependency{
		HealthCheckController: healthCheckControllerImpl,
		WalletController:      walletControllerImpl,
	}
}
