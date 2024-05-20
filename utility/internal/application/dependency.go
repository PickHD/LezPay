package application

import (
	"github.com/PickHD/LezPay/utility/internal/controller"
	"github.com/PickHD/LezPay/utility/internal/repository"
	"github.com/PickHD/LezPay/utility/internal/service"
)

type Dependency struct {
	HealthCheckController controller.HealthCheckController
	SendMailController    controller.SendMailController
}

func SetupDependencyInjection(app *App) *Dependency {
	// repository
	healthCheckRepoImpl := repository.NewHealthCheckRepository(app.Context, app.Config, app.Logger, app.Tracer, app.Redis)
	sendMailRepoImpl := repository.NewSendMailRepository(app.Context, app.Config, app.Logger, app.Tracer, app.Mailer)

	// service
	healthCheckSvcImpl := service.NewHealthCheckService(app.Context, app.Config, app.Tracer, healthCheckRepoImpl)
	sendMailSvcImpl := service.NewSendMailService(app.Context, app.Config, app.Logger, app.Tracer, sendMailRepoImpl)

	// controller
	healthCheckControllerImpl := controller.NewHealthCheckController(app.Context, app.Config, app.Tracer, healthCheckSvcImpl)
	sendMailControllerImpl := controller.NewSendMailController(app.Context, app.Config, app.Tracer, sendMailSvcImpl)

	return &Dependency{
		HealthCheckController: healthCheckControllerImpl,
		SendMailController:    sendMailControllerImpl,
	}
}
