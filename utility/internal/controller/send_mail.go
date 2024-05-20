package controller

import (
	"context"
	"net/http"
	"strings"

	"github.com/PickHD/LezPay/utility/internal/config"
	"github.com/PickHD/LezPay/utility/internal/helper"
	"github.com/PickHD/LezPay/utility/internal/model"
	"github.com/PickHD/LezPay/utility/internal/service"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/sdk/trace"
)

type (
	// SendMailController is an interface that has all the function to be implemented inside send mail controller
	SendMailController interface {
		Send(ctx *gin.Context)
	}

	// SendMailControllerImpl is an app send mail struct that consists of all the dependencies needed for send mail controller
	SendMailControllerImpl struct {
		Context     context.Context
		Config      *config.Configuration
		Tracer      *trace.TracerProvider
		SendMailSvc service.SendMailService
	}
)

// NewSendMailController return new instances send mail controller
func NewSendMailController(ctx context.Context, config *config.Configuration, tracer *trace.TracerProvider, sendMailSvc service.SendMailService) *SendMailControllerImpl {
	return &SendMailControllerImpl{
		Context:     ctx,
		Config:      config,
		Tracer:      tracer,
		SendMailSvc: sendMailSvc,
	}
}

// Check godoc
// @Summary      Send Mail
// @Tags         Mail
// @Accept       json
// @Produce      json
// @Param        user body model.SendMailRequest true "send mail"
// @Success      200  {object}  helper.BaseResponse
// @Failure      400  {object}  helper.BaseResponse
// @Failure      500  {object}  helper.BaseResponse
// @Router       /send [post]
func (sc *SendMailControllerImpl) Send(ctx *gin.Context) {
	var req model.SendMailRequest

	tr := sc.Tracer.Tracer("Utility-Send Controller")
	_, span := tr.Start(ctx, "Start Send")
	defer span.End()

	if err := ctx.BindJSON(&req); err != nil {
		helper.NewResponses[any](ctx, http.StatusBadRequest, "Invalid request", nil, err, nil)
		return
	}

	err := sc.SendMailSvc.Send(ctx, &req)
	if err != nil {
		if strings.Contains(err.Error(), string(model.Validation)) {
			helper.NewResponses[any](ctx, http.StatusBadRequest, err.Error(), nil, err, nil)
			return
		}

		helper.NewResponses[any](ctx, http.StatusInternalServerError, "Failed send mail", nil, err, nil)
		return
	}

	helper.NewResponses[any](ctx, http.StatusOK, "Success send mails", nil, nil, nil)
}
