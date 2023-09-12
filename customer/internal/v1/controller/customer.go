package controller

import (
	"context"

	"github.com/PickHD/LezPay/customer/internal/v1/config"
	"github.com/PickHD/LezPay/customer/internal/v1/model"
	"github.com/PickHD/LezPay/customer/internal/v1/service"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	customerpb "github.com/PickHD/LezPay/customer/pkg/proto/v1/customer"
)

type (
	// CustomerController is an interface that has all the function to be implemented inside customer controller
	CustomerController interface {
		CreateCustomer(ctx context.Context, req *customerpb.CustomerRequest) (*customerpb.CustomerResponse, error)
		UpdateVerifiedCustomer(ctx context.Context, req *customerpb.UpdateVerifiedCustomerRequest) (*customerpb.UpdateVerifiedCustomerResponse, error)
		GetCustomerIDByEmail(ctx context.Context, req *customerpb.GetCustomerIDByEmailRequest) (*customerpb.GetCustomerIDByEmailResponse, error)
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
