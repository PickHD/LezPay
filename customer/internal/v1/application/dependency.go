package application

import (
	"github.com/PickHD/LezPay/customer/internal/v1/controller"
	"github.com/PickHD/LezPay/customer/internal/v1/repository"
	"github.com/PickHD/LezPay/customer/internal/v1/service"

	walletpb "github.com/PickHD/LezPay/customer/pkg/proto/v1/wallet"
)

type Dependency struct {
	HealthCheckController controller.HealthCheckController
	CustomerController    *controller.CustomerControllerImpl
}

func SetupDependencyInjection(app *App) *Dependency {
	walletServiceClient := walletpb.NewWalletServiceClient(app.WalletGRPC)

	// repository
	healthCheckRepoImpl := repository.NewHealthCheckRepository(app.Context, app.Config, app.Logger, app.Tracer, app.DB, app.Redis)
	customerRepoImpl := repository.NewCustomerRepository(app.Context, app.Config, app.Logger, app.Tracer, app.DB, app.Redis)

	// service
	healthCheckSvcImpl := service.NewHealthCheckService(app.Context, app.Config, app.Tracer, healthCheckRepoImpl)
	customerSvcImpl := service.NewCustomerService(app.Context, app.Config, app.Tracer, customerRepoImpl, walletServiceClient)

	// controller
	healthCheckControllerImpl := controller.NewHealthCheckController(app.Context, app.Config, app.Tracer, healthCheckSvcImpl)
	customerControllerImpl := controller.NewCustomerController(app.Context, app.Config, app.Tracer, customerSvcImpl)

	return &Dependency{
		HealthCheckController: healthCheckControllerImpl,
		CustomerController:    customerControllerImpl,
	}
}
