package application

import (
	"github.com/PickHD/LezPay/merchant/internal/controller"
	"github.com/PickHD/LezPay/merchant/internal/repository"
	"github.com/PickHD/LezPay/merchant/internal/service"
)

type Dependency struct {
	HealthCheckController controller.HealthCheckController
	MerchantController    *controller.MerchantControllerImpl
}

func SetupDependencyInjection(app *App) *Dependency {
	// repository
	healthCheckRepoImpl := repository.NewHealthCheckRepository(app.Context, app.Config, app.Logger, app.Tracer, app.DB, app.Redis)
	merchantRepoImpl := repository.NewMerchantRepository(app.Context, app.Config, app.Logger, app.Tracer, app.DB, app.Redis)

	// service
	healthCheckSvcImpl := service.NewHealthCheckService(app.Context, app.Config, app.Tracer, healthCheckRepoImpl)
	merchantSvcImpl := service.NewMerchantService(app.Context, app.Config, app.Tracer, merchantRepoImpl)

	// controller
	healthCheckControllerImpl := controller.NewHealthCheckController(app.Context, app.Config, app.Tracer, healthCheckSvcImpl)
	merchantControllerImpl := controller.NewMerchantController(app.Context, app.Config, app.Tracer, merchantSvcImpl)

	return &Dependency{
		HealthCheckController: healthCheckControllerImpl,
		MerchantController:    merchantControllerImpl,
	}
}
