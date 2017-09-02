package main

import (
	"context"
	"fmt"

	argPkg "github.com/alexflint/go-arg"
	"google.golang.org/grpc"

	"github.com/ng-vu/go-grpc-sample/base"
	"github.com/ng-vu/go-grpc-sample/base/l"
	"github.com/ng-vu/go-grpc-sample/blue/config"
	pbAgency "github.com/ng-vu/go-grpc-sample/pb/agency"
	pb "github.com/ng-vu/go-grpc-sample/pb/sadmin"
)

var serviceProviders = []struct {
	Codename string
	Name     string
}{{
	Codename: "green",
	Name:     "GreenProvider",
}}

func addServiceProviders() {
	for _, d := range serviceProviders {
		req := &pb.CreateServiceProviderRequest{
			Codename: d.Codename,
			Name:     d.Name,
		}

		resp, err := client.CreateServiceProvider(ctx, req)
		if err == nil {
			ll.Info("Created provider", l.Object("resp", resp))
		} else {
			ll.Error("Error creating provider", l.String("name", req.Name), l.Error(err))
		}
	}
}

var agencyStaffs = []struct {
	Name     string
	Phone    string
	Password string
	Email    string
	Address  string
}{
	// Name, Phone, Password, Email, Address
	{"Đà Lạt 1", "0123456001", "123456", "", "Thung lũng Đà Lạt"},
	{"Lữ Gia 2", "0123456002", "123456", "", "70 Lữ Gia, Q.10, HCM"},
}

func addAgencyStaffs() {
	for _, d := range agencyStaffs {
		req := &pb.CreateAgencyStaffRequest{
			Password: d.Password,
			Info: &pbAgency.AgencyStaff{
				Name:    d.Name,
				Phone:   d.Phone,
				Address: d.Address,
				Email:   d.Email,
			},
		}

		resp, err := client.CreateAgencyStaff(ctx, req)
		if err == nil {
			ll.Info(fmt.Sprintf("Created agency: %v - %v", resp.UserId, req.Info.Name))
		} else {
			ll.Error("Error creating agency", l.String("name", req.Info.Name), l.Error(err))
		}
	}
}

var args = &struct {
	Address string
	Token   string
}{}

var (
	cfg = config.Default()
	ll  = l.New()
	ctx context.Context

	client pb.BlueSAdminClient
)

func main() {
	{
		cfg = config.Default()
		c := cfg.SAdminService
		args.Address = c.GRPC.Host + ":" + c.GRPC.Port
		args.Token = c.MagicToken
		argPkg.MustParse(args)

		conn, err := grpc.Dial(args.Address, grpc.WithInsecure())
		if err != nil {
			ll.Fatal("Unable to connect GRPC", l.Error(err), l.String("addr", args.Address))
		}
		client = pb.NewBlueSAdminClient(conn)
		ctx = base.AppendAccessToken(context.Background(), args.Token)
	}

	addServiceProviders()
	addAgencyStaffs()
}
