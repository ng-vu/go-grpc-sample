// Code generated by protoc-gen-go.
// source: agency/model.proto
// DO NOT EDIT!

package agency

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type Order_StatusDone int32

const (
	Order_process Order_StatusDone = 0
	Order_done    Order_StatusDone = 1
	Order_cancel  Order_StatusDone = -1
)

var Order_StatusDone_name = map[int32]string{
	0:  "process",
	1:  "done",
	-1: "cancel",
}
var Order_StatusDone_value = map[string]int32{
	"process": 0,
	"done":    1,
	"cancel":  -1,
}

func (x Order_StatusDone) String() string {
	return proto.EnumName(Order_StatusDone_name, int32(x))
}
func (Order_StatusDone) EnumDescriptor() ([]byte, []int) { return fileDescriptor1, []int{2, 0} }

type AgencyStaff struct {
	Id      string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Name    string `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	Phone   string `protobuf:"bytes,3,opt,name=phone" json:"phone,omitempty"`
	Email   string `protobuf:"bytes,4,opt,name=email" json:"email,omitempty"`
	Address string `protobuf:"bytes,5,opt,name=address" json:"address,omitempty"`
}

func (m *AgencyStaff) Reset()                    { *m = AgencyStaff{} }
func (m *AgencyStaff) String() string            { return proto.CompactTextString(m) }
func (*AgencyStaff) ProtoMessage()               {}
func (*AgencyStaff) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

func (m *AgencyStaff) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *AgencyStaff) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *AgencyStaff) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *AgencyStaff) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *AgencyStaff) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

type Customer struct {
	Id    string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Name  string `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	Phone string `protobuf:"bytes,3,opt,name=phone" json:"phone,omitempty"`
}

func (m *Customer) Reset()                    { *m = Customer{} }
func (m *Customer) String() string            { return proto.CompactTextString(m) }
func (*Customer) ProtoMessage()               {}
func (*Customer) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{1} }

func (m *Customer) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Customer) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Customer) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

