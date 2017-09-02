package service

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	pb "github.com/ng-vu/go-grpc-sample/pb/agency"
)

// AgencyService is alias to BlueAgencyServer
type AgencyService pb.BlueAgencyServer

type agencyService struct {
	InnerService
}

// NewAgencyService return new AgencyService
func NewAgencyService(s InnerService) AgencyService {
	return &agencyService{s}
}

func (s *agencyService) VersionInfo(ctx context.Context, req *pb.Empty) (*pb.VersionInfoResponse, error) {
	return &pb.VersionInfoResponse{
		Service:     "blue/agency",
		Version:     "0.1",
		UpdatedTime: 0,
	}, nil
}

func (s *agencyService) AccountLogin(ctx context.Context, req *pb.AccountLoginRequest) (*pb.AccountLoginResponse, error) {

	if req.Phone == "" || req.Password == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "Missing phone or password")
	}

	staff, tok, err := s.UserCtrl.StaffLoginByPhone(ctx, req.Phone, req.Password)
	if err != nil {
		return nil, err
	}

	return &pb.AccountLoginResponse{
		UserId:      string(staff.ID),
		AccessToken: tok.TokenStr,
		UserInfo: &pb.AgencyStaff{
			Id:      string(staff.ID),
			Name:    string(staff.Name),
			Phone:   string(staff.Phone),
			Email:   string(staff.Email),
			Address: string(staff.Address),
		},
	}, nil
}

func (s *agencyService) AccountLogout(ctx context.Context, req *pb.AccountLogoutRequest) (*pb.AccountLogoutResponse, error) {
	return nil, grpc.Errorf(codes.Unavailable, "TODO")
}

func (s *agencyService) CustomerLookup(ctx context.Context, req *pb.CustomerLookupRequest) (*pb.CustomerLookupResponse, error) {
	orders, err := s.OrderCtrl.DeliveryOrdersGetByCustomerPhone(ctx, req.Phone)
	if err != nil {
		return nil, err
	}
	if len(orders) == 0 {
		return nil, grpc.Errorf(codes.NotFound, "Not found")
	}

	customer := &pb.Customer{
		Name:  orders[0].Customer.Name,
		Phone: orders[0].Customer.Phone,
	}
	return &pb.CustomerLookupResponse{
		Customer: customer,
		Orders:   orders,
	}, nil
}

func (s *agencyService) CustomerAction(ctx context.Context, req *pb.CustomerActionRequest) (*pb.CustomerActionResponse, error) {
	return nil, grpc.Errorf(codes.Unavailable, "TODO")
}

func (s *agencyService) ReceiveFromSupplier(ctx context.Context, req *pb.ReceiveFromSupplierRequest) (*pb.ReceiveFromSupplierResponse, error) {
	return nil, grpc.Errorf(codes.Unavailable, "TODO")
}

func (s *agencyService) TransferToSupplier(ctx context.Context, req *pb.TransferToSupplierRequest) (*pb.TransferToSupplierResponse, error) {
	return nil, grpc.Errorf(codes.Unavailable, "TODO")
}
