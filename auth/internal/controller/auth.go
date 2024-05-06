package controller

import (
	"context"
	"net/http"
	"strings"

	"github.com/PickHD/LezPay/auth/internal/config"
	"github.com/PickHD/LezPay/auth/internal/helper"
	"github.com/PickHD/LezPay/auth/internal/model"
	"github.com/PickHD/LezPay/auth/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/sdk/trace"
)

type (
	// AuthController is an interface that has all the function to be implemented inside auth controller
	AuthController interface {
		Register(ctx *gin.Context)
		VerifyRegister(ctx *gin.Context)
		Login(ctx *gin.Context)
		ForgotPassword(ctx *gin.Context)
		VerifyForgotPassword(ctx *gin.Context)
		ResetPassword(ctx *gin.Context)
	}

	// AuthControllerImpl is an app auth struct that consists of all the dependencies needed for auth controller
	AuthControllerImpl struct {
		Context context.Context
		Config  *config.Configuration
		Logger  *logrus.Logger
		Tracer  *trace.TracerProvider
		AuthSvc service.AuthService
	}
)

// NewAuthController return new instances auth controller
func NewAuthController(ctx context.Context, config *config.Configuration, logger *logrus.Logger, tracer *trace.TracerProvider, authSvc service.AuthService) *AuthControllerImpl {
	return &AuthControllerImpl{
		Context: ctx,
		Config:  config,
		Logger:  logger,
		Tracer:  tracer,
		AuthSvc: authSvc,
	}
}

// Check godoc
// @Summary      Register customer/merchant
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        user body model.RegisterRequest true "register customer/merchant"
// @Success      201  {object}  helper.BaseResponse
// @Failure      400  {object}  helper.BaseResponse
// @Failure      500  {object}  helper.BaseResponse
// @Router       /register [post]
func (ac *AuthControllerImpl) Register(ctx *gin.Context) {
	var req model.RegisterRequest

	tr := ac.Tracer.Tracer("Auth-Register Controller")
	_, span := tr.Start(ctx, "Start Register")
	defer span.End()

	if err := ctx.BindJSON(&req); err != nil {
		helper.NewResponses[any](ctx, http.StatusBadRequest, "Invalid request", nil, err, nil)
		return
	}

	err := ac.AuthSvc.RegisterCustomerOrMerchant(&req)
	if err != nil {
		if strings.Contains(err.Error(), string(model.Validation)) {
			helper.NewResponses[any](ctx, http.StatusBadRequest, err.Error(), nil, err, nil)
			return
		}

		helper.NewResponses[any](ctx, http.StatusInternalServerError, "Failed register user", nil, err, nil)
		return
	}

	helper.NewResponses[any](ctx, http.StatusCreated, "Success register, please check email for further verification", nil, nil, nil)
}

// Check godoc
// @Summary      Verify Register customer/merchant
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        code  query   string  true  "Code Verification"
// @Param        user_type query string true "customer|merchant"
// @Success      200  {object}  helper.BaseResponse
// @Failure      400  {object}  helper.BaseResponse
// @Failure      404  {object}  helper.BaseResponse
// @Failure      500  {object}  helper.BaseResponse
// @Router       /register/verify [get]
func (ac *AuthControllerImpl) VerifyRegister(ctx *gin.Context) {
	tr := ac.Tracer.Tracer("Auth-VerifyRegister Controller")
	_, span := tr.Start(ctx, "Start VerifyRegister")
	defer span.End()

	getCode := ctx.Query("code")
	getUserType := ctx.Query("user_type")

	if getCode == "" || getUserType == "" {
		helper.NewResponses[any](ctx, http.StatusBadRequest, "Code Required", nil, nil, nil)
		return
	}

	if getUserType != string(model.Merchant) && getUserType != string(model.Customer) {
		helper.NewResponses[any](ctx, http.StatusBadRequest, "user_type must customer/merchant", nil, nil, nil)
		return
	}

	data, err := ac.AuthSvc.VerifyRegisterOrForgotPasswordCode(ctx, getCode, model.UserType(getUserType), model.RegisterVerification)
	if err != nil {
		helper.NewResponses[any](ctx, http.StatusInternalServerError, "Failed Verify Code", nil, err, nil)
		return
	}

	helper.NewResponses[any](ctx, http.StatusOK, "Verify success, Redirecting to Login Page....", data, err, nil)
}

