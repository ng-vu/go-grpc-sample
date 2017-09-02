package main

import (
	"context"
	"time"

	argPkg "github.com/alexflint/go-arg"
	"google.golang.org/grpc"

	"github.com/ng-vu/go-grpc-sample/base"
	"github.com/ng-vu/go-grpc-sample/base/l"
	"github.com/ng-vu/go-grpc-sample/blue/config"
	pb "github.com/ng-vu/go-grpc-sample/pb/partner"
)

var orders = []struct {
	OrderID      string
	OrderCode    string
	OrderTime    string // DD/MM/YYYY HH:mm
	ExpectedTime string
	Note         string

	CustomerPhone   string
	CustomerName    string
	CustomerAddress string

	TotalFee    int64
	TotalAmount int64
}{
	{
		"ID_1001", "CODE_1001", "26/07/2017 15:10", "30/07/2017 09:00", "",
		"01234560001", "Alice", "Somewhere in her home town",
		20000, 420000,
	},
	{
		"ID_1002", "CODE_1002", "25/07/2017 9:45", "28/07/2017 18:00", "",
		"01234560001", "Alice", "Somewhere in her home town",
		80000, 80000,
	},
	{
		"ID_2001", "CODE_2001", "25/07/2017 9:45", "28/07/2017 18:00", "",
		"01234560002", "Johny", "His company address",
		35000, 1035000,
	},
	{
		"ID_2002", "CODE_2002", "29/07/2017 12:30", "01/08/2017 14:00", "",
		"01234560002", "Johny", "His company address",
		15000, 15000,
	},
	{
		"ID_2003", "CODE_2003", "30/07/2017 10:30", "02/08/2017 18:00", "",
		"01234560002", "Johny", "His company address",
		40000, 40000,
	},
}

func parseTime(s string) int64 {
	t, err := time.Parse("02/01/2006 15:04", s)
	if err != nil {
		ll.Fatal("Error", l.Error(err))
	}
	return base.Millis(t)
}

func createOrders() {
	pbOrders := make([]*pb.DeliveryOrder, len(orders))
	for i, d := range orders {
		pbOrders[i] = &pb.DeliveryOrder{
			OrderId:     d.OrderID,
			TotalFee:    d.TotalFee,
			TotalAmount: d.TotalAmount,

			Info: &pb.DeliveryOrderInfo{
				OrderCode:      d.OrderCode,
				ExtraOrderCode: "",
				OrderTime:      parseTime(d.OrderTime),
				ExpectedTime:   0,
				SenderName:     d.CustomerName,
				SenderPhone:    d.CustomerPhone,
				SenderAddress: &pb.Address{
					Address: d.CustomerAddress,
				},
				ReceiverName: d.CustomerName,
			},

			Customer: &pb.OrderCustomer{
				Name:    d.CustomerName,
				Phone:   d.CustomerPhone,
				Address: d.CustomerAddress,
			},
			Location: &pb.DeliveryOrder_LocationAddress{
				LocationAddress: &pb.Address{
					Address: d.CustomerAddress,
				},
			},
			Service: &pb.DeliveryOrder_ServicePickup{
				ServicePickup: &pb.OrderPickup{
					Fee: d.TotalFee,
					Cod: d.TotalAmount - d.TotalFee,
				},
			},
		}
	}
	req := &pb.DeliveryOrdersCreateRequest{
		Partner: "green",
		Orders:  pbOrders,
	}
	resp, err := client.DeliveryOrdersCreate(ctx, req)
	if err != nil {
		ll.Error("Error", l.Error(err))
	}
	ll.Info("Response", l.Object("resp", resp))
}

var (
	cfg    *config.Config
	ll     = l.New()
	client pb.BluePartnerClient
	ctx    context.Context
)

var args = &struct {
	Address string `arg:"env:PARTNER_ADDRESS,help:GRPC address"`
	APIKey  string `arg:"env:PARTNER_APIKEY,required,help:Partner APIKEY"`
}{}

func main() {
	{
		cfg = config.Default()
		c := cfg.PartnerService
		args.Address = c.GRPC.Host + ":" + c.GRPC.Port
		args.APIKey = ""
		argPkg.MustParse(args)

		conn, err := grpc.Dial(args.Address, grpc.WithInsecure())
		if err != nil {
			ll.Fatal("Unable to connect GRPC", l.Error(err), l.String("addr", args.Address))
		}
		client = pb.NewBluePartnerClient(conn)
		ctx = base.AppendAccessToken(context.Background(), args.APIKey)
	}

	createOrders()
}
