package service

import (
	"context"

	"github.com/PickHD/LezPay/merchant/internal/v1/config"
	"github.com/PickHD/LezPay/merchant/internal/v1/model"
	"github.com/PickHD/LezPay/merchant/internal/v1/repository"
	"go.opentelemetry.io/otel/sdk/trace"
)

type (
	// MerchantService is an interface that has all the function to be implemented inside merchant service
	MerchantService interface {
		CreateMerchant(req *model.CreateMerchantRequest) (*model.CreateMerchantResponse, error)
		UpdateVerifiedMerchant(email string) (*model.UpdateVerifiedMerchantResponse, error)
		GetMerchantDetailsByEmail(req *model.GetMerchantDetailsByEmailRequest) (*model.GetMerchantDetailsByEmailResponse, error)
	}

	// MerchantServiceImpl is an app Merchant struct that consists of all the dependencies needed for merchant service
	MerchantServiceImpl struct {
		Context      context.Context
		Config       *config.Configuration
		Tracer       *trace.TracerProvider
		MerchantRepo repository.MerchantRepository
	}
)

// NewMerchantService return new instances merchant service
func NewMerchantService(ctx context.Context, config *config.Configuration, tracer *trace.TracerProvider, merchantRepo repository.MerchantRepository) *MerchantServiceImpl {
	return &MerchantServiceImpl{
		Context:      ctx,
		Config:       config,
		Tracer:       tracer,
		MerchantRepo: merchantRepo,
	}
}

func (ms *MerchantServiceImpl) CreateMerchant(req *model.CreateMerchantRequest) (*model.CreateMerchantResponse, error) {
	tr := ms.Tracer.Tracer("Merchant-CreateMerchant Service")
	_, span := tr.Start(ms.Context, "Start CreateMerchant")
	defer span.End()

	Id, isVerified, err := ms.MerchantRepo.CreateMerchant(req)
	if err != nil {
		return nil, err
	}

	return &model.CreateMerchantResponse{
		ID:         Id,
		IsVerified: isVerified,
	}, nil
}

func (ms *MerchantServiceImpl) UpdateVerifiedMerchant(email string) (*model.UpdateVerifiedMerchantResponse, error) {
	tr := ms.Tracer.Tracer("Merchant-UpdateVerifiedMerchant Service")
	_, span := tr.Start(ms.Context, "Start UpdateVerifiedMerchant")
	defer span.End()

	isVerified, err := ms.MerchantRepo.UpdateVerifiedMerchant(email)
	if err != nil {
		return nil, err
	}

	return &model.UpdateVerifiedMerchantResponse{
		IsVerified: isVerified,
	}, nil
}

func (ms *MerchantServiceImpl) GetMerchantDetailsByEmail(req *model.GetMerchantDetailsByEmailRequest) (*model.GetMerchantDetailsByEmailResponse, error) {
	tr := ms.Tracer.Tracer("Merchant-GetMerchantDetailsByEmail Service")
	_, span := tr.Start(ms.Context, "Start GetMerchantDetailsByEmail")
	defer span.End()

	getMerchant, err := ms.MerchantRepo.GetMerchantDetailsByEmail(req)
	if err != nil {
		return nil, err
	}

	return &model.GetMerchantDetailsByEmailResponse{
		ID:          getMerchant.ID,
		FullName:    getMerchant.FullName,
		PhoneNumber: getMerchant.PhoneNumber,
		Email:       getMerchant.Email,
		Password:    getMerchant.Password,
	}, nil
}
