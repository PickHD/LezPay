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
		TopupWalletCustomer(customerID uint64, req *model.TopupWalletCustomerRequest) (*model.TopupWalletCustomerResponse, error)
	}

	// CustomerServiceImpl is an app customer struct that consists of all the dependencies needed for customer service
	CustomerServiceImpl struct {
		Context          context.Context
		Config           *config.Configuration
		Tracer           *trace.TracerProvider
		CustomerRepo     repository.CustomerRepository
		WalletClients    walletpb.WalletServiceClient
		WalletTrxClients walletpb.TransactionServiceClient
	}
)

// NewCustomerService return new instances customer service
func NewCustomerService(ctx context.Context, config *config.Configuration, tracer *trace.TracerProvider, customerRepo repository.CustomerRepository, walletClient walletpb.WalletServiceClient, walletTrxClient walletpb.TransactionServiceClient) *CustomerServiceImpl {
	return &CustomerServiceImpl{
		Context:          ctx,
		Config:           config,
		Tracer:           tracer,
		CustomerRepo:     customerRepo,
		WalletClients:    walletClient,
		WalletTrxClients: walletTrxClient,
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
	tr := cs.Tracer.Tracer("Customer-GetCustomerDashboard Service")
	_, span := tr.Start(cs.Context, "Start GetCustomerDashboard")
	defer span.End()

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

func (cs *CustomerServiceImpl) TopupWalletCustomer(customerID uint64, req *model.TopupWalletCustomerRequest) (*model.TopupWalletCustomerResponse, error) {
	/*
		TODO:
		1. validate topup request x
		2. implement grpc & call grpc to create topup transaction with status pending, return transaction_id x
		3. construct proto message TopupTransactionMessage with payload : transaction_id,customer_id,amount,payment_channel_id,EventName (PaymentRequested)
		4. create repo ProduceTopupTransaction : convert proto message to byte and send data to kafka topic (topup-transaction)
		5. return response transaction_id with status pending
	*/
	tr := cs.Tracer.Tracer("Customer-TopupWalletCustomer Service")
	_, span := tr.Start(cs.Context, "Start TopupWalletCustomer")
	defer span.End()

	err := cs.validateTopupWalletCustomerRequest(req)
	if err != nil {
		return nil, err
	}

	data, err := cs.WalletTrxClients.CreateTransaction(cs.Context, &walletpb.CreateTransactionRequest{
		CustomerId:       customerID,
		Amount:           req.Amount,
		PaymentChannelId: req.PaymentChannelID,
		TypeTransaction:  string(model.TopupTransaction),
	})
	if err != nil {
		return nil, err
	}

	err = cs.CustomerRepo.ProduceCustomerTopupTransaction(
		cs.prepareTopupTransactionMessage(customerID, data.GetTransactionId(), req))
	if err != nil {
		return nil, err
	}

	return &model.TopupWalletCustomerResponse{
		TransactionID: data.GetTransactionId(),
		Status:        data.GetStatus(),
	}, nil
}

func (cs *CustomerServiceImpl) validateTopupWalletCustomerRequest(req *model.TopupWalletCustomerRequest) error {
	if req.Amount < 1 && req.Amount > 10000000 {
		return model.NewError(model.Validation, "amount cant be less than 0 & amount cant be more than 10mio")
	}

	if req.PaymentChannelID < 1 {
		return model.NewError(model.Validation, "payment_channel_id cant be less than 0")
	}

	return nil
}

func (cs *CustomerServiceImpl) prepareTopupTransactionMessage(customerID, transactionID uint64, req *model.TopupWalletCustomerRequest) *walletpb.TopupTransactionMessage {
	return &walletpb.TopupTransactionMessage{
		TransactionId:    transactionID,
		CustomerId:       customerID,
		Amount:           req.Amount,
		PaymentChannelId: req.PaymentChannelID,
		EventName:        string(model.EventPaymentRequested),
	}
}