// Check godoc
// @Summary      Login customer/merchant
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        user body model.LoginRequest true "login"
// @Success      200  {object}  helper.BaseResponse
// @Failure      400  {object}  helper.BaseResponse
// @Failure      404  {object}  helper.BaseResponse
// @Failure      500  {object}  helper.BaseResponse
// @Router       /login [post]
func (ac *AuthControllerImpl) Login(ctx *gin.Context) {
	var req model.LoginRequest

	tr := ac.Tracer.Tracer("Auth-Login Controller")
	_, span := tr.Start(ctx, "Start Login")
	defer span.End()

	if err := ctx.BindJSON(&req); err != nil {
		helper.NewResponses[any](ctx, http.StatusBadRequest, "Invalid request", nil, err, nil)
		return
	}

	data, err := ac.AuthSvc.LoginCustomerOrMerchant(ctx, &req)
	if err != nil {
		if strings.Contains(err.Error(), string(model.Validation)) {
			helper.NewResponses[any](ctx, http.StatusBadRequest, err.Error(), nil, err, nil)
			return
		}

		if strings.Contains(err.Error(), string(model.NotFound)) {
			helper.NewResponses[any](ctx, http.StatusNotFound, err.Error(), nil, err, nil)
			return
		}

		helper.NewResponses[any](ctx, http.StatusInternalServerError, "Failed login customer/merchant", nil, err, nil)
		return
	}

	helper.NewResponses[any](ctx, http.StatusOK, "Success login", data, nil, nil)
}

// Check godoc
// @Summary      Forgot Password customer/merchant
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        forgotPassword body model.ForgotPasswordRequest true "forgot password customer/merchant"
// @Success      200  {object}  helper.BaseResponse
// @Failure      400  {object}  helper.BaseResponse
// @Failure      404  {object}  helper.BaseResponse
// @Failure      500  {object}  helper.BaseResponse
// @Router       /forgot-password [post]
func (ac *AuthControllerImpl) ForgotPassword(ctx *gin.Context) {
	var req model.ForgotPasswordRequest

	tr := ac.Tracer.Tracer("Auth-ForgotPassword Controller")
	_, span := tr.Start(ctx, "Start ForgotPassword")
	defer span.End()

	if err := ctx.BindJSON(&req); err != nil {
		helper.NewResponses[any](ctx, http.StatusBadRequest, "Invalid request", req, err, nil)
		return
	}

	if string(req.UserType) != string(model.Merchant) && string(req.UserType) != string(model.Customer) {
		helper.NewResponses[any](ctx, http.StatusBadRequest, "user_type must customer/merchant", nil, nil, nil)
		return
	}

	err := ac.AuthSvc.ForgotPasswordCustomerOrMerchant(ctx, &req)
	if err != nil {
		if strings.Contains(err.Error(), string(model.Validation)) {
			helper.NewResponses[any](ctx, http.StatusBadRequest, err.Error(), req.Email, err, nil)
			return
		}

		if strings.Contains(err.Error(), string(model.NotFound)) {
			helper.NewResponses[any](ctx, http.StatusNotFound, err.Error(), req.Email, err, nil)
			return
		}

		helper.NewResponses[any](ctx, http.StatusInternalServerError, "Failed request forgot password", nil, err, nil)
		return
	}

	helper.NewResponses[any](ctx, http.StatusOK, "Success sent forgot password verification ,please check your email", nil, nil, nil)
}

