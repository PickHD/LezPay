package controller

import (
	"context"
	"strings"

	"github.com/PickHD/LezPay/customer/internal/v1/config"
	"github.com/PickHD/LezPay/customer/internal/v1/helper"
	"github.com/PickHD/LezPay/customer/internal/v1/middleware"
	"github.com/PickHD/LezPay/customer/internal/v1/model"
	"github.com/PickHD/LezPay/customer/internal/v1/service"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	customerpb "github.com/PickHD/LezPay/customer/pkg/proto/v1/customer"
)

type (
	// CustomerController is an interface that has all the function to be implemented inside customer controller
	CustomerController interface {
		// GRPC
		CreateCustomer(ctx context.Context, req *customerpb.CustomerRequest) (*customerpb.CustomerResponse, error)
		UpdateVerifiedCustomer(ctx context.Context, req *customerpb.UpdateVerifiedCustomerRequest) (*customerpb.UpdateVerifiedCustomerResponse, error)
		GetCustomerIDByEmail(ctx context.Context, req *customerpb.GetCustomerIDByEmailRequest) (*customerpb.GetCustomerIDByEmailResponse, error)
		GetCustomerDetailsByEmail(ctx context.Context, req *customerpb.GetCustomerDetailsByEmailRequest) (*customerpb.GetCustomerDetailsByEmailResponse, error)
		UpdateCustomerPasswordByEmail(ctx context.Context, req *customerpb.UpdateCustomerPasswordByEmailRequest) (*customerpb.UpdateCustomerPasswordByEmailResponse, error)

		// REST API
		GetCustomerDashboard(ctx *fiber.Ctx) error
	}

	// CustomerControllerImpl is an app customer struct that consists of all the dependencies needed for customer controller
	CustomerControllerImpl struct {
		Context     context.Context
		Config      *config.Configuration
		Tracer      *trace.TracerProvider
		CustomerSvc service.CustomerService
		customerpb.UnimplementedCustomerServiceServer
	}
)

// NewCustomerController return new instances customer controller
func NewCustomerController(ctx context.Context, config *config.Configuration, tracer *trace.TracerProvider, customerSvc service.CustomerService) *CustomerControllerImpl {
	return &CustomerControllerImpl{
		Context:     ctx,
		Config:      config,
		Tracer:      tracer,
		CustomerSvc: customerSvc,
	}
}

func (cc *CustomerControllerImpl) CreateCustomer(ctx context.Context, req *customerpb.CustomerRequest) (*customerpb.CustomerResponse, error) {
	tr := cc.Tracer.Tracer("Customer-CreateCustomer Controller")
	_, span := tr.Start(ctx, "Start CreateCustomer")
	defer span.End()

	newCustomer := model.CreateCustomerRequest{
		FullName:    req.GetFullName(),
		Email:       req.GetEmail(),
		PhoneNumber: req.GetPhoneNumber(),
		Password:    req.GetPassword(),
		Pin:         req.GetPin(),
	}

	data, err := cc.CustomerSvc.CreateCustomer(&newCustomer)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed Create Customer %s", err.Error())
	}

	return &customerpb.CustomerResponse{
		Id:         uint64(data.ID),
		IsVerified: data.IsVerified,
	}, nil
}

func (cc *CustomerControllerImpl) UpdateVerifiedCustomer(ctx context.Context, req *customerpb.UpdateVerifiedCustomerRequest) (*customerpb.UpdateVerifiedCustomerResponse, error) {
	tr := cc.Tracer.Tracer("Customer-UpdateVerifiedCustomer Controller")
	_, span := tr.Start(ctx, "Start UpdateVerifiedCustomer")
	defer span.End()

	data, err := cc.CustomerSvc.UpdateVerifiedCustomer(req.GetEmail())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed Update Verified Customer %s", err.Error())
	}

	return &customerpb.UpdateVerifiedCustomerResponse{
		IsVerified: data.IsVerified,
	}, nil
}