type Order struct {
	Id          string           `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Code        string           `protobuf:"bytes,2,opt,name=code" json:"code,omitempty"`
	ServiceId   string           `protobuf:"bytes,3,opt,name=service_id,json=serviceId" json:"service_id,omitempty"`
	Type        string           `protobuf:"bytes,4,opt,name=type" json:"type,omitempty"`
	StatusDone  Order_StatusDone `protobuf:"varint,5,opt,name=status_done,json=statusDone,enum=agency.Order_StatusDone" json:"status_done,omitempty"`
	TotalAmount int32            `protobuf:"varint,6,opt,name=total_amount,json=totalAmount" json:"total_amount,omitempty"`
	TotalFee    int32            `protobuf:"varint,7,opt,name=total_fee,json=totalFee" json:"total_fee,omitempty"`
	Info        *OrderInfo       `protobuf:"bytes,10,opt,name=info" json:"info,omitempty"`
	Customer    *OrderCustomer   `protobuf:"bytes,11,opt,name=customer" json:"customer,omitempty"`
	Service     *OrderService    `protobuf:"bytes,12,opt,name=service" json:"service,omitempty"`
}

func (m *Order) Reset()                    { *m = Order{} }
func (m *Order) String() string            { return proto.CompactTextString(m) }
func (*Order) ProtoMessage()               {}
func (*Order) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{2} }

func (m *Order) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Order) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *Order) GetServiceId() string {
	if m != nil {
		return m.ServiceId
	}
	return ""
}

func (m *Order) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *Order) GetStatusDone() Order_StatusDone {
	if m != nil {
		return m.StatusDone
	}
	return Order_process
}

func (m *Order) GetTotalAmount() int32 {
	if m != nil {
		return m.TotalAmount
	}
	return 0
}

func (m *Order) GetTotalFee() int32 {
	if m != nil {
		return m.TotalFee
	}
	return 0
}

func (m *Order) GetInfo() *OrderInfo {
	if m != nil {
		return m.Info
	}
	return nil
}

func (m *Order) GetCustomer() *OrderCustomer {
	if m != nil {
		return m.Customer
	}
	return nil
}

func (m *Order) GetService() *OrderService {
	if m != nil {
		return m.Service
	}
	return nil
}

type OrderInfo struct {
	ProviderCode string `protobuf:"bytes,1,opt,name=provider_code,json=providerCode" json:"provider_code,omitempty"`
	CustomerCode string `protobuf:"bytes,2,opt,name=customer_code,json=customerCode" json:"customer_code,omitempty"`
	OrderTime    int64  `protobuf:"varint,3,opt,name=order_time,json=orderTime" json:"order_time,omitempty"`
	ExpectedTime int64  `protobuf:"varint,4,opt,name=expected_time,json=expectedTime" json:"expected_time,omitempty"`
	Note         string `protobuf:"bytes,10,opt,name=note" json:"note,omitempty"`
}

func (m *OrderInfo) Reset()                    { *m = OrderInfo{} }
func (m *OrderInfo) String() string            { return proto.CompactTextString(m) }
func (*OrderInfo) ProtoMessage()               {}
func (*OrderInfo) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{3} }

func (m *OrderInfo) GetProviderCode() string {
	if m != nil {
		return m.ProviderCode
	}
	return ""
}

func (m *OrderInfo) GetCustomerCode() string {
	if m != nil {
		return m.CustomerCode
	}
	return ""
}

func (m *OrderInfo) GetOrderTime() int64 {
	if m != nil {
		return m.OrderTime
	}
	return 0
}

func (m *OrderInfo) GetExpectedTime() int64 {
	if m != nil {
		return m.ExpectedTime
	}
	return 0
}

func (m *OrderInfo) GetNote() string {
	if m != nil {
		return m.Note
	}
	return ""
}

type OrderCustomer struct {
	Id      string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Name    string `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	Phone   string `protobuf:"bytes,3,opt,name=phone" json:"phone,omitempty"`
	Email   string `protobuf:"bytes,4,opt,name=email" json:"email,omitempty"`
	Address string `protobuf:"bytes,5,opt,name=address" json:"address,omitempty"`
}

func (m *OrderCustomer) Reset()                    { *m = OrderCustomer{} }
func (m *OrderCustomer) String() string            { return proto.CompactTextString(m) }
func (*OrderCustomer) ProtoMessage()               {}
func (*OrderCustomer) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{4} }

func (m *OrderCustomer) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *OrderCustomer) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *OrderCustomer) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *OrderCustomer) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *OrderCustomer) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

type OrderService struct {
	Fee string `protobuf:"bytes,1,opt,name=fee" json:"fee,omitempty"`
	Cod string `protobuf:"bytes,2,opt,name=cod" json:"cod,omitempty"`
}

func (m *OrderService) Reset()                    { *m = OrderService{} }
func (m *OrderService) String() string            { return proto.CompactTextString(m) }
func (*OrderService) ProtoMessage()               {}
func (*OrderService) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{5} }

func (m *OrderService) GetFee() string {
	if m != nil {
		return m.Fee
	}
	return ""
}

func (m *OrderService) GetCod() string {
	if m != nil {
		return m.Cod
	}
	return ""
}

