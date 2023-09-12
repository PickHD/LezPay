package application

import (
	"github.com/PickHD/LezPay/customer/internal/v1/controller"
	"github.com/PickHD/LezPay/customer/internal/v1/repository"
	"github.com/PickHD/LezPay/customer/internal/v1/service"
)

type Dependency struct {
	HealthCheckController controller.HealthCheckController
	CustomerController    *controller.CustomerControllerImpl
}

func SetupDependencyInjection(app *App) *Dependency {
	// repository
	healthCheckRepoImpl := repository.NewHealthCheckRepository(app.Context, app.Config, app.Logger, app.Tracer, app.DB, app.Redis)
	customerRepoImpl := repository.NewCustomerRepository(app.Context, app.Config, app.Logger, app.Tracer, app.DB, app.Redis)

	// service
	healthCheckSvcImpl := service.NewHealthCheckService(app.Context, app.Config, app.Tracer, healthCheckRepoImpl)
	customerSvcImpl := service.NewCustomerService(app.Context, app.Config, app.Tracer, customerRepoImpl)

	// controller
	healthCheckControllerImpl := controller.NewHealthCheckController(app.Context, app.Config, app.Tracer, healthCheckSvcImpl)
	customerControllerImpl := controller.NewCustomerController(app.Context, app.Config, app.Tracer, customerSvcImpl)

	return &Dependency{
		HealthCheckController: healthCheckControllerImpl,
		CustomerController:    customerControllerImpl,
	}
}
