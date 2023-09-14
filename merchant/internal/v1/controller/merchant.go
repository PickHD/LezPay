package controller

import (
	"context"
	"strings"

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
		GetMerchantDetailsByEmail(ctx context.Context, req *merchantpb.GetMerchantDetailsByEmailRequest) (*merchantpb.GetMerchantDetailsByEmailResponse, error)
		UpdateMerchantPasswordByEmail(ctx context.Context, req *merchantpb.UpdateMerchantPasswordByEmailRequest) (*merchantpb.UpdateMerchantPasswordByEmailResponse, error)
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

func (mc *MerchantControllerImpl) GetMerchantDetailsByEmail(ctx context.Context, req *merchantpb.GetMerchantDetailsByEmailRequest) (*merchantpb.GetMerchantDetailsByEmailResponse, error) {
	tr := mc.Tracer.Tracer("Merchant-GetMerchantDetailsByEmail Controller")
	_, span := tr.Start(ctx, "Start GetMerchantDetailsByEmail")
	defer span.End()

	getMerchant := &model.GetMerchantDetailsByEmailRequest{
		Email: req.GetEmail(),
	}

	data, err := mc.MerchantSvc.GetMerchantDetailsByEmail(getMerchant)
	if err != nil {
		if strings.Contains(err.Error(), string(model.NotFound)) {
			return nil, status.Error(codes.NotFound, "Email Not Found")
		}

		return nil, status.Errorf(codes.Internal, "Failed Get Merchant Details By Email %s", err.Error())
	}

	return &merchantpb.GetMerchantDetailsByEmailResponse{
		Id:          uint64(data.ID),
		FullName:    data.FullName,
		Email:       data.Email,
		PhoneNumber: data.PhoneNumber,
		Password:    data.Password,
	}, nil
}

func (mc *MerchantControllerImpl) UpdateMerchantPasswordByEmail(ctx context.Context, req *merchantpb.UpdateMerchantPasswordByEmailRequest) (*merchantpb.UpdateMerchantPasswordByEmailResponse, error) {
	tr := mc.Tracer.Tracer("Merchant-UpdateMerchantPasswordByEmail Controller")
	_, span := tr.Start(ctx, "Start UpdateMerchantPasswordByEmail")
	defer span.End()

	updatePasswordReq := model.UpdateMerchantPasswordByEmailRequest{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	data, err := mc.MerchantSvc.UpdateMerchantPasswordByEmail(&updatePasswordReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed Update Password %s", err.Error())
	}

	return &merchantpb.UpdateMerchantPasswordByEmailResponse{
		Email: data.Email,
	}, nil
}
