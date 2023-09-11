package application

import (
	"github.com/PickHD/LezPay/wallet/internal/v1/controller"
	"github.com/PickHD/LezPay/wallet/internal/v1/repository"
	"github.com/PickHD/LezPay/wallet/internal/v1/service"
)

type Dependency struct {
	HealthCheckController controller.HealthCheckController
}

func SetupDependencyInjection(app *App) *Dependency {
	// repository
	healthCheckRepoImpl := repository.NewHealthCheckRepository(app.Context, app.Config, app.Logger, app.Tracer, app.DB, app.Redis)

	// service
	healthCheckSvcImpl := service.NewHealthCheckService(app.Context, app.Config, app.Tracer, healthCheckRepoImpl)

	// controller
	healthCheckControllerImpl := controller.NewHealthCheckController(app.Context, app.Config, app.Tracer, healthCheckSvcImpl)

	return &Dependency{
		HealthCheckController: healthCheckControllerImpl,
	}
}