type Service struct {
	Id          string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Name        string `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	Description string `protobuf:"bytes,3,opt,name=description" json:"description,omitempty"`
}

func (m *Service) Reset()                    { *m = Service{} }
func (m *Service) String() string            { return proto.CompactTextString(m) }
func (*Service) ProtoMessage()               {}
func (*Service) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{6} }

func (m *Service) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Service) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Service) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func init() {
	proto.RegisterType((*AgencyStaff)(nil), "agency.AgencyStaff")
	proto.RegisterType((*Customer)(nil), "agency.Customer")
	proto.RegisterType((*Order)(nil), "agency.Order")
	proto.RegisterType((*OrderInfo)(nil), "agency.OrderInfo")
	proto.RegisterType((*OrderCustomer)(nil), "agency.OrderCustomer")
	proto.RegisterType((*OrderService)(nil), "agency.OrderService")
	proto.RegisterType((*Service)(nil), "agency.Service")
	proto.RegisterEnum("agency.Order_StatusDone", Order_StatusDone_name, Order_StatusDone_value)
}

func init() { proto.RegisterFile("agency/model.proto", fileDescriptor1) }

var fileDescriptor1 = []byte{
	// 529 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xbc, 0x94, 0x4d, 0xab, 0xd3, 0x4c,
	0x14, 0xc7, 0x9f, 0xf4, 0xbd, 0x27, 0xe9, 0xa5, 0xcf, 0x78, 0x85, 0x01, 0x11, 0x6a, 0x2e, 0x42,
	0x45, 0x9a, 0x62, 0xdd, 0xe8, 0xf2, 0xda, 0x8b, 0x70, 0x57, 0x17, 0x52, 0x57, 0x6e, 0x42, 0x3a,
	0x73, 0xda, 0x3b, 0xd0, 0x64, 0x42, 0x32, 0xad, 0xf6, 0xc3, 0xf8, 0x11, 0xfc, 0x8c, 0xca, 0x9c,
	0x49, 0x6a, 0x0b, 0x2e, 0xc4, 0x85, 0x5d, 0x9d, 0xf9, 0xcd, 0x7f, 0xce, 0x7b, 0x03, 0x2c, 0xdd,
	0x62, 0x2e, 0x8e, 0xf3, 0x4c, 0x4b, 0xdc, 0x45, 0x45, 0xa9, 0x8d, 0x66, 0x3d, 0xc7, 0xc2, 0x2f,
	0xe0, 0xdf, 0x92, 0xb5, 0x32, 0xe9, 0x66, 0xc3, 0xae, 0xa0, 0xa5, 0x24, 0xf7, 0x26, 0xde, 0x74,
	0x18, 0xb7, 0x94, 0x64, 0x0c, 0x3a, 0x79, 0x9a, 0x21, 0x6f, 0x11, 0x21, 0x9b, 0x5d, 0x43, 0xb7,
	0x78, 0xd4, 0x39, 0xf2, 0x36, 0x41, 0x77, 0xb0, 0x14, 0xb3, 0x54, 0xed, 0x78, 0xc7, 0x51, 0x3a,
	0x30, 0x0e, 0xfd, 0x54, 0xca, 0x12, 0xab, 0x8a, 0x77, 0x89, 0x37, 0xc7, 0xf0, 0x0e, 0x06, 0xcb,
	0x7d, 0x65, 0x74, 0x86, 0xe5, 0xdf, 0x47, 0x0d, 0xbf, 0xb5, 0xa1, 0xfb, 0x50, 0xca, 0xdf, 0xfb,
	0x10, 0x5a, 0x9e, 0x7c, 0x58, 0x9b, 0x3d, 0x07, 0xa8, 0xb0, 0x3c, 0x28, 0x81, 0x89, 0x92, 0xb5,
	0xa3, 0x61, 0x4d, 0xee, 0xe9, 0x89, 0x39, 0x16, 0x58, 0x57, 0x40, 0x36, 0x7b, 0x0f, 0x7e, 0x65,
	0x52, 0xb3, 0xaf, 0x12, 0x69, 0x83, 0xdb, 0x22, 0xae, 0x16, 0x3c, 0x72, 0xdd, 0x8b, 0x28, 0x74,
	0xb4, 0x22, 0xc1, 0x9d, 0xce, 0x31, 0x86, 0xea, 0x64, 0xb3, 0x17, 0x10, 0x18, 0x6d, 0xd2, 0x5d,
	0x92, 0x66, 0x7a, 0x9f, 0x1b, 0xde, 0x9b, 0x78, 0xd3, 0x6e, 0xec, 0x13, 0xbb, 0x25, 0xc4, 0x9e,
	0xc1, 0xd0, 0x49, 0x36, 0x88, 0xbc, 0x4f, 0xf7, 0x03, 0x02, 0x1f, 0x11, 0xd9, 0x4b, 0xe8, 0xa8,
	0x7c, 0xa3, 0x39, 0x4c, 0xbc, 0xa9, 0xbf, 0xf8, 0xff, 0x22, 0xe6, 0x7d, 0xbe, 0xd1, 0x31, 0x5d,
	0xb3, 0x37, 0x30, 0x10, 0x75, 0x23, 0xb9, 0x4f, 0xd2, 0xa7, 0x17, 0xd2, 0xa6, 0xcb, 0xf1, 0x49,
	0xc6, 0x22, 0xe8, 0xd7, 0x55, 0xf3, 0x80, 0x5e, 0x5c, 0x5f, 0xbc, 0x58, 0xb9, 0xbb, 0xb8, 0x11,
	0x85, 0xef, 0x00, 0x7e, 0xd5, 0xc8, 0x7c, 0xe8, 0x17, 0xa5, 0x16, 0x58, 0x55, 0xe3, 0xff, 0xd8,
	0x00, 0x3a, 0xb6, 0x31, 0x63, 0x8f, 0x3d, 0x81, 0x9e, 0x48, 0x73, 0x81, 0xbb, 0xf1, 0x8f, 0xe6,
	0xe7, 0x85, 0xdf, 0x3d, 0x18, 0x9e, 0x12, 0x66, 0x37, 0x30, 0x2a, 0x4a, 0x7d, 0x50, 0x12, 0xcb,
	0x84, 0x86, 0xe3, 0xc6, 0x15, 0x34, 0x70, 0x69, 0x87, 0x74, 0x03, 0xa3, 0x26, 0xd1, 0xe4, 0x6c,
	0x82, 0x41, 0x03, 0x97, 0xf5, 0x24, 0xb5, 0x75, 0x9b, 0x18, 0x95, 0xb9, 0x95, 0x68, 0xc7, 0x43,
	0x22, 0x9f, 0x54, 0x46, 0x3e, 0xf0, 0x6b, 0x81, 0xc2, 0xa0, 0x74, 0x8a, 0x0e, 0x29, 0x82, 0x06,
	0x92, 0xc8, 0x6e, 0x99, 0x36, 0x48, 0xfd, 0xb5, 0x5b, 0xa6, 0x0d, 0x86, 0x47, 0x18, 0x5d, 0x34,
	0xed, 0x1f, 0xfe, 0x21, 0x16, 0x10, 0x9c, 0x77, 0x9f, 0x8d, 0xa1, 0x6d, 0xb7, 0xc2, 0x85, 0xb6,
	0xa6, 0x25, 0x42, 0xcb, 0x3a, 0xb4, 0x35, 0xc3, 0x07, 0xe8, 0x37, 0xf2, 0x3f, 0x49, 0x74, 0x02,
	0xbe, 0xc4, 0x4a, 0x94, 0xaa, 0x30, 0x4a, 0xe7, 0x75, 0xba, 0xe7, 0xe8, 0xc3, 0xeb, 0xcf, 0xaf,
	0xb6, 0xca, 0x3c, 0xee, 0xd7, 0x91, 0xd0, 0xd9, 0x3c, 0xdf, 0xce, 0x0e, 0xfb, 0xf9, 0x56, 0xcf,
	0xb6, 0x65, 0x21, 0x66, 0x55, 0x9a, 0x15, 0x3b, 0x9c, 0x17, 0xeb, 0xb9, 0x5b, 0x96, 0x75, 0x8f,
	0x3e, 0x25, 0x6f, 0x7f, 0x06, 0x00, 0x00, 0xff, 0xff, 0xe3, 0xfd, 0xf7, 0x7f, 0x60, 0x04, 0x00,
	0x00,
}