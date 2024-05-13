package application

import (
	"github.com/PickHD/LezPay/utility/internal/controller"
	"github.com/PickHD/LezPay/utility/internal/repository"
	"github.com/PickHD/LezPay/utility/internal/service"
)

type Dependency struct {
	HealthCheckController controller.HealthCheckController
}

func SetupDependencyInjection(app *App) *Dependency {
	// repository
	healthCheckRepoImpl := repository.NewHealthCheckRepository(app.Context, app.Config, app.Logger, app.Tracer, app.Redis)

	// service
	healthCheckSvcImpl := service.NewHealthCheckService(app.Context, app.Config, app.Tracer, healthCheckRepoImpl)

	// controller
	healthCheckControllerImpl := controller.NewHealthCheckController(app.Context, app.Config, app.Tracer, healthCheckSvcImpl)

	return &Dependency{
		HealthCheckController: healthCheckControllerImpl,
	}
}
