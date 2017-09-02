package ctrl

import (
	"context"

	google_rpc "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/ng-vu/go-grpc-sample/base"
	"github.com/ng-vu/go-grpc-sample/base/l"
	"github.com/ng-vu/go-grpc-sample/base/validate"
	pbAgency "github.com/ng-vu/go-grpc-sample/pb/agency"
	pbPartner "github.com/ng-vu/go-grpc-sample/pb/partner"
	"github.com/ng-vu/go-grpc-sample/blue/internal/auth"
	"github.com/ng-vu/go-grpc-sample/blue/internal/model"
	"github.com/ng-vu/go-grpc-sample/blue/internal/store"
)

// OrderCtrl ...
type OrderCtrl struct {
	OrderStore            *store.OrderStore            `inject:""`
	OrderTransactionStore *store.OrderTransactionStore `inject:""`
}

// DeliveryOrdersCreate ...
func (c *OrderCtrl) DeliveryOrdersCreate(ctx context.Context, sp auth.ServiceProviderClaim, req *pbPartner.DeliveryOrdersCreateRequest) (*pbPartner.DeliveryOrdersCreateResponse, error) {

	orders := make([]*model.Order, len(req.Orders))
	for i, pbOrder := range req.Orders {
		if pbOrder.Info == nil || pbOrder.Customer == nil || pbOrder.Location == nil || pbOrder.Service == nil {
			ll.Error("Invalid order info", l.Object("order", pbOrder))
			return nil, grpc.Errorf(codes.InvalidArgument, "Invalid order info")
		}
		order := &model.Order{
			ProviderID:        sp.ID,
			ProviderOrderID:   model.String(pbOrder.OrderId),
			ProviderOrderCode: model.String(pbOrder.Info.OrderCode),

			OrderTime:    base.FromMillis(pbOrder.Info.OrderTime),
			ExpectedTime: base.FromMillisP(pbOrder.Info.ExpectedTime),
			Note:         model.String(pbOrder.Info.Note),

			CustomerReferenceID: model.String(pbOrder.Customer.ReferenceId),
			CustomerName:        model.String(pbOrder.Customer.Name),
			CustomerPhone:       model.String(pbOrder.Customer.Phone),
			CustomerEmail:       model.String(pbOrder.Customer.Email),
			CustomerAddress:     model.String(pbOrder.Customer.Address),

			TotalFee:    int32(pbOrder.TotalFee),
			TotalAmount: int32(pbOrder.TotalAmount),
		}
		orders[i] = order
	}

	errors := c.OrderStore.BulkCreate(orders)
	mapSuccess := make(map[string]string)
	mapError := make(map[string]*google_rpc.Status)
	for i, err := range errors {
		order := orders[i]
		if err == nil {
			mapSuccess[string(order.ID)] = "ok"
		} else {
			mapError[string(order.ID)] = &google_rpc.Status{
				Code:    int32(codes.Unknown),
				Message: err.Error(),
			}
		}
	}

	if len(mapSuccess) == 0 {
		return nil, grpc.Errorf(codes.Unknown, "Unable to create orders")
	}
	if len(mapError) > 0 {
		ll.Error("Create orders", l.Object("errors", mapError))
	}
	return &pbPartner.DeliveryOrdersCreateResponse{
		OrderErrors:  mapError,
		OrderSuccess: mapSuccess,
	}, nil
}

// DeliveryOrdersGetByCustomerPhone ...
func (c *OrderCtrl) DeliveryOrdersGetByCustomerPhone(ctx context.Context, inputPhone string) ([]*pbAgency.Order, error) {
	if inputPhone == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "Phone required")
	}

	phone, ok := validate.NormalizePhone(inputPhone)
	if !ok {
		ll.Error("Invalid phone number", l.String("inputPhone", inputPhone))
		return nil, grpc.Errorf(codes.InvalidArgument, "Invalid phone")
	}

	ll.Info("Query orders by phone", l.String("phone", phone))
	orders, err := c.OrderStore.GetByCustomerPhone(phone)
	if err != nil {
		ll.Error("Order.GetByCustomerPhone", l.Error(err), l.String("phone", phone))
		return nil, grpc.Errorf(codes.NotFound, err.Error())
	}

	pbOrders := make([]*pbAgency.Order, len(orders))
	for i, order := range orders {
		pbOrders[i] = &pbAgency.Order{
			Id:         string(order.ID),
			Code:       string(order.Code),
			ServiceId:  "",
			Type:       "",
			StatusDone: 0,
			Info: &pbAgency.OrderInfo{
				ProviderCode: string(order.ProviderOrderCode),
				CustomerCode: string(order.ExtraOrderCode),
				OrderTime:    base.Millis(order.OrderTime),
				ExpectedTime: base.MillisP(order.ExpectedTime),
			},
			Customer: &pbAgency.OrderCustomer{
				Name:    string(order.CustomerName),
				Phone:   string(order.CustomerPhone),
				Address: string(order.CustomerAddress),
				Email:   string(order.CustomerEmail),
			},
			TotalAmount: order.TotalAmount,
			TotalFee:    order.TotalFee,
		}
	}
	return pbOrders, nil
}
