package service

import (
	"context"

	"github.com/PickHD/LezPay/customer/internal/v1/config"
	"github.com/PickHD/LezPay/customer/internal/v1/model"
	"github.com/PickHD/LezPay/customer/internal/v1/repository"
	"go.opentelemetry.io/otel/sdk/trace"

	walletpb "github.com/PickHD/LezPay/customer/pkg/proto/v1/wallet"
)

type (
	// CustomerService is an interface that has all the function to be implemented inside customer service
	CustomerService interface {
		CreateCustomer(req *model.CreateCustomerRequest) (*model.CreateCustomerResponse, error)
		UpdateVerifiedCustomer(email string) (*model.UpdateVerifiedCustomerResponse, error)
		GetCustomerIDByEmail(req *model.GetCustomerIDByEmailRequest) (*model.GetCustomerIDByEmailResponse, error)
		GetCustomerDetailsByEmail(req *model.GetCustomerDetailsByEmailRequest) (*model.GetCustomerDetailsByEmailResponse, error)
		UpdateCustomerPasswordByEmail(req *model.UpdateCustomerPasswordByEmailRequest) (*model.UpdateCustomerPasswordByEmailResponse, error)
		GetCustomerDashboard(customerID uint64) (*model.GetCustomerDashboardResponse, error)
	}

	// CustomerServiceImpl is an app customer struct that consists of all the dependencies needed for customer service
	CustomerServiceImpl struct {
		Context       context.Context
		Config        *config.Configuration
		Tracer        *trace.TracerProvider
		CustomerRepo  repository.CustomerRepository
		WalletClients walletpb.WalletServiceClient
	}
)

// NewCustomerService return new instances customer service
func NewCustomerService(ctx context.Context, config *config.Configuration, tracer *trace.TracerProvider, customerRepo repository.CustomerRepository, walletClient walletpb.WalletServiceClient) *CustomerServiceImpl {
	return &CustomerServiceImpl{
		Context:       ctx,
		Config:        config,
		Tracer:        tracer,
		CustomerRepo:  customerRepo,
		WalletClients: walletClient,
	}
}

func (cs *CustomerServiceImpl) CreateCustomer(req *model.CreateCustomerRequest) (*model.CreateCustomerResponse, error) {
	tr := cs.Tracer.Tracer("Customer-CreateCustomer Service")
	_, span := tr.Start(cs.Context, "Start CreateCustomer")
	defer span.End()

	Id, isVerified, err := cs.CustomerRepo.CreateCustomer(req)
	if err != nil {
		return nil, err
	}

	return &model.CreateCustomerResponse{
		ID:         Id,
		IsVerified: isVerified,
	}, nil
}

func (cs *CustomerServiceImpl) UpdateVerifiedCustomer(email string) (*model.UpdateVerifiedCustomerResponse, error) {
	tr := cs.Tracer.Tracer("Customer-UpdateVerifiedCustomer Service")
	_, span := tr.Start(cs.Context, "Start UpdateVerifiedCustomer")
	defer span.End()

	isVerified, err := cs.CustomerRepo.UpdateVerifiedCustomer(email)
	if err != nil {
		return nil, err
	}

	return &model.UpdateVerifiedCustomerResponse{
		IsVerified: isVerified,
	}, nil
}

func (cs *CustomerServiceImpl) GetCustomerIDByEmail(req *model.GetCustomerIDByEmailRequest) (*model.GetCustomerIDByEmailResponse, error) {
	tr := cs.Tracer.Tracer("Customer-GetCustomerIDByEmail Service")
	_, span := tr.Start(cs.Context, "Start GetCustomerIDByEmail")
	defer span.End()

	getCustomerID, err := cs.CustomerRepo.GetCustomerIDByEmail(req)
	if err != nil {
		return nil, err
	}

	return &model.GetCustomerIDByEmailResponse{
		CustomerID: getCustomerID,
	}, nil
}

func (cs *CustomerServiceImpl) GetCustomerDetailsByEmail(req *model.GetCustomerDetailsByEmailRequest) (*model.GetCustomerDetailsByEmailResponse, error) {
	tr := cs.Tracer.Tracer("Customer-GetCustomerDetailsByEmail Service")
	_, span := tr.Start(cs.Context, "Start GetCustomerDetailsByEmail")
	defer span.End()

	getCustomer, err := cs.CustomerRepo.GetCustomerDetailsByEmail(req)
	if err != nil {
		return nil, err
	}

	return &model.GetCustomerDetailsByEmailResponse{
		ID:          getCustomer.ID,
		FullName:    getCustomer.FullName,
		PhoneNumber: getCustomer.PhoneNumber,
		Email:       getCustomer.Email,
		Password:    getCustomer.Password,
		Pin:         getCustomer.Pin,
	}, nil
}

func (cs *CustomerServiceImpl) UpdateCustomerPasswordByEmail(req *model.UpdateCustomerPasswordByEmailRequest) (*model.UpdateCustomerPasswordByEmailResponse, error) {
	tr := cs.Tracer.Tracer("Customer-UpdateCustomerPasswordByEmail Service")
	_, span := tr.Start(cs.Context, "Start UpdateCustomerPasswordByEmail")
	defer span.End()

	data, err := cs.CustomerRepo.UpdateCustomerPasswordByEmail(req)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (cs *CustomerServiceImpl) GetCustomerDashboard(customerID uint64) (*model.GetCustomerDashboardResponse, error) {
	data, err := cs.WalletClients.GetCustomerWallet(cs.Context, &walletpb.GetCustomerWalletRequest{
		CustomerId: customerID})
	if err != nil {
		return nil, err
	}

	return &model.GetCustomerDashboardResponse{
		WalletID: data.GetId(),
		Balance:  data.GetBalance(),
	}, nil
}
