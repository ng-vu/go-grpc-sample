package service

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/ng-vu/go-grpc-sample/blue/internal/model"
	pb "github.com/ng-vu/go-grpc-sample/pb/sadmin"
)

// SAdminService ...
type SAdminService pb.BlueSAdminServer

type sadminService struct {
	InnerService
}

// NewSAdminService returns new SAdminService
func NewSAdminService(s InnerService) SAdminService {
	return &sadminService{s}
}

func (s *sadminService) VersionInfo(context.Context, *pb.Empty) (*pb.VersionInfoResponse, error) {
	return &pb.VersionInfoResponse{
		Service:     "blue/internal",
		Version:     "0.1",
		UpdatedTime: 0,
	}, nil
}

func (s *sadminService) CreateAgencyStaff(ctx context.Context, req *pb.CreateAgencyStaffRequest) (*pb.CreateAgencyStaffResponse, error) {
	info := req.Info
	if info == nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "Missing info")
	}
	staff := &model.AgencyStaff{
		Name:    model.String(info.Name),
		Email:   model.String(info.Email),
		Phone:   model.String(info.Phone),
		Address: model.String(info.Address),
	}
	userID, err := s.UserCtrl.CreateStaff(staff, req.Password)
	return &pb.CreateAgencyStaffResponse{
		UserId: string(userID),
	}, err
}

func (s *sadminService) CreateServiceProvider(ctx context.Context, req *pb.CreateServiceProviderRequest) (*pb.CreateServiceProviderResponse, error) {
	sp := &model.ServiceProvider{
		Codename: model.String(req.Codename),
		Name:     model.String(req.Name),
	}
	id, apikey, err := s.ProviderCtrl.Create(sp)
	return &pb.CreateServiceProviderResponse{
		Id:     string(id),
		ApiKey: apikey,
	}, err
}

func (s *sadminService) CreateService(ctx context.Context, req *pb.CreateServiceRequest) (*pb.CreateServiceResponse, error) {
	return nil, ErrTODO
}

func (s *sadminService) GenerateServiceProviderSecret(ctx context.Context, req *pb.GenerateServiceProviderSecretRequest) (*pb.GenerateServiceProviderSecretResponse, error) {
	return nil, ErrTODO
}

func (s *sadminService) AddWebhook(ctx context.Context, req *pb.AddWebhookRequest) (*pb.AddWebhookResponse, error) {
	return nil, ErrTODO
}

func (s *sadminService) RemoveWebhook(ctx context.Context, req *pb.RemoveWebhookRequest) (*pb.RemoveWebhookResponse, error) {
	return nil, ErrTODO
}
