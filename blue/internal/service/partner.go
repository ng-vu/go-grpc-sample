package service

import (
	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/ng-vu/go-grpc-sample/base/l"
	"github.com/ng-vu/go-grpc-sample/blue/internal/auth"
	pb "github.com/ng-vu/go-grpc-sample/pb/partner"
)

// PartnerService ...
type PartnerService pb.BluePartnerServer

type partnerService struct {
	InnerService
}

// NewPartnerService returns new PartnerService
func NewPartnerService(s InnerService) PartnerService {
	return &partnerService{s}
}

// Implement auth.Validator interface
func (s *partnerService) Validate(tokenStr string) (auth.ServiceProviderClaim, error) {
	sp, err := s.ProviderCtrl.ValidateAPIKey(tokenStr)
	return sp, err
}

func (s *partnerService) VersionInfo(context.Context, *pb.Empty) (*pb.VersionInfoResponse, error) {
	return &pb.VersionInfoResponse{
		Service:     "blue/partner",
		Version:     "0.1",
		UpdatedTime: 0,
	}, nil
}

func (s *partnerService) AccountInfo(ctx context.Context, _ *pb.Empty) (*pb.AccountInfoResponse, error) {
	sp, ok := auth.ProviderFromContext(ctx)
	if !ok {
		ll.Error("Unexpected context value", l.Object("sp", sp))
		return nil, grpc.Errorf(codes.Internal, "Unexpected")
	}
	return &pb.AccountInfoResponse{
		Id:       string(sp.ID),
		Codename: sp.Codename,
		Name:     sp.Name,
	}, nil
}

func (s *partnerService) DeliveryOrdersCreate(ctx context.Context, req *pb.DeliveryOrdersCreateRequest) (*pb.DeliveryOrdersCreateResponse, error) {
	sp, ok := auth.ProviderFromContext(ctx)
	if !ok {
		ll.Error("Unexpected context value", l.Object("sp", sp))
		return nil, grpc.Errorf(codes.Internal, "Unexpected")
	}
	if sp.Codename != req.Partner {
		ll.Error("Invalid partner")
		return nil, grpc.Errorf(codes.PermissionDenied, "Permission denied")
	}
	return s.OrderCtrl.DeliveryOrdersCreate(ctx, sp, req)
}

func (s *partnerService) DeliveryOrdersUpdateStatus(ctx context.Context, req *pb.DeliveryOrdersUpdateStatusRequest) (*pb.DeliveryOrdersUpdateStatusResponse, error) {
	return nil, ErrTODO
}

func (s *partnerService) DeliveryOrdersCancel(ctx context.Context, req *pb.DeliveryOrdersCancelRequest) (*pb.DeliveryOrdersCancelResponse, error) {
	return nil, ErrTODO
}
