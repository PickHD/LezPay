package controller

import (
	"context"

	"github.com/PickHD/LezPay/merchant/internal/v1/config"
	"github.com/PickHD/LezPay/merchant/internal/v1/model"
	"github.com/PickHD/LezPay/merchant/internal/v1/service"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	merchantpb "github.com/PickHD/LezPay/merchant/pkg/proto/v1/merchant"
)

type (
	// MerchantController is an interface that has all the function to be implemented inside Merchant controller
	MerchantController interface {
		CreateMerchant(ctx context.Context, req *merchantpb.MerchantRequest) (*merchantpb.MerchantResponse, error)
		UpdateVerifiedMerchant(ctx context.Context, req *merchantpb.UpdateVerifiedMerchantRequest) (*merchantpb.UpdateVerifiedMerchantResponse, error)
	}

	// MerchantControllerImpl is an app Merchant struct that consists of all the dependencies needed for Merchant controller
	MerchantControllerImpl struct {
		Context     context.Context
		Config      *config.Configuration
		Tracer      *trace.TracerProvider
		MerchantSvc service.MerchantService
		merchantpb.UnimplementedMerchantServiceServer
	}
)

// NewMerchantController return new instances Merchant controller
func NewMerchantController(ctx context.Context, config *config.Configuration, tracer *trace.TracerProvider, merchantSvc service.MerchantService) *MerchantControllerImpl {
	return &MerchantControllerImpl{
		Context:     ctx,
		Config:      config,
		Tracer:      tracer,
		MerchantSvc: merchantSvc,
	}
}

func (mc *MerchantControllerImpl) CreateMerchant(ctx context.Context, req *merchantpb.MerchantRequest) (*merchantpb.MerchantResponse, error) {
	tr := mc.Tracer.Tracer("Merchant-CreateMerchant Controller")
	_, span := tr.Start(ctx, "Start CreateMerchant")
	defer span.End()

	newMerchant := model.CreateMerchantRequest{
		FullName:    req.GetFullName(),
		Email:       req.GetEmail(),
		PhoneNumber: req.GetPhoneNumber(),
		Password:    req.GetPassword(),
	}

	data, err := mc.MerchantSvc.CreateMerchant(&newMerchant)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed Create Merchant %s", err.Error())
	}

	return &merchantpb.MerchantResponse{
		Id:         uint64(data.ID),
		IsVerified: data.IsVerified,
	}, nil
}

func (mc *MerchantControllerImpl) UpdateVerifiedMerchant(ctx context.Context, req *merchantpb.UpdateVerifiedMerchantRequest) (*merchantpb.UpdateVerifiedMerchantResponse, error) {
	tr := mc.Tracer.Tracer("Merchant-UpdateVerifiedMerchant Controller")
	_, span := tr.Start(ctx, "Start UpdateVerifiedMerchant")
	defer span.End()

	data, err := mc.MerchantSvc.UpdateVerifiedMerchant(req.GetEmail())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed Update Verified Merchant %s", err.Error())
	}

	return &merchantpb.UpdateVerifiedMerchantResponse{
		IsVerified: data.IsVerified,
	}, nil
}