// Check godoc
// @Summary      Verify Forgot Password customer/merchant
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        code  query   string  true  "Code Verification"
// @Param        user_type query string true "customer|merchant"
// @Success      200  {object}  helper.BaseResponse
// @Failure      400  {object}  helper.BaseResponse
// @Failure      404  {object}  helper.BaseResponse
// @Failure      500  {object}  helper.BaseResponse
// @Router       /forgot-password/verify [get]
func (ac *AuthControllerImpl) VerifyForgotPassword(ctx *gin.Context) {
	tr := ac.Tracer.Tracer("Auth-VerifyForgotPassword Controller")
	_, span := tr.Start(ctx, "Start VerifyForgotPassword")
	defer span.End()

	getCode := ctx.Query("code")
	getUserType := ctx.Query("user_type")

	if getCode == "" || getUserType == "" {
		helper.NewResponses[any](ctx, http.StatusBadRequest, "Code Required", nil, nil, nil)
		return
	}

	if getUserType != string(model.Merchant) && getUserType != string(model.Customer) {
		helper.NewResponses[any](ctx, http.StatusBadRequest, "user_type must customer/merchant", nil, nil, nil)
		return
	}

	data, err := ac.AuthSvc.VerifyRegisterOrForgotPasswordCode(ctx, getCode, model.UserType(getUserType), model.ForgotPasswordVerification)
	if err != nil {
		helper.NewResponses[any](ctx, http.StatusInternalServerError, "Failed Verify Code", nil, err, nil)
		return
	}

	helper.NewResponses[any](ctx, http.StatusOK, "Verify success, Redirecting to Reset Password Page....", data, err, nil)
}

// Check godoc
// @Summary      Reset Password customer/merchant
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        code  query   string  true  "Code Verification"
// @Param        user_type query string true "customer|merchant"
// @Param        forgotPassword body model.ResetPasswordRequest true "reset password customer/merchant"
// @Success      200  {object}  helper.BaseResponse
// @Failure      400  {object}  helper.BaseResponse
// @Failure      404  {object}  helper.BaseResponse
// @Failure      500  {object}  helper.BaseResponse
// @Router       /reset-password [put]
func (ac *AuthControllerImpl) ResetPassword(ctx *gin.Context) {
	tr := ac.Tracer.Tracer("Auth-ResetPassword Controller")
	_, span := tr.Start(ctx, "Start ResetPassword")
	defer span.End()

	var req model.ResetPasswordRequest

	getCode := ctx.Query("code")
	getUserType := ctx.Query("user_type")

	if getCode == "" || getUserType == "" {
		helper.NewResponses[any](ctx, http.StatusBadRequest, "Code Required", nil, nil, nil)
		return
	}

	if getUserType != string(model.Merchant) && getUserType != string(model.Customer) {
		helper.NewResponses[any](ctx, http.StatusBadRequest, "user_type must customer/merchant", nil, nil, nil)
		return
	}

	if err := ctx.BindJSON(&req); err != nil {
		helper.NewResponses[any](ctx, http.StatusBadRequest, "Invalid request", req, err, nil)
		return
	}

	err := ac.AuthSvc.ResetPasswordCustomerOrMerchant(ctx, &req, getCode, model.UserType(getUserType))
	if err != nil {
		if strings.Contains(err.Error(), string(model.Validation)) {
			helper.NewResponses[any](ctx, http.StatusBadRequest, err.Error(), nil, err, nil)
			return
		}

		if strings.Contains(err.Error(), string(model.NotFound)) {
			helper.NewResponses[any](ctx, http.StatusNotFound, err.Error(), nil, err, nil)
			return
		}

		helper.NewResponses[any](ctx, http.StatusInternalServerError, "Failed reset password", nil, err, nil)
		return
	}

	helper.NewResponses[any](ctx, http.StatusOK, "Success reset password", nil, nil, nil)
}