func (cc *CustomerControllerImpl) GetCustomerIDByEmail(ctx context.Context, req *customerpb.GetCustomerIDByEmailRequest) (*customerpb.GetCustomerIDByEmailResponse, error) {
	tr := cc.Tracer.Tracer("Customer-GetCustomerIDByEmail Controller")
	_, span := tr.Start(ctx, "Start GetCustomerIDByEmail")
	defer span.End()

	getCustomer := &model.GetCustomerIDByEmailRequest{
		Email: req.GetEmail(),
	}

	data, err := cc.CustomerSvc.GetCustomerIDByEmail(getCustomer)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed Get Customer ID By Email %s", err.Error())
	}

	return &customerpb.GetCustomerIDByEmailResponse{
		CustomerId: uint64(data.CustomerID),
	}, nil
}

func (cc *CustomerControllerImpl) GetCustomerDetailsByEmail(ctx context.Context, req *customerpb.GetCustomerDetailsByEmailRequest) (*customerpb.GetCustomerDetailsByEmailResponse, error) {
	tr := cc.Tracer.Tracer("Customer-GetCustomerDetailsByEmail Controller")
	_, span := tr.Start(ctx, "Start GetCustomerDetailsByEmail")
	defer span.End()

	getCustomer := &model.GetCustomerDetailsByEmailRequest{
		Email: req.GetEmail(),
	}

	data, err := cc.CustomerSvc.GetCustomerDetailsByEmail(getCustomer)
	if err != nil {
		if strings.Contains(err.Error(), string(model.NotFound)) {
			return nil, status.Error(codes.NotFound, "Email Not Found")
		}

		return nil, status.Errorf(codes.Internal, "Failed Get Customer Details By Email %s", err.Error())
	}

	return &customerpb.GetCustomerDetailsByEmailResponse{
		Id:          uint64(data.ID),
		FullName:    data.FullName,
		Email:       data.Email,
		PhoneNumber: data.PhoneNumber,
		Password:    data.Password,
		Pin:         data.Pin,
	}, nil
}

func (cc *CustomerControllerImpl) UpdateCustomerPasswordByEmail(ctx context.Context, req *customerpb.UpdateCustomerPasswordByEmailRequest) (*customerpb.UpdateCustomerPasswordByEmailResponse, error) {
	tr := cc.Tracer.Tracer("Customer-UpdateCustomerPasswordByEmail Controller")
	_, span := tr.Start(ctx, "Start UpdateCustomerPasswordByEmail")
	defer span.End()

	updatePasswordReq := model.UpdateCustomerPasswordByEmailRequest{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	data, err := cc.CustomerSvc.UpdateCustomerPasswordByEmail(&updatePasswordReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed Update Password %s", err.Error())
	}

	return &customerpb.UpdateCustomerPasswordByEmailResponse{
		Email: data.Email,
	}, nil
}

// Check godoc
// @Summary      Get Customer Dashboard
// @Tags         Customer
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Authorization Bearer <Place Access Token Here>"
// @Success      200  {object}  helper.BaseResponse
// @Failure      404  {object}  helper.BaseResponse
// @Failure      500  {object}  helper.BaseResponse
// @Router       /dashboard [get]
func (cc *CustomerControllerImpl) GetCustomerDashboard(ctx *fiber.Ctx) error {
	tr := cc.Tracer.Tracer("Customer-GetCustomerDashboard Controller")
	_, span := tr.Start(cc.Context, "Start GetCustomerDashboard")
	defer span.End()

	data := ctx.Locals(model.KeyJWTValidAccess)
	decodedData, err := middleware.Extract(data)
	if err != nil {
		return helper.NewResponses[any](ctx, fiber.StatusInternalServerError, err.Error(), nil, err, nil)
	}

	result, err := cc.CustomerSvc.GetCustomerDashboard(decodedData.UserID)
	if err != nil {
		return helper.NewResponses[any](ctx, fiber.StatusInternalServerError, "Failed get Dashboard", nil, err, nil)
	}

	return helper.NewResponses[any](ctx, fiber.StatusOK, "Success get Dashboard", result, nil, nil)
}
